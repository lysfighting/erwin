import type { TerminalStep } from '../components/TerminalPlayground'

/**
 * Mock terminal data — 1:1 reproduction of REAL Teamwork skill output.
 *
 * Visual encoding (Claude Code rendering):
 *   **bold** → white bold (section headers, labels)
 *   `backticks` → purple-blue rgb(177,185,249) (commands, paths, counts)
 *   *italic* → dim italic (branch names, timestamps)
 *   **#NNN** → light-blue bold clickable (issue/PR numbers)
 *   ──── → box-drawing dividers (44 chars)
 *
 * Virtual project: acme-corp/taskflow
 * Team: alice (tech-lead), bob (backend), charlie (frontend)
 * Base branch: pre-launch | Production: main
 * Active mission: #042 "Add user authentication with JWT"
 */

// ─────────────────────────────────────────────
// CH3: Setup — /team init, /team, /team config
// ─────────────────────────────────────────────

export const ch3SetupMocks: TerminalStep[][] = [
  // Case 1: /team (first time — triggers init)
  [
    {
      command: '/team',
      output: [
        '⚡ No team config found. Running first-time setup...',
        '',
        '🔍 Scanning acme-corp/taskflow...',
        '  Language: TypeScript · Package manager: npm',
        '  Remote:   github.com/acme-corp/taskflow',
        '  Branches: *main*, *pre-launch*',
        '',
        '👥 3 contributors detected from git log:',
        '  alice    47 commits  →  tech-lead',
        '  bob      31 commits  →  backend',
        '  charlie  28 commits  →  frontend',
        '',
        '📝 Writing config...',
        '  ✓ `.teamwork/config.yml`',
        '  ✓ GitHub labels synced (`mission`, `status:queued`, `status:wip`, `status:review`, `status:done`, `priority:P0`–`P3`, `size:S`–`XL`)',
        '  ✓ Milestone `V1.0` created',
        '  ✓ `.github/ISSUE_TEMPLATE/mission.yml`',
        '  ✓ `.github/workflows/ci.yml`',
        '  ✓ `.githooks/pre-commit` + `.githooks/pre-push`',
        '',
        '────────────────────────────────────────────',
        '✅ Team initialized! Showing dashboard...',
        '',
      ],
      delay: 60,
    },
    {
      command: '',
      output: [
        '**MY DASHBOARD** — acme-corp/taskflow ──────────────────',
        'Teamwork v3.8.3 · **V1.0**  `░░░░░░░░░░`  0%  `0`/`0`',
        '',
        '👥 **TEAM** (3 members)',
        '  alice (tech-lead)',
        '    ⚪ *(idle)*',
        '  bob (backend)',
        '    ⚪ *(idle)*',
        '  charlie (frontend)',
        '    ⚪ *(idle)*',
        '',
        '📋 **Unassigned** (0)',
        '  No issues yet. Run `/team-issue` to create one.',
        '',
        '💡 **Next:** `/team-issue` to create your first mission',
        '',
        '`/team help` · `/team doctor` · `/team config`',
        '────────────────────────────────────────────',
      ],
      delay: 40,
    },
  ],

  // Case 2: /team config
  [
    {
      command: '/team config',
      output: [
        '⚙️ **TEAMWORK CONFIG** ── v3.8.3 ────────────────',
        '**Config:**   `.teamwork/config.yml`',
        '**Worktree mode:**  `false` *(local, per-user)*',
        '',
        '────────── **Shared Config** (committed) ──────────',
        'project: acme-corp/taskflow',
        'version: 3',
        '',
        'conventions:',
        '  base_branch: "pre-launch"',
        '  production_branch: "main"',
        '  branch_pattern: "mission/{issue}-{slug}-{user}"',
        '  commit_style: "conventional"',
        '',
        'milestone: "V1.0"',
        '',
        'team:',
        '  - github: alice',
        '    role: tech-lead',
        '  - github: bob',
        '    role: backend',
        '  - github: charlie',
        '    role: frontend',
        '',
        'labels:  [mission, status:queued, status:wip, status:review,',
        '          status:done, priority:P0, priority:P1, priority:P2,',
        '          priority:P3, size:S, size:M, size:L, size:XL]',
        '',
        '────────── **Local Settings** (per-user) ──────────',
        'Worktree: `false`',
        '',
        '────────────────────────────────────────────',
        '**HOW TO EDIT:**',
        '  Shared config: edit `.teamwork/config.yml` directly',
        '                 (changes take effect immediately, no re-init needed)',
        '  Worktree mode: `git config --local teamwork.worktree true/false`',
        '  Full re-init:  `/team init` (upgrades version, adds missing sections)',
        '────────────────────────────────────────────',
      ],
    },
  ],

  // Case 3: /team (dashboard with active missions)
  [
    {
      command: '/team',
      output: [
        '**MY DASHBOARD** — acme-corp/taskflow ──────────────────',
        'Teamwork v3.8.3 · **V1.0**  `████░░░░░░`  25%  `1`/`4`',
        '',
        '📂 **Worktrees** *(worktree mode)* 🟥 main · 🟦 develop · 🟧 rc · 🟩 mission · ⬜ other',
        '  ⭐ 🟩 `~/Projects/taskflow.worktrees/042-add-auth-bob` ← *mission/042-add-auth-bob* → **#042** Add user auth',
        '  · 🟩 `~/Projects/taskflow.worktrees/048-fix-payment-alice` ← *mission/048-fix-payment-alice* → **#048** Fix payment crash',
        '',
        '👥 **TEAM** (3 members)',
        '  alice (tech-lead)',
        '    🟢 **#048** Fix payment crash on expired cards 🟠 P0 → *wip*',
        '    *mission/048-fix-payment-alice* → `/team-ship`',
        '    🔀 **#045** — ✅',
        '  bob (backend)',
        '    🟢 **#042** Add user authentication with JWT 🟠 P1 → *wip*',
        '    *mission/042-add-auth-bob* → `/team-drive`',
        '  charlie (frontend)',
        '    ⚪ *(idle)*',
        '',
        '📋 **Unassigned** (`2`)',
        '  🟡 P2 **#043**  Add rate limiting',
        '  🟡 P3 **#044**  Write API docs',
        '',
        '📊 **Merges**  alice: `12` · bob: `8` · charlie: `5`',
        '',
        '💡 **Next:** `/team-claim` **#043** to start',
        '',
        '`/team help` · `/team doctor` · `/team #N` · `/team config`',
        '────────────────────────────────────────────',
      ],
    },
  ],
]

// ─────────────────────────────────────────────
// CH4: Create — /team-issue
// ─────────────────────────────────────────────

export const ch4CreateMocks: TerminalStep[][] = [
  // Case 1: /team-issue (single issue creation)
  [
    {
      command: '/team-issue "Add user authentication with JWT"',
      output: [
        '📋 **PREVIEW** ── Add user authentication with JWT ──────────────────',
        '  **Assignee**   @bob',
        '  🟠 `P1`  **Size** `M`  **Domain** `backend`',
        '  **Labels**     `mission`, `status:queued`, `priority:P1`, `domain:backend`, `size:M`',
        '  **Milestone**  `V1.0`',
        '────────────────────────────────────────────',
        '',
        '## Mission Contract',
        '',
        '**Objective:** Implement JWT-based authentication for the TaskFlow API.',
        '',
        '**Sub-tasks:**',
        '- [ ] Install jsonwebtoken + bcrypt dependencies',
        '- [ ] Create auth service (generate, verify, refresh tokens)',
        '- [ ] Implement POST /auth/login and POST /auth/register',
        '- [ ] Add JWT validation middleware to protected routes',
        '- [ ] Write unit tests (≥80% coverage for auth module)',
        '',
        '**Acceptance Criteria:**',
        '- POST /auth/login returns valid JWT with 15m expiry',
        '- POST /auth/register hashes password with bcrypt',
        '- Protected routes return 401 without valid token',
        '- Token refresh endpoint with rotation',
        '',
        '**Context Files:**',
        '  `src/routes/index.ts` · `src/middleware/` · `package.json`',
        '',
        '**Test Command:** `npm test -- --grep "auth"`',
        '────────────────────────────────────────────',
      ],
      delay: 50,
    },
    {
      command: '',
      output: [
        '📦 **CREATED** ── **#042** Add user authentication with JWT ─────────────',
        '  **Assignee**  @bob',
        '  🟠 `P1`  **Size** `M`  **Milestone** `V1.0`',
        '  **URL**       `https://github.com/acme-corp/taskflow/issues/42`',
        '',
        '────────────────────────────────────────────',
        '`/team-claim` **#042** to generate Contract + branch',
        '',
        '💡 Tip: Use `/team-issue batch` to create multiple issues at once.',
      ],
      delay: 40,
    },
  ],

  // Case 2: /team-issue batch
  [
    {
      command: '/team-issue batch V1.0 -- "Auth with JWT" "Rate limiting" "API docs" "Fix payment crash"',
      output: [
        '🏁 **MILESTONE + MC PREVIEW**',
        '────────────────────────────────────────────',
        '  **Milestone**    `V1.0`',
        '  **Description**  First release milestone',
        '  **MCs to create** `4`',
        '',
        '  🟠 MC1: Add user authentication with JWT → @bob',
        '  🟡 MC2: Add rate limiting to API endpoints → @—',
        '  🟡 MC3: Write API documentation → @—',
        '  🟠 MC4: Fix payment crash on expired cards → @alice',
        '────────────────────────────────────────────',
        '',
      ],
      delay: 40,
    },
    {
      command: '',
      output: [
        '🏁 **MILESTONE** ── V1.0 ────────────────────',
        '  **URL**  `https://github.com/acme-corp/taskflow/milestone/1`',
        '  **MCs**  `4` created',
        '',
        '  **#042**  Add user authentication with JWT  @bob',
        '  **#043**  Add rate limiting to API endpoints  @—',
        '  **#044**  Write API documentation  @—',
        '  **#048**  Fix payment crash on expired cards  @alice',
        '',
        '────────────────────────────────────────────',
        '`/team-claim` **#042** to start working',
        '',
        '💡 Tip: `/team` to see your updated dashboard.',
      ],
      delay: 60,
    },
  ],

  // Case 3: View single issue details (/team #42)
  [
    {
      command: '/team #42',
      output: [
        '**MC #042** ── Add user authentication with JWT ────────────────',
        '**Assignee:**   @bob',
        '**Status:**     *wip*',
        '**Priority:**   🟠 P1',
        '**Milestone:**  `V1.0`',
        '**URL:**        `https://github.com/acme-corp/taskflow/issues/42`',
        '────────────────────────────────────────────',
        '',
        '## Mission Contract',
        '**Objective:** Implement JWT-based authentication for the TaskFlow API.',
        '',
        '**Sub-tasks:**',
        '  ✅ Install jsonwebtoken + bcrypt dependencies',
        '  ✅ Create auth service',
        '  ○  Implement login/register routes',
        '  ○  Add JWT middleware',
        '  ○  Write unit tests',
        '',
        '────────────────────────────────────────────',
        '',
        '▸ **BRANCH**  *mission/042-add-auth-bob*',
        '  Commits: `3` on *mission/042-add-auth-bob*',
        '',
        '▸ **RELATED PRs**',
        '  (none)',
        '',
        '────────────────────────────────────────────',
        '`/team-claim` **#042** to generate Contract + branch',
      ],
    },
  ],
]

// ─────────────────────────────────────────────
// CH5: Claim — /team-claim
// ─────────────────────────────────────────────

export const ch5ClaimMocks: TerminalStep[][] = [
  // Case 1: /team-claim #42
  [
    {
      command: '/team-claim #42',
      output: [
        '**CLAIMED** ──── **#042** Add user authentication with JWT ────────────',
        '🟠  `P1`  Branch: *mission/042-add-auth-bob*',
        'Contract:  `.teamwork/active/MISSION-42.md`',
        '',
        '**OBJECTIVE**',
        '  Implement JWT-based authentication for the TaskFlow API',
        '',
        '**SUB-TASKS**',
        '  [ ] Install jsonwebtoken + bcrypt dependencies',
        '  [ ] Create auth service (generate, verify, refresh tokens)',
        '  [ ] Implement POST /auth/login and POST /auth/register',
        '  [ ] Add JWT validation middleware to protected routes',
        '  [ ] Write unit tests (≥80% coverage for auth module)',
        '',
        '**SUCCESS CRITERIA**',
        '  POST /auth/login returns valid JWT with 15m expiry',
        '  POST /auth/register hashes password with bcrypt',
        '  Protected routes return 401 without valid token',
        '  Token refresh endpoint with rotation',
        '',
        '**CONTEXT FILES**',
        '  `src/routes/index.ts`',
        '  `src/middleware/`',
        '  `package.json`',
        '',
        '────────────────────────────────────────────',
        '`/team-drive` to execute │ `/team-ship` when done',
        '',
        '💡 Tip: The contract is your source of truth during execution.',
      ],
      delay: 60,
    },
  ],

  // Case 2: /team-claim (worktree mode)
  [
    {
      command: '/team-claim #42',
      output: [
        '**CLAIMED** ──── **#042** Add user authentication with JWT ────────────',
        '🟠  `P1`  Branch: *mission/042-add-auth-bob*',
        'Contract:  `.teamwork/active/MISSION-42.md`',
        'Worktree:  `~/Projects/taskflow.worktrees/042-add-auth-bob`',
        '',
        '**OBJECTIVE**',
        '  Implement JWT-based authentication for the TaskFlow API',
        '',
        '**SUB-TASKS**',
        '  [ ] Install jsonwebtoken + bcrypt dependencies',
        '  [ ] Create auth service (generate, verify, refresh tokens)',
        '  [ ] Implement POST /auth/login and POST /auth/register',
        '  [ ] Add JWT validation middleware to protected routes',
        '  [ ] Write unit tests (≥80% coverage for auth module)',
        '',
        '**SUCCESS CRITERIA**',
        '  POST /auth/login returns valid JWT with 15m expiry',
        '  POST /auth/register hashes password with bcrypt',
        '  Protected routes return 401 without valid token',
        '',
        '**CONTEXT FILES**',
        '  `src/routes/index.ts`',
        '  `src/middleware/`',
        '  `package.json`',
        '',
        '⚠️ **DO NOT run `/team-drive` in this terminal**',
        '────────────────────────────────────────────',
        '当前终端在主仓库。Worktree 模式下，所有操作必须在 worktree 目录中执行。',
        '在这里运行 `/team-drive` 会操作错误的代码。',
        '',
        '**打开新终端 tab，执行：**',
        '  `cd ~/Projects/taskflow.worktrees/042-add-auth-bob && claude`',
        '  Then: `/team-drive` to execute │ `/team-ship` when done',
        '────────────────────────────────────────────',
        '',
        '💡 Tip: Each worktree gets its own iTerm tab with branch color.',
      ],
      delay: 60,
    },
  ],

  // Case 3: /team-claim list
  [
    {
      command: '/team-claim list',
      output: [
        '**Your assigned missions:**',
        '  🟠 `1.` **#042** Add user authentication with JWT',
        '  🟠 `2.` **#048** Fix payment crash on expired cards',
        '',
        '**Unassigned missions:**',
        '  🟡 `3.` **#043** Add rate limiting',
        '  🟡 `4.` **#044** Write API docs',
        '',
        '────────────────────────────────────────────',
        '`/team-claim #042` to claim',
      ],
    },
  ],
]

// ─────────────────────────────────────────────
// CH6: Execute — /team-drive
// ─────────────────────────────────────────────

export const ch6ExecuteMocks: TerminalStep[][] = [
  // Case 1: /team-drive (full mission execution)
  [
    {
      command: '/team-drive',
      output: [
        '🎯 **MISSION** ── **#042** Add user authentication with JWT ─────────────',
        '🟠 P1  🔀 *mission/042-add-auth-bob*',
        'Progress:  `░░░░░░░░░░`  `0`/`5`',
        '',
        '🎯 **OBJECTIVE**',
        '  Implement JWT-based authentication for the TaskFlow API',
        '',
        '📋 **SUB-TASKS**',
        '  ○ Install jsonwebtoken + bcrypt dependencies',
        '  ○ Create auth service (generate, verify, refresh tokens)',
        '  ○ Implement POST /auth/login and POST /auth/register',
        '  ○ Add JWT validation middleware to protected routes',
        '  ○ Write unit tests (≥80% coverage for auth module)',
        '',
        '✅ **ACCEPTANCE CRITERIA**',
        '  POST /auth/login returns valid JWT with 15m expiry',
        '  Protected routes return 401 without valid token',
        '',
        '📁 **CONTEXT FILES**',
        '  `src/routes/index.ts` · `src/middleware/` · `package.json`',
        '────────────────────────────────────────────',
      ],
      delay: 50,
    },
    {
      command: '▸ Task 1/5: Install dependencies',
      output: [
        '  `npm install jsonwebtoken bcrypt @types/jsonwebtoken @types/bcrypt`',
        '  added 4 packages',
        '  ✅ Task complete',
      ],
      delay: 60,
    },
    {
      command: '▸ Task 2/5: Create auth service',
      output: [
        '  Creating `src/services/auth.ts`...',
        '  → generateToken(userId): signs JWT with 15m expiry',
        '  → verifyToken(token): validates and decodes',
        '  → refreshToken(token): rotation with new jti',
        '  → hashPassword / comparePassword: bcrypt cost=12',
        '  ✅ Task complete',
      ],
      delay: 70,
    },
    {
      command: '▸ Task 3/5: Login & register routes',
      output: [
        '  Creating `src/routes/auth.ts`...',
        '  → POST /auth/register: validate input, hash password, store, return token',
        '  → POST /auth/login: find user, compare password, return token',
        '  Running tests...',
        '  ✗ FAIL: "should hash password before storing"',
        '    Expected bcrypt.hash to be called — got plaintext insert',
        '',
        '  🔧 Fixing: added await bcrypt.hash(password, 12) in register handler',
        '  Running tests... 6 passed ✅',
        '  ✅ Task complete',
      ],
      delay: 80,
    },
    {
      command: '▸ Task 4/5: JWT middleware',
      output: [
        '  Creating `src/middleware/auth.ts`...',
        '  → Extract token from Authorization: Bearer header',
        '  → Verify token, attach user to req.user',
        '  → 401 if missing, 403 if invalid',
        '  Applied to: /api/tasks, /api/users, /api/payments',
        '  ✅ Task complete',
      ],
      delay: 60,
    },
    {
      command: '▸ Task 5/5: Write unit tests',
      output: [
        '  Creating `tests/auth.test.ts`...',
        '  → 12 test cases covering login, register, middleware, refresh',
        '  Running full suite... 12 passed, 0 failed ✅',
        '  Coverage: 87% (target ≥80%) ✅',
        '',
        '────────────────────────────────────────────',
        '**🏁 COMPLETE** ── **#042** Add user authentication with JWT ────────────',
        'Mode:     sequential  `5` tasks  `5` waves',
        'Commits:  `5` on 🔀 *mission/042-add-auth-bob*',
        'Tests:    ✅ passing  Coverage: 87%',
        '',
        '**📊 PLAN VALIDATION**',
        '  CONFIRMED: all 5 sub-tasks mapped to code changes',
        '',
        '**📋 SUB-TASKS**',
        '  ✅ Install jsonwebtoken + bcrypt dependencies',
        '  ✅ Create auth service',
        '  ✅ Implement login/register routes',
        '  ✅ Add JWT middleware',
        '  ✅ Write unit tests (87% coverage)',
        '',
        '**✅ ACCEPTANCE CRITERIA**',
        '  POST /auth/login returns valid JWT — verified by test',
        '  Protected routes return 401 — verified by test',
        '  Token refresh with rotation — verified by test',
        '',
        '────────────────────────────────────────────',
        '`/team-ship` to create PR',
        '',
        '💡 Tip: Review your changes before shipping.',
      ],
      delay: 50,
    },
  ],

  // Case 2: /team-drive resume after interruption
  [
    {
      command: '/team-drive',
      output: [
        '🎯 **MISSION** ── **#042** Add user authentication with JWT ─────────────',
        '🟠 P1  🔀 *mission/042-add-auth-bob*',
        'Progress:  `██████░░░░`  `3`/`5`',
        '',
        '📋 **SUB-TASKS**',
        '  ✅ Install jsonwebtoken + bcrypt dependencies',
        '  ✅ Create auth service',
        '  ✅ Implement POST /auth/login and POST /auth/register',
        '  ▸ Add JWT validation middleware  ← resuming here',
        '  ○ Write unit tests',
        '',
        '💾 State recovered from `.teamwork/active/MISSION-42.md`',
        '────────────────────────────────────────────',
      ],
      delay: 50,
    },
    {
      command: '▸ Task 4/5: JWT middleware (resumed)',
      output: [
        '  Creating `src/middleware/auth.ts`...',
        '  → Token extraction, verification, user attachment',
        '  Applied to protected routes',
        '  ✅ Task complete',
      ],
      delay: 60,
    },
    {
      command: '▸ Task 5/5: Write unit tests',
      output: [
        '  Creating `tests/auth.test.ts`...',
        '  → 12 test cases',
        '  Running full suite... 12 passed ✅',
        '  Coverage: 87% ✅',
        '',
        '────────────────────────────────────────────',
        '**🏁 COMPLETE** ── **#042** Add user authentication with JWT ────────────',
        'Mode:     sequential  `5` tasks  ·  Resumed from task 4',
        'Commits:  `5` on 🔀 *mission/042-add-auth-bob*',
        'Tests:    ✅ passing  Coverage: 87%',
        '',
        '────────────────────────────────────────────',
        '`/team-ship` to create PR',
      ],
      delay: 50,
    },
  ],
]

// ─────────────────────────────────────────────
// CH7: Deliver — /team-ship
// ─────────────────────────────────────────────

export const ch7DeliverMocks: TerminalStep[][] = [
  // Case 1: /team-ship (create PR)
  [
    {
      command: '/team-ship',
      output: [
        '**📦 SHIPPED** ── **#042** Add user authentication with JWT ─────────────',
        '🔀 *mission/042-add-auth-bob*',
        '🔀 PR: `https://github.com/acme-corp/taskflow/pull/87`',
        '✅ CI: passing',
        'Review:    requested from @alice',
        'Labels:    `status:review`',
        '',
        '**📋 SUB-TASKS DELIVERED**',
        '  ✅ Install jsonwebtoken + bcrypt dependencies',
        '  ✅ Create auth service',
        '  ✅ Implement login/register routes',
        '  ✅ Add JWT middleware',
        '  ✅ Write unit tests (87% coverage)',
        '',
        '────────────────────────────────────────────',
        '`/team-ship done` (after merge) │ `/team-ship review` │ `/team`',
        '',
        '💡 Tip: Run `/team-ship review` for AI code review before human review.',
      ],
      delay: 70,
    },
  ],

  // Case 2: /team-ship review
  [
    {
      command: '/team-ship review',
      output: [
        '**AI REVIEW** ── PR **#087** ────────────────────────',
        '',
        '`src/services/auth.ts`',
        '  ✓ Token expiry set correctly (15m access, 7d refresh)',
        '  ⚠ Consider adding jti (JWT ID) for token revocation',
        '  ✓ Refresh rotation prevents replay attacks',
        '',
        '`src/routes/auth.ts`',
        '  ✓ Input validation on login/register endpoints',
        '  ✓ Password hashed with bcrypt (cost=12)',
        '  ⚠ Add rate limiting to /auth/login (brute force protection)',
        '',
        '`src/middleware/auth.ts`',
        '  ✓ Proper 401/403 error responses',
        '  ✓ Token extracted from Authorization header',
        '',
        '────────────────────────────────────────────',
        'Summary: `2` suggestions, `0` blockers. ✅ Approve recommended.',
        '',
        '💡 Tip: Address suggestions before or after merge — your call.',
      ],
      delay: 50,
    },
  ],

  // Case 3: /team-ship done (post-merge cleanup)
  [
    {
      command: '/team-ship done',
      output: [
        '**🏁 DONE** ── **#042** Add user authentication with JWT ────────────────',
        'Issue:     closed',
        '🔀 PR: `https://github.com/acme-corp/taskflow/pull/87` (merged)',
        'Labels:    `status:done`',
        'Worktree:  removed',
        'Branch:    *pre-launch*',
        '',
        '────────────────────────────────────────────',
        '`/team-claim` (next mission) │ `/team` (dashboard)',
        '',
        '💡 Tip: Run `/team` to see who needs help.',
      ],
      delay: 50,
    },
  ],
]

// ─────────────────────────────────────────────
// CH8: Shortcuts — /team auto, /team wrap
// ─────────────────────────────────────────────

export const ch8ShortcutsMocks: TerminalStep[][] = [
  // Case 1: /team auto (one-shot: claim → drive → ship)
  [
    {
      command: '/team auto #43',
      output: [
        '🤖 **AUTO** ── **#043** Add rate limiting ────────────────',
        '🟡 P2   Size: `M`',
        'Milestone: `V1.0`',
        '',
        '🎯 **OBJECTIVE**',
        '  Add API rate limiting to prevent abuse',
        '',
        '⚡ **AUTO PLAN**',
        '  Phase 1: Claim → Contract + branch',
        '  Phase 2: Drive → execute all sub-tasks',
        '  Phase 3: Ship → push + PR',
        '  Stops at: PR created (you merge manually)',
        '────────────────────────────────────────────',
      ],
      delay: 50,
    },
    {
      command: '',
      output: [
        '✅ **PHASE 1: CLAIMED**',
        '   Contract: `.teamwork/active/MISSION-43.md`',
        '   Branch:   *mission/043-rate-limiting-charlie*',
        '   Proceeding to drive...',
        '────────────────────────────────────────────',
      ],
      delay: 60,
    },
    {
      command: '',
      output: [
        '✅ **PHASE 2: DRIVEN**',
        '   Code tasks: `2`/`2` complete',
        '   Commits:   `2` commits on *mission/043-rate-limiting-charlie*',
        '   Proceeding to ship...',
        '────────────────────────────────────────────',
      ],
      delay: 60,
    },
    {
      command: '',
      output: [
        '🤖 **AUTO COMPLETE** ── **#043** Add rate limiting ────────────',
        '────────────────────────────────────────────',
        '',
        '✅ **PHASE 1: CLAIMED**',
        '   Branch: *mission/043-rate-limiting-charlie*',
        '',
        '✅ **PHASE 2: DRIVEN**',
        '   Code tasks: `2`/`2`',
        '',
        '✅ **PHASE 3: SHIPPED**',
        '   🔀 PR: `https://github.com/acme-corp/taskflow/pull/88`',
        '   ✅ CI: passing',
        '   Review: requested from @alice',
        '',
        '────────────────────────────────────────────',
        '⏳ Waiting for human merge.',
        'After merge: `/team-ship done #043`',
        '────────────────────────────────────────────',
      ],
      delay: 70,
    },
  ],

  // Case 2: /team wrap (retroactive)
  [
    {
      command: '/team wrap',
      output: [
        '🔍 Detecting uncommitted changes...',
        '  4 files changed in working tree:',
        '    `src/utils/logger.ts`         (new)',
        '    `src/config/logging.ts`       (modified)',
        '    `tests/logger.test.ts`        (new)',
        '    `package.json`                (added: winston)',
        '',
        '🤖 Analyzing diff...',
        '  Inferred: "Add structured logging with Winston"',
        '',
        '**WRAP** ── retroactive mission ────────────────',
        '  ✓ Issue **#049** created: "Add structured logging with Winston"',
        '  ✓ Branch: *mission/049-structured-logging-bob*',
        '  ✓ Contract generated from existing changes',
        '  ✓ Changes committed: 4 files (+187 -12)',
        '  ✓ PR **#089** created',
        '  Labels: `status:review`',
        '',
        '────────────────────────────────────────────',
        '`/team-ship done` (after merge) │ `/team`',
        '',
        '💡 Tip: `/team wrap` captures work you did without a mission.',
      ],
      delay: 60,
    },
  ],
]

// ─────────────────────────────────────────────
// CH9: Health — /team doctor
// ─────────────────────────────────────────────

export const ch9HealthMocks: TerminalStep[][] = [
  // Case 1: /team doctor (with findings — table format)
  [
    {
      command: '/team doctor',
      output: [
        '🏥 **LOCAL HEALTH** — acme-corp/taskflow',
        '   Current: *pre-launch* · Base: *pre-launch*',
        '────────────────────────────────────────────',
        '',
        '📋 **BRANCHES** (`5` local)',
        '  *main*              🟢 clean    base branch',
        '  *pre-launch*        🟢 clean    integration branch',
        '  *mission/042-*      🔴 dirty    +6, -1 vs base  **#042** OPEN │ PR **#087** ✅',
        '  *mission/039-*      💀 orphan   +0, -38 vs base **#039** CLOSED │ PR **#040** merged',
        '  *mission/041-*      💀 orphan   merged, 32 days old',
        '',
        '📋 **OPEN PRs** (`1`)',
        '  **#087**  feat: user auth   ✅ CI   ✅ Mergeable  🟢 SHIP',
        '',
        '📋 **STASHES** (`2`)',
        '  stash@{0}   3 days   "WIP on mission/042"   14 files   🟡 LARGE',
        '  stash@{1}   42 days  "WIP on mission/018"   1 file     🔴 DROP',
        '',
        '📋 **WORKTREES** (`2` extra)',
        '  `~/taskflow.worktrees/042-auth`    *mission/042-*    🟢 active',
        '  `~/taskflow.worktrees/039-old`     (missing)         🟡 prunable',
        '',
        '📋 **ISSUES** (`2` problems in `6` open)',
        '  **#050**  fix: Safari camera   @alice   *wip*   🟡 DUP?',
        '  **#052**  fix: iOS media       @alice   *wip*   🟡 DUP?',
        '',
        '🔍 **DIAGNOSIS**',
        '────────────────────────────────────────────',
        '  🔴 Branch *mission/042-* has uncommitted changes — stash or commit',
        '  🔴 Stash stash@{1} is 42 days old on deleted branch — drop it',
        '  🟡 **#087** is ready to ship — merge it',
        '  🟡 2 orphan branches can be deleted',
        '  🟡 Issues **#050** and **#052** may be duplicates',
        '  🟢 Worktree `039-old` is prunable',
        '',
        '📌 **ACTION PLAN** (`6` actions)',
        '────────────────────────────────────────────',
        '  🔴 [1] Commit or stash changes on *mission/042-*',
        '  🔴 [2] Drop stale stash: `git stash drop stash@{1}`',
        '  🟡 [3] Merge PR **#087**: approve and merge',
        '  🟡 [4] Delete orphan branches: `git branch -d mission/039-* mission/041-*`',
        '  🟡 [5] Review duplicate issues **#050** **#052**',
        '  🟢 [6] Prune worktree: `git worktree remove 039-old`',
        '',
        '────────────────────────────────────────────',
        '**SUMMARY:** 🔴 `2` critical · 🟡 `3` warnings · 🟢 `1` cleanup',
        '`/team doctor fix` to execute interactively',
      ],
      delay: 40,
    },
  ],

  // Case 2: /team doctor fix (auto-cleanup)
  [
    {
      command: '/team doctor fix',
      output: [
        '🟢 **AUTO-CLEANUP** (`4` actions):',
        '',
        '  [4] Delete orphan *mission/039-old-feature*: `git branch -d mission/039-*`',
        '      ✅ Deleted',
        '  [4] Delete orphan *mission/041-typo-fix*: `git branch -d mission/041-*`',
        '      ✅ Deleted',
        '  [2] Drop stale stash: `git stash drop stash@{1}`',
        '      ✅ Dropped',
        '  [6] Prune worktree: `git worktree remove 039-old`',
        '      ✅ Pruned',
        '',
        '🏥 **DOCTOR FIX COMPLETE**',
        '────────────────────────────────────────────',
        '   ✅ Executed: `4` actions',
        '   ⏭️ Skipped:  `2` actions (need manual review)',
        '────────────────────────────────────────────',
      ],
      delay: 60,
    },
  ],
]

// ─────────────────────────────────────────────
// CH10: Release — /team-rc
// ─────────────────────────────────────────────

export const ch10ReleaseMocks: TerminalStep[][] = [
  // Case 1: /team-rc (cut RC from pre-launch)
  [
    {
      command: '/team-rc',
      output: [
        '**PREPARE RC** ── **V1.0.0** ────────────────────────',
        '*rc/V1.0.0* ← *pre-launch*',
        'Develop CI: green ✅',
        '',
        '🔍 Pre-release checks:',
        '  ✓ lint — 0 errors',
        '  ✓ typecheck — 0 errors',
        '  ✓ unit tests — 47 passed',
        '  ✓ integration — 12 passed',
        '  ✓ coverage — 91% (target ≥80%)',
        '',
      ],
      delay: 40,
    },
    {
      command: '',
      output: [
        '**RC PREPARED ── V1.0.0** ────────────────────',
        '*rc/V1.0.0*  ←  *pre-launch*',
        'Staging:  triggered',
        '',
        '**NEXT STEPS**',
        '  1. 部署 staging，验证',
        '  2. 发现 bug? → 在 *pre-launch* 修复 → `/team-rc cancel` → `/team-rc` 重新切',
        '  3. 验证通过 → `/team-rc promote`',
        '',
        '────────────────────────────────────────────',
        '',
        '💡 Tip: Never commit directly to the RC branch.',
      ],
      delay: 60,
    },
  ],

  // Case 2: /team-rc promote
  [
    {
      command: '/team-rc promote',
      output: [
        '**PROMOTED ── V1.0.0** ────────────────────────',
        'Tag: **V1.0.0** (on *main*)',
        'PR:  **#156** (squash merged)',
        'Release: `https://github.com/acme-corp/taskflow/releases/tag/V1.0.0`',
        'RC Branch:  *rc/V1.0.0* (deleted)',
        'Milestone:  **V1.0** — `4` closed, `0` open',
        '',
        '────────────────────────────────────────────',
        '`/team-rc` to prepare next RC',
        '',
        '💡 Tip: *pre-launch* has been rebased on *main*.',
      ],
      delay: 70,
    },
  ],
]

// ─────────────────────────────────────────────
// CH11: Help & Learn — /team help, /team learn
// ─────────────────────────────────────────────

export const ch11HelpMocks: TerminalStep[][] = [
  // Case 1: /team help
  [
    {
      command: '/team help',
      output: [
        '**TEAMWORK v3.8.3** — AI-Native Team Coordination (Push Model)',
        '────────────────────────────────────────────',
        '',
        '**SETUP** (one time)',
        '  `/team`                 Auto-detect project, create config + CI + labels',
        '  `/team init`            Force re-initialize (upgrade existing config)',
        '',
        '**DAILY WORKFLOW**',
        '  `/team`                 See role-based dashboard (leader=team view, member=my view)',
        '  `/team #42`             View MC **#042** details (branch, commits, PRs)',
        '',
        '**ISSUE MANAGEMENT**',
        '  `/team-issue <desc>`    Create AI-enriched Issue with labels + contract structure',
        '  `/team-issue #42`       View Issue details',
        '  `/team-issue batch`     Create milestone + multiple issues at once',
        '',
        '**MISSION EXECUTION**',
        '  `/team-claim #42`       Generate Contract + create branch (+ worktree)',
        '  `/team-drive`           Autonomous execution loop against Contract',
        '  `/team-ship`            Push + create PR + CI check + cleanup',
        '  `/team-ship done`       Post-merge: close Issue, clean worktree, update labels',
        '',
        '**SHORTCUTS**',
        '  `/team auto #42`        One-shot: claim → drive → ship',
        '  `/team wrap`            Retroactive: detect changes → issue + branch + PR',
        '',
        '**RELEASE**',
        '  `/team-rc`              Cut RC branch from *pre-launch*',
        '  `/team-rc promote`      Squash-merge RC → *main*, tag, GitHub Release',
        '  `/team-rc cancel`       Delete RC, fix on *pre-launch*, re-cut',
        '',
        '**UTILITIES**',
        '  `/team doctor`          12-check project health diagnostic',
        '  `/team doctor fix`      Auto-repair found issues',
        '  `/team config`          View/edit team configuration',
        '  `/team learn`           Interactive tutorial & concepts',
        '  `/team help`            This guide',
        '',
        '────────────────────────────────────────────',
      ],
    },
  ],

  // Case 2: /team learn
  [
    {
      command: '/team learn',
      output: [
        '**TEAMWORK CONCEPTS** ────────────────────────',
        '',
        '**1. GitHub-First Philosophy**',
        '  Everything lives in GitHub. Issues = tasks, PRs = delivery,',
        '  Actions = CI, Milestones = versions. Zero external tools.',
        '',
        '**2. Mission Contract (MC)**',
        '  An MC is a GitHub Issue enriched by AI with:',
        '  • Objective, sub-tasks, acceptance criteria',
        '  • Context files (which code to read/modify)',
        '  • Test commands, AI notes',
        '  Local copy: `.teamwork/active/MISSION-N.md`',
        '',
        '**3. Four-Branch Model**',
        '  🟥 *main*         Production, tagged releases only',
        '  🟦 *pre-launch*   Integration hub, all PRs merge here',
        '  🟩 *mission/**    One per MC, short-lived work branches',
        '  🟧 *rc/**         Release candidates, frozen for QA',
        '',
        '**4. Worktree Mode**',
        '  Each mission gets its own git worktree + iTerm tab.',
        '  Isolation: no stashing, no branch-switching, parallel work.',
        '  Enable: `git config --local teamwork.worktree true`',
        '',
        '**5. Push Model**',
        '  No polling. Teamwork pushes notifications via iTerm2 alerts.',
        '  Team lead sees team-wide dashboard; members see their own view.',
        '',
        '**6. The Loop**',
        '  `/team-issue` → `/team-claim` → `/team-drive` → `/team-ship` → `/team-ship done`',
        '  Or shortcut: `/team auto #N` (one command, full cycle)',
        '',
        '────────────────────────────────────────────',
        '`/team help` full command reference · `/team` dashboard',
      ],
    },
  ],
]

// ─────────────────────────────────────────────
// CH12: Worktree — Isolated parallel execution
// ─────────────────────────────────────────────

export const ch12WorktreeMocks: TerminalStep[][] = [
  // Case 1: Enable worktree mode
  [
    {
      command: 'git config --local teamwork.worktree true',
      output: [
        '',
      ],
    },
    {
      command: '/team config',
      output: [
        '⚙️ **TEAMWORK CONFIG** ── v3.8.3 ────────────────',
        '**Config:**   `.teamwork/config.yml`',
        '**Worktree mode:**  `true` *(local, per-user)*',
        '',
        '────────── **Local Settings** (per-user) ──────────',
        'Worktree: `true`',
        '',
        '────────────────────────────────────────────',
        '**HOW TO EDIT:**',
        '  Worktree mode: `git config --local teamwork.worktree true/false`',
        '────────────────────────────────────────────',
      ],
    },
  ],

  // Case 2: /team-claim with worktree
  [
    {
      command: '/team-claim #42',
      output: [
        '**CLAIMED** ──── **#042** Add user authentication with JWT ────────────',
        '🟠  `P1`  Branch: *mission/042-add-auth-bob*',
        'Contract:  `.teamwork/active/MISSION-42.md`',
        'Worktree:  `~/Projects/taskflow.worktrees/042-add-auth-bob`',
        '',
        '**OBJECTIVE**',
        '  Implement JWT-based authentication for the TaskFlow API',
        '',
        '**SUB-TASKS**',
        '  [ ] Install jsonwebtoken + bcrypt dependencies',
        '  [ ] Create auth service',
        '  [ ] Implement login/register routes',
        '  [ ] Add JWT middleware',
        '  [ ] Write unit tests',
        '',
        '⚠️ **DO NOT run `/team-drive` in this terminal**',
        '────────────────────────────────────────────',
        '当前终端在主仓库。Worktree 模式下，所有操作必须在 worktree 目录中执行。',
        '在这里运行 `/team-drive` 会操作错误的代码。',
        '',
        '**打开新终端 tab，执行：**',
        '  `cd ~/Projects/taskflow.worktrees/042-add-auth-bob && claude`',
        '  Then: `/team-drive` to execute │ `/team-ship` when done',
        '────────────────────────────────────────────',
        '',
        '💡 Tip: Each worktree gets its own iTerm tab with branch color.',
      ],
      delay: 60,
    },
  ],

  // Case 3: Multiple worktrees active
  [
    {
      command: 'git worktree list',
      output: [
        '~/Projects/taskflow                                        abc1234 [pre-launch]',
        '~/Projects/taskflow.worktrees/042-add-auth-bob             def5678 [mission/042-add-auth-bob]',
        '~/Projects/taskflow.worktrees/048-fix-payment-alice        ghi9012 [mission/048-fix-payment-alice]',
      ],
    },
    {
      command: '/team',
      output: [
        '**MY DASHBOARD** — acme-corp/taskflow ──────────────────',
        'Teamwork v3.8.3 · **V1.0**  `████░░░░░░`  25%  `1`/`4`',
        '',
        '👥 **TEAM** (3 members)',
        '  alice (tech-lead)',
        '    🟢 **#048** Fix payment crash 🟠 P0 → *wip*',
        '    *mission/048-fix-payment-alice*',
        '    Worktree: `~/Projects/taskflow.worktrees/048-fix-payment-alice`',
        '  bob (backend)',
        '    🟢 **#042** Add user authentication 🟠 P1 → *wip*',
        '    *mission/042-add-auth-bob*',
        '    Worktree: `~/Projects/taskflow.worktrees/042-add-auth-bob`',
        '  charlie (frontend)',
        '    ⚪ *(idle)*',
        '',
        '📊 **Worktrees**  `2` active',
        '',
        '`/team help` · `/team doctor` · `/team #N`',
        '────────────────────────────────────────────',
      ],
    },
  ],
]
