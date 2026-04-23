# E2E Evaluation Template: Runtime/Web System

> Task Type: runtime-web
> Version: 1.0
> For use by: Eval Architect (harness.md generation) + Evaluator Orchestrator (execution)

**Scope:** Projects whose deliverable is a startable, interactive runtime system — web services, APIs, WebSocket servers, desktop apps, HTML/CSS/JS sites. The core evaluation verb is "start it, use it, verify behavior."

**Representative Projects:** agent-core@2.3.0, aws-deploy@2.3.0, glasses-e2e-flow@2.0.0, realestate-pitch@2.0.0, eval-system@2.0.0

---

## §1 Agent Roster

### env-agent
- **Role:** System bootstrap and health verification — starts the system, confirms it is alive, and captures baseline state.
- **Tools:** Bash
- **Evidence:** `{evidence_dir}/env-health.json`, `{evidence_dir}/startup-log.txt`
- **Prompt Template:**
  "You are the environment agent for a runtime/web system evaluation. Your sole responsibility is to start the system under test and verify it reaches a healthy state.

  Execute the following sequence using Bash:

  1. Navigate to `{mission_dir}` and identify the startup command. Check for: package.json scripts (npm start, npm run dev), Makefile targets, shell scripts (start.sh, run.sh), Python entry points (main.py, app.py, manage.py runserver), or Docker Compose files. If {startup_command} is provided, use it directly.

  2. Start the system in background. Capture the startup log to `{evidence_dir}/startup-log.txt`. Record the start timestamp.

  3. Wait up to {startup_timeout_s|30} seconds for the health endpoint to respond. Poll `{health_endpoint}` every 2 seconds using curl. Record time-to-healthy.

  4. Once healthy (HTTP 200), capture the full response body. If the system does not reach healthy state within the timeout, record the failure and the last 50 lines of the startup log.

  5. Write results to `{evidence_dir}/env-health.json` in this exact format:
  ```json
  {
    "startup_command": "the command used",
    "startup_time_s": 0.0,
    "health_endpoint": "{health_endpoint}",
    "health_status": 200,
    "health_body_snippet": "first 200 chars of response",
    "startup_success": true,
    "error": null
  }
  ```

  If startup fails, set startup_success to false, populate error with the failure reason, and still write the JSON. Do not abandon the evaluation — downstream agents need this file to know the system state."

---

### browser-agent
- **Role:** Browser interaction sequence — navigates pages, interacts with UI elements, captures screenshots and console state.
- **Tools:** browser_navigate, browser_snapshot, browser_take_screenshot, browser_click, browser_fill_form, browser_console_messages, browser_network_requests
- **Evidence:** `{evidence_dir}/browser-*.png`, `{evidence_dir}/browser-interaction.json`, `{evidence_dir}/console-log.json`
- **Prompt Template:**
  "You are the browser agent for a runtime/web system evaluation. You interact with the system through a real browser to verify that UI features work correctly and capture visual evidence.

  Precondition: Read `{evidence_dir}/env-health.json`. If startup_success is false, write a skip report and exit.

  Execute the following interaction sequence:

  1. Navigate to `{system_url}` using browser_navigate. Take a screenshot immediately after load and save to `{evidence_dir}/browser-01-initial.png`. Call browser_snapshot to capture the accessibility tree.

  2. For each feature in {feature_list}:
     a. Identify the UI element that triggers or represents this feature (button, link, form, navigation item).
     b. Interact with it: use browser_click for buttons/links, browser_fill_form for inputs, browser_select_option for dropdowns.
     c. Wait for the page to settle (browser_wait_for if needed).
     d. Take a screenshot: `{evidence_dir}/browser-{NN}-{feature_slug}.png`.
     e. Call browser_snapshot to capture the post-interaction DOM state.
     f. Record whether the interaction succeeded (expected element appeared, expected state change occurred, no error message displayed).

  3. After all feature interactions, collect:
     a. Console messages via browser_console_messages. Save to `{evidence_dir}/console-log.json`.
     b. Network requests via browser_network_requests. Save to `{evidence_dir}/network-log.json`.
     c. Count JS errors (messages with level 'error') and failed network requests (status >= 400).

  4. Write the interaction summary to `{evidence_dir}/browser-interaction.json`:
  ```json
  {
    "features_tested": [
      {"name": "feature_name", "interacted": true, "screenshot": "browser-02-feature.png", "passed": true, "notes": ""}
    ],
    "features_passed": 0,
    "features_total": 0,
    "js_error_count": 0,
    "failed_network_requests": 0,
    "screenshots_captured": []
  }
  ```

  Be thorough: every feature in {feature_list} must be attempted. If a feature is unreachable (element not found), record it as failed with a clear note explaining what element was expected and what the DOM actually contained."

---

### api-agent
- **Role:** HTTP/WebSocket/API endpoint testing — validates REST responses, status codes, payloads, and real-time connections.
- **Tools:** Bash (curl, websocat), WebFetch
- **Evidence:** `{evidence_dir}/api-results.json`
- **Prompt Template:**
  "You are the API agent for a runtime/web system evaluation. You test HTTP endpoints and WebSocket connections to verify backend behavior independently of the browser UI.

  Precondition: Read `{evidence_dir}/env-health.json`. If startup_success is false, write a skip report and exit.

  Execute the following tests:

  1. For each endpoint in {api_endpoints}:
     a. Send the appropriate HTTP request (GET/POST/PUT/DELETE) using Bash curl with -w to capture timing info, or WebFetch for complex requests.
     b. Record: HTTP status code, response time (TTFB), response body (first 500 chars), content-type header.
     c. Validate against expected behavior: correct status code, expected fields in response body, response time within {ttfb_threshold_ms|2000}ms.

  2. Error handling tests — for each endpoint, send one malformed request:
     a. Missing required fields (empty body for POST endpoints).
     b. Invalid content type (send text/plain to a JSON endpoint).
     c. Verify the system returns a proper error response (4xx status code with error message) and does NOT crash (check that subsequent valid requests still succeed).

  3. If {websocket_endpoints} is defined, for each WebSocket endpoint:
     a. Establish connection using Bash websocat or equivalent.
     b. Send a test message, wait up to 10s for response.
     c. Record: connection success, message roundtrip time, response content.

  4. Write results to `{evidence_dir}/api-results.json`:
  ```json
  {
    "endpoints_tested": [
      {
        "method": "GET",
        "url": "/api/resource",
        "status": 200,
        "ttfb_ms": 45,
        "body_snippet": "first 200 chars",
        "passed": true,
        "error_handling_tested": true,
        "error_handling_passed": true
      }
    ],
    "websocket_tested": [],
    "total_passed": 0,
    "total_tested": 0,
    "ttfb_p95_ms": 0,
    "error_handling_rate": 0.0
  }
  ```

  Critical: TTFB measurement must use curl's time_starttransfer, not total time. Error handling rate = endpoints where malformed request got a proper 4xx / total endpoints tested."

---

### log-agent
- **Role:** Post-interaction log analysis — parses server logs, aggregates performance metrics, detects anomalies.
- **Tools:** Bash, Read
- **Evidence:** `{evidence_dir}/log-analysis.json`
- **Prompt Template:**
  "You are the log agent for a runtime/web system evaluation. After env-agent, browser-agent, and api-agent have completed, you analyze the system's runtime logs and aggregate performance data.

  Execute the following analysis:

  1. Read the startup log from `{evidence_dir}/startup-log.txt`. Check for:
     a. Warning or error messages during startup (lines containing WARN, ERROR, FATAL, Exception, Traceback).
     b. Deprecation notices.
     c. Unhandled promise rejections or uncaught exceptions.

  2. Read server runtime logs. Check `{mission_dir}` for common log locations: logs/, *.log, stdout captured during evaluation, nohup.out, pm2 logs. Use Bash to search for log files.
     a. Count total error-level log entries during the evaluation period.
     b. Count unique error types (deduplicate by message pattern).
     c. Identify any stack traces or crash dumps.

  3. Read `{evidence_dir}/console-log.json` (from browser-agent) and `{evidence_dir}/api-results.json` (from api-agent). Aggregate:
     a. Total JS console errors from browser session.
     b. API response time distribution: min, max, p50, p95, mean.
     c. Failed request count from both browser and API testing.

  4. Write consolidated analysis to `{evidence_dir}/log-analysis.json`:
  ```json
  {
    "startup_warnings": 0,
    "startup_errors": 0,
    "runtime_error_entries": 0,
    "unique_error_types": 0,
    "stack_traces_found": 0,
    "console_js_errors": 0,
    "api_ttfb_p50_ms": 0,
    "api_ttfb_p95_ms": 0,
    "api_ttfb_max_ms": 0,
    "total_failed_requests": 0,
    "anomalies": ["description of any unusual pattern"],
    "health_summary": "CLEAN | WARNINGS | ERRORS"
  }
  ```

  Focus on facts, not interpretation. Count what is there. The judge-agent will interpret the numbers."

---

### judge-agent
- **Role:** Evidence aggregation and final verdict — reads all evidence files, scores each checkpoint, produces the evaluation judgment.
- **Tools:** Read
- **Evidence:** `{evidence_dir}/judge-verdict.json`
- **Prompt Template:**
  "You are the judge agent for a runtime/web system evaluation. You do not run any tests yourself. You read the evidence files produced by the other agents and render a structured judgment.

  Read the following evidence files:
  - `{evidence_dir}/env-health.json` (from env-agent)
  - `{evidence_dir}/browser-interaction.json` (from browser-agent)
  - `{evidence_dir}/api-results.json` (from api-agent)
  - `{evidence_dir}/log-analysis.json` (from log-agent)

  For each checkpoint metric below, extract the relevant value from the evidence and score it:

  CP1 startup_success: env-health.json -> startup_success == true AND startup_time_s < {startup_timeout_s|30}
  CP2 ttfb_ms: api-results.json -> ttfb_p95_ms < {ttfb_threshold_ms|2000}
  CP3 feature_coverage: browser-interaction.json -> features_passed / features_total >= {feature_coverage_threshold|0.8}
  CP4 error_handling_rate: api-results.json -> error_handling_rate >= {error_handling_threshold|0.8}
  CP5 console_error_count: browser-interaction.json -> js_error_count == 0
  CP6 perf_metrics: log-analysis.json -> api_ttfb_p95_ms < {perf_p95_threshold_ms|3000} AND health_summary != 'ERRORS'

  For each checkpoint, record: metric name, observed value, threshold, PASS/FAIL.

  Calculate overall: checkpoints_passed / checkpoints_total.

  Write verdict to `{evidence_dir}/judge-verdict.json`:
  ```json
  {
    "checkpoints": [
      {"id": "CP1", "name": "startup_success", "observed": true, "threshold": true, "result": "PASS"},
      {"id": "CP2", "name": "ttfb_ms", "observed": 145, "threshold": 2000, "result": "PASS"}
    ],
    "checkpoints_passed": 0,
    "checkpoints_total": 6,
    "pass_rate": 0.0,
    "evidence_files_read": [],
    "evidence_files_missing": [],
    "verdict": "PASS",
    "verdict_reason": "All 6 checkpoints passed. System starts cleanly, features are interactive, API responses are fast, no JS errors."
  }
  ```

  Verdict rules:
  - PASS: pass_rate >= 0.8 AND CP1 (startup_success) is PASS
  - CONDITIONAL PASS: pass_rate >= 0.6 AND CP1 is PASS
  - FAIL: pass_rate < 0.6 OR CP1 is FAIL
  
  If any evidence file is missing, note it in evidence_files_missing and score that checkpoint as FAIL with observed = 'MISSING_EVIDENCE'. Never infer or assume results."

---

## §2 E2E Flow & Checkpoint Metrics

### Execution Flow

```
Stage 1: env-agent starts system
    ↓ (gate: startup_success == true)
Stage 2: browser-agent + api-agent run in PARALLEL
    ↓ (both complete)
Stage 3: log-agent analyzes all evidence
    ↓
Stage 4: judge-agent renders verdict
```

### CP1: startup_success
- **Stage:** 1
- **Agent:** env-agent
- **Metrics:** startup_success (boolean, threshold: true), startup_time_s (float, threshold: < {startup_timeout_s|30})
- **Agent Dispatch Spec:**
  Agent: env-agent
  Prompt: "Start the system at {mission_dir} using {startup_command}. Poll {health_endpoint} every 2s for up to {startup_timeout_s|30}s. Write env-health.json with startup_success and startup_time_s."
  Evidence: {evidence_dir}/env-health.json
  Extract: startup_success, startup_time_s from JSON
- **Gate:** If startup_success == false, skip Stage 2-3. Judge-agent still runs with available evidence.

### CP2: ttfb_ms
- **Stage:** 2
- **Agent:** api-agent
- **Metrics:** ttfb_p95_ms (integer, threshold: < {ttfb_threshold_ms|2000})
- **Agent Dispatch Spec:**
  Agent: api-agent
  Prompt: "Test all endpoints in {api_endpoints}. Use curl -w '%{time_starttransfer}' for each request. Compute p95 across all requests. Write api-results.json."
  Evidence: {evidence_dir}/api-results.json
  Extract: ttfb_p95_ms from JSON

### CP3: feature_coverage
- **Stage:** 2
- **Agent:** browser-agent
- **Metrics:** feature_coverage (float, threshold: >= {feature_coverage_threshold|0.8})
- **Agent Dispatch Spec:**
  Agent: browser-agent
  Prompt: "Navigate to {system_url}. For each feature in {feature_list}, interact with UI and verify outcome. Capture screenshot per feature. Write browser-interaction.json with features_passed/features_total."
  Evidence: {evidence_dir}/browser-interaction.json
  Extract: features_passed / features_total

### CP4: error_handling_rate
- **Stage:** 2
- **Agent:** api-agent
- **Metrics:** error_handling_rate (float, threshold: >= {error_handling_threshold|0.8})
- **Agent Dispatch Spec:**
  Agent: api-agent
  Prompt: "For each endpoint, send one malformed request. Verify the system returns 4xx with error body (not 5xx crash, not 200 silent ignore). Compute error_handling_rate = proper_4xx / total_malformed."
  Evidence: {evidence_dir}/api-results.json
  Extract: error_handling_rate from JSON

### CP5: console_error_count
- **Stage:** 2
- **Agent:** browser-agent
- **Metrics:** console_error_count (integer, threshold: == 0)
- **Agent Dispatch Spec:**
  Agent: browser-agent
  Prompt: "After all interactions, call browser_console_messages. Count messages with level 'error'. Write to browser-interaction.json."
  Evidence: {evidence_dir}/browser-interaction.json
  Extract: js_error_count from JSON

### CP6: perf_metrics
- **Stage:** 3
- **Agent:** log-agent
- **Metrics:** api_ttfb_p95_ms (integer, threshold: < {perf_p95_threshold_ms|3000}), health_summary (string, threshold: != 'ERRORS')
- **Agent Dispatch Spec:**
  Agent: log-agent
  Prompt: "Aggregate api-results.json TTFB data. Compute p50, p95, max. Scan server logs for errors. Write log-analysis.json with health_summary."
  Evidence: {evidence_dir}/log-analysis.json
  Extract: api_ttfb_p95_ms, health_summary from JSON

---

## §3 Scenario Templates

### S1: Happy Path Full Stack Walkthrough [positive, critical]
**Agent Mission:**
  Agent: browser-agent
  Prompt: "Navigate to {system_url}. Execute the primary user journey: {primary_user_journey}. At each step, take a screenshot and verify the expected outcome. This scenario tests that the golden path works end-to-end without errors. Record every step in browser-interaction.json. After the journey, collect console messages and verify js_error_count == 0."
  Pass Criteria: All steps in {primary_user_journey} complete successfully AND js_error_count == 0 AND all screenshots captured.

### S2: API Contract Verification [positive, high]
**Agent Mission:**
  Agent: api-agent
  Prompt: "For each endpoint in {api_endpoints}, verify the response matches the documented contract: correct HTTP status code, correct content-type, response body contains required fields as specified in {api_contract}. Measure TTFB for every request. A contract violation is any response that deviates from the documented schema. Write per-endpoint results to api-results.json."
  Pass Criteria: 100% of endpoints match their documented contract AND ttfb_p95_ms < {ttfb_threshold_ms|2000}.

### S3: Error Resilience Under Bad Input [adversarial, high]
**Agent Mission:**
  Agent: api-agent
  Prompt: "You are probing the system's error handling. For every endpoint in {api_endpoints}, send three categories of bad input: (1) empty/missing required body, (2) wrong content-type header, (3) oversized payload (1MB random string). For each, verify: the system returns a 4xx error with a meaningful error message in the response body, the system does NOT crash (subsequent valid request still succeeds), and no stack trace is leaked in the response body. Record each test in api-results.json under the error_handling fields."
  Pass Criteria: error_handling_rate >= {error_handling_threshold|0.8} AND zero stack traces leaked AND system remains responsive after all bad inputs.

### S4: Performance Under Sequential Load [adversarial, medium]
**Agent Mission:**
  Agent: api-agent
  Prompt: "Send {load_request_count|20} sequential requests to {primary_endpoint} as fast as possible using a Bash loop with curl. Record TTFB for every request. Compute: min, max, p50, p95, mean. Check for response time degradation: is the p95 of the last 5 requests more than 2x the p95 of the first 5? This detects memory leaks and connection exhaustion. Write results to {evidence_dir}/load-test.json."
  Pass Criteria: p95 TTFB < {perf_p95_threshold_ms|3000} AND no degradation pattern (last-5-p95 / first-5-p95 < 2.0).

### S5: Visual Regression Baseline [positive, medium]
**Agent Mission:**
  Agent: browser-agent
  Prompt: "Navigate to each key page in {page_list}. Take a full-page screenshot of each. Call browser_snapshot to capture the accessibility tree. Verify: no broken images (check network requests for failed image loads), no layout overflow (check for horizontal scrollbar via snapshot), page title and main heading are present and match expectations. Save screenshots as {evidence_dir}/visual-{page_slug}.png."
  Pass Criteria: All pages load with zero broken images AND zero layout overflow AND all titles/headings match.

---

## §4 Evidence Protocol

| # | Agent | File | Format | Min Quality |
|---|-------|------|--------|-------------|
| 1 | env-agent | env-health.json | JSON with startup_success, startup_time_s, health_status | All fields non-null |
| 2 | env-agent | startup-log.txt | Raw text, last 200 lines of startup output | Non-empty if startup attempted |
| 3 | browser-agent | browser-interaction.json | JSON with features array, js_error_count, screenshots list | features_total > 0 |
| 4 | browser-agent | browser-{NN}-{slug}.png | PNG screenshot per feature interaction | File size > 1KB (not blank) |
| 5 | browser-agent | console-log.json | JSON array of console messages with levels | Array (may be empty) |
| 6 | browser-agent | network-log.json | JSON array of network requests with status codes | Array (may be empty) |
| 7 | api-agent | api-results.json | JSON with endpoints_tested, ttfb_p95_ms, error_handling_rate | endpoints_tested length > 0 |
| 8 | log-agent | log-analysis.json | JSON with error counts, TTFB aggregates, health_summary | health_summary is one of CLEAN/WARNINGS/ERRORS |
| 9 | judge-agent | judge-verdict.json | JSON with checkpoints array, pass_rate, verdict, verdict_reason | verdict is one of PASS/CONDITIONAL PASS/FAIL |

### Evidence Integrity Rules:
- Every JSON file must be valid JSON (parseable without error).
- Every screenshot must be a real PNG (file size > 1KB rules out blank/error captures).
- Missing evidence files are scored as FAIL for the corresponding checkpoint — never inferred.
- The judge-agent must list all evidence files it read and any that were missing.

---

## §5 Placeholder Reference

| Placeholder | Description | Example | Required |
|-------------|-------------|---------|----------|
| {mission_dir} | Project root directory | ~/drive-projects/v2/my-project@2.3.0 | Yes |
| {evidence_dir} | Evidence output directory | {mission_dir}/grading/evidence | Yes |
| {system_url} | Base URL of running system | http://localhost:3000 | Yes |
| {health_endpoint} | Health check URL | http://localhost:3000/health | Yes |
| {startup_command} | Command to start system | npm start | No (auto-detect) |
| {startup_timeout_s} | Max wait for healthy state | 30 | No (default: 30) |
| {feature_list} | JSON array of features to test | ["login", "dashboard", "search"] | Yes |
| {api_endpoints} | JSON array of endpoint specs | [{"method":"GET","url":"/api/users"}] | No (skip api-agent) |
| {websocket_endpoints} | WebSocket URLs to test | ["ws://localhost:3000/ws"] | No (skip WS tests) |
| {ttfb_threshold_ms} | Max acceptable TTFB p95 | 2000 | No (default: 2000) |
| {perf_p95_threshold_ms} | Performance threshold | 3000 | No (default: 3000) |
| {feature_coverage_threshold} | Min feature pass rate | 0.8 | No (default: 0.8) |
| {error_handling_threshold} | Min error handling rate | 0.8 | No (default: 0.8) |
| {primary_user_journey} | Step-by-step main user flow | "1. Click Login 2. Enter creds 3. See dashboard" | Yes |
| {primary_endpoint} | Main API endpoint for load test | /api/main | No |
| {page_list} | Pages for visual check | ["/", "/about", "/dashboard"] | No |
| {api_contract} | API response schema docs | "GET /users returns {users: [...]}" | No |
| {load_request_count} | Requests for load scenario | 20 | No (default: 20) |
