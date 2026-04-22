#!/usr/bin/env bash
# tw-pr.sh — PR lifecycle operations for teamwork skills
#
# USAGE:
#   bash tw-pr.sh exists BRANCH [BASE]              # Check if open PR exists, print number
#   bash tw-pr.sh commit-type TITLE                   # Infer feat/fix/refactor/test/docs
#   bash tw-pr.sh create ISSUE TITLE BODY BASE BRANCH # Create PR, return number + url
#   bash tw-pr.sh comment ISSUE TEXT                   # Post comment on Issue (non-fatal)
#   bash tw-pr.sh watch PR                               # Watch CI checks (blocks until done)
#   bash tw-pr.sh wait-merged PR [REPO]                 # Poll until PR is merged (10 min max)
#   bash tw-pr.sh add-reviewer PR REVIEWER             # Request review (non-fatal)
#   bash tw-pr.sh verify-merged BRANCH                 # Check if PR merged, print number+url
#   bash tw-pr.sh merged-since BASE [PREV_TAG] [--repo R] [--format F] [--version V] [--base-sha S] [--env E]
#                                                       # List merged PRs since PREV_TAG's date, categorized.
#                                                       # FMT: json (default) | markdown | feishu
#                                                       # Used by /team-rc Prepare (staging changelog) + Promote (release notes).
#
# EXIT CODES:
#   0 — success
#   1 — not found / usage error
#   2 — operation failed

set -euo pipefail

usage() {
  echo "Usage: tw-pr.sh <subcommand> [args...]" >&2
  echo "Subcommands: exists, commit-type, create, comment, watch, wait-merged," >&2
  echo "             add-reviewer, verify-merged, merged-since" >&2
  exit 1
}

cmd_exists() {
  local branch="${1:-}" base="${2:-}"
  if [ -z "$branch" ]; then
    echo "ERROR: exists requires BRANCH" >&2
    exit 1
  fi

  local -a base_args=()
  [ -n "$base" ] && base_args=(--base "$base")

  local result
  result=$(gh pr list --head "$branch" "${base_args[@]}" --state open --json number,url --limit 1 2>/dev/null) || {
    echo "ERROR: Could not query PRs" >&2
    exit 2
  }

  local count
  count=$(echo "$result" | jq 'length' 2>/dev/null || echo "0")

  if [ "$count" -gt 0 ]; then
    local number url
    number=$(echo "$result" | jq -r '.[0].number')
    url=$(echo "$result" | jq -r '.[0].url')
    echo "$number $url"
    return 0
  else
    return 1
  fi
}

cmd_commit_type() {
  local title="${1:-}"
  if [ -z "$title" ]; then
    echo "feat"  # default
    return 0
  fi

  local lower
  lower=$(echo "$title" | tr '[:upper:]' '[:lower:]')

  # Use word-boundary matching to avoid false positives (e.g., "prefix" → fix)
  if echo "$lower" | grep -qEw '(fix|bug|hotfix|patch)'; then
    echo "fix"
  elif echo "$lower" | grep -qEw '(refactor|clean|restructure)'; then
    echo "refactor"
  elif echo "$lower" | grep -qEw '(test|spec)'; then
    echo "test"
  elif echo "$lower" | grep -qEw '(doc|docs|readme)'; then
    echo "docs"
  else
    echo "feat"
  fi
}

cmd_create() {
  local issue="${1:-}" title="${2:-}" body="${3:-}" base="${4:-}" branch="${5:-}"
  if [ -z "$issue" ] || [ -z "$title" ] || [ -z "$base" ] || [ -z "$branch" ]; then
    echo "ERROR: create requires ISSUE, TITLE, BASE, BRANCH (BODY optional)" >&2
    exit 1
  fi

  # Safety: read base_branch from config and override if caller got it wrong
  local script_dir
  script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  local config_base
  config_base=$(bash "$script_dir/tw-config.sh" conventions.base_branch "" 2>/dev/null) || true
  if [ -n "$config_base" ] && [ "$config_base" != "$base" ]; then
    echo "WARNING: caller passed base='$base' but config says base_branch='$config_base'. Using config value." >&2
    base="$config_base"
  fi

  local url
  url=$(gh pr create \
    --title "$title" \
    --body "$body" \
    --base "$base" \
    --head "$branch" 2>/dev/null) || {
    echo "ERROR: PR creation failed" >&2
    exit 2
  }

  # gh pr create outputs the URL on stdout; extract PR number from last line
  local number
  url=$(echo "$url" | tail -1)
  number=$(echo "$url" | grep -oE '[0-9]+$' || echo "")

  echo "$number $url"
}

cmd_comment() {
  local issue="${1:-}" text="${2:-}"
  if [ -z "$issue" ] || [ -z "$text" ]; then
    echo "WARNING: comment requires ISSUE and TEXT, skipping" >&2
    return 0
  fi

  gh issue comment "$issue" --body "$text" 2>/dev/null || {
    echo "WARNING: Could not post comment to Issue #$issue" >&2
    return 0  # Non-fatal
  }
  echo "Comment posted to Issue #$issue"
}

cmd_watch() {
  local pr="${1:-}"
  if [ -z "$pr" ]; then
    echo "ERROR: watch requires PR number" >&2
    exit 1
  fi

  local interval=10
  if gh pr checks "$pr" --watch --interval "$interval" 2>&1; then
    echo "CI passed"
    return 0
  else
    echo "CI failed or timed out" >&2
    return 2
  fi
}

cmd_add_reviewer() {
  local pr="${1:-}" reviewer="${2:-}"
  if [ -z "$pr" ] || [ -z "$reviewer" ]; then
    echo "WARNING: add-reviewer requires PR and REVIEWER, skipping" >&2
    return 0
  fi

  gh pr edit "$pr" --add-reviewer "$reviewer" 2>/dev/null || {
    echo "WARNING: Could not add reviewer $reviewer to PR #$pr" >&2
    return 0
  }
  echo "Review requested from $reviewer on PR #$pr"
}

cmd_wait_merged() {
  local pr="${1:-}" repo="${2:-}"
  if [ -z "$pr" ]; then
    echo "ERROR: wait-merged requires PR number" >&2
    exit 1
  fi

  local -a repo_args=()
  [ -n "$repo" ] && repo_args=(--repo "$repo")

  local max_polls=60  # 10 minutes max (60 × 10s)
  local poll_count=0
  echo "⏳ Waiting for PR #$pr to merge..."

  while true; do
    local state
    state=$(gh pr view "$pr" --json state --jq '.state' "${repo_args[@]}" 2>/dev/null) || {
      echo "WARNING: Could not query PR state" >&2
      state="UNKNOWN"
    }

    if [ "$state" = "MERGED" ]; then
      echo "PR #$pr merged"
      return 0
    fi

    if [ "$state" = "CLOSED" ]; then
      echo "ERROR: PR #$pr was closed without merging" >&2
      echo "Check: gh pr view $pr ${repo_args[*]}" >&2
      return 1
    fi

    poll_count=$((poll_count + 1))
    if [ "$poll_count" -ge "$max_polls" ]; then
      echo "ERROR: Timed out waiting for PR #$pr to merge (10 min)" >&2
      echo "Check CI: gh pr checks $pr ${repo_args[*]}" >&2
      return 3
    fi

    local elapsed=$((poll_count * 10))
    echo "  ⏳ ${elapsed}s / 600s — state: $state"
    sleep 10
  done
}

cmd_verify_merged() {
  local branch="${1:-}"
  if [ -z "$branch" ]; then
    echo "ERROR: verify-merged requires BRANCH" >&2
    exit 1
  fi

  local result
  result=$(gh pr list --head "$branch" --state merged --json number,url --limit 1 2>/dev/null) || {
    echo "ERROR: Could not query PRs" >&2
    exit 2
  }

  local count
  count=$(echo "$result" | jq 'length' 2>/dev/null || echo "0")

  if [ "$count" -gt 0 ]; then
    local number url
    number=$(echo "$result" | jq -r '.[0].number')
    url=$(echo "$result" | jq -r '.[0].url')
    echo "$number $url"
    return 0
  else
    echo "No merged PR found for branch: $branch" >&2
    return 1
  fi
}

# ───────────────────────────────────────────────────────────────────
# merged-since — list merged PRs on BASE since PREV_TAG's date,
# categorized by conventional-commit prefix. Shared between
# /team-rc Prepare (staging changelog preview) and Promote (release notes).
#
# Three output formats:
#   json      — structured PR list with category field (AI/tool post-processing)
#   markdown  — GitHub-style release notes (# What's in this release, ### Features, etc.)
#   feishu    — 飞书-friendly rich-text with header, metadata block, categorized PRs
#
# Args (positional then named):
#   $1 BASE           — base branch (e.g. pre-launch, develop, main)
#   $2 PREV_TAG       — prior version tag (optional; if empty, includes all PRs back to repo start)
#   --repo REPO       — gh repo slug (default: detect from current)
#   --format FMT      — json | markdown | feishu (default json)
#   --version VER     — current version (for feishu header, optional)
#   --base-sha SHA    — base commit sha (for feishu metadata, optional)
#   --env ENV         — deployment env name for feishu (staging | production)
#   --release-url URL — GitHub Release URL for feishu (promote only)
#   --pr-number N     — release PR number for feishu (promote only)
#   --tag-sha SHA     — tagged merge commit for feishu (promote only)
# ───────────────────────────────────────────────────────────────────
cmd_merged_since() {
  local base="" prev_tag=""
  local repo="" format="json" version="" base_sha="" env=""
  local release_url="" pr_number="" tag_sha=""

  # Positional first, then flags. Accept empty "" as explicit "no prev_tag".
  if [ $# -gt 0 ] && [[ "${1:-}" != --* ]]; then
    base="${1:-}"; shift
  fi
  if [ $# -gt 0 ] && [[ "${1:-}" != --* ]]; then
    prev_tag="${1:-}"; shift
  fi
  while [ $# -gt 0 ]; do
    case "$1" in
      --repo)        repo="${2:-}"; shift 2 ;;
      --format)      format="${2:-json}"; shift 2 ;;
      --version)     version="${2:-}"; shift 2 ;;
      --base-sha)    base_sha="${2:-}"; shift 2 ;;
      --env)         env="${2:-}"; shift 2 ;;
      --release-url) release_url="${2:-}"; shift 2 ;;
      --pr-number)   pr_number="${2:-}"; shift 2 ;;
      --tag-sha)     tag_sha="${2:-}"; shift 2 ;;
      *) echo "ERROR: unknown flag '$1'" >&2; exit 1 ;;
    esac
  done

  if [ -z "$base" ]; then
    echo "ERROR: merged-since requires BASE" >&2
    exit 1
  fi
  case "$format" in
    json|markdown|feishu) ;;
    *) echo "ERROR: --format must be json | markdown | feishu (got '$format')" >&2; exit 1 ;;
  esac
  if [ -z "$repo" ]; then
    repo=$(gh repo view --json nameWithOwner --jq '.nameWithOwner' 2>/dev/null || echo "")
    if [ -z "$repo" ]; then
      echo "ERROR: --repo not given and gh could not detect (gh auth login?)" >&2
      exit 2
    fi
  fi

  # Derive PREV_DATE in UTC Z format (matches GitHub API mergedAt format).
  # Uses %at (author unix timestamp) + date -u -r/-d for locale-independent conversion.
  # Fix history: 3.13.0 ship fixed the %aI local-tz bug that silently dropped PRs merged
  # in the offset window (e.g. UTC 12:38 – local 20:38+08:00 vanished from notes).
  local prev_date="2000-01-01T00:00:00Z"
  if [ -n "$prev_tag" ]; then
    local prev_ts
    prev_ts=$(git log -1 --format='%at' "$prev_tag" 2>/dev/null)
    if [ -n "$prev_ts" ]; then
      prev_date=$(date -u -r "$prev_ts" '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null \
                 || date -u -d "@$prev_ts" '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null \
                 || echo "2000-01-01T00:00:00Z")
    fi
  fi

  # Fetch merged PRs (limit 300 handles large release windows; GitHub caps at 1000 via gh)
  local pr_json
  pr_json=$(gh pr list --base "$base" --state merged \
    --json number,title,author,mergedAt,url \
    --limit 300 --repo "$repo" 2>/dev/null) || {
    echo "ERROR: gh pr list failed (auth? network?)" >&2
    exit 2
  }

  # Categorize: feat:/feat(x): → Features; fix:/fix(x): → Bug Fixes; else Other Changes.
  # Filter mergedAt > prev_date (both UTC Z, valid instant comparison).
  local filtered
  filtered=$(echo "$pr_json" | jq --arg since "$prev_date" '
    [ .[] | select(.mergedAt > $since) | {
        number, title, url, mergedAt,
        author: .author.login,
        category: (
          if (.title | test("^feat(\\(.+\\))?!?:")) then "Features"
          elif (.title | test("^fix(\\(.+\\))?!?:")) then "Bug Fixes"
          else "Other Changes"
          end
        )
    } ]
  ') || {
    echo "ERROR: jq filter failed" >&2
    exit 2
  }

  case "$format" in
    json)     echo "$filtered" ;;
    markdown) _render_markdown "$filtered" "$prev_tag" "$version" ;;
    feishu)   _render_feishu   "$filtered" "$prev_tag" "$version" "$base" "$base_sha" "$env" "$release_url" "$pr_number" "$tag_sha" "$repo" ;;
  esac
}

# GitHub-style release notes (same as former P5 inline output)
_render_markdown() {
  local j="$1" prev_tag="$2" version="$3"
  local counts
  counts=$(echo "$j" | jq 'group_by(.category) | map({key: .[0].category, value: length}) | from_entries')

  echo "## What's in this release"
  echo ""
  for cat in "Features" "Bug Fixes" "Other Changes"; do
    local n
    n=$(echo "$counts" | jq --arg k "$cat" '.[$k] // 0')
    [ "$n" -eq 0 ] && continue
    local label
    case "$cat" in
      "Features")      label="### Features" ;;
      "Bug Fixes")     label="### Bug Fixes" ;;
      "Other Changes") label="### Other Changes" ;;
    esac
    echo "$label"
    echo "$j" | jq -r --arg c "$cat" '
      .[] | select(.category == $c) |
      "* \(.title) by @\(.author) in \(.url) (**#\(.number)**)"
    '
    echo ""
  done

  # Contributors
  local contribs
  contribs=$(echo "$j" | jq -r '
    group_by(.author) | map({a: .[0].author, n: length}) |
    sort_by(-.n) | map("@\(.a) (\(.n))") | join(" · ")
  ')
  if [ -n "$contribs" ] && [ "$contribs" != "" ]; then
    echo "### Contributors"
    echo "$contribs"
    echo ""
  fi

  if [ -n "$prev_tag" ] && [ -n "$version" ]; then
    echo "**Full Changelog**: ${prev_tag}...${version}"
  elif [ -z "$prev_tag" ]; then
    echo "**Initial release**"
  fi
}

# 飞书-friendly markdown (rich text renders in 飞书; plain text also readable)
_render_feishu() {
  local j="$1" prev_tag="$2" version="$3" base="$4" base_sha="$5" env="$6"
  local release_url="$7" pr_number="$8" tag_sha="$9" repo="${10}"

  local now_stamp
  now_stamp=$(date '+%Y-%m-%d %H:%M %z')
  # Pretty tz as +08:00 instead of +0800
  now_stamp=$(echo "$now_stamp" | sed -E 's/([+-][0-9]{2})([0-9]{2})$/\1:\2/')

  local header_emoji="🚀"
  local kind_word="Release"
  if [ "$env" = "production" ]; then
    header_emoji="🎉"
    kind_word="Production Release"
  elif [ "$env" = "staging" ]; then
    header_emoji="🚀"
    kind_word="Staging Release"
  fi

  local title="${version:-(no version)}"
  echo "# ${header_emoji} ${title} ${kind_word}"
  echo ""
  echo "**环境**: \`${env:-未知}\`"
  echo "**时间**: ${now_stamp}"
  if [ -n "$version" ] && [ -n "$base" ]; then
    if [ "$env" = "staging" ]; then
      echo "**分支**: \`rc/${version}\` ← \`${base}\`"
    else
      echo "**分支**: \`${base}\`"
    fi
  fi
  if [ -n "$base_sha" ]; then
    local short_sha="${base_sha:0:12}"
    echo "**基线 SHA**: \`${short_sha}\`"
  fi
  if [ -n "$tag_sha" ]; then
    local short_tag_sha="${tag_sha:0:12}"
    echo "**Tag commit**: \`${short_tag_sha}\`"
  fi
  if [ -n "$pr_number" ]; then
    echo "**Release PR**: #${pr_number}"
  fi
  if [ -n "$release_url" ]; then
    echo "**GitHub Release**: ${release_url}"
  fi
  if [ -n "$prev_tag" ] && [ -n "$version" ] && [ -n "$repo" ]; then
    echo "**Diff**: https://github.com/${repo}/compare/${prev_tag}...${version}"
  fi
  echo ""
  echo "## 本次变更"
  echo ""

  # Categories with emoji
  local counts
  counts=$(echo "$j" | jq 'group_by(.category) | map({key: .[0].category, value: length}) | from_entries')

  for cat in "Features" "Bug Fixes" "Other Changes"; do
    local n
    n=$(echo "$counts" | jq --arg k "$cat" '.[$k] // 0')
    [ "$n" -eq 0 ] && continue
    local label
    case "$cat" in
      "Features")      label="### ✨ Features" ;;
      "Bug Fixes")     label="### 🐛 Bug Fixes" ;;
      "Other Changes") label="### 📦 Other Changes" ;;
    esac
    echo "$label"
    # Prefer short URL-less list items for 飞书 density; Issue auto-links via #NNN
    echo "$j" | jq -r --arg c "$cat" '
      .[] | select(.category == $c) |
      "- \(.title) — @\(.author) (#\(.number))"
    '
    echo ""
  done

  # Contributors
  local contribs
  contribs=$(echo "$j" | jq -r '
    group_by(.author) | map({a: .[0].author, n: length}) |
    sort_by(-.n) | map("@\(.a) (\(.n))") | join(" · ")
  ')
  if [ -n "$contribs" ] && [ "$contribs" != "" ]; then
    echo "### 👥 Contributors"
    echo "$contribs"
    echo ""
  fi

  local total
  total=$(echo "$j" | jq 'length')
  echo "---"
  echo "合计 **${total}** PR 进入本次发布"
}

# --- Main ---
SUBCOMMAND="${1:-}"
shift || true

case "$SUBCOMMAND" in
  exists)         cmd_exists "$@" ;;
  commit-type)    cmd_commit_type "$@" ;;
  create)         cmd_create "$@" ;;
  comment)        cmd_comment "$@" ;;
  watch)          cmd_watch "$@" ;;
  wait-merged)    cmd_wait_merged "$@" ;;
  add-reviewer)   cmd_add_reviewer "$@" ;;
  verify-merged)  cmd_verify_merged "$@" ;;
  merged-since)   cmd_merged_since "$@" ;;
  *)              usage ;;
esac
