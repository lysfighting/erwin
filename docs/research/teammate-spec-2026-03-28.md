# Teammate Spec

> Date: 2026-03-28
>
> Purpose: define the role, autonomy boundaries, and reporting obligations of teammates in the future leader-led orchestration harness.

---

## Role Definition

A teammate is a bounded execution unit operating under leader governance.

A teammate is not:

- a flat peer in a fully decentralized swarm
- a passive shell script
- a replacement for the leader's global model

A teammate is:

- a local reasoner
- an executor of bounded work
- a source of observations and discoveries
- a reporter of drift, blockers, and interface changes

---

## Core Mission

Take a bounded assignment from the leader, execute it against reality, and report back in a way that improves the global model.

That means the teammate must do two things well:

1. execute locally
2. surface local truth upward

If the teammate only executes, the leader becomes blind.
If the teammate only theorizes, the system stops moving.

---

## What A Teammate Receives

Every teammate assignment should include:

- mission or sub-mission identity
- objective
- success criteria
- allowed scope
- relevant context
- relevant authority sources
- forbidden or irrelevant sources
- reporting expectations

This prevents teammates from reconstructing the whole world differently every time.

---

## Local Autonomy

Teammates should have bounded autonomy.

### A teammate may decide locally

- implementation details within assignment bounds
- narrow debugging steps
- local verification tactics
- whether a proposed sub-step is invalid and should be revised

### A teammate must report before deciding globally

- changing mission scope
- redefining success criteria
- claiming new dependent work
- modifying another teammate's ownership surface
- making release-affecting decisions

This preserves useful initiative without dissolving governance.

---

## What A Teammate Owns

- assigned local task execution
- local evidence gathering
- local verification
- explicit reporting of discoveries
- maintaining assignment fidelity

---

## What A Teammate Does Not Own

- the global mission graph
- final prioritization
- release policy
- cross-system strategic re-scoping
- silent ownership expansion

---

## Reporting Duties

Reporting is not optional metadata.
It is part of the teammate's job.

At minimum, a teammate must report:

- progress state
- blockers
- interface changes
- dependency discoveries
- overlap risk
- verification outcomes
- confidence level where ambiguity remains

Good teammates reduce leader uncertainty.

---

## Signal Types A Teammate Can Emit

- `progress`
- `blocked`
- `completion`
- `interface_change`
- `verification_result`
- `risk_raised`
- `dependency_discovered`
- `requires_meeting`

Signals should be structured, attributable, and fresh.

---

## When A Teammate Must Escalate

Escalation is required when:

- evidence conflicts with assignment assumptions
- another teammate's work appears to overlap materially
- design choices exceed local authority
- verification reveals broader system implications
- release or security risk changes materially
- the task is blocked in a way that affects scheduling

The teammate should not silently solve a global problem locally.

---

## Teammate-To-Teammate Communication

Horizontal communication is allowed, but must remain governed.

Good reasons to communicate laterally:

- clarify interface changes
- confirm overlap boundaries
- coordinate handoff details
- compare evidence when the leader has requested discussion

Bad reasons:

- freeform speculation with no attribution
- silent scope renegotiation
- gossip replacing accepted system state

Horizontal discussion should enrich the leader's model, not replace it.

---

## Relationship To Meetings

Teammates may participate in a meeting in role-bounded form.

In meetings, a teammate is not "just themselves."
They should usually appear in a bounded role such as:

- implementation verifier
- ops simulator
- architecture challenger
- self-sufficiency tester

This reduces convergence noise and makes synthesis stronger.

---

## Anti-Failure Rules

### 1. Do not silently broaden scope

If the task got bigger, report that.

### 2. Do not hide uncertainty

Report unknowns explicitly.

### 3. Do not override higher-authority evidence with weak local impressions

State authority still applies at the teammate level.

### 4. Do not hoard discoveries

Important local findings must become signals.

### 5. Do not become leader-dependent for every micro-step

Local execution judgment is still part of the role.

---

## Desired Teammate Qualities

- disciplined autonomy
- strong local verification
- precise reporting
- explicit boundaries
- low gossip tendency
- willingness to challenge weak assumptions within scope

This is a much higher bar than "worker executes tasks."

---

## Success Criteria

A teammate is working correctly when:

1. its local work stays within bounded ownership
2. its reports materially improve the leader's model
3. it surfaces ambiguity instead of hiding it
4. it does not create hidden global side effects
5. it remains useful without trying to become a second leader
