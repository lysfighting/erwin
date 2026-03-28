#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

case "${1:-}" in
  current)
    git branch --show-current
    ;;
  worktree-main-repo)
    git worktree list --porcelain 2>/dev/null | head -1 | sed 's/^worktree //'
    ;;
  worktree-find-by-branch)
    branch="${2:-}"
    if [ -z "$branch" ]; then
      echo "ERROR: worktree-find-by-branch requires BRANCH" >&2
      exit 1
    fi
    git worktree list --porcelain 2>/dev/null | awk -v b="$branch" '
      /^worktree / { wt = substr($0, 10) }
      /^branch refs\/heads\// { br = substr($0, 18); if (br == b) { print wt; exit 0 } }
    '
    ;;
  *)
    echo "Usage: harness-git.sh {current|worktree-main-repo|worktree-find-by-branch BRANCH}" >&2
    exit 1
    ;;
esac
