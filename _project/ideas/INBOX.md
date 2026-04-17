# Idea Inbox

> Dump raw ideas here immediately. Process weekly (promote to skills/ or meta/ or discard).
> Format: `- 💡 [idea title] — [one sentence context] — [date]`

---

## Inbox

- 💡 `version-management-skill` — A universal `/erwin-init` skill that initializes any project with the `_project/` cognitive framework (CLAUDE.md, CHANGELOG, GOALS, ideas/, decisions/, sessions/) — 2026-03-02 → **PROMOTED** to meta/project-management-template.md

- 💡 `visual-test-skill` — A `/visual-test` skill that enables Claude to see the running frontend page, identify visual bugs, and verify fixes without human screenshots. See design spec in `ideas/visual-test-skill.md` — 2026-03-02

- 💡 `teamwork-label-lifecycle-fix` — Audit found 3 label gaps in teamwork skills v2.3.0: (1) `status:blocked` created during `/team` init but no skill sets/unsets it; (2) `domain:*` / `size:*` labels used by `/team-issue` but not created during `/team` init; (3) `version:*` labels have commented-out creation code and no skill references them. Fix: add missing label creation to `/team` init, add block/unblock commands to `/team-drive`, decide whether to remove or implement version labels. — 2026-03-03

- 💡 `intelligent-deferral` — 为 drive-skill 增加"智慧延迟"机制：当 Claude 判断某件事不适合在当前空间/session 中执行（上下文不对、scope 溢出、需要不同工具链、或当前任务优先级更高），不应该遗忘或草率处理，而应触发一个 **Deferral Protocol**，将意图完整捕获到 `_project/ideas/INBOX.md`（或指定目标文件）。捕获内容必须包含：(1) 要做什么（具体行动描述）；(2) 为什么要做（动机和预期价值）；(3) 为什么现在不做（触发延迟的判断理由）；(4) 最佳执行时机/条件；(5) 相关上下文（文件路径、代码片段、对话要点）。哲学锚点：**塞涅卡（Seneca）**"论生命之短促"——真正的时间智慧不是永远行动，而是在对的时机做对的事；结合 drive-skill 已有的 **知行合一（王阳明）**——知行合一的前提是「知」的完整保存，当行动时机未到，「知」必须被高保真地记录下来，否则延迟就变成了遗忘，遗忘就是认知损失。这个机制本质上是知行合一的补集：**知而暂不行，行需待时机，但知不可失。** — 2026-03-03

- 💡 `team-fleet` — `/team fleet` sub-command for multi-mission teamwork orchestration: team-lead + feature-lead × N + persistent worktrees + Mode-A first-checkpoint gate. Fills the gap between `/team auto #N` (single-mission) and flat swarm dispatchers (claude-code-swarm, Anthropic Agent Teams) — nested lead hierarchy + contract discipline + checkpoint gate is the combination nobody else ships. n=1 validated on vi_agent #1298 / #1299 / #1302 (all merged, 2026-04-16). Sprout spec: `ideas/skills/team-fleet.md`; session evidence pack: `ideas/skills/team-fleet-appendix-session-2026-04-16.md`. — 2026-04-16

---

## Promoted (removed from inbox, now in skills/ or meta/)

| Idea | Promoted To | Date |
|------|------------|------|
| version-management-skill | meta/project-management-template.md | 2026-03-02 |
