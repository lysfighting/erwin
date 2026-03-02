---
# AI session workflow rules for erwin project
---

# AI Session Workflow Protocol

## 何时读 CURRENT.md — 原则优先，不罗列场景

**一个判断测试**：「我现在的行动，如果方向错了，需要撤销多少工作？」

- **撤销成本高** → 先读 `_project/CURRENT.md`（30秒），再行动
- **撤销成本低** → 直接做，无需读

**必须先读**（行动产出不可逆结果）：
- 即将写代码、创建或修改文件
- 即将做架构或技术决策
- 用户说「继续」「下一步」「做吧」「drive」
- 对话开始后第一次实质性行动

**无需读**（纯对话，不留痕迹）：
- 解释概念、回答技术问题
- 讨论方案但尚未决定动手

> 核心逻辑来自 drive-skill 的 Cardinal Rule：永不停止，除非检测到方向偏差。
> 读 CURRENT.md 是在行动前排除「方向偏差」的最小成本手段。

## 同步读取（与 CURRENT.md 一起）

当读 CURRENT.md 时，同步读：
- `_project/GOALS.md` — 验证与北极星对齐
- 相关 ADR — 如果即将改动某个有决策记录的组件

## During Session

- Discovered bugs/tasks → create `_project/issues/ISS-NNNN-title.md`
- Major decisions → create `_project/decisions/NNNN-title.md`
- Good skill ideas → add to `_project/ideas/INBOX.md` immediately

## Session End (mandatory, 10 min)

1. Update `_project/CURRENT.md`
2. Write `_project/sessions/YYYY-MM-DD.md`
3. Update `CHANGELOG.md` if features completed
4. Update `~/.claude/projects/-Users-liyasong-Project-collov-erwin/memory/MEMORY.md`
   — 新发现的规律、用户偏好、技术陷阱、项目特有知识（非项目状态，是 Claude 的认知积累）
5. Run: `git add -A && git commit -m "type(scope): description"`

## File Size Limits

| File | Target | Max |
|------|--------|-----|
| CLAUDE.md | <150 lines | 200 lines |
| _project/CURRENT.md | <80 lines | 120 lines |
| ADR files | 50-100 lines | 150 lines |
| Session logs | 50-150 lines | 300 lines |

## What Not To Do

- Don't skip CURRENT.md before irreversible actions — stale context = direction error
- Don't skip the session-end update — stale CURRENT.md = broken continuity
- Don't add large images or binary files to docs/ or _project/
- Don't edit reference/ directory files (read-only source material)
