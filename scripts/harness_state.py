#!/usr/bin/env python3
"""Shared state engine for the emerging orchestration harness."""

from __future__ import annotations

import json
import re
import subprocess
import sys
from contextlib import contextmanager
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

import fcntl


def now_iso() -> str:
    return datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")


def now_dt() -> datetime:
    return datetime.now(timezone.utc)


def parse_iso(value: str) -> datetime | None:
    if not value:
        return None
    try:
        return datetime.fromisoformat(value.replace("Z", "+00:00"))
    except ValueError:
        return None


def repo_root() -> Path:
    script_dir = Path(__file__).resolve().parent
    try:
        output = subprocess.check_output(
            ["git", "rev-parse", "--show-toplevel"],
            cwd=script_dir.parent,
            stderr=subprocess.DEVNULL,
            text=True,
        ).strip()
        if output:
            return Path(output)
    except Exception:
        pass
    return script_dir.parent


def definitions_root() -> Path:
    return repo_root() / "harness"


def runtime_root() -> Path:
    return repo_root() / ".harness"


def state_lock_path() -> Path:
    return runtime_root() / ".lock"


def ensure_runtime_dirs() -> None:
    root = runtime_root()
    dirs = [
        root / "runs",
        root / "signals",
        root / "contracts",
        root / "meetings" / "active",
        root / "meetings" / "archive",
        root / "rc" / "active",
        root / "rc" / "archive",
        root / "scheduler" / "snapshots",
        root / "maintenance",
        root / "cache",
    ]
    for path in dirs:
        path.mkdir(parents=True, exist_ok=True)


@contextmanager
def harness_lock() -> Any:
    ensure_runtime_dirs()
    lock_path = state_lock_path()
    lock_path.parent.mkdir(parents=True, exist_ok=True)
    with lock_path.open("a+", encoding="utf-8") as handle:
        fcntl.flock(handle.fileno(), fcntl.LOCK_EX)
        try:
            yield
        finally:
            fcntl.flock(handle.fileno(), fcntl.LOCK_UN)


def _strip_quotes(value: str) -> str:
    if len(value) >= 2 and value[0] == value[-1] and value[0] in ("'", '"'):
        return value[1:-1]
    return value


def parse_scalar(value: str) -> Any:
    value = value.strip()
    if value == "":
        return ""
    if value == "[]":
        return []
    if value == "{}":
        return {}
    if value in ("true", "false"):
        return value == "true"
    if value in ("null", "~"):
        return None
    if value.startswith(("'", '"')) and value.endswith(("'", '"')):
        return _strip_quotes(value)
    try:
        if value.startswith("0") and value != "0":
            raise ValueError
        return int(value)
    except ValueError:
        return value


def parse_simple_yaml(text: str) -> Any:
    lines = []
    for raw in text.splitlines():
        if not raw.strip():
            continue
        if raw.lstrip().startswith("#"):
            continue
        indent = len(raw) - len(raw.lstrip(" "))
        content = raw.strip()
        lines.append((indent, content))

    def parse_block(index: int, expected_indent: int) -> tuple[Any, int]:
        if index >= len(lines):
            return {}, index

        indent, content = lines[index]
        if indent < expected_indent:
            return {}, index

        if content.startswith("- "):
            items = []
            while index < len(lines):
                indent, content = lines[index]
                if indent != expected_indent or not content.startswith("- "):
                    break
                item = content[2:].strip()
                index += 1
                if item == "":
                    child, index = parse_block(index, expected_indent + 2)
                    items.append(child)
                elif ":" in item and not item.startswith(("'", '"')):
                    key, raw_val = item.split(":", 1)
                    raw_val = raw_val.strip()
                    entry = {key.strip(): parse_scalar(raw_val)} if raw_val else {}
                    if not raw_val:
                        child, index = parse_block(index, expected_indent + 2)
                        if isinstance(child, dict):
                            entry.update(child)
                    items.append(entry)
                else:
                    items.append(parse_scalar(item))
            return items, index

        mapping: dict[str, Any] = {}
        while index < len(lines):
            indent, content = lines[index]
            if indent != expected_indent or content.startswith("- "):
                break
            if ":" not in content:
                index += 1
                continue
            key, raw_val = content.split(":", 1)
            key = key.strip()
            raw_val = raw_val.strip()
            index += 1
            if raw_val:
                mapping[key] = parse_scalar(raw_val)
            else:
                if index < len(lines) and lines[index][0] > expected_indent:
                    child, index = parse_block(index, expected_indent + 2)
                    mapping[key] = child
                else:
                    mapping[key] = {}
        return mapping, index

    parsed, _ = parse_block(0, 0)
    return parsed


def load_data(path: Path, default: Any) -> Any:
    if not path.exists():
        return default
    text = path.read_text(encoding="utf-8").strip()
    if not text:
        return default
    try:
        return json.loads(text)
    except json.JSONDecodeError:
        return parse_simple_yaml(text)


def dump_data(path: Path, data: Any) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(data, indent=2, ensure_ascii=True) + "\n", encoding="utf-8")


def run_command(cmd: list[str]) -> str:
    return subprocess.check_output(cmd, cwd=repo_root(), stderr=subprocess.DEVNULL, text=True).strip()


def run_command_with_stderr(cmd: list[str], cwd: Path | None = None) -> str:
    return subprocess.check_output(cmd, cwd=cwd or repo_root(), stderr=subprocess.STDOUT, text=True).strip()


def slugify(value: str) -> str:
    slug = re.sub(r"[^a-zA-Z0-9]+", "-", value.strip().lower()).strip("-")
    slug = re.sub(r"-{2,}", "-", slug)
    return slug[:48] or "mission"


def branch_from_pattern(pattern: str, issue_number: str, slug: str) -> str:
    return pattern.replace("{issue}", str(issue_number)).replace("{slug}", slug)


def merge_dicts(base: Any, override: Any) -> Any:
    if isinstance(base, dict) and isinstance(override, dict):
        result = dict(base)
        for key, value in override.items():
            if key in result:
                result[key] = merge_dicts(result[key], value)
            else:
                result[key] = value
        return result
    return override


def fill_if_blank(mapping: dict[str, Any], key: str, value: Any) -> None:
    if not mapping.get(key):
        mapping[key] = value


def get_key(data: Any, key: str, default: Any = "") -> Any:
    current = data
    for part in key.split("."):
        if not isinstance(current, dict) or part not in current:
            return default
        current = current[part]
    return current


def read_stdin_json() -> dict[str, Any]:
    raw = sys.stdin.read().strip()
    if not raw:
        return {}
    return json.loads(raw)


def profile_paths(name: str | None = None) -> list[Path]:
    paths = []
    local_override = runtime_root() / "profile.local.yml"
    if local_override.exists():
        paths.append(local_override)
    profiles_dir = definitions_root() / "profiles"
    if name:
        candidate = profiles_dir / f"{name}.yml"
        if candidate.exists():
            paths.append(candidate)
    else:
        yml_files = sorted(profiles_dir.glob("*.yml"))
        if len(yml_files) == 1:
            paths.append(yml_files[0])
        default_hint = profiles_dir / "default.yml"
        if default_hint.exists():
            paths.append(default_hint)
    return paths


def resolve_profile(name: str | None = None) -> tuple[dict[str, Any], Path | None]:
    profiles_dir = definitions_root() / "profiles"
    base_path = None
    if name:
        explicit = profiles_dir / f"{name}.yml"
        if explicit.exists():
            base_path = explicit
    else:
        yml_files = sorted(profiles_dir.glob("*.yml"))
        if len(yml_files) == 1:
            base_path = yml_files[0]
        elif (profiles_dir / "default.yml").exists():
            base_path = profiles_dir / "default.yml"

    base = load_data(base_path, {}) if base_path else {}
    local_override = load_data(runtime_root() / "profile.local.yml", {})
    merged = merge_dicts(base, local_override)
    return merged, base_path


def registry_path() -> Path:
    return runtime_root() / "runs" / "registry.yml"


def signal_log_path() -> Path:
    return runtime_root() / "signals" / "log.yml"


def contract_template_path() -> Path:
    return definitions_root() / "templates" / "contract.yml"


def run_template_path() -> Path:
    return definitions_root() / "templates" / "run.yml"


def signal_template_path() -> Path:
    return definitions_root() / "templates" / "signal.yml"


def meeting_template_path() -> Path:
    return definitions_root() / "templates" / "meeting.yml"


def rc_template_path() -> Path:
    return definitions_root() / "templates" / "rc.yml"


def issue_template_path() -> Path:
    return definitions_root() / "templates" / "issue.md"


def pr_template_path() -> Path:
    return definitions_root() / "templates" / "pr.md"


def policies_path() -> Path:
    return definitions_root() / "policies" / "defaults.yml"


def maintenance_backlog_path() -> Path:
    return runtime_root() / "maintenance" / "backlog.yml"


def meetings_active_root() -> Path:
    return runtime_root() / "meetings" / "active"


def meetings_archive_root() -> Path:
    return runtime_root() / "meetings" / "archive"


def rc_active_root() -> Path:
    return runtime_root() / "rc" / "active"


def rc_archive_root() -> Path:
    return runtime_root() / "rc" / "archive"


def meeting_path(meeting_id: str, archived: bool = False) -> Path:
    root = meetings_archive_root() if archived else meetings_active_root()
    return root / f"{meeting_id}.yml"


def rc_path(rc_id: str, archived: bool = False) -> Path:
    root = rc_archive_root() if archived else rc_active_root()
    return root / f"{rc_id}.yml"


def read_registry() -> dict[str, Any]:
    return load_data(registry_path(), {"version": 1, "updated_at": now_iso(), "runs": []})


def write_registry(data: dict[str, Any]) -> None:
    data["updated_at"] = now_iso()
    dump_data(registry_path(), data)


def read_signal_log() -> dict[str, Any]:
    return load_data(signal_log_path(), {"version": 1, "updated_at": now_iso(), "signals": []})


def write_signal_log(data: dict[str, Any]) -> None:
    data["updated_at"] = now_iso()
    dump_data(signal_log_path(), data)


def read_maintenance_backlog() -> dict[str, Any]:
    return load_data(
        maintenance_backlog_path(),
        {"version": 1, "updated_at": now_iso(), "meetings": [], "tasks": []},
    )


def write_maintenance_backlog(data: dict[str, Any]) -> None:
    data["updated_at"] = now_iso()
    dump_data(maintenance_backlog_path(), data)


def read_meeting(meeting_id: str) -> dict[str, Any] | None:
    active_path = meeting_path(meeting_id)
    if active_path.exists():
        return load_data(active_path, {})
    archived_path = meeting_path(meeting_id, archived=True)
    if archived_path.exists():
        return load_data(archived_path, {})
    return None


def read_active_meeting(meeting_id: str) -> dict[str, Any] | None:
    active_path = meeting_path(meeting_id)
    if active_path.exists():
        return load_data(active_path, {})
    return None


def read_rc(rc_id: str) -> dict[str, Any] | None:
    active_path = rc_path(rc_id)
    if active_path.exists():
        return load_data(active_path, {})
    archived_path = rc_path(rc_id, archived=True)
    if archived_path.exists():
        return load_data(archived_path, {})
    return None


def read_active_rc(rc_id: str) -> dict[str, Any] | None:
    active_path = rc_path(rc_id)
    if active_path.exists():
        return load_data(active_path, {})
    return None


def list_meetings(include_archived: bool = False) -> list[dict[str, Any]]:
    ensure_runtime_dirs()
    paths = sorted(meetings_active_root().glob("*.yml"))
    if include_archived:
        paths.extend(sorted(meetings_archive_root().glob("*.yml")))
    meetings = []
    for path in paths:
        data = load_data(path, {})
        meeting = data.get("meeting", {})
        if meeting:
            meetings.append(meeting)
    meetings.sort(key=lambda item: item.get("created_at", ""), reverse=True)
    return meetings


def list_rcs(include_archived: bool = False) -> list[dict[str, Any]]:
    ensure_runtime_dirs()
    paths = sorted(rc_active_root().glob("*.yml"))
    if include_archived:
        paths.extend(sorted(rc_archive_root().glob("*.yml")))
    rcs = []
    for path in paths:
        data = load_data(path, {})
        rc = data.get("rc", {})
        if rc:
            rcs.append(rc)
    rcs.sort(key=lambda item: item.get("created_at", ""), reverse=True)
    return rcs


def active_rc() -> dict[str, Any] | None:
    rcs = [rc for rc in list_rcs() if rc.get("state") == "active"]
    return rcs[0] if rcs else None


def latest_closed_meeting_for_run(run_id: str) -> dict[str, Any] | None:
    meetings = [
        meeting
        for meeting in list_meetings(include_archived=True)
        if meeting.get("run_id") == run_id and meeting.get("state") == "closed"
    ]
    return meetings[0] if meetings else None


def open_meeting_actions() -> list[dict[str, Any]]:
    actions: list[dict[str, Any]] = []
    for meeting in list_meetings(include_archived=True):
        for item in meeting.get("action_items", []):
            if isinstance(item, dict) and item.get("state", "open") != "open":
                continue
            summary = _text_value(item)
            if not summary:
                continue
            actions.append(
                {
                    "meeting_id": meeting.get("meeting_id", ""),
                    "meeting_state": meeting.get("state", ""),
                    "run_id": meeting.get("run_id", ""),
                    "mission_id": meeting.get("mission_id", ""),
                    "summary": summary,
                    "action_id": item.get("action_id", "") if isinstance(item, dict) else "",
                }
            )
    return actions


def find_run_index(registry: dict[str, Any], run_id: str) -> int:
    for index, run in enumerate(registry.get("runs", [])):
        if run.get("run_id") == run_id:
            return index
    return -1


def find_signal_index(log: dict[str, Any], signal_id: str) -> int:
    for index, signal in enumerate(log.get("signals", [])):
        if signal.get("signal_id") == signal_id:
            return index
    return -1


def append_unique_text(items: list[Any], value: str) -> list[Any]:
    if value and value not in items:
        items.append(value)
    return items


def active_run(run: dict[str, Any]) -> bool:
    return run.get("state") not in ("completed", "abandoned")


def active_signal(signal: dict[str, Any]) -> bool:
    return signal.get("state") == "active"


def allowed_run_transitions() -> dict[str, set[str]]:
    return {
        "draft": {"ready", "abandoned"},
        "ready": {"executing", "waiting", "blocked", "completed", "abandoned"},
        "executing": {"ready", "waiting", "blocked", "completed", "abandoned"},
        "waiting": {"ready", "executing", "blocked", "abandoned"},
        "blocked": {"ready", "waiting", "abandoned"},
        "completed": set(),
        "abandoned": set(),
    }


def validate_run_transition(run: dict[str, Any], new_state: str) -> None:
    current = run.get("state", "draft")
    if current == new_state:
        return
    allowed = allowed_run_transitions().get(current, set())
    if new_state not in allowed:
        raise SystemExit(f"invalid run state transition: {current} -> {new_state}")


def policy_value(key: str, default: Any) -> Any:
    policies = load_data(policies_path(), {}).get("policies", {})
    return get_key(policies, key, default)


def record_verification(run: dict[str, Any], ref: str) -> None:
    refs = run.setdefault("verification_refs", [])
    append_unique_text(refs, ref)
    run["verification_status"] = "passed"
    run["updated_at"] = now_iso()


def record_ship_ref(run: dict[str, Any], ref: str) -> None:
    refs = run.setdefault("ship_refs", [])
    append_unique_text(refs, ref)
    run["updated_at"] = now_iso()


def issue_body_from_template(
    *,
    title: str,
    objective: str,
    success_criteria: list[str],
    notes: list[str],
) -> str:
    template = issue_template_path().read_text(encoding="utf-8")
    criteria_block = "\n".join(f"- {item}" for item in success_criteria) if success_criteria else "- TBD"
    notes_block = "\n".join(f"- {item}" for item in notes) if notes else "- None"
    return (
        template.replace("__TITLE__", title)
        .replace("__OBJECTIVE__", objective or title)
        .replace("__SUCCESS_CRITERIA__", criteria_block)
        .replace("__NOTES__", notes_block)
    )


def pr_payload_for_run(run: dict[str, Any], base_branch: str, override_title: str = "") -> dict[str, str]:
    title = override_title.strip() or f"[{run.get('mission_id', '')}] {run.get('objective', '').strip()}".strip()
    summary = run.get("latest_summary", "").strip() or run.get("objective", "").strip() or "Mission delivery"
    verification_refs = run.get("verification_refs", [])
    verification_block = "\n".join(f"- {item}" for item in verification_refs) if verification_refs else "- Not yet recorded"
    template = pr_template_path().read_text(encoding="utf-8")
    body = (
        template.replace("__SUMMARY__", summary)
        .replace("__VERIFICATION__", verification_block)
        .replace("__BRANCH__", run.get("branch", ""))
        .replace("__RUN_ID__", run.get("run_id", ""))
        .replace("__MISSION_ID__", run.get("mission_id", ""))
    )
    return {
        "title": title,
        "body": body,
        "base_branch": base_branch,
        "branch": run.get("branch", ""),
    }


def open_github_pr(
    *,
    run: dict[str, Any],
    base_branch: str,
    title: str,
    body: str,
    draft: bool,
    reviewers: list[str],
) -> dict[str, str]:
    cmd = [
        "gh",
        "pr",
        "create",
        "--base",
        base_branch,
        "--head",
        run.get("branch", ""),
        "--title",
        title,
        "--body",
        body,
    ]
    if draft:
        cmd.append("--draft")
    for reviewer in reviewers:
        cmd.extend(["--reviewer", reviewer])
    output = run_command_with_stderr(cmd)
    pr_url = output.strip().splitlines()[-1].strip()
    pr_ref = pr_url.rstrip("/").split("/")[-1] if pr_url else ""
    return {
        "pr_ref": pr_ref,
        "pr_url": pr_url,
    }


def current_branch() -> str:
    try:
        return run_command(["git", "branch", "--show-current"])
    except Exception:
        return ""


def branch_exists(branch: str) -> bool:
    try:
        run_command(["git", "show-ref", "--verify", f"refs/heads/{branch}"])
        return True
    except Exception:
        return False


def git_status_porcelain() -> str:
    try:
        return run_command(["git", "status", "--porcelain"])
    except Exception:
        return ""


def git_is_dirty() -> bool:
    return bool(git_status_porcelain().strip())


def ensure_clean_checkout_context() -> None:
    if git_is_dirty():
        raise RuntimeError("current worktree has uncommitted changes; refusing branch checkout")


def worktree_root(profile: dict[str, Any]) -> Path:
    configured = get_key(profile, "profile.worktree_root", "")
    if configured:
        root = Path(configured)
        if not root.is_absolute():
            root = repo_root() / root
        return root
    return repo_root().parent / f"{repo_root().name}-worktrees"


def provision_worktree_path(profile: dict[str, Any], slug: str) -> Path:
    return worktree_root(profile) / slug


def git_checkout_new_branch(branch: str, base_branch: str) -> None:
    current = current_branch()
    if current == branch:
        return
    ensure_clean_checkout_context()
    try:
        run_command(["git", "fetch", "origin", base_branch])
    except Exception:
        pass
    if branch_exists(branch):
        run_command(["git", "checkout", branch])
        return
    run_command(["git", "checkout", base_branch])
    try:
        run_command(["git", "pull", "--ff-only", "origin", base_branch])
    except Exception:
        pass
    run_command(["git", "checkout", "-b", branch])


def git_create_worktree(branch: str, path: Path, base_branch: str) -> None:
    if path.exists():
        raise RuntimeError(f"worktree path already exists: {path}")
    path.parent.mkdir(parents=True, exist_ok=True)
    try:
        run_command(["git", "fetch", "origin", base_branch])
    except Exception:
        pass
    if branch_exists(branch):
        run_command(["git", "worktree", "add", str(path), branch])
    else:
        try:
            run_command(["git", "worktree", "add", "-b", branch, str(path), f"origin/{base_branch}"])
        except Exception:
            run_command(["git", "worktree", "add", "-b", branch, str(path), base_branch])


def git_delete_branch(branch: str, fallback_branch: str = "main") -> None:
    if not branch:
        return
    current = current_branch()
    if current == branch:
        ensure_clean_checkout_context()
        run_command(["git", "checkout", fallback_branch])
    if not branch_exists(branch):
        return
    run_command(["git", "branch", "-D", branch])


def git_remove_worktree(path: Path) -> None:
    if not path.exists():
        return
    main_repo = ""
    try:
        main_repo = run_command(["git", "worktree", "list", "--porcelain"])
    except Exception:
        pass
    if str(path) and str(path) in main_repo:
        run_command(["git", "worktree", "remove", "--force", str(path)])
        return
    raise RuntimeError(f"path exists but is not a registered worktree: {path}")


def default_authority_hierarchy() -> list[str]:
    return ["L1: live git/GitHub state", "L2: active runtime state", "L3: committed harness definitions"]


def default_meeting_roles() -> list[str]:
    return ["fact-verifier", "architecture-challenger", "completeness-checker"]


def materialize_meeting(
    *,
    meeting_id: str,
    run: dict[str, Any],
    question: str,
    requested_by: str,
) -> dict[str, Any]:
    template = load_data(meeting_template_path(), {}).get("meeting", {})
    timestamp = now_iso()
    meeting = merge_dicts(
        template,
        {
            "meeting_id": meeting_id,
            "mission_id": run.get("mission_id", ""),
            "run_id": run.get("run_id", ""),
            "question": question,
            "scope": run.get("objective", "") or run.get("mission_id", ""),
            "leader_id": run.get("leader_id", ""),
            "requested_by": requested_by,
            "state": "active",
            "authority_hierarchy": default_authority_hierarchy(),
            "locked_decisions": [],
            "roles": default_meeting_roles(),
            "packet": {
                "expected_output": "Produce accepted findings, rejected findings, unresolved questions, and action items.",
                "context_refs": [run.get("contract_ref", "")],
                "run_summary": run.get("latest_summary", ""),
                "status_reason": run.get("status_reason", ""),
            },
            "created_at": timestamp,
            "updated_at": timestamp,
            "closed_at": "",
        },
    )
    dump_data(meeting_path(meeting_id), {"version": 1, "meeting": meeting})
    return meeting


def update_contract_from_run(run: dict[str, Any]) -> None:
    contract_ref = run.get("contract_ref", "")
    if not contract_ref:
        return
    contract_path = Path(contract_ref)
    if not contract_path.exists():
        return
    data = load_data(contract_path, {})
    contract = data.get("contract", {})
    contract["branch"] = run.get("branch", "")
    contract["worktree_path"] = run.get("worktree_path", "")
    freshness = contract.setdefault("freshness", {})
    freshness["source_checked_at"] = now_iso()
    dump_data(contract_path, {"version": data.get("version", 1), "contract": contract})


def reprovision_run(
    *,
    run: dict[str, Any],
    profile: dict[str, Any],
    mode: str,
    apply_changes: bool,
) -> dict[str, Any]:
    if mode not in ("none", "branch", "worktree"):
        raise SystemExit("reprovision mode must be one of: none, branch, worktree")
    branch = run.get("branch", "")
    if not branch:
        raise SystemExit("run has no branch; cannot reprovision")
    base_branch = get_key(profile, "profile.base_branch", "main")
    current = current_branch()
    objective = run.get("objective", "") or run.get("mission_id", "")
    slug = slugify(objective)
    worktree_path = ""
    if mode == "worktree":
        worktree_path = str(provision_worktree_path(profile, slug))
        run["worktree_path"] = worktree_path
    elif mode == "none":
        run["worktree_path"] = ""

    provisioning = {
        "mode": mode,
        "apply_requested": apply_changes,
        "applied": False,
        "base_branch": base_branch,
        "previous_branch": current,
        "target_branch": branch,
    }
    if worktree_path:
        provisioning["target_worktree_path"] = worktree_path

    if apply_changes and mode != "none":
        if mode == "branch":
            git_checkout_new_branch(branch, base_branch)
            provisioning["applied"] = True
            provisioning["current_branch_after_apply"] = current_branch()
        elif mode == "worktree":
            git_create_worktree(branch, Path(worktree_path), base_branch)
            mission_marker = Path(worktree_path) / ".mission"
            mission_marker.write_text(f"{run.get('issue_number', '')}\n", encoding="utf-8")
            provisioning["applied"] = True

    run["mode"] = "solo" if mode == "none" else "managed"
    run["provisioning"] = provisioning
    update_contract_from_run(run)
    return run


def update_backlog_meeting(meeting_id: str, updates: dict[str, Any]) -> dict[str, Any]:
    backlog = read_maintenance_backlog()
    meetings = backlog.setdefault("meetings", [])
    for index, entry in enumerate(meetings):
        if entry.get("meeting_id") == meeting_id:
            meetings[index] = merge_dicts(entry, updates)
            write_maintenance_backlog(backlog)
            return meetings[index]
    new_entry = merge_dicts({"meeting_id": meeting_id, "created_at": now_iso()}, updates)
    meetings.append(new_entry)
    write_maintenance_backlog(backlog)
    return new_entry


def write_active_meeting(meeting: dict[str, Any]) -> None:
    meeting["updated_at"] = now_iso()
    dump_data(meeting_path(meeting["meeting_id"]), {"version": 1, "meeting": meeting})


def meeting_has_output(meeting: dict[str, Any]) -> bool:
    findings = meeting.get("findings", {})
    return bool(findings.get("accepted") or findings.get("rejected") or findings.get("unresolved") or meeting.get("action_items"))


def _text_value(item: Any) -> str:
    if isinstance(item, dict):
        return str(item.get("summary", "")).strip()
    if isinstance(item, str):
        return item.strip()
    return ""


def meeting_outcome(meeting: dict[str, Any]) -> dict[str, str]:
    findings = meeting.get("findings", {})
    accepted = [_text_value(item) for item in findings.get("accepted", [])]
    accepted = [item for item in accepted if item]
    actions = []
    for item in meeting.get("action_items", []):
        if isinstance(item, dict) and item.get("state", "open") != "open":
            continue
        text = _text_value(item)
        if text:
            actions.append(text)

    summary = accepted[0] if accepted else (meeting.get("resolution_note", "").strip() or meeting.get("question", "").strip())
    next_action = actions[0] if actions else "Resume execution with the meeting outcome applied."
    return {
        "summary": summary,
        "next_action": next_action,
    }


def materialize_rc(
    *,
    rc_id: str,
    version: str,
    branch: str,
    base_branch: str,
    production_branch: str,
    run_ids: list[str],
) -> dict[str, Any]:
    template = load_data(rc_template_path(), {}).get("rc", {})
    timestamp = now_iso()
    rc = merge_dicts(
        template,
        {
            "rc_id": rc_id,
            "version": version,
            "state": "active",
            "branch": branch,
            "base_branch": base_branch,
            "production_branch": production_branch,
            "run_ids": run_ids,
            "verification_status": "unknown",
            "verification_refs": [],
            "notes": [],
            "created_at": timestamp,
            "updated_at": timestamp,
            "promoted_at": "",
            "closed_at": "",
        },
    )
    dump_data(rc_path(rc_id), {"version": 1, "rc": rc})
    return rc


def issue_record(
    *,
    severity: str,
    kind: str,
    ref: str,
    summary: str,
    recommendation: str,
) -> dict[str, str]:
    return {
        "severity": severity,
        "kind": kind,
        "ref": ref,
        "summary": summary,
        "recommendation": recommendation,
    }


def age_hours(value: str) -> float | None:
    parsed = parse_iso(value)
    if not parsed:
        return None
    return round((now_dt() - parsed).total_seconds() / 3600, 2)


def doctor_report() -> dict[str, Any]:
    ensure_runtime_dirs()
    profile, _ = resolve_profile()
    registry = read_registry()
    signals = read_signal_log()
    meetings = [meeting for meeting in list_meetings() if meeting.get("state") == "active"]
    meeting_actions = open_meeting_actions()
    active_rc_record = active_rc()
    active_rcs = [rc for rc in list_rcs() if rc.get("state") == "active"]
    active_runs = [run for run in registry.get("runs", []) if active_run(run)]
    active_signals = [signal for signal in signals.get("signals", []) if active_signal(signal)]
    active_run_ids = {run.get("run_id", "") for run in active_runs}
    run_stale_after = float(policy_value("cleanup.stale_run_after_hours", 24))
    signal_stale_after = float(policy_value("cleanup.stale_signal_after_hours", 12))
    meeting_stale_after = float(policy_value("cleanup.stale_meeting_after_hours", 6))
    issues: list[dict[str, str]] = []

    for run in active_runs:
        run_id = run.get("run_id", "")
        if not run.get("created_at"):
            issues.append(
                issue_record(
                    severity="warning",
                    kind="run-missing-created-at",
                    ref=run_id,
                    summary="Run is active but missing created_at.",
                    recommendation="Backfill timestamps or recreate the run cleanly.",
                )
            )
        report_age = age_hours(run.get("last_reported_at", ""))
        if report_age is None:
            issues.append(
                issue_record(
                    severity="warning",
                    kind="run-missing-last-reported-at",
                    ref=run_id,
                    summary="Run is active but missing last_reported_at.",
                    recommendation="Have the owner emit a report or repair the run state.",
                )
            )
        elif report_age > run_stale_after:
            issues.append(
                issue_record(
                    severity="warning",
                    kind="run-stale",
                    ref=run_id,
                    summary=f"Run has not reported for {report_age} hours.",
                    recommendation="Review whether the run is still active, blocked, or abandoned.",
                )
            )
        if run.get("state") == "blocked":
            issues.append(
                issue_record(
                    severity="warning",
                    kind="run-blocked",
                    ref=run_id,
                    summary=run.get("status_reason", "") or "Run is blocked.",
                    recommendation="Resolve the blocker or convene a meeting.",
                )
            )
        latest_closed_meeting = latest_closed_meeting_for_run(run_id)
        if latest_closed_meeting and meeting_has_output(latest_closed_meeting):
            outcome = meeting_outcome(latest_closed_meeting)
            if outcome["summary"] and run.get("latest_summary", "").strip() != outcome["summary"]:
                issues.append(
                    issue_record(
                        severity="warning",
                        kind="run-missing-meeting-summary",
                        ref=run_id,
                        summary="Run has a closed meeting with accepted output that is not reflected in latest_summary.",
                        recommendation="Adopt the meeting outcome into the run state.",
                    )
                )
            if outcome["next_action"] and run.get("next_expected_action", "").strip() != outcome["next_action"]:
                issues.append(
                    issue_record(
                        severity="warning",
                        kind="run-missing-meeting-next-action",
                        ref=run_id,
                        summary="Run has a closed meeting with action items that are not reflected in next_expected_action.",
                        recommendation="Adopt the meeting outcome into the run state.",
                    )
                )
        worktree_path = run.get("worktree_path", "")
        if worktree_path and not Path(worktree_path).exists():
            issues.append(
                issue_record(
                    severity="warning",
                    kind="worktree-missing",
                    ref=run_id,
                    summary=f"Run references missing worktree path: {worktree_path}",
                    recommendation="Re-provision the worktree or update the run record.",
                )
            )

    for signal in active_signals:
        signal_id = signal.get("signal_id", "")
        run_id = signal.get("run_id", "")
        if run_id and run_id not in active_run_ids:
            issues.append(
                issue_record(
                    severity="critical",
                    kind="signal-orphaned",
                    ref=signal_id,
                    summary="Active signal references a missing run.",
                    recommendation="Run maintenance reconcile before proceeding.",
                )
            )
        if not signal.get("emitted_at"):
            issues.append(
                issue_record(
                    severity="warning",
                    kind="signal-missing-emitted-at",
                    ref=signal_id,
                    summary="Signal is active but missing emitted_at.",
                    recommendation="Repair signal metadata or re-emit the signal.",
                )
            )
        emitted_age = age_hours(signal.get("emitted_at", ""))
        if emitted_age is not None and emitted_age > signal_stale_after:
            issues.append(
                issue_record(
                    severity="warning",
                    kind="signal-stale",
                    ref=signal_id,
                    summary=f"Signal has been active for {emitted_age} hours.",
                    recommendation="Resolve or refresh the signal if it still matters.",
                )
            )
        fresh_until = parse_iso(signal.get("fresh_until", ""))
        if fresh_until and fresh_until < now_dt():
            issues.append(
                issue_record(
                    severity="warning",
                    kind="signal-expired",
                    ref=signal_id,
                    summary="Signal freshness window has expired.",
                    recommendation="Resolve or refresh the signal.",
                )
            )

    for meeting in meetings:
        meeting_id = meeting.get("meeting_id", "")
        run_id = meeting.get("run_id", "")
        if run_id and run_id not in active_run_ids:
            issues.append(
                issue_record(
                    severity="critical",
                    kind="meeting-orphaned",
                    ref=meeting_id,
                    summary="Active meeting references a missing run.",
                    recommendation="Run maintenance reconcile before proceeding.",
                )
            )

    if len(active_rcs) > 1:
        issues.append(
            issue_record(
                severity="critical",
                kind="rc-multiple-active",
                ref="rc-control-plane",
                summary="More than one active RC exists.",
                recommendation="Close or promote surplus RCs before proceeding.",
            )
        )

    for rc in active_rcs:
        rc_id = rc.get("rc_id", "")
        if not rc.get("run_ids", []):
            issues.append(
                issue_record(
                    severity="warning",
                    kind="rc-missing-runs",
                    ref=rc_id,
                    summary="Active RC has no associated shipped runs.",
                    recommendation="Attach shipped runs or close the RC.",
                )
            )
        if rc.get("verification_status", "unknown") != "passed":
            issues.append(
                issue_record(
                    severity="warning",
                    kind="rc-unverified",
                    ref=rc_id,
                    summary="Active RC has not been marked verified.",
                    recommendation="Record RC verification before promotion.",
                )
            )
        if not meeting.get("authority_hierarchy"):
            issues.append(
                issue_record(
                    severity="warning",
                    kind="meeting-missing-authority",
                    ref=meeting_id,
                    summary="Meeting has no authority hierarchy.",
                    recommendation="Repair the meeting packet before using its conclusions.",
                )
            )
        if not meeting.get("roles"):
            issues.append(
                issue_record(
                    severity="warning",
                    kind="meeting-missing-roles",
                    ref=meeting_id,
                    summary="Meeting has no assigned roles.",
                    recommendation="Assign bounded roles before continuing the meeting.",
                )
            )
        meeting_age = age_hours(meeting.get("created_at", ""))
        if meeting_age is not None and meeting_age > meeting_stale_after:
            issues.append(
                issue_record(
                    severity="warning",
                    kind="meeting-stale",
                    ref=meeting_id,
                    summary=f"Meeting has been open for {meeting_age} hours.",
                    recommendation="Close the meeting or refresh its packet and participants.",
                )
            )
        if meeting_age is not None and meeting_age > 0.5 and not meeting_has_output(meeting):
            issues.append(
                issue_record(
                    severity="warning",
                    kind="meeting-no-output",
                    ref=meeting_id,
                    summary="Meeting is active but has not recorded findings or action items yet.",
                    recommendation="Record council findings or close the meeting if it is no longer needed.",
                )
            )

    severity_order = {"critical": 0, "warning": 1, "info": 2}
    issues.sort(key=lambda item: (severity_order.get(item["severity"], 9), item["kind"], item["ref"]))
    critical = sum(1 for item in issues if item["severity"] == "critical")
    warning = sum(1 for item in issues if item["severity"] == "warning")
    info = sum(1 for item in issues if item["severity"] == "info")

    if critical:
        next_action = "Repair critical runtime drift before new execution."
    elif warning:
        next_action = "Triage warnings, then continue execution."
    else:
        next_action = "Runtime state looks healthy."

    return {
        "profile_id": get_key(profile, "profile.profile_id", ""),
        "issue_count": len(issues),
        "critical_count": critical,
        "warning_count": warning,
        "info_count": info,
        "issues": issues,
        "next_action": next_action,
    }


def cmd_paths(args: list[str]) -> int:
    if not args:
        raise SystemExit("paths requires a subcommand")
    subcommand = args[0]
    if subcommand == "repo-root":
        print(repo_root())
    elif subcommand == "definitions-root":
        print(definitions_root())
    elif subcommand == "runtime-root":
        print(runtime_root())
    elif subcommand == "ensure-runtime-dirs":
        ensure_runtime_dirs()
        print(runtime_root())
    elif subcommand == "profile-path":
        _, base_path = resolve_profile(args[1] if len(args) > 1 else None)
        print(base_path or "")
    else:
        raise SystemExit(f"unknown paths subcommand: {subcommand}")
    return 0


def cmd_profile(args: list[str]) -> int:
    if not args:
        raise SystemExit("profile requires a subcommand")
    subcommand = args[0]
    profile, base_path = resolve_profile(args[1] if len(args) > 1 and subcommand == "path" else None)
    if subcommand == "show":
        print(json.dumps(profile, indent=2, ensure_ascii=True))
    elif subcommand == "path":
        print(base_path or "")
    elif subcommand == "id":
        print(get_key(profile, "profile.profile_id", ""))
    elif subcommand == "get":
        if len(args) < 2:
            raise SystemExit("profile get requires KEY")
        key = args[1]
        default = args[2] if len(args) > 2 else ""
        target = profile
        if not key.startswith("profile."):
            key = f"profile.{key}"
        value = get_key(target, key, default)
        if isinstance(value, (dict, list)):
            print(json.dumps(value, ensure_ascii=True))
        else:
            print(value)
    else:
        raise SystemExit(f"unknown profile subcommand: {subcommand}")
    return 0


def cmd_runs(args: list[str]) -> int:
    if not args:
        raise SystemExit("runs requires a subcommand")
    ensure_runtime_dirs()
    subcommand = args[0]
    registry = read_registry()
    if subcommand == "init":
        with harness_lock():
            write_registry(registry)
        print(registry_path())
        return 0
    if subcommand == "list":
        print(json.dumps(registry, indent=2, ensure_ascii=True))
        return 0
    if subcommand == "active":
        mission_id = args[1] if len(args) > 1 else None
        runs = [run for run in registry.get("runs", []) if active_run(run)]
        if mission_id:
            runs = [run for run in runs if run.get("mission_id") == mission_id]
        print(json.dumps(runs, indent=2, ensure_ascii=True))
        return 0
    if subcommand == "get":
        if len(args) < 2:
            raise SystemExit("runs get requires RUN_ID")
        index = find_run_index(registry, args[1])
        if index < 0:
            raise SystemExit(f"run not found: {args[1]}")
        print(json.dumps(registry["runs"][index], indent=2, ensure_ascii=True))
        return 0
    if subcommand == "add":
        with harness_lock():
            registry = read_registry()
            template = load_data(run_template_path(), {}).get("run", {})
            payload = merge_dicts(template, read_stdin_json())
            if not payload.get("run_id"):
                raise SystemExit("runs add requires run_id")
            if find_run_index(registry, payload["run_id"]) >= 0:
                raise SystemExit(f"run already exists: {payload['run_id']}")
            timestamp = now_iso()
            fill_if_blank(payload, "created_at", timestamp)
            if payload.get("state") in ("ready", "executing", "paused", "waiting", "blocked"):
                fill_if_blank(payload, "started_at", timestamp)
                fill_if_blank(payload, "last_reported_at", timestamp)
            payload["updated_at"] = timestamp
            registry.setdefault("runs", []).append(payload)
            write_registry(registry)
        print(json.dumps(payload, indent=2, ensure_ascii=True))
        return 0
    if subcommand == "update":
        if len(args) < 2:
            raise SystemExit("runs update requires RUN_ID")
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, args[1])
            if index < 0:
                raise SystemExit(f"run not found: {args[1]}")
            registry["runs"][index] = merge_dicts(registry["runs"][index], read_stdin_json())
            registry["runs"][index]["updated_at"] = now_iso()
            write_registry(registry)
        print(json.dumps(registry["runs"][index], indent=2, ensure_ascii=True))
        return 0
    if subcommand == "set-state":
        if len(args) < 3:
            raise SystemExit("runs set-state requires RUN_ID STATE [REASON]")
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, args[1])
            if index < 0:
                raise SystemExit(f"run not found: {args[1]}")
            validate_run_transition(registry["runs"][index], args[2])
            registry["runs"][index]["state"] = args[2]
            if len(args) > 3:
                registry["runs"][index]["status_reason"] = args[3]
            registry["runs"][index]["updated_at"] = now_iso()
            if args[2] in ("completed", "abandoned"):
                registry["runs"][index]["completed_at"] = now_iso()
            write_registry(registry)
        print(json.dumps(registry["runs"][index], indent=2, ensure_ascii=True))
        return 0
    raise SystemExit(f"unknown runs subcommand: {subcommand}")


def cmd_signals(args: list[str]) -> int:
    if not args:
        raise SystemExit("signals requires a subcommand")
    ensure_runtime_dirs()
    subcommand = args[0]
    log = read_signal_log()
    if subcommand == "init":
        with harness_lock():
            write_signal_log(log)
        print(signal_log_path())
        return 0
    if subcommand == "list":
        print(json.dumps(log, indent=2, ensure_ascii=True))
        return 0
    if subcommand == "active":
        mission_id = args[1] if len(args) > 1 else None
        signals = [sig for sig in log.get("signals", []) if active_signal(sig)]
        if mission_id:
            signals = [sig for sig in signals if sig.get("mission_id") == mission_id]
        print(json.dumps(signals, indent=2, ensure_ascii=True))
        return 0
    if subcommand == "emit":
        with harness_lock():
            log = read_signal_log()
            template = load_data(signal_template_path(), {}).get("signal", {})
            payload = merge_dicts(template, read_stdin_json())
            if not payload.get("signal_id"):
                raise SystemExit("signals emit requires signal_id")
            if find_signal_index(log, payload["signal_id"]) >= 0:
                raise SystemExit(f"signal already exists: {payload['signal_id']}")
            fill_if_blank(payload, "emitted_at", now_iso())
            fill_if_blank(payload, "state", "active")
            log.setdefault("signals", []).append(payload)
            write_signal_log(log)
        print(json.dumps(payload, indent=2, ensure_ascii=True))
        return 0
    if subcommand == "set-state":
        if len(args) < 3:
            raise SystemExit("signals set-state requires SIGNAL_ID STATE [NOTE]")
        with harness_lock():
            log = read_signal_log()
            index = find_signal_index(log, args[1])
            if index < 0:
                raise SystemExit(f"signal not found: {args[1]}")
            log["signals"][index]["state"] = args[2]
            if len(args) > 3:
                log["signals"][index]["resolution_note"] = args[3]
            write_signal_log(log)
        print(json.dumps(log["signals"][index], indent=2, ensure_ascii=True))
        return 0
    raise SystemExit(f"unknown signals subcommand: {subcommand}")


def cmd_contracts(args: list[str]) -> int:
    if not args:
        raise SystemExit("contracts requires a subcommand")
    ensure_runtime_dirs()
    subcommand = args[0]
    if subcommand == "materialize":
        if len(args) < 2:
            raise SystemExit("contracts materialize requires OUTPUT_PATH")
        output_path = Path(args[1])
        if not output_path.is_absolute():
            output_path = repo_root() / output_path
        template = load_data(contract_template_path(), {}).get("contract", {})
        payload = merge_dicts(template, read_stdin_json())
        freshness = payload.setdefault("freshness", {})
        if not freshness.get("created_at"):
            freshness["created_at"] = now_iso()
        dump_data(output_path, {"version": 1, "contract": payload})
        print(output_path)
        return 0
    if subcommand == "show-template":
        template = load_data(contract_template_path(), {})
        print(json.dumps(template, indent=2, ensure_ascii=True))
        return 0
    raise SystemExit(f"unknown contracts subcommand: {subcommand}")


def cmd_home(args: list[str]) -> int:
    del args
    ensure_runtime_dirs()
    profile, _ = resolve_profile()
    registry = read_registry()
    signals = read_signal_log()
    meetings = [meeting for meeting in list_meetings() if meeting.get("state") == "active"]
    meeting_actions = open_meeting_actions()
    active_rc_record = active_rc()
    active_runs = [run for run in registry.get("runs", []) if active_run(run)]
    active_signals = [signal for signal in signals.get("signals", []) if active_signal(signal)]
    blocked_runs = [run for run in active_runs if run.get("state") == "blocked"]
    waiting_runs = [run for run in active_runs if run.get("state") == "waiting"]
    active_run_ids = {run.get("run_id", "") for run in active_runs}
    orphan_signals = [
        signal for signal in active_signals if signal.get("run_id") and signal.get("run_id") not in active_run_ids
    ]
    orphan_meetings = [
        meeting for meeting in meetings if meeting.get("run_id") and meeting.get("run_id") not in active_run_ids
    ]

    if orphan_signals or orphan_meetings:
        next_action = "Reconcile orphaned runtime objects before proceeding."
    elif not active_runs:
        next_action = "No active runs. Claim a mission."
    elif meeting_actions:
        next_action = "Review open meeting action items and translate them into active execution."
    elif blocked_runs:
        next_action = "Review blockers or convene a meeting."
    elif meetings:
        next_action = "Review active meetings and convert conclusions into run updates."
    elif active_rc_record:
        next_action = "Review the active RC and decide whether to promote, fix, or close it."
    elif waiting_runs:
        next_action = "Review waiting runs and dependency state."
    else:
        next_action = "Continue active execution and monitor signals."

    summary = {
        "profile_id": get_key(profile, "profile.profile_id", ""),
        "active_run_count": len(active_runs),
        "blocked_run_count": len(blocked_runs),
        "active_signal_count": len(active_signals),
        "active_meeting_count": len(meetings),
        "active_rc_count": 1 if active_rc_record else 0,
        "open_meeting_action_count": len(meeting_actions),
        "orphan_signal_count": len(orphan_signals),
        "orphan_meeting_count": len(orphan_meetings),
        "runs": active_runs,
        "signals": active_signals,
        "meetings": meetings,
        "active_rc": active_rc_record or {},
        "meeting_actions": meeting_actions,
        "orphan_signals": orphan_signals,
        "orphan_meetings": orphan_meetings,
        "next_action": next_action,
    }
    print(json.dumps(summary, indent=2, ensure_ascii=True))
    return 0


def cmd_doctor(args: list[str]) -> int:
    del args
    print(json.dumps(doctor_report(), indent=2, ensure_ascii=True))
    return 0


def cmd_issue(args: list[str]) -> int:
    ensure_runtime_dirs()
    if not args:
        raise SystemExit("issue requires a subcommand")
    subcommand = args[0]
    if subcommand not in ("body", "draft", "create"):
        raise SystemExit(f"unknown issue subcommand: {subcommand}")

    title_parts: list[str] = []
    objective = ""
    criteria: list[str] = []
    notes: list[str] = []
    labels: list[str] = []
    assignees: list[str] = []
    draft_only = subcommand in ("body", "draft")

    for arg in args[1:]:
        if arg.startswith("--objective="):
            objective = arg.split("=", 1)[1].strip()
        elif arg.startswith("--criteria="):
            raw = arg.split("=", 1)[1].strip()
            criteria.extend([item.strip() for item in raw.split("|") if item.strip()])
        elif arg.startswith("--note="):
            note = arg.split("=", 1)[1].strip()
            if note:
                notes.append(note)
        elif arg.startswith("--label="):
            label = arg.split("=", 1)[1].strip()
            if label:
                labels.append(label)
        elif arg.startswith("--assignee="):
            assignee = arg.split("=", 1)[1].strip()
            if assignee:
                assignees.append(assignee)
        elif arg == "--draft":
            draft_only = True
        else:
            title_parts.append(arg)

    title = " ".join(title_parts).strip()
    if not title:
        raise SystemExit(f"issue {subcommand} requires TITLE")

    body = issue_body_from_template(
        title=title,
        objective=objective or title,
        success_criteria=criteria,
        notes=notes,
    )

    payload = {
        "title": title,
        "objective": objective or title,
        "success_criteria": criteria,
        "notes": notes,
        "labels": labels,
        "assignees": assignees,
        "body": body,
        "mode": "draft" if draft_only else "create",
    }

    if subcommand == "body":
        print(body)
        return 0

    if draft_only:
        print(json.dumps(payload, indent=2, ensure_ascii=True))
        return 0

    cmd = ["gh", "issue", "create", "--title", title, "--body", body]
    for label in labels:
        cmd.extend(["--label", label])
    for assignee in assignees:
        cmd.extend(["--assignee", assignee])

    try:
        output = run_command_with_stderr(cmd)
    except Exception as exc:
        raise SystemExit(f"issue create failed: {exc}")

    issue_url = output.strip().splitlines()[-1].strip()
    issue_number = 0
    match = re.search(r"/issues/(\\d+)$", issue_url)
    if match:
        issue_number = int(match.group(1))

    payload["issue_url"] = issue_url
    payload["issue_number"] = issue_number
    print(json.dumps(payload, indent=2, ensure_ascii=True))
    return 0


def cmd_claim(args: list[str]) -> int:
    if not args:
        raise SystemExit("claim requires ISSUE_NUMBER [TITLE_OVERRIDE]")
    ensure_runtime_dirs()
    provision_mode = "none"
    apply_changes = False
    positionals: list[str] = []
    for arg in args:
        if arg.startswith("--provision="):
            provision_mode = arg.split("=", 1)[1]
        elif arg == "--apply":
            apply_changes = True
        else:
            positionals.append(arg)
    if not positionals:
        raise SystemExit("claim requires ISSUE_NUMBER [TITLE_OVERRIDE]")
    issue_number = positionals[0]
    title_override = " ".join(positionals[1:]).strip() if len(positionals) > 1 else ""

    profile, _ = resolve_profile()
    issue_title = title_override
    issue_url = ""
    issue_body = ""
    if not issue_title:
        try:
            issue_data = json.loads(
                run_command(["gh", "issue", "view", issue_number, "--json", "title,body,url"])
            )
            issue_title = issue_data.get("title", "").strip()
            issue_body = issue_data.get("body", "").strip()
            issue_url = issue_data.get("url", "").strip()
        except Exception:
            raise SystemExit("claim could not fetch issue data; provide a title override")

    with harness_lock():
        registry = read_registry()
        mission_id = f"issue-{issue_number}"
        active_existing = [run for run in registry.get("runs", []) if active_run(run) and run.get("mission_id") == mission_id]
        if active_existing:
            raise SystemExit(f"active run already exists for {mission_id}")

        slug = slugify(issue_title)
        branch_pattern = get_key(profile, "profile.branch_pattern", "mission/{issue}-{slug}")
        branch = branch_from_pattern(branch_pattern, issue_number, slug)
        contract_path = runtime_root() / "contracts" / f"issue-{issue_number}.yml"
        run_id = f"run-{datetime.now(timezone.utc).strftime('%Y%m%d%H%M%S')}-claim-{issue_number}"
        base_branch = get_key(profile, "profile.base_branch", "main")
        planned_worktree_path = ""
        current = current_branch()

        if provision_mode == "worktree":
            planned_worktree_path = str(provision_worktree_path(profile, slug))
        elif provision_mode not in ("none", "branch"):
            raise SystemExit("claim --provision must be one of: none, branch, worktree")

        contract_payload = merge_dicts(
            load_data(contract_template_path(), {}).get("contract", {}),
            {
                "mission_id": mission_id,
                "issue_number": int(issue_number),
                "title": issue_title,
                "objective": issue_title,
                "success_criteria": [],
                "sub_tasks": [],
                "branch": branch,
                "worktree_path": planned_worktree_path,
                "context_files": [],
                "verification_commands": get_key(profile, "profile.default_verify_commands", []),
                "notes": {
                    "assumptions_invalidated": [],
                    "manual_followups": [],
                    "discoveries": [item for item in [issue_url, issue_body[:240] if issue_body else ""] if item],
                },
            },
        )
        freshness = contract_payload.setdefault("freshness", {})
        if not freshness.get("created_at"):
            freshness["created_at"] = now_iso()
        dump_data(contract_path, {"version": 1, "contract": contract_payload})

        run_payload = merge_dicts(
            load_data(run_template_path(), {}).get("run", {}),
            {
                "run_id": run_id,
                "mission_id": mission_id,
                "issue_number": int(issue_number),
                "leader_id": "leader-main",
                "owner_type": "leader",
                "owner_id": "leader-main",
                "state": "ready",
                "mode": "solo" if provision_mode == "none" else "managed",
                "execution_surface": "foreground-skill",
                "branch": branch,
                "worktree_path": planned_worktree_path,
                "contract_ref": str(contract_path),
                "objective": issue_title,
                "latest_summary": "Claimed and ready for execution",
                "next_expected_action": "Start drive flow and run the narrowest useful verification.",
            },
        )
        timestamp = now_iso()
        fill_if_blank(run_payload, "created_at", timestamp)
        fill_if_blank(run_payload, "started_at", timestamp)
        fill_if_blank(run_payload, "last_reported_at", timestamp)
        run_payload["updated_at"] = timestamp

        provisioning = {
            "mode": provision_mode,
            "apply_requested": apply_changes,
            "applied": False,
            "base_branch": base_branch,
            "previous_branch": current,
        }
        if provision_mode == "branch":
            provisioning["target_branch"] = branch
        elif provision_mode == "worktree":
            provisioning["target_branch"] = branch
            provisioning["target_worktree_path"] = planned_worktree_path

        if apply_changes and provision_mode != "none":
            try:
                if provision_mode == "branch":
                    git_checkout_new_branch(branch, base_branch)
                elif provision_mode == "worktree":
                    git_create_worktree(branch, Path(planned_worktree_path), base_branch)
                    mission_marker = Path(planned_worktree_path) / ".mission"
                    mission_marker.write_text(f"{issue_number}\n", encoding="utf-8")
                provisioning["applied"] = True
                if provision_mode == "branch":
                    provisioning["current_branch_after_apply"] = current_branch()
            except Exception as exc:
                provisioning["error"] = str(exc)
                run_payload["state"] = "blocked"
                run_payload["status_reason"] = f"Provisioning failed: {exc}"
                run_payload["next_expected_action"] = "Clean the current worktree or switch to worktree provisioning."
                log = read_signal_log()
                signal = merge_dicts(
                    load_data(signal_template_path(), {}).get("signal", {}),
                    {
                        "signal_id": f"sig-{datetime.now(timezone.utc).strftime('%Y%m%d%H%M%S')}-claim-failed-{issue_number}",
                        "kind": "blocked",
                        "source_type": "leader",
                        "source_id": "leader-main",
                        "run_id": run_id,
                        "mission_id": mission_id,
                        "applies_to": [mission_id],
                        "scope_type": "local",
                        "title": "Claim provisioning failed",
                        "summary": str(exc),
                        "authority_level": "L1",
                        "confidence": "high",
                        "fresh_until": now_iso(),
                    },
                )
                fill_if_blank(signal, "emitted_at", now_iso())
                fill_if_blank(signal, "state", "active")
                log.setdefault("signals", []).append(signal)
                write_signal_log(log)

        run_payload["provisioning"] = provisioning
        registry.setdefault("runs", []).append(run_payload)
        write_registry(registry)

    summary = {
        "mission_id": mission_id,
        "issue_number": int(issue_number),
        "title": issue_title,
        "branch": branch,
        "contract_path": str(contract_path),
        "run_id": run_id,
        "issue_url": issue_url,
        "provisioning": provisioning,
        "next_action": run_payload["next_expected_action"],
    }
    print(json.dumps(summary, indent=2, ensure_ascii=True))
    return 0


def cmd_drive(args: list[str]) -> int:
    if not args:
        raise SystemExit("drive requires a subcommand")
    ensure_runtime_dirs()
    subcommand = args[0]
    registry = read_registry()

    if subcommand == "report":
        if len(args) < 3:
            raise SystemExit("drive report requires RUN_ID SUMMARY")
        run_id = args[1]
        summary = " ".join(args[2:]).strip()
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            registry["runs"][index]["latest_summary"] = summary
            registry["runs"][index]["last_reported_at"] = now_iso()
            registry["runs"][index]["updated_at"] = now_iso()
            write_registry(registry)
        print(json.dumps(registry["runs"][index], indent=2, ensure_ascii=True))
        return 0

    if subcommand == "start":
        if len(args) < 2:
            raise SystemExit("drive start requires RUN_ID [SUMMARY]")
        run_id = args[1]
        summary = " ".join(args[2:]).strip() if len(args) > 2 else "Execution started"
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            validate_run_transition(registry["runs"][index], "executing")
            registry["runs"][index]["state"] = "executing"
            registry["runs"][index]["latest_summary"] = summary
            registry["runs"][index]["status_reason"] = ""
            registry["runs"][index]["last_reported_at"] = now_iso()
            registry["runs"][index]["updated_at"] = now_iso()
            write_registry(registry)
        print(json.dumps(registry["runs"][index], indent=2, ensure_ascii=True))
        return 0

    if subcommand == "resume":
        if len(args) < 2:
            raise SystemExit("drive resume requires RUN_ID [SUMMARY]")
        run_id = args[1]
        summary = " ".join(args[2:]).strip() if len(args) > 2 else "Execution resumed"
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            validate_run_transition(registry["runs"][index], "ready")
            registry["runs"][index]["state"] = "ready"
            registry["runs"][index]["latest_summary"] = summary
            registry["runs"][index]["status_reason"] = ""
            registry["runs"][index]["last_reported_at"] = now_iso()
            registry["runs"][index]["updated_at"] = now_iso()
            write_registry(registry)
        print(json.dumps(registry["runs"][index], indent=2, ensure_ascii=True))
        return 0

    if subcommand == "complete":
        if len(args) < 2:
            raise SystemExit("drive complete requires RUN_ID [SUMMARY]")
        run_id = args[1]
        summary = " ".join(args[2:]).strip() if len(args) > 2 else "Execution completed"
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            validate_run_transition(registry["runs"][index], "completed")
            registry["runs"][index]["state"] = "completed"
            registry["runs"][index]["latest_summary"] = summary
            registry["runs"][index]["status_reason"] = ""
            registry["runs"][index]["completed_at"] = now_iso()
            registry["runs"][index]["last_reported_at"] = now_iso()
            registry["runs"][index]["updated_at"] = now_iso()
            write_registry(registry)
        print(json.dumps(registry["runs"][index], indent=2, ensure_ascii=True))
        return 0

    if subcommand == "adopt-meeting":
        if len(args) < 3:
            raise SystemExit("drive adopt-meeting requires RUN_ID MEETING_ID")
        run_id = args[1]
        meeting_id = args[2]
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            data = read_meeting(meeting_id)
            if not data:
                raise SystemExit(f"meeting not found: {meeting_id}")
            meeting = data.get("meeting", {})
            if meeting.get("state") != "closed":
                raise SystemExit("drive adopt-meeting requires a closed meeting")
            outcome = meeting_outcome(meeting)
            if outcome["summary"]:
                registry["runs"][index]["latest_summary"] = outcome["summary"]
            if outcome["next_action"]:
                registry["runs"][index]["next_expected_action"] = outcome["next_action"]
            registry["runs"][index]["status_reason"] = meeting.get("resolution_note", "") or registry["runs"][index].get("status_reason", "")
            registry["runs"][index]["updated_at"] = now_iso()
            write_registry(registry)
        print(
            json.dumps(
                {
                    "run": registry["runs"][index],
                    "meeting": {
                        "meeting_id": meeting.get("meeting_id", ""),
                        "resolution_note": meeting.get("resolution_note", ""),
                    },
                },
                indent=2,
                ensure_ascii=True,
            )
        )
        return 0

    if subcommand == "reprovision":
        if len(args) < 2:
            raise SystemExit("drive reprovision requires RUN_ID [--mode=none|branch|worktree] [--apply]")
        run_id = ""
        mode = ""
        apply_changes = False
        for arg in args[1:]:
            if arg.startswith("--mode="):
                mode = arg.split("=", 1)[1]
            elif arg == "--apply":
                apply_changes = True
            elif not run_id:
                run_id = arg
        if not run_id:
            raise SystemExit("drive reprovision requires RUN_ID")
        with harness_lock():
            profile, _ = resolve_profile()
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            current_mode = get_key(registry["runs"][index], "provisioning.mode", "") or ("worktree" if registry["runs"][index].get("worktree_path") else "none")
            target_mode = mode or current_mode or "none"
            try:
                registry["runs"][index] = reprovision_run(
                    run=registry["runs"][index],
                    profile=profile,
                    mode=target_mode,
                    apply_changes=apply_changes,
                )
                registry["runs"][index]["status_reason"] = ""
                registry["runs"][index]["next_expected_action"] = "Start drive flow and run the narrowest useful verification."
            except Exception as exc:
                registry["runs"][index]["status_reason"] = f"Reprovisioning failed: {exc}"
                registry["runs"][index]["next_expected_action"] = "Review reprovision failure and choose a safer provisioning mode."
                raise
            finally:
                registry["runs"][index]["updated_at"] = now_iso()
                write_registry(registry)
        print(json.dumps(registry["runs"][index], indent=2, ensure_ascii=True))
        return 0

    if subcommand == "set-next":
        if len(args) < 3:
            raise SystemExit("drive set-next requires RUN_ID NEXT_ACTION")
        run_id = args[1]
        next_action = " ".join(args[2:]).strip()
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            registry["runs"][index]["next_expected_action"] = next_action
            registry["runs"][index]["updated_at"] = now_iso()
            write_registry(registry)
        print(json.dumps(registry["runs"][index], indent=2, ensure_ascii=True))
        return 0

    if subcommand == "block":
        if len(args) < 3:
            raise SystemExit("drive block requires RUN_ID REASON")
        run_id = args[1]
        reason = " ".join(args[2:]).strip()
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            validate_run_transition(registry["runs"][index], "blocked")
            registry["runs"][index]["state"] = "blocked"
            registry["runs"][index]["status_reason"] = reason
            registry["runs"][index]["last_reported_at"] = now_iso()
            registry["runs"][index]["updated_at"] = now_iso()
            write_registry(registry)

            log = read_signal_log()
            signal = merge_dicts(
                load_data(signal_template_path(), {}).get("signal", {}),
                {
                    "signal_id": f"sig-{datetime.now(timezone.utc).strftime('%Y%m%d%H%M%S')}-blocked-{run_id}",
                    "kind": "blocked",
                    "source_type": "teammate",
                    "source_id": registry["runs"][index].get("owner_id", ""),
                    "run_id": run_id,
                    "mission_id": registry["runs"][index].get("mission_id", ""),
                    "applies_to": [registry["runs"][index].get("mission_id", "")],
                    "scope_type": "local",
                    "title": "Run blocked",
                    "summary": reason,
                    "authority_level": "L2",
                    "confidence": "high",
                    "fresh_until": now_iso(),
                },
            )
            fill_if_blank(signal, "emitted_at", now_iso())
            fill_if_blank(signal, "state", "active")
            log.setdefault("signals", []).append(signal)
            write_signal_log(log)
        print(json.dumps({"run": registry["runs"][index], "signal": signal}, indent=2, ensure_ascii=True))
        return 0

    if subcommand == "meeting":
        if len(args) < 3:
            raise SystemExit("drive meeting requires RUN_ID QUESTION")
        run_id = args[1]
        question = " ".join(args[2:]).strip()
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            validate_run_transition(registry["runs"][index], "waiting")
            backlog = read_maintenance_backlog()
            meeting_id = f"meeting-{datetime.now(timezone.utc).strftime('%Y%m%d%H%M%S')}-{registry['runs'][index].get('mission_id','')}"
            meeting = materialize_meeting(
                meeting_id=meeting_id,
                run=registry["runs"][index],
                question=question,
                requested_by=registry["runs"][index].get("owner_id", ""),
            )
            backlog.setdefault("meetings", []).append(
                {
                    "meeting_id": meeting_id,
                    "run_id": run_id,
                    "mission_id": registry["runs"][index].get("mission_id", ""),
                    "question": question,
                    "state": "active",
                    "meeting_ref": str(meeting_path(meeting_id)),
                    "created_at": meeting.get("created_at", now_iso()),
                    "updated_at": meeting.get("updated_at", now_iso()),
                }
            )
            write_maintenance_backlog(backlog)
            registry["runs"][index]["state"] = "waiting"
            registry["runs"][index]["status_reason"] = f"Meeting active: {question}"
            registry["runs"][index]["next_expected_action"] = "Chair the meeting and convert conclusions into run updates."
            registry["runs"][index]["updated_at"] = now_iso()
            write_registry(registry)
        print(
            json.dumps(
                {
                    "meeting": meeting,
                    "backlog_entry": backlog["meetings"][-1],
                    "run": registry["runs"][index],
                },
                indent=2,
                ensure_ascii=True,
            )
        )
        return 0

    if subcommand == "meeting-close":
        if len(args) < 3:
            raise SystemExit("drive meeting-close requires RUN_ID MEETING_ID [NOTE]")
        run_id = args[1]
        meeting_id = args[2]
        note = " ".join(args[3:]).strip() if len(args) > 3 else "Meeting closed"
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            data = read_meeting(meeting_id)
            if not data:
                raise SystemExit(f"meeting not found: {meeting_id}")
            meeting = data.get("meeting", {})
            outcome = meeting_outcome(meeting)
            meeting["state"] = "closed"
            meeting["updated_at"] = now_iso()
            meeting["closed_at"] = now_iso()
            meeting["resolution_note"] = note
            dump_data(meeting_path(meeting_id, archived=True), {"version": 1, "meeting": meeting})
            active_path = meeting_path(meeting_id)
            if active_path.exists():
                active_path.unlink()
            update_backlog_meeting(
                meeting_id,
                {
                    "state": "closed",
                    "closed_at": meeting["closed_at"],
                    "resolution_note": note,
                    "meeting_ref": str(meeting_path(meeting_id, archived=True)),
                    "updated_at": now_iso(),
                },
            )
            validate_run_transition(registry["runs"][index], "ready")
            registry["runs"][index]["state"] = "ready"
            registry["runs"][index]["status_reason"] = note
            if outcome["summary"]:
                registry["runs"][index]["latest_summary"] = outcome["summary"]
            registry["runs"][index]["next_expected_action"] = outcome["next_action"]
            registry["runs"][index]["updated_at"] = now_iso()
            write_registry(registry)
        print(json.dumps({"meeting": meeting, "run": registry["runs"][index]}, indent=2, ensure_ascii=True))
        return 0

    raise SystemExit(f"unknown drive subcommand: {subcommand}")


def cmd_ship(args: list[str]) -> int:
    ensure_runtime_dirs()
    if not args:
        raise SystemExit("ship requires a subcommand")
    subcommand = args[0]

    if subcommand == "status":
        if len(args) < 2:
            raise SystemExit("ship status requires RUN_ID")
        registry = read_registry()
        index = find_run_index(registry, args[1])
        if index < 0:
            raise SystemExit(f"run not found: {args[1]}")
        run = registry["runs"][index]
        summary = {
            "run_id": run.get("run_id", ""),
            "mission_id": run.get("mission_id", ""),
            "verification_status": run.get("verification_status", "unknown"),
            "verification_refs": run.get("verification_refs", []),
            "ship_status": run.get("ship_status", "unstarted"),
            "ship_refs": run.get("ship_refs", []),
            "pr_ref": run.get("pr_ref", ""),
            "pr_url": run.get("pr_url", ""),
            "rc_id": run.get("rc_id", ""),
            "next_expected_action": run.get("next_expected_action", ""),
        }
        print(json.dumps(summary, indent=2, ensure_ascii=True))
        return 0

    if subcommand == "verify":
        if len(args) < 3:
            raise SystemExit("ship verify requires RUN_ID REF")
        run_id = args[1]
        ref = " ".join(args[2:]).strip()
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            record_verification(registry["runs"][index], ref)
            registry["runs"][index]["latest_summary"] = "Verification recorded for ship readiness"
            registry["runs"][index]["next_expected_action"] = "Mark the run ship-ready or keep iterating if more work remains."
            write_registry(registry)
        print(json.dumps(registry["runs"][index], indent=2, ensure_ascii=True))
        return 0

    if subcommand == "ready":
        if len(args) < 2:
            raise SystemExit("ship ready requires RUN_ID [SUMMARY]")
        run_id = args[1]
        summary = " ".join(args[2:]).strip() if len(args) > 2 else "Run is ready to ship"
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            run = registry["runs"][index]
            if policy_value("admissions.require_pre_ship_validation", True) and run.get("verification_status") != "passed":
                raise SystemExit("ship ready requires verification_status=passed")
            run["ship_status"] = "ready"
            run["latest_summary"] = summary
            run["next_expected_action"] = "Push the branch and create or update the PR."
            run["updated_at"] = now_iso()
            write_registry(registry)
        print(json.dumps(registry["runs"][index], indent=2, ensure_ascii=True))
        return 0

    if subcommand == "push":
        if len(args) < 2:
            raise SystemExit("ship push requires RUN_ID [REMOTE]")
        run_id = args[1]
        remote = args[2] if len(args) > 2 else "origin"
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            run = registry["runs"][index]
            branch = run.get("branch", "")
            if not branch:
                raise SystemExit("ship push requires a run branch")
            try:
                run_command_with_stderr(["git", "push", "-u", remote, branch])
            except Exception as exc:
                raise SystemExit(f"ship push failed: {exc}")
            run["ship_status"] = "pushed"
            record_ship_ref(run, f"push:{remote}/{branch}")
            run["latest_summary"] = f"Branch pushed to {remote}: {branch}"
            run["next_expected_action"] = "Create or update the PR."
            write_registry(registry)
        print(json.dumps(registry["runs"][index], indent=2, ensure_ascii=True))
        return 0

    if subcommand == "pr-draft":
        if len(args) < 2:
            raise SystemExit("ship pr-draft requires RUN_ID [--title=...] [--base=...]")
        run_id = args[1]
        override_title = ""
        override_base = ""
        for arg in args[2:]:
            if arg.startswith("--title="):
                override_title = arg.split("=", 1)[1].strip()
            elif arg.startswith("--base="):
                override_base = arg.split("=", 1)[1].strip()
        registry = read_registry()
        index = find_run_index(registry, run_id)
        if index < 0:
            raise SystemExit(f"run not found: {run_id}")
        run = registry["runs"][index]
        base_branch = override_base or get_key(resolve_profile()[0], "profile.base_branch", "main")
        payload = pr_payload_for_run(run, base_branch, override_title)
        payload["run_id"] = run_id
        payload["mode"] = "draft"
        print(json.dumps(payload, indent=2, ensure_ascii=True))
        return 0

    if subcommand == "pr-open":
        if len(args) < 2:
            raise SystemExit("ship pr-open requires RUN_ID [--title=...] [--base=...] [--draft] [--reviewer=...]")
        run_id = args[1]
        override_title = ""
        override_base = ""
        reviewers: list[str] = []
        draft = False
        for arg in args[2:]:
            if arg.startswith("--title="):
                override_title = arg.split("=", 1)[1].strip()
            elif arg.startswith("--base="):
                override_base = arg.split("=", 1)[1].strip()
            elif arg.startswith("--reviewer="):
                reviewer = arg.split("=", 1)[1].strip()
                if reviewer:
                    reviewers.append(reviewer)
            elif arg == "--draft":
                draft = True
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            run = registry["runs"][index]
            if run.get("ship_status", "unstarted") not in ("ready", "pushed", "in_review"):
                raise SystemExit("ship pr-open requires ship ready or pushed state")
            base_branch = override_base or get_key(resolve_profile()[0], "profile.base_branch", "main")
            payload = pr_payload_for_run(run, base_branch, override_title)
            try:
                pr_result = open_github_pr(
                    run=run,
                    base_branch=base_branch,
                    title=payload["title"],
                    body=payload["body"],
                    draft=draft,
                    reviewers=reviewers,
                )
            except Exception as exc:
                raise SystemExit(f"ship pr-open failed: {exc}")
            run["ship_status"] = "in_review"
            run["pr_ref"] = pr_result["pr_ref"]
            run["pr_url"] = pr_result["pr_url"]
            record_ship_ref(run, f"pr:{pr_result['pr_ref']}")
            run["latest_summary"] = f"PR opened via GitHub: {pr_result['pr_ref']}"
            run["next_expected_action"] = "Track review, CI, and merge readiness."
            write_registry(registry)
        print(
            json.dumps(
                {
                    "run": registry["runs"][index],
                    "pr": {
                        "title": payload["title"],
                        "base_branch": payload["base_branch"],
                        "pr_ref": pr_result["pr_ref"],
                        "pr_url": pr_result["pr_url"],
                        "draft": draft,
                        "reviewers": reviewers,
                    },
                },
                indent=2,
                ensure_ascii=True,
            )
        )
        return 0

    if subcommand == "pr":
        if len(args) < 3:
            raise SystemExit("ship pr requires RUN_ID PR_REF [PR_URL]")
        run_id = args[1]
        pr_ref = args[2]
        pr_url = args[3] if len(args) > 3 else ""
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            run = registry["runs"][index]
            if run.get("ship_status", "unstarted") == "unstarted":
                raise SystemExit("ship pr requires ship ready state first")
            run["ship_status"] = "in_review"
            run["pr_ref"] = pr_ref
            if pr_url:
                run["pr_url"] = pr_url
            record_ship_ref(run, f"pr:{pr_ref}")
            run["latest_summary"] = f"PR opened or updated: {pr_ref}"
            run["next_expected_action"] = "Track review, CI, and merge readiness."
            write_registry(registry)
        print(json.dumps(registry["runs"][index], indent=2, ensure_ascii=True))
        return 0

    if subcommand == "merged":
        if len(args) < 2:
            raise SystemExit("ship merged requires RUN_ID [SUMMARY]")
        run_id = args[1]
        summary = " ".join(args[2:]).strip() if len(args) > 2 else "PR merged"
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            run = registry["runs"][index]
            if policy_value("admissions.require_pre_ship_validation", True) and run.get("verification_status") != "passed":
                raise SystemExit("ship merged requires verification_status=passed")
            run["ship_status"] = "merged"
            run["shipped_at"] = now_iso()
            record_ship_ref(run, "merged")
            if run.get("state") != "completed":
                validate_run_transition(run, "completed")
                run["state"] = "completed"
                run["completed_at"] = now_iso()
            run["latest_summary"] = summary
            run["next_expected_action"] = "Consider whether this shipped work should enter an RC."
            write_registry(registry)
        print(json.dumps(registry["runs"][index], indent=2, ensure_ascii=True))
        return 0

    if subcommand == "cleanup":
        if len(args) < 2:
            raise SystemExit("ship cleanup requires RUN_ID [--delete-branch] [--delete-worktree]")
        run_id = args[1]
        delete_branch = "--delete-branch" in args[2:]
        delete_worktree = "--delete-worktree" in args[2:]
        if not delete_branch and not delete_worktree:
            delete_worktree = True
        with harness_lock():
            registry = read_registry()
            index = find_run_index(registry, run_id)
            if index < 0:
                raise SystemExit(f"run not found: {run_id}")
            run = registry["runs"][index]
            profile, _ = resolve_profile()
            cleaned: list[str] = []
            if delete_worktree and run.get("worktree_path"):
                try:
                    git_remove_worktree(Path(run["worktree_path"]))
                    cleaned.append(f"worktree:{run['worktree_path']}")
                    run["worktree_path"] = ""
                except Exception as exc:
                    raise SystemExit(f"ship cleanup failed removing worktree: {exc}")
            if delete_branch and run.get("branch"):
                try:
                    git_delete_branch(run["branch"], get_key(profile, "profile.base_branch", "main"))
                    cleaned.append(f"branch:{run['branch']}")
                except Exception as exc:
                    raise SystemExit(f"ship cleanup failed deleting branch: {exc}")
            for ref in cleaned:
                record_ship_ref(run, f"cleanup:{ref}")
            run["latest_summary"] = "Post-merge cleanup completed" if cleaned else "No cleanup action was needed"
            run["next_expected_action"] = "Run is cleanly delivered." if run.get("ship_status") == "merged" else run.get("next_expected_action", "")
            run["updated_at"] = now_iso()
            update_contract_from_run(run)
            write_registry(registry)
        print(json.dumps({"run": registry["runs"][index], "cleaned": cleaned}, indent=2, ensure_ascii=True))
        return 0

    raise SystemExit(f"unknown ship subcommand: {subcommand}")


def cmd_rc(args: list[str]) -> int:
    ensure_runtime_dirs()
    if not args:
        raise SystemExit("rc requires a subcommand")
    subcommand = args[0]

    if subcommand == "list":
        include_archived = "--all" in args[1:]
        print(json.dumps(list_rcs(include_archived=include_archived), indent=2, ensure_ascii=True))
        return 0

    if subcommand == "get":
        if len(args) < 2:
            raise SystemExit("rc get requires RC_ID")
        data = read_rc(args[1])
        if not data:
            raise SystemExit(f"rc not found: {args[1]}")
        print(json.dumps(data.get("rc", {}), indent=2, ensure_ascii=True))
        return 0

    if subcommand == "create":
        if len(args) < 3:
            raise SystemExit("rc create requires VERSION RUN_ID...")
        version = args[1]
        run_ids = [arg for arg in args[2:] if not arg.startswith("--")]
        apply_changes = "--apply" in args[2:]
        if not run_ids:
            raise SystemExit("rc create requires at least one RUN_ID")
        with harness_lock():
            if active_rc():
                raise SystemExit("an active RC already exists")
            profile, _ = resolve_profile()
            registry = read_registry()
            selected_runs = []
            for run_id in run_ids:
                index = find_run_index(registry, run_id)
                if index < 0:
                    raise SystemExit(f"run not found: {run_id}")
                run = registry["runs"][index]
                if run.get("ship_status") != "merged":
                    raise SystemExit(f"run is not merged and cannot enter RC: {run_id}")
                selected_runs.append((index, run))
            branch_pattern = get_key(profile, "profile.rc_branch_pattern", "rc/{version}")
            branch = branch_pattern.replace("{version}", version)
            base_branch = get_key(profile, "profile.base_branch", "main")
            production_branch = get_key(profile, "profile.production_branch", base_branch)
            rc_id = f"rc-{datetime.now(timezone.utc).strftime('%Y%m%d%H%M%S')}-{slugify(version)}"
            if apply_changes:
                git_checkout_new_branch(branch, base_branch)
            rc = materialize_rc(
                rc_id=rc_id,
                version=version,
                branch=branch,
                base_branch=base_branch,
                production_branch=production_branch,
                run_ids=run_ids,
            )
            note = {
                "recorded_at": now_iso(),
                "summary": f"RC created for version {version}",
            }
            rc.setdefault("notes", []).append(note)
            rc["updated_at"] = now_iso()
            dump_data(rc_path(rc_id), {"version": 1, "rc": rc})
            for index, run in selected_runs:
                registry["runs"][index]["rc_id"] = rc_id
                registry["runs"][index]["next_expected_action"] = f"Verify active RC {rc_id} and decide on promote or fix."
                registry["runs"][index]["updated_at"] = now_iso()
            write_registry(registry)
        print(json.dumps(rc, indent=2, ensure_ascii=True))
        return 0

    if subcommand == "verify":
        if len(args) < 3:
            raise SystemExit("rc verify requires RC_ID REF")
        rc_id = args[1]
        ref = " ".join(args[2:]).strip()
        with harness_lock():
            data = read_active_rc(rc_id)
            if not data:
                raise SystemExit(f"active rc not found: {rc_id}")
            rc = data.get("rc", {})
            refs = rc.setdefault("verification_refs", [])
            append_unique_text(refs, ref)
            rc["verification_status"] = "passed"
            rc["updated_at"] = now_iso()
            rc.setdefault("notes", []).append({"recorded_at": now_iso(), "summary": f"Verification recorded: {ref}"})
            dump_data(rc_path(rc_id), {"version": 1, "rc": rc})
        print(json.dumps(rc, indent=2, ensure_ascii=True))
        return 0

    if subcommand in ("promote", "close"):
        if len(args) < 2:
            raise SystemExit(f"rc {subcommand} requires RC_ID [NOTE]")
        rc_id = args[1]
        note = " ".join(args[2:]).strip() if len(args) > 2 else ("RC promoted" if subcommand == "promote" else "RC closed")
        with harness_lock():
            data = read_active_rc(rc_id)
            if not data:
                raise SystemExit(f"active rc not found: {rc_id}")
            rc = data.get("rc", {})
            if subcommand == "promote" and policy_value("admissions.require_rc_validation", True) and rc.get("verification_status") != "passed":
                raise SystemExit("rc promote requires verification_status=passed")
            rc["state"] = "promoted" if subcommand == "promote" else "closed"
            rc["updated_at"] = now_iso()
            rc["closed_at"] = now_iso()
            if subcommand == "promote":
                rc["promoted_at"] = now_iso()
            rc.setdefault("notes", []).append({"recorded_at": now_iso(), "summary": note})
            dump_data(rc_path(rc_id, archived=True), {"version": 1, "rc": rc})
            active_path = rc_path(rc_id)
            if active_path.exists():
                active_path.unlink()
        print(json.dumps(rc, indent=2, ensure_ascii=True))
        return 0

    raise SystemExit(f"unknown rc subcommand: {subcommand}")


def cmd_meetings(args: list[str]) -> int:
    ensure_runtime_dirs()
    if not args:
        raise SystemExit("meetings requires a subcommand")
    subcommand = args[0]
    if subcommand == "list":
        include_archived = "--all" in args[1:]
        print(json.dumps(list_meetings(include_archived=include_archived), indent=2, ensure_ascii=True))
        return 0
    if subcommand == "active":
        meetings = [meeting for meeting in list_meetings() if meeting.get("state") == "active"]
        mission_id = args[1] if len(args) > 1 and not args[1].startswith("--") else None
        if mission_id:
            meetings = [meeting for meeting in meetings if meeting.get("mission_id") == mission_id]
        print(json.dumps(meetings, indent=2, ensure_ascii=True))
        return 0
    if subcommand == "get":
        if len(args) < 2:
            raise SystemExit("meetings get requires MEETING_ID")
        data = read_meeting(args[1])
        if not data:
            raise SystemExit(f"meeting not found: {args[1]}")
        print(json.dumps(data.get("meeting", {}), indent=2, ensure_ascii=True))
        return 0
    if subcommand == "summary":
        if len(args) < 2:
            raise SystemExit("meetings summary requires MEETING_ID")
        data = read_meeting(args[1])
        if not data:
            raise SystemExit(f"meeting not found: {args[1]}")
        meeting = data.get("meeting", {})
        outcome = meeting_outcome(meeting)
        summary = {
            "meeting_id": meeting.get("meeting_id", ""),
            "state": meeting.get("state", ""),
            "question": meeting.get("question", ""),
            "outcome_summary": outcome["summary"],
            "next_action": outcome["next_action"],
            "accepted_count": len(meeting.get("findings", {}).get("accepted", [])),
            "rejected_count": len(meeting.get("findings", {}).get("rejected", [])),
            "unresolved_count": len(meeting.get("findings", {}).get("unresolved", [])),
            "open_action_count": len(
                [
                    item
                    for item in meeting.get("action_items", [])
                    if not isinstance(item, dict) or item.get("state", "open") == "open"
                ]
            ),
        }
        print(json.dumps(summary, indent=2, ensure_ascii=True))
        return 0
    if subcommand == "add-role":
        if len(args) < 3:
            raise SystemExit("meetings add-role requires MEETING_ID ROLE")
        meeting_id = args[1]
        role = " ".join(args[2:]).strip()
        with harness_lock():
            data = read_active_meeting(meeting_id)
            if not data:
                raise SystemExit(f"active meeting not found: {meeting_id}")
            meeting = data.get("meeting", {})
            roles = meeting.setdefault("roles", [])
            if role not in roles:
                roles.append(role)
            write_active_meeting(meeting)
        print(json.dumps(meeting, indent=2, ensure_ascii=True))
        return 0
    if subcommand == "add-context":
        if len(args) < 3:
            raise SystemExit("meetings add-context requires MEETING_ID REF")
        meeting_id = args[1]
        ref = " ".join(args[2:]).strip()
        with harness_lock():
            data = read_active_meeting(meeting_id)
            if not data:
                raise SystemExit(f"active meeting not found: {meeting_id}")
            meeting = data.get("meeting", {})
            packet = meeting.setdefault("packet", {})
            refs = packet.setdefault("context_refs", [])
            if ref not in refs:
                refs.append(ref)
            write_active_meeting(meeting)
        print(json.dumps(meeting, indent=2, ensure_ascii=True))
        return 0
    if subcommand == "set-expected-output":
        if len(args) < 3:
            raise SystemExit("meetings set-expected-output requires MEETING_ID TEXT")
        meeting_id = args[1]
        expected_output = " ".join(args[2:]).strip()
        with harness_lock():
            data = read_active_meeting(meeting_id)
            if not data:
                raise SystemExit(f"active meeting not found: {meeting_id}")
            meeting = data.get("meeting", {})
            packet = meeting.setdefault("packet", {})
            packet["expected_output"] = expected_output
            write_active_meeting(meeting)
        print(json.dumps(meeting, indent=2, ensure_ascii=True))
        return 0
    if subcommand == "add-finding":
        if len(args) < 4:
            raise SystemExit("meetings add-finding requires MEETING_ID accepted|rejected|unresolved TEXT")
        meeting_id = args[1]
        finding_kind = args[2]
        summary = " ".join(args[3:]).strip()
        if finding_kind not in ("accepted", "rejected", "unresolved"):
            raise SystemExit("meetings add-finding kind must be accepted, rejected, or unresolved")
        with harness_lock():
            data = read_active_meeting(meeting_id)
            if not data:
                raise SystemExit(f"active meeting not found: {meeting_id}")
            meeting = data.get("meeting", {})
            findings = meeting.setdefault("findings", {})
            bucket = findings.setdefault(finding_kind, [])
            bucket.append(
                {
                    "finding_id": f"finding-{datetime.now(timezone.utc).strftime('%Y%m%d%H%M%S')}",
                    "summary": summary,
                    "recorded_at": now_iso(),
                }
            )
            write_active_meeting(meeting)
        print(json.dumps(meeting, indent=2, ensure_ascii=True))
        return 0
    if subcommand == "add-action":
        if len(args) < 3:
            raise SystemExit("meetings add-action requires MEETING_ID TEXT")
        meeting_id = args[1]
        action_text = " ".join(args[2:]).strip()
        with harness_lock():
            data = read_active_meeting(meeting_id)
            if not data:
                raise SystemExit(f"active meeting not found: {meeting_id}")
            meeting = data.get("meeting", {})
            actions = meeting.setdefault("action_items", [])
            actions.append(
                {
                    "action_id": f"action-{datetime.now(timezone.utc).strftime('%Y%m%d%H%M%S')}",
                    "summary": action_text,
                    "recorded_at": now_iso(),
                    "state": "open",
                }
            )
            write_active_meeting(meeting)
        print(json.dumps(meeting, indent=2, ensure_ascii=True))
        return 0
    if subcommand == "resolve-action":
        if len(args) < 3:
            raise SystemExit("meetings resolve-action requires MEETING_ID ACTION_ID [STATE]")
        meeting_id = args[1]
        action_id = args[2]
        action_state = args[3] if len(args) > 3 else "done"
        if action_state not in ("done", "cancelled", "open"):
            raise SystemExit("meetings resolve-action state must be one of: open, done, cancelled")
        with harness_lock():
            data = read_meeting(meeting_id)
            if not data:
                raise SystemExit(f"meeting not found: {meeting_id}")
            meeting = data.get("meeting", {})
            target = None
            for item in meeting.setdefault("action_items", []):
                if isinstance(item, dict) and item.get("action_id") == action_id:
                    target = item
                    break
            if target is None:
                raise SystemExit(f"action not found: {action_id}")
            target["state"] = action_state
            target["updated_at"] = now_iso()
            if action_state in ("done", "cancelled"):
                target["closed_at"] = now_iso()
            if meeting_path(meeting_id).exists():
                write_active_meeting(meeting)
            else:
                meeting["updated_at"] = now_iso()
                dump_data(meeting_path(meeting_id, archived=True), {"version": 1, "meeting": meeting})
        print(json.dumps(meeting, indent=2, ensure_ascii=True))
        return 0
    raise SystemExit(f"unknown meetings subcommand: {subcommand}")


def cmd_maintenance(args: list[str]) -> int:
    ensure_runtime_dirs()
    if not args:
        raise SystemExit("maintenance requires a subcommand")
    subcommand = args[0]
    if subcommand == "normalize":
        with harness_lock():
            registry = read_registry()
            log = read_signal_log()
            normalized_runs: list[str] = []
            normalized_signals: list[str] = []

            for run in registry.get("runs", []):
                changed = False
                timestamp = run.get("updated_at", "") or now_iso()
                if not run.get("created_at"):
                    run["created_at"] = timestamp
                    changed = True
                if run.get("state") in ("ready", "executing", "waiting", "blocked") and not run.get("started_at"):
                    run["started_at"] = run.get("created_at", timestamp)
                    changed = True
                if run.get("state") in ("ready", "executing", "waiting", "blocked") and not run.get("last_reported_at"):
                    run["last_reported_at"] = run.get("updated_at", "") or run.get("started_at", "") or now_iso()
                    changed = True
                if changed:
                    run["updated_at"] = now_iso()
                    normalized_runs.append(run.get("run_id", ""))

            for signal in log.get("signals", []):
                if not signal.get("emitted_at"):
                    signal["emitted_at"] = signal.get("fresh_until", "") or now_iso()
                    normalized_signals.append(signal.get("signal_id", ""))

            if normalized_runs:
                write_registry(registry)
            if normalized_signals:
                write_signal_log(log)

            summary = {
                "normalized_run_ids": [item for item in normalized_runs if item],
                "normalized_signal_ids": [item for item in normalized_signals if item],
                "run_count": len(normalized_runs),
                "signal_count": len(normalized_signals),
            }
        print(json.dumps(summary, indent=2, ensure_ascii=True))
        return 0
    if subcommand != "reconcile":
        raise SystemExit(f"unknown maintenance subcommand: {subcommand}")

    with harness_lock():
        registry = read_registry()
        log = read_signal_log()
        active_run_ids = {run.get("run_id", "") for run in registry.get("runs", []) if active_run(run)}

        reconciled_signals: list[str] = []
        for signal in log.get("signals", []):
            if signal.get("state") != "active":
                continue
            run_id = signal.get("run_id", "")
            if run_id and run_id not in active_run_ids:
                signal["state"] = "resolved"
                signal["resolution_note"] = "Reconciled orphaned signal: referenced run is missing."
                signal["resolved_at"] = now_iso()
                reconciled_signals.append(signal.get("signal_id", ""))

        reconciled_meetings: list[str] = []
        for meeting in list_meetings():
            if meeting.get("state") != "active":
                continue
            run_id = meeting.get("run_id", "")
            if run_id and run_id not in active_run_ids:
                meeting["state"] = "closed"
                meeting["updated_at"] = now_iso()
                meeting["closed_at"] = now_iso()
                meeting["resolution_note"] = "Reconciled orphaned meeting: referenced run is missing."
                dump_data(meeting_path(meeting["meeting_id"], archived=True), {"version": 1, "meeting": meeting})
                active_path = meeting_path(meeting["meeting_id"])
                if active_path.exists():
                    active_path.unlink()
                update_backlog_meeting(
                    meeting["meeting_id"],
                    {
                        "state": "closed",
                        "closed_at": meeting["closed_at"],
                        "resolution_note": meeting["resolution_note"],
                        "meeting_ref": str(meeting_path(meeting["meeting_id"], archived=True)),
                        "updated_at": now_iso(),
                    },
                )
                reconciled_meetings.append(meeting.get("meeting_id", ""))

        if reconciled_signals:
            write_signal_log(log)

        summary = {
            "reconciled_signal_ids": [item for item in reconciled_signals if item],
            "reconciled_meeting_ids": [item for item in reconciled_meetings if item],
            "signal_count": len(reconciled_signals),
            "meeting_count": len(reconciled_meetings),
        }
    print(json.dumps(summary, indent=2, ensure_ascii=True))
    return 0


def main() -> int:
    if len(sys.argv) < 2:
        raise SystemExit("usage: harness_state.py <paths|profile|runs|signals|contracts|home|doctor|issue|claim|drive|ship|rc|meetings|maintenance> ...")
    command = sys.argv[1]
    args = sys.argv[2:]
    if command == "paths":
        return cmd_paths(args)
    if command == "profile":
        return cmd_profile(args)
    if command == "runs":
        return cmd_runs(args)
    if command == "signals":
        return cmd_signals(args)
    if command == "contracts":
        return cmd_contracts(args)
    if command == "home":
        return cmd_home(args)
    if command == "doctor":
        return cmd_doctor(args)
    if command == "issue":
        return cmd_issue(args)
    if command == "claim":
        return cmd_claim(args)
    if command == "drive":
        return cmd_drive(args)
    if command == "ship":
        return cmd_ship(args)
    if command == "rc":
        return cmd_rc(args)
    if command == "meetings":
        return cmd_meetings(args)
    if command == "maintenance":
        return cmd_maintenance(args)
    raise SystemExit(f"unknown command: {command}")


if __name__ == "__main__":
    sys.exit(main())
