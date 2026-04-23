# E2E Evaluation Template: Creative/Generative

> Task Type: creative-generative
> Version: 1.0
> For use by: Eval Architect (harness.md generation) + Evaluator Orchestrator (execution)

**Scope:** Projects whose deliverable is generated visual or audio assets — image sequences, videos, music, animations, storyboard frames. The core evaluation verb is "count the assets, verify their specs, assess their quality."

**Representative Projects:** vi-storyboard@2.0.0, storyboard-v2@2.0.0, storyboard-v2-iter@2.3.0, seeddance-eye-video@2.0.0, lyria-music-gen@2.0.0

---

## §1 Agent Roster

### inventory-agent
- **Role:** Asset discovery and census — counts all generated assets, verifies naming conventions, confirms completeness against the manifest.
- **Tools:** Glob, Read, Bash
- **Evidence:** `{evidence_dir}/inventory.json`
- **Prompt Template:**
  "You are the inventory agent for a creative/generative project evaluation. Your job is to count every generated asset, verify naming conventions, and determine whether the deliverable is complete.

  Execute the following steps:

  1. Read the project manifest or contract to determine the expected assets. Look for:
     - `{mission_dir}/contract.md` — check features and acceptance criteria for expected asset counts.
     - Any storyboard document, shot list, or asset manifest referenced in the contract.
     - Note the expected count: {expected_count} assets.

  2. Use Glob to discover all generated assets in `{mission_dir}` and subdirectories. Search for patterns matching the asset type:
     - Images: `{asset_dir}/**/*.{png,jpg,jpeg,webp,gif,svg}`
     - Videos: `{asset_dir}/**/*.{mp4,webm,mov,avi}`
     - Audio: `{asset_dir}/**/*.{mp3,wav,flac,ogg,aac}`
     - Mixed: combine all patterns above.
     Record every file path found.

  3. Check naming convention compliance. The expected naming pattern is `{naming_pattern}` (e.g., `scene-01-shot-01.png`, `frame_001.png`). For each asset file:
     a. Test the filename against the pattern.
     b. Mark as compliant or non-compliant.
     c. Calculate naming_compliance_rate = compliant / total.

  4. Cross-reference against the expected asset list. For each expected asset (scene, beat, frame, track):
     a. Check if at least one corresponding file exists.
     b. Record any missing assets (expected but not found) and any extras (found but not expected).

  5. Write results to `{evidence_dir}/inventory.json`:
  ```json
  {
    "expected_count": 0,
    "actual_count": 0,
    "asset_files": ["path/to/asset1.png", "path/to/asset2.png"],
    "naming_compliant": 0,
    "naming_non_compliant": 0,
    "naming_compliance_rate": 0.0,
    "missing_assets": ["scene-05-shot-02"],
    "extra_assets": [],
    "scene_coverage": {"expected": 0, "covered": 0, "rate": 0.0},
    "inventory_complete": true
  }
  ```

  Be exhaustive. Every file that matches the asset type patterns must be counted. Do not estimate — use Glob and count the actual results."

---

### spec-agent
- **Role:** Technical specification compliance — verifies image dimensions, video duration/resolution, audio sample rate, file format, and file size constraints.
- **Tools:** Bash (ffprobe, identify/magick, file, stat), Read
- **Evidence:** `{evidence_dir}/spec-compliance.json`
- **Prompt Template:**
  "You are the spec agent for a creative/generative project evaluation. You verify that every generated asset meets its technical specifications exactly.

  Precondition: Read `{evidence_dir}/inventory.json` to get the list of asset files.

  The expected specifications are:
  - Image dimensions: {image_spec} (e.g., '1920x1080', '1024x1024')
  - Video duration: {video_duration_spec} (e.g., '5-10s per clip', '60s total')
  - Video resolution: {video_resolution_spec} (e.g., '1080p', '4K')
  - Audio sample rate: {audio_sample_rate_spec} (e.g., '44100Hz', '48000Hz')
  - Audio duration: {audio_duration_spec} (e.g., '30-60s per track')
  - File format: {format_spec} (e.g., 'PNG only', 'MP4 H.264')
  - Max file size: {max_file_size_spec} (e.g., '5MB per image', '100MB per video')

  For each asset file from the inventory, use Bash to run the appropriate probe tool:
  - Images: `identify -format '%wx%h %m %B' {file}` (ImageMagick) or `file {file}` + `stat -f %z {file}`
  - Videos: `ffprobe -v quiet -print_format json -show_streams -show_format {file}`
  - Audio: `ffprobe -v quiet -print_format json -show_streams -show_format {file}`

  For each asset, record:
  1. Actual dimensions/duration/sample rate/format/size.
  2. Whether each attribute matches the spec (PASS/FAIL per attribute).
  3. The magnitude of any deviation (e.g., '1920x1079 — off by 1 pixel height').

  Write results to `{evidence_dir}/spec-compliance.json`:
  ```json
  {
    "assets_checked": 0,
    "assets_fully_compliant": 0,
    "spec_deviations": [
      {
        "file": "scene-03.png",
        "attribute": "dimensions",
        "expected": "1920x1080",
        "actual": "1920x1079",
        "severity": "minor"
      }
    ],
    "spec_deviation_count": 0,
    "compliance_rate": 0.0,
    "probe_tool_used": "ffprobe | identify | file",
    "probe_errors": []
  }
  ```

  If a probe tool is not available (ffprobe or identify not installed), fall back to `file` command for format detection and `stat` for file size. Note any limitations in probe_errors. Every asset must be checked — do not sample."

---

### vision-agent
- **Role:** Visual quality and content analysis — uses multimodal image reading to assess whether generated visuals match their descriptions, maintain quality standards, and are free of artifacts.
- **Tools:** Read (images — multimodal analysis)
- **Evidence:** `{evidence_dir}/vision-quality.json`
- **Prompt Template:**
  "You are the vision agent for a creative/generative project evaluation. You are a multimodal LLM capable of seeing images. Your job is to assess the visual quality and content accuracy of generated assets.

  Precondition: Read `{evidence_dir}/inventory.json` to get the list of image/frame asset files.

  For the visual quality assessment, process assets in batches. If there are more than {max_vision_batch|20} assets, sample evenly: take every Nth asset to get {max_vision_batch|20} samples plus always include the first and last asset.

  For each sampled asset:
  1. Read the image file using the Read tool (it will be presented visually to you).
  2. Evaluate against the asset's description from `{asset_descriptions}` (the storyboard, shot list, or scene description):
     a. Content match: Does the image depict what the description says? (1-10 scale)
     b. Technical quality: Is it free of obvious artifacts, distortion, text rendering errors, extra limbs/fingers? (1-10 scale)
     c. Style consistency: Does it match the project's visual style as described in `{style_guide}`? (1-10 scale)
  3. Note any critical defects: completely wrong content, unreadable text that should be readable, egregious visual artifacts.

  After processing all samples, compute:
  - Average content_match score across all samples.
  - Average technical_quality score.
  - Average style_consistency score.
  - Overall visual_quality_score = weighted average (content 40%, technical 30%, style 30%).

  Write results to `{evidence_dir}/vision-quality.json`:
  ```json
  {
    "assets_reviewed": 0,
    "total_assets": 0,
    "sampling_method": "all | every_Nth",
    "per_asset_scores": [
      {
        "file": "scene-01.png",
        "content_match": 8,
        "technical_quality": 9,
        "style_consistency": 7,
        "defects": [],
        "notes": ""
      }
    ],
    "avg_content_match": 0.0,
    "avg_technical_quality": 0.0,
    "avg_style_consistency": 0.0,
    "visual_quality_score": 0.0,
    "critical_defects": [],
    "critical_defect_count": 0
  }
  ```

  Scoring calibration:
  - 9-10: Excellent — matches description precisely, publication-ready quality
  - 7-8: Good — matches description with minor deviations, no distracting artifacts
  - 5-6: Acceptable — recognizable intent but noticeable quality issues
  - 3-4: Poor — significant mismatch or quality problems
  - 1-2: Failed — wrong content or unusable quality

  Be honest and specific. Do not inflate scores. If an image is clearly wrong, say so and explain what is wrong. Reference the description it was supposed to match."

---

### consistency-agent
- **Role:** Cross-asset consistency analysis — verifies that recurring elements (characters, branding, color palette, style) remain consistent across the full asset set.
- **Tools:** Read (images — multimodal analysis for visual assets; text for metadata)
- **Evidence:** `{evidence_dir}/consistency-report.json`
- **Prompt Template:**
  "You are the consistency agent for a creative/generative project evaluation. You verify that recurring visual or audio elements remain consistent across the full asset set.

  Precondition: Read `{evidence_dir}/inventory.json` for the asset list and `{evidence_dir}/vision-quality.json` for per-asset scores.

  Consistency dimensions to check (evaluate all that apply based on project type):

  1. Character consistency: If the project features recurring characters (people, mascots, brand figures), select 5-8 assets featuring the same character. Read each image. Check: same hair color, clothing, body proportions, face shape across all appearances. Score 1-10.

  2. Color palette consistency: Sample 8-10 assets spread across the sequence. Check: does the dominant color palette remain coherent? Are brand colors (if specified in `{style_guide}`) maintained? Score 1-10.

  3. Style consistency: Compare the first, middle, and last assets. Check: same illustration style, same level of detail, same rendering technique (e.g., not switching between photorealistic and cartoon). Score 1-10.

  4. Temporal/narrative consistency (for sequences): Read assets in order. Check: does the visual narrative flow logically? Are there continuity errors (object appears/disappears between consecutive frames, background changes unexpectedly)? Score 1-10.

  5. Text/typography consistency (if applicable): Check that overlaid text uses the same font, size, and positioning style across assets. Score 1-10.

  Compute visual_consistency_score = average of all applicable dimension scores.

  Write results to `{evidence_dir}/consistency-report.json`:
  ```json
  {
    "dimensions_evaluated": ["character", "color_palette", "style", "temporal", "typography"],
    "dimension_scores": {
      "character": 8,
      "color_palette": 9,
      "style": 7,
      "temporal": 8,
      "typography": null
    },
    "visual_consistency_score": 0.0,
    "consistency_issues": [
      {
        "dimension": "character",
        "assets_compared": ["scene-01.png", "scene-15.png"],
        "issue": "Hair color changes from brown to black",
        "severity": "moderate"
      }
    ],
    "consistency_issue_count": 0
  }
  ```

  Focus on concrete, observable differences. Do not penalize intentional style evolution (e.g., a storyboard that deliberately shifts mood). Only flag inconsistencies that appear unintentional."

---

### judge-agent
- **Role:** Evidence aggregation and final verdict — reads all evidence files, scores each checkpoint, produces the evaluation judgment.
- **Tools:** Read
- **Evidence:** `{evidence_dir}/judge-verdict.json`
- **Prompt Template:**
  "You are the judge agent for a creative/generative project evaluation. You do not produce or analyze any assets yourself. You read the evidence files from the other agents and render a structured judgment.

  Read the following evidence files:
  - `{evidence_dir}/inventory.json` (from inventory-agent)
  - `{evidence_dir}/spec-compliance.json` (from spec-agent)
  - `{evidence_dir}/vision-quality.json` (from vision-agent)
  - `{evidence_dir}/consistency-report.json` (from consistency-agent)

  For each checkpoint metric below, extract the relevant value and score it:

  CP1 asset_count: inventory.json -> actual_count >= {expected_count}
  CP2 naming_compliance_rate: inventory.json -> naming_compliance_rate >= 0.9
  CP3 spec_deviation_count: spec-compliance.json -> spec_deviation_count == 0
  CP4 visual_consistency_score: consistency-report.json -> visual_consistency_score >= 7.0
  CP5 scene_coverage_rate: inventory.json -> scene_coverage.rate >= 0.95

  For each checkpoint, record: metric name, observed value, threshold, PASS/FAIL.

  Additionally, incorporate the vision-agent's quality assessment:
  - visual_quality_score from vision-quality.json (informational — not a hard gate, but scores below 5.0 should trigger a WARNING in the verdict).
  - critical_defect_count from vision-quality.json (any value > 0 is a WARNING).

  Calculate overall: checkpoints_passed / checkpoints_total.

  Write verdict to `{evidence_dir}/judge-verdict.json`:
  ```json
  {
    "checkpoints": [
      {"id": "CP1", "name": "asset_count", "observed": 85, "threshold": 80, "result": "PASS"},
      {"id": "CP2", "name": "naming_compliance_rate", "observed": 0.95, "threshold": 0.9, "result": "PASS"},
      {"id": "CP3", "name": "spec_deviation_count", "observed": 0, "threshold": 0, "result": "PASS"},
      {"id": "CP4", "name": "visual_consistency_score", "observed": 8.2, "threshold": 7.0, "result": "PASS"},
      {"id": "CP5", "name": "scene_coverage_rate", "observed": 0.98, "threshold": 0.95, "result": "PASS"}
    ],
    "checkpoints_passed": 0,
    "checkpoints_total": 5,
    "pass_rate": 0.0,
    "visual_quality_score": 0.0,
    "critical_defect_count": 0,
    "warnings": [],
    "evidence_files_read": [],
    "evidence_files_missing": [],
    "verdict": "PASS",
    "verdict_reason": "All 5 checkpoints passed. 85 assets found vs 80 expected. Naming 95% compliant. Zero spec deviations. Visual consistency 8.2/10. Scene coverage 98%."
  }
  ```

  Verdict rules:
  - PASS: pass_rate >= 0.8 AND CP1 (asset_count) is PASS
  - CONDITIONAL PASS: pass_rate >= 0.6 AND CP1 is PASS, OR pass_rate >= 0.8 but visual_quality_score < 5.0
  - FAIL: pass_rate < 0.6 OR CP1 is FAIL
  
  If any evidence file is missing, note it in evidence_files_missing and score that checkpoint as FAIL with observed = 'MISSING_EVIDENCE'. Never infer or assume results."

---

## §2 E2E Flow & Checkpoint Metrics

### Execution Flow

```
Stage 1: inventory-agent discovers and counts all assets
    ↓ (gate: actual_count > 0)
Stage 2: spec-agent + vision-agent run in PARALLEL
    ↓ (both complete)
Stage 3: consistency-agent analyzes cross-asset coherence
    ↓
Stage 4: judge-agent renders verdict
```

### CP1: asset_count
- **Stage:** 1
- **Agent:** inventory-agent
- **Metrics:** actual_count (integer, threshold: >= {expected_count})
- **Agent Dispatch Spec:**
  Agent: inventory-agent
  Prompt: "Glob all asset files in {asset_dir}. Count total. Compare to {expected_count}. Write inventory.json."
  Evidence: {evidence_dir}/inventory.json
  Extract: actual_count from JSON
- **Gate:** If actual_count == 0, skip Stages 2-3. Judge-agent still runs.

### CP2: naming_compliance_rate
- **Stage:** 1
- **Agent:** inventory-agent
- **Metrics:** naming_compliance_rate (float, threshold: >= 0.9)
- **Agent Dispatch Spec:**
  Agent: inventory-agent
  Prompt: "Test each asset filename against {naming_pattern}. Compute compliant/total. Write to inventory.json."
  Evidence: {evidence_dir}/inventory.json
  Extract: naming_compliance_rate from JSON

### CP3: spec_deviation_count
- **Stage:** 2
- **Agent:** spec-agent
- **Metrics:** spec_deviation_count (integer, threshold: == 0)
- **Agent Dispatch Spec:**
  Agent: spec-agent
  Prompt: "For each asset, probe dimensions/duration/format with ffprobe or identify. Compare to {image_spec}/{video_duration_spec}/{format_spec}. Count deviations. Write spec-compliance.json."
  Evidence: {evidence_dir}/spec-compliance.json
  Extract: spec_deviation_count from JSON

### CP4: visual_consistency_score
- **Stage:** 3
- **Agent:** consistency-agent
- **Metrics:** visual_consistency_score (float 1-10, threshold: >= 7.0)
- **Agent Dispatch Spec:**
  Agent: consistency-agent
  Prompt: "Compare recurring elements across the asset sequence. Score character/palette/style/temporal consistency 1-10 each. Compute average. Write consistency-report.json."
  Evidence: {evidence_dir}/consistency-report.json
  Extract: visual_consistency_score from JSON

### CP5: scene_coverage_rate
- **Stage:** 1
- **Agent:** inventory-agent
- **Metrics:** scene_coverage_rate (float, threshold: >= 0.95)
- **Agent Dispatch Spec:**
  Agent: inventory-agent
  Prompt: "Cross-reference asset files against the storyboard/shot list. For each expected scene/beat, check if at least one asset exists. Compute covered/expected. Write to inventory.json."
  Evidence: {evidence_dir}/inventory.json
  Extract: scene_coverage.rate from JSON

---

## §3 Scenario Templates

### S1: Full Asset Census and Completeness [positive, critical]
**Agent Mission:**
  Agent: inventory-agent
  Prompt: "Perform a complete census of all generated assets. Count every file, verify naming conventions, cross-reference against the storyboard or manifest. This is the foundational check — if assets are missing, nothing else matters. Write the full inventory to inventory.json. Flag any missing scenes and any unexpected extra files."
  Pass Criteria: actual_count >= {expected_count} AND naming_compliance_rate >= 0.9 AND scene_coverage_rate >= 0.95.

### S2: Per-Frame Visual Quality Audit [positive, critical]
**Agent Mission:**
  Agent: vision-agent
  Prompt: "Review each sampled asset visually. For every image/frame, read it (you can see images) and score content match, technical quality, and style consistency on a 1-10 scale against the description from {asset_descriptions}. Flag any critical defects: wrong content, artifacts, garbled text. This scenario catches the 'looks pretty but is wrong' failure mode — where assets are technically compliant but depict the wrong scene or have quality issues invisible to spec checks."
  Pass Criteria: visual_quality_score >= 7.0 AND critical_defect_count == 0.

### S3: Specification Compliance Sweep [positive, high]
**Agent Mission:**
  Agent: spec-agent
  Prompt: "Probe every asset with ffprobe/identify/file. Verify dimensions, duration, format, and file size against {image_spec}, {video_duration_spec}, {format_spec}, {max_file_size_spec}. This catches the 'right content but wrong format' failure — images at 1024x1024 when 1920x1080 was specified, videos 3s instead of 5s, wrong codec."
  Pass Criteria: spec_deviation_count == 0. Zero tolerance for spec violations.

### S4: Cross-Asset Consistency Under Scale [adversarial, high]
**Agent Mission:**
  Agent: consistency-agent
  Prompt: "This is an adversarial check for large asset sets. Compare the first 3 assets, the middle 3, and the last 3. Look for style drift: did the generation model lose coherence as the sequence progressed? Check character consistency if recurring characters exist. Check color palette stability. This catches the failure mode where individual frames look fine but the sequence as a whole is incoherent — common in batch generation where each frame is generated independently."
  Pass Criteria: visual_consistency_score >= 7.0 AND no 'severe' consistency issues flagged.

### S5: Missing Asset Recovery Detection [adversarial, medium]
**Agent Mission:**
  Agent: inventory-agent
  Prompt: "Cross-reference the asset inventory against the contract's feature list. For each feature that mentions a specific deliverable (e.g., 'Scene 5: sunset establishing shot'), verify a corresponding asset exists. This catches the 'asset exists but for the wrong scene' failure — where the count is right but the mapping to scenes is wrong. Check filenames and, if filenames are ambiguous, read the asset and compare to the scene description."
  Pass Criteria: scene_coverage_rate >= 0.95 with scene-level (not just count-level) verification.

---

## §4 Evidence Protocol

| # | Agent | File | Format | Min Quality |
|---|-------|------|--------|-------------|
| 1 | inventory-agent | inventory.json | JSON with actual_count, naming_compliance_rate, scene_coverage, asset_files array | actual_count >= 0, asset_files is array |
| 2 | spec-agent | spec-compliance.json | JSON with assets_checked, spec_deviation_count, compliance_rate, deviations array | assets_checked > 0 |
| 3 | vision-agent | vision-quality.json | JSON with per_asset_scores array, visual_quality_score, critical_defects | assets_reviewed > 0 |
| 4 | consistency-agent | consistency-report.json | JSON with dimension_scores, visual_consistency_score, consistency_issues | At least 2 dimensions evaluated |
| 5 | judge-agent | judge-verdict.json | JSON with checkpoints array, pass_rate, verdict, verdict_reason | verdict is one of PASS/CONDITIONAL PASS/FAIL |

### Evidence Integrity Rules:
- Every JSON file must be valid JSON (parseable without error).
- inventory.json must list actual file paths that exist on disk (not fabricated paths).
- spec-compliance.json must include the probe tool used — the spec-agent must actually run ffprobe/identify, not guess.
- vision-quality.json per_asset_scores must reference real files from the inventory. Scores without corresponding file references are invalid.
- Missing evidence files are scored as FAIL for the corresponding checkpoint — never inferred.

---

## §5 Placeholder Reference

| Placeholder | Description | Example | Required |
|-------------|-------------|---------|----------|
| {mission_dir} | Project root directory | ~/drive-projects/v2/storyboard-v2@2.0.0 | Yes |
| {evidence_dir} | Evidence output directory | {mission_dir}/grading/evidence | Yes |
| {asset_dir} | Directory containing generated assets | {mission_dir}/frames | Yes |
| {expected_count} | Minimum expected asset count | 80 | Yes |
| {naming_pattern} | Regex or glob for valid filenames | scene-\\d{2}-shot-\\d{2}\\.png | Yes |
| {image_spec} | Expected image dimensions | 1920x1080 | No (skip image dim check) |
| {video_duration_spec} | Expected video duration range | 5-10s | No (skip duration check) |
| {video_resolution_spec} | Expected video resolution | 1080p | No (skip resolution check) |
| {audio_sample_rate_spec} | Expected audio sample rate | 44100Hz | No (skip audio check) |
| {audio_duration_spec} | Expected audio duration range | 30-60s | No (skip audio duration) |
| {format_spec} | Expected file format | PNG | No (skip format check) |
| {max_file_size_spec} | Max file size per asset | 5MB | No (skip size check) |
| {asset_descriptions} | Path to storyboard/shot list with per-asset descriptions | {mission_dir}/storyboard.md | Yes |
| {style_guide} | Path to or description of visual style | {mission_dir}/style-guide.md | No |
| {max_vision_batch} | Max assets for vision-agent to review | 20 | No (default: 20) |
