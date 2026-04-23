#!/usr/bin/env python3
"""Generate a comprehensive Markdown evaluation report for a drive-skill mission.

Reads eval.json (v2.3) or eval-data.json (v2.2 legacy) and evidence/ from
a mission's grading/ directory. Outputs a terminal-readable Markdown report.

Usage:
    python generate_eval_report.py <mission-dir>
    python generate_eval_report.py <mission-dir> --output <path>
    python generate_eval_report.py <mission-dir> --stdout  # print to terminal
"""

import argparse
import json
import sys
from datetime import datetime, timezone
from pathlib import Path


def load_json(path: Path) -> dict | list | None:
    if path.exists():
        try:
            return json.loads(path.read_text())
        except (json.JSONDecodeError, OSError):
            pass
    return None


def load_eval_data(mission_dir: Path) -> dict:
    # Try v2.3 format first, fall back to v2.2 legacy
    return (load_json(mission_dir / "grading" / "eval.json")
            or load_json(mission_dir / "grading" / "eval-data.json")
            or {})


def load_comparison(mission_dir: Path) -> dict | None:
    """Load or generate comparison data."""
    grading = mission_dir / "grading"

    # Try inline aggregation
    try:
        from aggregate_eval import build_comparison
        result = build_comparison(grading)
        if result and result.get("iteration_count", 0) >= 1:
            return result
    except ImportError:
        pass

    # Fall back to pre-generated file
    return load_json(grading / "eval-comparison.json")


def load_eval_reports(mission_dir: Path) -> list[tuple[str, str]]:
    """Load all eval-N.md reports. Returns [(name, content), ...]"""
    grading = mission_dir / "grading"
    reports = []
    if grading.is_dir():
        for f in sorted(grading.iterdir()):
            if f.is_file() and f.name.startswith("eval-") and f.suffix == ".md":
                try:
                    reports.append((f.name, f.read_text(errors="replace")))
                except OSError:
                    pass
    return reports


def list_dir_files(directory: Path) -> list[str]:
    if directory.is_dir():
        return sorted(f.name for f in directory.iterdir() if f.is_file() and not f.name.startswith("."))
    return []


def fmt_delta(val: int | float) -> str:
    if val > 0:
        return f"**+{val}** ↑"
    elif val < 0:
        return f"**{val}** ↓"
    return "0"


def fmt_delta_str(s: str) -> str:
    """Format a delta string like '+3' or '-2' with arrow."""
    try:
        n = float(s)
        if n > 0:
            return f"**{s}** ↑"
        elif n < 0:
            return f"**{s}** ↓"
    except (ValueError, TypeError):
        pass
    return str(s)


def generate_report(mission_dir: Path) -> str:
    eval_data = load_eval_data(mission_dir)
    comparison = load_comparison(mission_dir)
    contract_path = mission_dir / "contract.md"
    contract = contract_path.read_text(errors="replace") if contract_path.exists() else ""
    evidence_files = list_dir_files(mission_dir / "grading" / "evidence")
    # v2.3: screenshots are in evidence/; v2.2 legacy: separate screenshots/
    screenshot_files = list_dir_files(mission_dir / "grading" / "screenshots")
    eval_reports = load_eval_reports(mission_dir)

    mission_name = eval_data.get("mission", mission_dir.name)
    version = eval_data.get("version", "?")
    timestamp = eval_data.get("timestamp", "")
    iteration = eval_data.get("iteration", "?")
    verdict = eval_data.get("verdict", "UNKNOWN")
    verdict_reason = eval_data.get("verdict_reason", "")

    features = eval_data.get("features", [])
    dims = eval_data.get("quality_dimensions", [])
    scenarios = eval_data.get("scenarios", [])

    passed_features = sum(1 for f in features if f.get("passed"))
    total_features = len(features)
    passed_scenarios = sum(1 for s in scenarios if s.get("status") == "pass")
    total_scenarios = len(scenarios)

    lines = []

    # ── Header ──
    verdict_icon = {"PASS": "✅", "FAIL": "❌", "PARTIAL": "⚠️"}.get(verdict, "❓")
    lines.append(f"# Evaluation Report — {mission_name}")
    lines.append("")
    lines.append(f"> **Version:** {version} | **Iteration:** {iteration} | **Time:** {timestamp}")
    lines.append(f"> **Verdict:** {verdict_icon} **{verdict}** — {verdict_reason}")
    lines.append("")

    # ── Summary ──
    lines.append("## Summary")
    lines.append("")
    lines.append(f"| Metric | Value |")
    lines.append(f"|--------|-------|")
    lines.append(f"| Features Passed | {passed_features}/{total_features} |")
    lines.append(f"| Scenarios Passed | {passed_scenarios}/{total_scenarios} |")
    lines.append(f"| Evidence Files | {len(evidence_files)} |")
    lines.append(f"| Screenshots | {len(screenshot_files)} |")
    lines.append(f"| Verdict | {verdict_icon} {verdict} |")
    lines.append("")

    # ── Quality Scores ──
    if dims:
        lines.append("## Quality Scores")
        lines.append("")
        lines.append("| Dimension | Score | Threshold | Status | Evidence |")
        lines.append("|-----------|-------|-----------|--------|----------|")
        for d in dims:
            score = d.get("score", 0)
            threshold = d.get("threshold", 0)
            ok = score >= threshold
            status = "✅ PASS" if ok else "❌ FAIL"
            ev = d.get("evidence", "")
            lines.append(f"| {d.get('name', '?')} | **{score}**/10 | ≥{threshold} | {status} | {ev} |")
        lines.append("")

    # ── Features ──
    if features:
        lines.append("## Features")
        lines.append("")
        lines.append("| ID | Title | Status |")
        lines.append("|----|-------|--------|")
        for f in features:
            status = "✅ Pass" if f.get("passed") else "❌ Fail"
            lines.append(f"| {f.get('id', '?')} | {f.get('title', '?')} | {status} |")
        lines.append("")

        # Detailed acceptance criteria
        for f in features:
            acc = f.get("acceptance", [])
            if acc:
                status_icon = "✅" if f.get("passed") else "❌"
                lines.append(f"### {f.get('id')}: {f.get('title')} {status_icon}")
                lines.append("")
                if f.get("description"):
                    lines.append(f"{f['description']}")
                    lines.append("")
                lines.append("**Acceptance Criteria:**")
                for a in acc:
                    lines.append(f"- {a}")
                if f.get("verify_command"):
                    lines.append(f"- Verify: `{f['verify_command']}`")
                lines.append("")

    # ── Scenarios ──
    if scenarios:
        lines.append("## Scenarios")
        lines.append("")
        lines.append("| Scenario | Status | Notes |")
        lines.append("|----------|--------|-------|")
        for s in scenarios:
            status_map = {"pass": "✅ Pass", "fail": "❌ Fail", "partial": "⚠️ Partial"}
            status = status_map.get(s.get("status", ""), s.get("status", "?"))
            lines.append(f"| {s.get('name', '?')} | {status} | {s.get('notes', '')} |")
        lines.append("")

    # ══════════════════════════════════════════════════════════════════
    # ── Cross-Iteration Comparison ──
    # ══════════════════════════════════════════════════════════════════
    if comparison and comparison.get("iteration_count", 0) >= 2:
        iters = comparison.get("iterations", [])
        f_matrix = comparison.get("feature_matrix", {})
        d_trend = comparison.get("dimension_trend", {})
        d_thresholds = comparison.get("dimension_thresholds", {})
        s_prog = comparison.get("scenario_progression", {})
        delta = comparison.get("delta", {})

        lines.append("---")
        lines.append("")
        lines.append("## Cross-Iteration Comparison")
        lines.append("")

        # Verdict progression
        verdict_flow = " → ".join(
            f"**Iter {it['iteration']}: {it['verdict']}**" for it in iters
        )
        lines.append(f"**Verdict Progression:** {verdict_flow}")
        lines.append("")

        # Summary delta table
        if len(iters) >= 2:
            prev_it = iters[-2]
            curr_it = iters[-1]

            lines.append("### Summary Delta")
            lines.append("")
            lines.append(f"| Metric | Iter {prev_it['iteration']} | Iter {curr_it['iteration']} | Delta |")
            lines.append("|--------|------------|------------|-------|")

            # Features
            pf = prev_it["features"]["passed"]
            cf = curr_it["features"]["passed"]
            pt = prev_it["features"]["total"]
            ct = curr_it["features"]["total"]
            lines.append(f"| Features Passed | {pf}/{pt} | {cf}/{ct} | {fmt_delta(cf - pf)} |")

            # Scenarios
            ps = prev_it["scenarios"]["passed"]
            cs = curr_it["scenarios"]["passed"]
            pst = prev_it["scenarios"]["total"]
            cst = curr_it["scenarios"]["total"]
            lines.append(f"| Scenarios Passed | {ps}/{pst} | {cs}/{cst} | {fmt_delta(cs - ps)} |")

            # Quality dimensions
            prev_dims = {d["name"]: d["score"] for d in prev_it.get("quality_dimensions", [])}
            curr_dims = {d["name"]: d["score"] for d in curr_it.get("quality_dimensions", [])}
            for name in dict.fromkeys(list(prev_dims) + list(curr_dims)):
                pv = prev_dims.get(name)
                cv = curr_dims.get(name)
                pstr = f"{pv}/10" if pv is not None else "—"
                cstr = f"{cv}/10" if cv is not None else "—"
                dstr = fmt_delta(cv - pv) if pv is not None and cv is not None else "—"
                lines.append(f"| {name} | {pstr} | {cstr} | {dstr} |")

            # Verdict
            lines.append(f"| Verdict | {prev_it['verdict']} | {curr_it['verdict']} | {delta.get('verdict', '')} |")
            lines.append("")

            # Feature flips callout
            flipped_pass = delta.get("features_flipped_to_pass", [])
            flipped_fail = delta.get("features_flipped_to_fail", [])
            if flipped_pass:
                lines.append(f"**🟢 Newly Passing:** {', '.join(f'`{f}`' for f in flipped_pass)}")
                lines.append("")
            if flipped_fail:
                lines.append(f"**🔴 Regressions:** {', '.join(f'`{f}`' for f in flipped_fail)}")
                lines.append("")

        # Feature matrix
        if f_matrix:
            lines.append("### Feature Matrix")
            lines.append("")
            iter_headers = " | ".join(f"Iter {it['iteration']}" for it in iters)
            lines.append(f"| Feature | {iter_headers} |")
            lines.append(f"|---------|{'|'.join('---------' for _ in iters)}|")

            # Build title lookup
            title_lookup = {f.get("id", ""): f.get("title", "") for f in features}

            for fid, row in f_matrix.items():
                title = title_lookup.get(fid, "")
                label = f"**{fid}** {title}" if title else f"**{fid}**"
                cells = []
                for val in row:
                    if val is True:
                        cells.append("✅")
                    elif val is False:
                        cells.append("❌")
                    else:
                        cells.append("—")
                lines.append(f"| {label} | {' | '.join(cells)} |")
            lines.append("")

        # Quality dimension trend
        if d_trend:
            lines.append("### Quality Dimension Trend")
            lines.append("")
            iter_headers = " | ".join(f"Iter {it['iteration']}" for it in iters)
            lines.append(f"| Dimension | {iter_headers} | Threshold |")
            lines.append(f"|-----------|{'|'.join('---------' for _ in iters)}|-----------|")

            for name, scores in d_trend.items():
                threshold = d_thresholds.get(name, "?")
                cells = []
                for s in scores:
                    if s is None:
                        cells.append("—")
                    elif s >= (threshold if isinstance(threshold, (int, float)) else 0):
                        cells.append(f"**{s}** ✅")
                    else:
                        cells.append(f"**{s}** ❌")
                lines.append(f"| {name} | {' | '.join(cells)} | ≥{threshold} |")
            lines.append("")

        # Scenario progression
        if s_prog:
            lines.append("### Scenario Progression")
            lines.append("")
            iter_headers = " | ".join(f"Iter {it['iteration']}" for it in iters)
            lines.append(f"| Scenario | {iter_headers} |")
            lines.append(f"|----------|{'|'.join('---------' for _ in iters)}|")

            status_icons = {"pass": "✅", "fail": "❌", "partial": "⚠️"}
            for name, row in s_prog.items():
                cells = [status_icons.get(v, "—") if v else "—" for v in row]
                lines.append(f"| {name} | {' | '.join(cells)} |")
            lines.append("")

    # ── Evidence ──
    if evidence_files:
        lines.append("## Evidence Artifacts")
        lines.append("")
        for f in evidence_files:
            lines.append(f"- `grading/evidence/{f}`")
        lines.append("")

    # ── Screenshots ──
    if screenshot_files:
        lines.append("## Screenshots")
        lines.append("")
        for f in screenshot_files:
            lines.append(f"- `grading/screenshots/{f}`")
        lines.append("")

    # ── Appended Eval Reports ──
    if eval_reports:
        lines.append("---")
        lines.append("")
        lines.append("## Full Evaluation Reports")
        lines.append("")
        for name, content in eval_reports:
            lines.append(f"<details><summary>{name}</summary>")
            lines.append("")
            lines.append(content)
            lines.append("")
            lines.append("</details>")
            lines.append("")

    return "\n".join(lines)


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Generate Markdown evaluation report for a drive-skill mission"
    )
    parser.add_argument("mission_dir", type=Path, help="Path to mission directory")
    parser.add_argument(
        "--output", "-o", type=Path, default=None,
        help="Output path (default: <mission_dir>/grading/eval-report.md)"
    )
    parser.add_argument(
        "--stdout", action="store_true",
        help="Print to stdout instead of writing file"
    )
    args = parser.parse_args()

    mission_dir = args.mission_dir.resolve()
    if not mission_dir.is_dir():
        print(f"Error: {mission_dir} is not a directory", file=sys.stderr)
        sys.exit(1)

    report = generate_report(mission_dir)

    if args.stdout:
        print(report)
    else:
        output_path = args.output or (mission_dir / "grading" / "eval-report.md")
        output_path.parent.mkdir(parents=True, exist_ok=True)
        output_path.write_text(report)
        print(f"Generated: {output_path}")


if __name__ == "__main__":
    main()
