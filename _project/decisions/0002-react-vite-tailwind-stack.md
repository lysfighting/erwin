# ADR-0002: React + Vite + Tailwind + Framer Motion Stack

- **Date**: 2026-03-02
- **Status**: Accepted (inherited from v1.0.0)
- **Session**: docs/sessions/2026-03-02.md

## Context

The drive-showcase needed a frontend stack when it was created at v1.0.0.

## Decision

Use React 19.2 + TypeScript + Vite 7.3 + Tailwind CSS 4.2 + Framer Motion 12.34.

## Consequences

**Positive**:
- React + TypeScript: familiar, large ecosystem, good for component-based UI
- Vite: fast HMR for development iteration
- Tailwind 4: utility-first, no runtime overhead, v4's new CSS variables approach
- Framer Motion: declarative animations already integrated throughout

**Negative**:
- Framer Motion adds ~50KB bundle weight
- React 19 is relatively new (concurrent features, not all patterns settled)

## Alternatives Considered

- **Vue/Nuxt**: Rejected at inception (React preferred for this team)
- **Next.js**: Rejected (no SSR needed for a static showcase; Vite simpler)
- **CSS animations only**: Rejected (Framer Motion provides better DX for complex sequences)

## Note for V2

V2 adds only `d3-force` (15KB, force-directed layout). No other new libraries. See ADR-0004.
