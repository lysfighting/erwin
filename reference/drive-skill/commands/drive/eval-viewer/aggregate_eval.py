#!/usr/bin/env python3
"""Aggregate eval iterations into comparison data.

Supports both v2.3 format (single eval.json with iterations_history)
and legacy v2.2 format (eval-data.json + eval-data-iter-*.json).

Outputs eval-comparison.json to the grading/ directory.

Usage:
    python aggregate_eval.py <mission-dir>
    python aggregate_eval.py <mission-dir> --output <path>

Can also be imported and used as a library:
    from aggregate_eval import build_comparison
"""

import argparse
import json
import sys
from pathlib import Path


def load_iterations(grading_dir: Path) -> list[dict]:
    """Load all eval iterations from grading/, sorted by iteration number.

    Supports two formats:
    - v2.3: eval.json with iterations_history array (cumulative)
    - v2.2 legacy: eval-data.json + eval-data-iter-*.json (multiple files)
    """
    iterations = []

    # Try v2.3 format first: eval.json with iterations_history
    eval_json = grading_dir / "eval.json"
    if eval_json.exists():
        try:
            data = json.loads(eval_json.read_text())
            if isinstance(data, dict):
                # The current iteration IS the main data
                iterations.append(data)
                # Also check iterations_history for previous iterations
                for hist in data.get("iterations_history", []):
                    # History entries are summaries; wrap them to match expected format
                    if isinstance(hist, dict) and "iteration" in hist:
                        iterations.append(hist)
        except (json.JSONDecodeError, OSError) as e:
            print(f"Warning: skipping eval.json: {e}", file=sys.stderr)

    # Fall back to v2.2 legacy format
    if not iterations:
        for f in sorted(grading_dir.glob("eval-data-iter-*.json")):
            try:
                data = json.loads(f.read_text())
                if isinstance(data, dict) and "iteration" in data:
                    iterations.append(data)
            except (json.JSONDecodeError, OSError) as e:
                print(f"Warning: skipping {f.name}: {e}", file=sys.stderr)

        current = grading_dir / "eval-data.json"
        if current.exists():
            try:
                data = json.loads(current.read_text())
                if isinstance(data, dict) and "iteration" in data:
                    iterations.append(data)
            except (json.JSONDecodeError, OSError) as e:
                print(f"Warning: skipping eval-data.json: {e}", file=sys.stderr)

    # Sort by iteration number, deduplicate
    iterations.sort(key=lambda d: d.get("iteration", 0))
    seen = {}
    for it in iterations:
        seen[it.get("iteration", 0)] = it
    return list(seen.values())


def build_feature_matrix(iterations: list[dict]) -> dict[str, list[bool | None]]:
    """Build feature pass/fail matrix across iterations.

    Returns: {"F1": [true, false, true], "F2": [null, true, true], ...}
    null means the feature didn't exist in that iteration.
    """
    # Collect all feature IDs across all iterations
    all_ids = []
    seen = set()
    for it in iterations:
        for f in it.get("features", []):
            fid = f.get("id", "")
            if fid and fid not in seen:
                all_ids.append(fid)
                seen.add(fid)

    matrix = {}
    for fid in all_ids:
        row = []
        for it in iterations:
            feature = next((f for f in it.get("features", []) if f.get("id") == fid), None)
            if feature is None:
                row.append(None)
            else:
                row.append(bool(feature.get("passed", False)))
        matrix[fid] = row

    return matrix


def build_dimension_trend(iterations: list[dict]) -> dict[str, list[float | None]]:
    """Build quality dimension score trend across iterations.

    Returns: {"Essence": [5, 7, 8], "Craft": [6, 6, 7], ...}
    """
    # Collect all dimension names
    all_names = []
    seen = set()
    for it in iterations:
        for d in it.get("quality_dimensions", []):
            name = d.get("name", "")
            if name and name not in seen:
                all_names.append(name)
                seen.add(name)

    trend = {}
    for name in all_names:
        row = []
        for it in iterations:
            dim = next((d for d in it.get("quality_dimensions", []) if d.get("name") == name), None)
            if dim is None:
                row.append(None)
            else:
                row.append(dim.get("score"))
        trend[name] = row

    return trend


def build_scenario_progression(iterations: list[dict]) -> dict[str, list[str | None]]:
    """Build scenario status progression across iterations.

    Returns: {"Login flow": ["fail", "pass", "pass"], ...}
    """
    all_names = []
    seen = set()
    for it in iterations:
        for s in it.get("scenarios", []):
            name = s.get("name", "")
            if name and name not in seen:
                all_names.append(name)
                seen.add(name)

    progression = {}
    for name in all_names:
        row = []
        for it in iterations:
            sc = next((s for s in it.get("scenarios", []) if s.get("name") == name), None)
            if sc is None:
                row.append(None)
            else:
                row.append(sc.get("status", "unknown"))
        progression[name] = row

    return progression


def compute_delta(iterations: list[dict], feature_matrix: dict) -> dict:
    """Compute delta between last two iterations.

    Positive delta = improvement. Convention:
    - features_passed: +N means more features pass
    - quality dimensions: +N means higher score
    - verdict: "FAIL -> PASS" means improvement
    """
    if len(iterations) < 2:
        return {}

    prev = iterations[-2]
    curr = iterations[-1]

    delta = {}

    # Feature delta
    prev_passed = sum(1 for f in prev.get("features", []) if f.get("passed"))
    curr_passed = sum(1 for f in curr.get("features", []) if f.get("passed"))
    diff = curr_passed - prev_passed
    delta["features_passed"] = f"+{diff}" if diff >= 0 else str(diff)

    # Quality dimension deltas
    prev_dims = {d["name"]: d["score"] for d in prev.get("quality_dimensions", []) if "name" in d and "score" in d}
    curr_dims = {d["name"]: d["score"] for d in curr.get("quality_dimensions", []) if "name" in d and "score" in d}
    for name in set(list(prev_dims.keys()) + list(curr_dims.keys())):
        if name in prev_dims and name in curr_dims:
            diff = round(curr_dims[name] - prev_dims[name], 2)
            delta[name] = f"+{diff}" if diff >= 0 else str(diff)

    # Scenario delta
    prev_sc_pass = sum(1 for s in prev.get("scenarios", []) if s.get("status") == "pass")
    curr_sc_pass = sum(1 for s in curr.get("scenarios", []) if s.get("status") == "pass")
    diff = curr_sc_pass - prev_sc_pass
    delta["scenarios_passed"] = f"+{diff}" if diff >= 0 else str(diff)

    # Verdict delta
    prev_verdict = prev.get("verdict", "UNKNOWN")
    curr_verdict = curr.get("verdict", "UNKNOWN")
    if prev_verdict != curr_verdict:
        delta["verdict"] = f"{prev_verdict} -> {curr_verdict}"
    else:
        delta["verdict"] = f"{curr_verdict} (unchanged)"

    # Feature flips
    flipped_to_pass = []
    flipped_to_fail = []
    for fid, row in feature_matrix.items():
        if len(row) >= 2 and row[-2] is not None and row[-1] is not None:
            if not row[-2] and row[-1]:
                flipped_to_pass.append(fid)
            elif row[-2] and not row[-1]:
                flipped_to_fail.append(fid)
    delta["features_flipped_to_pass"] = flipped_to_pass
    delta["features_flipped_to_fail"] = flipped_to_fail

    return delta


def build_comparison(grading_dir: Path) -> dict | None:
    """Build complete comparison data from all iterations in grading/.

    Returns None if fewer than 1 iteration found.
    """
    iterations_data = load_iterations(grading_dir)
    if not iterations_data:
        return None

    feature_matrix = build_feature_matrix(iterations_data)
    dimension_trend = build_dimension_trend(iterations_data)
    scenario_progression = build_scenario_progression(iterations_data)
    delta = compute_delta(iterations_data, feature_matrix)

    # Build iteration summaries
    iteration_summaries = []
    for it in iterations_data:
        features = it.get("features", [])
        scenarios = it.get("scenarios", [])
        dims = it.get("quality_dimensions", [])
        iteration_summaries.append({
            "iteration": it.get("iteration"),
            "timestamp": it.get("timestamp", ""),
            "verdict": it.get("verdict", "UNKNOWN"),
            "verdict_reason": it.get("verdict_reason", ""),
            "features": {
                "passed": sum(1 for f in features if f.get("passed")),
                "total": len(features),
            },
            "quality_dimensions": [
                {"name": d.get("name"), "score": d.get("score"), "threshold": d.get("threshold")}
                for d in dims
            ],
            "scenarios": {
                "passed": sum(1 for s in scenarios if s.get("status") == "pass"),
                "total": len(scenarios),
            },
        })

    # Build dimension thresholds (from most recent iteration)
    dim_thresholds = {}
    if iterations_data:
        for d in iterations_data[-1].get("quality_dimensions", []):
            if "name" in d and "threshold" in d:
                dim_thresholds[d["name"]] = d["threshold"]

    return {
        "iteration_count": len(iterations_data),
        "iterations": iteration_summaries,
        "feature_matrix": feature_matrix,
        "dimension_trend": dimension_trend,
        "dimension_thresholds": dim_thresholds,
        "scenario_progression": scenario_progression,
        "delta": delta,
    }


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Aggregate eval iterations into comparison data"
    )
    parser.add_argument("mission_dir", type=Path, help="Path to mission directory")
    parser.add_argument(
        "--output", "-o", type=Path, default=None,
        help="Output path (default: <mission_dir>/grading/eval-comparison.json)"
    )
    args = parser.parse_args()

    mission_dir = args.mission_dir.resolve()
    grading_dir = mission_dir / "grading"

    if not grading_dir.is_dir():
        print(f"Error: {grading_dir} is not a directory", file=sys.stderr)
        sys.exit(1)

    comparison = build_comparison(grading_dir)
    if comparison is None:
        print("No eval-data files found in grading/", file=sys.stderr)
        sys.exit(1)

    output_path = args.output or (grading_dir / "eval-comparison.json")
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(json.dumps(comparison, indent=2))
    print(f"Generated: {output_path}")

    # Print summary
    n = comparison["iteration_count"]
    delta = comparison.get("delta", {})
    print(f"\nIterations: {n}")
    if n >= 2:
        print(f"Features delta: {delta.get('features_passed', '—')}")
        print(f"Verdict: {delta.get('verdict', '—')}")
        flipped = delta.get("features_flipped_to_pass", [])
        regressed = delta.get("features_flipped_to_fail", [])
        if flipped:
            print(f"Newly passing: {', '.join(flipped)}")
        if regressed:
            print(f"Regressions: {', '.join(regressed)}")


if __name__ == "__main__":
    main()
