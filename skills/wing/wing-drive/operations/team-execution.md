# Team Execution — Parallel & Wave-Based Execution

> Loaded when wing-drive selects Parallel or Full Team execution mode.

---

## Team Assembly

Count remaining unchecked **code** sub-tasks (exclude `🔧 MANUAL`).

### Medium missions (4-7 tasks): Parallel Mode

1. Analyze dependencies → group into **waves** (parallel-safe batches)
2. Create team: `TeamCreate: team_name: "mission-{issue}", agent_type: "team-lead"`
3. Spawn 1-2 teammates. **ALL teammates MUST use `model: "opus"`.** No exceptions.
   - `mode: "bypassPermissions"` — teammates don't need AskUserQuestion
   - `isolation: "worktree"` — when teammates modify same files as you

   **Leaf implementer prompt:**
   ```
   You are an implementer on team "mission-{issue}".
   Mission: #{issue} — {title}.
   YOUR tasks: {list with descriptions and target file paths}.
   Work: read code → implement → verify (run tests) → commit → report done to "team-lead".
   Source of truth: $WING_DIR/active/MISSION-{issue}.md
   Files always win over memory. Re-read after any compaction.
   ```

4. You take critical-path tasks, teammates take parallel tasks.

### Large missions (8+ tasks): Full Team Mode

1. Full wave decomposition with critical path analysis
2. Create team with STL (Sub-Team-Lead) hierarchy:
   - Domain with 3+ tasks → spawn STL who manages inner subagents
   - Domain with 1-2 tasks → spawn leaf implementer
3. Spawn all in parallel. All `model: "opus"`, `mode: "bypassPermissions"`.

   **STL prompt:**
   ```
   You are STL for {Domain} on team "mission-{issue}".
   YOUR domain tasks: {list with descriptions}.
   You OWN this domain: analyze → spawn 2-4 Task subagents → review output → report to "team-lead".
   Escalate ONLY: cross-domain conflicts, ambiguous requirements, blockers.
   Source of truth: $WING_DIR/active/MISSION-{issue}.md
   ```

### Sizing

`team_size = min(max_parallel_tasks_in_widest_wave, 5)`

| Scale | Topology | Streams |
|---|---|---|
| 4-7 tasks, 1 domain | Flat: you + 1-2 leaf | 2-3 |
| 8-12 tasks, 2 domains | 1-2 STLs + 1 leaf | 5-8 |
| 13+ tasks, 3+ domains | 2-3 STLs (you = pure orchestrator) | 10-15 |

### Teammate Isolation Principle

Brief teammates from **Contract file ONLY**. Never from conversation history. Include file paths. Teammates have zero discussion context — this is their superpower against context pollution.

---

## Wave Decomposition

Think in WAVES, not lists. Wave = tasks that run simultaneously.

```
Zero-dependency tasks → Wave 0
Depend only on Wave 0 → Wave 1
Continue until all assigned
```

### Wave Map (mandatory)

**📊 WAVE MAP** — Mission **#{issue}**
────────────────────────────────────────────
⚡ Wave 0 (Foundation) — Width: `{N}`
├── [T1] {task} → {you/teammate}
└── [T2] {task} → {you/teammate}
⚡ Wave 1 (Core) — Width: `{N}`  ← PEAK
├── [T3] {task} → {you/teammate}
🔗 Critical Path: T1 → T3 → T5
Speedup: `{X}x` via parallelism
────────────────────────────────────────────

### Emit Topology (mandatory)

**🗺️ TEAM TOPOLOGY:** mission-{issue}
YOU (Team Lead) ─── orchestrating + critical path
├── 🎖️ stl-{domain} ─── [`{N}` tasks, spawns `{M}` inner]
├── 🔧 implementer-1 ─── [`{N}` tasks]
Active: `{N}` teammates + ~`{M}` inner = `{total}` streams

---

## The Execution Loop — Swarm Orchestration

```
while (mission != COMPLETE) {

    // WAVE MANAGEMENT (highest priority)
    → All wave-N tasks dispatched? If not → dispatch NOW
    → Critical-path complete? → Pipeline wave N+1
    → All wave-N complete? → "WAVE N COMPLETE" → Launch N+1

    // YOUR WORK (critical-path tasks)
    execute → verify → if error: root cause, fix, verify
    → self-review → passes: broadcast → complete → next

    // SWARM ORCHESTRATION
    on STL domain_report → if complete: review + mark done
    on STL escalation → INSTANT response. Unblock NOW.
    on leaf complete → check work → toggle Contract → assign next

    // DYNAMIC HIERARCHY (every 2 waves)
    >=3 tasks in domain? → PROMOTE leaf to STL
    <=2 tasks left? → DEMOTE STL to leaf
    After changes → re-emit TOPOLOGY

    // CHECKPOINT (every 3 tasks / every wave)
    → RE-READ Contract FROM DISK
    → Aligned with Objective?
    → Append key decisions to AI Notes incrementally
}
```

### Team Lead Principles

1. **Swarm throughput > personal output.** Unblocking one STL > finishing your own task.
2. **Never idle while teammates work** — review, pipeline-prep, or take the next task.
3. **STL escalations = INSTANT response.**
4. **Delegate DOMAINS to STLs, not individual tasks.**
5. **Scale hierarchy dynamically.** The team is living, not planned-once.

### Broadcast Protocol

**📡 DISPATCH:** [{agent}] → {task} | Parallel with: [{others}]
**✅ RETURN:** [{agent}] ← {task} | Files: [{list}]
**🔄 RE-DISPATCH:** [{agent}] → {next-task}
**🏁 DOMAIN COMPLETE:** [stl-{domain}] ← delivered
**📊 WAVE STATUS:** Wave `{N}/{total}` | Tasks: `{done}/{total}` | Agents: `{active}`

### Checkpoint Ownership

**Only the Team Lead (you) writes to the Contract and syncs to GitHub.** Teammates report completion via `SendMessage` to "team-lead". Lead reviews → toggles checkbox → batch-syncs after wave.
