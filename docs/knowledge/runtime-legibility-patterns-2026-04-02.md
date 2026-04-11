# Runtime Legibility Patterns

> Date: 2026-04-02
>
> Purpose: capture what recent `Harness` runtime work taught the repo about leader-facing legibility, especially for delegated follow-up lifecycle work and live council execution.

---

## Bottom Line

Recent runtime work keeps reinforcing the same pattern:

**`Harness` becomes more governable when the control plane surfaces the bounded object that is actually in motion, rather than flattening everything into generic run, meeting, or next-action wording.**

Two repeated examples stand out:

1. delegated adopted follow-up needed to become a first-class lifecycle object
2. active councils needed to project their current live work packet, not just generic progress wording

These are not isolated UX improvements.
They are concrete examples of the broader `harness engineering as legibility + eval pressure` thread.

---

## Pattern 1: Delegated Follow-Up Needed Its Own Lifecycle Semantics

The repo repeatedly hit the same problem:

- follow-up work had been adopted from a council
- that follow-up was delegated to a teammate
- the packet then became stale, cancelled, completed, or reopened
- generic run-level wording could not tell the leader what kind of recovery or quieting behavior was actually needed

The repeated fixes are now captured in:

- [../feedback/cases/HU-005-delegated-follow-up-needed-bounded-lifecycle-objects.md](../feedback/cases/HU-005-delegated-follow-up-needed-bounded-lifecycle-objects.md)

What this taught the repo:

- stale delegated ownership should not look like generic stale-run cleanup
- cancelled adopted follow-up should not look like ordinary completion
- completed delegated follow-up should stop producing false audit or doctor noise
- reopened delegated follow-up should regain bounded control-plane visibility instead of depending on human memory

The durable lesson is that a delegated adopted follow-up packet is not just “some task on a teammate run.”
It is a bounded governance object with a lifecycle that leader-facing surfaces must model explicitly.

---

## Pattern 2: Active Councils Needed Current Work Packets

The repo also kept finding that a live council is too opaque if the control plane only exposes:

- a generic `next_requirement`
- a batch id
- a phase label

The repeated fixes are now captured in:

- [../feedback/cases/HU-006-active-council-needed-current-work-visibility.md](../feedback/cases/HU-006-active-council-needed-current-work-visibility.md)

What this taught the repo:

- leaders need to see the current live council packet
- they also need the concrete current-wave action packets
- and they benefit from blocked-batch summaries so queued bounded work stays visible

The durable lesson is that a live council should project the bounded work object currently in flight, not just abstract process state.

---

## Shared Interpretation

These two patterns look different on the surface, but they reinforce the same control-plane principle:

- generic summary wording is not enough when a bounded governance object is carrying the real state transition
- the leader needs the smallest truthful object that explains the current obligation
- once that bounded object is visible, scheduler, `home`, `doctor`, and audit surfaces become easier to align

This is the runtime-side version of the same lesson already visible in the broader harness-engineering thread:

- legibility is not cosmetic
- visible structure is part of control-plane reliability
- stronger observability often means surfacing the right object, not only adding more prose

---

## Why This Matters For The Knowledge System

Without notes like this, these lessons stay scattered across many implementation sessions and look like unrelated bug-fix history.

With this note in place, the repo can now connect:

- runtime feedback cases
- the external harness-engineering thread
- proposal posture around `encode-as-eval`
- future review gates about traceability and maturity claims

This is exactly what the knowledge layer is supposed to do.

---

## Suggested Next Moves

1. Add more runtime-facing feedback cases when leader-facing surfaces repeatedly need a more precise bounded object.
2. Use this note when sharpening `HLP-003` or future eval candidates about control-plane legibility.
3. Prefer checks that ask “is the right bounded object visible?” over checks that only ask whether some summary text exists.
