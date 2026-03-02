# Plan: Bundle Code Splitting
## Contract: .claude/drive/bundle-splitting/contract.md

## Approach
Convert 15 ACT III section imports + SkillBuilder from eager to React.lazy() dynamic imports.
Add Suspense boundaries in ActIIIWrapper (for sections) and App.tsx (for SkillBuilder).
ActIIIWrapper already gates rendering on `isExpanded` — lazy sections won't load until user clicks "Expand".

## Implementation Sequence
1. Convert 16 imports in App.tsx to React.lazy() with named-export wrapper
2. Add Suspense boundary inside ActIIIWrapper around {children}
3. Add Suspense boundary around SkillBuilder in App.tsx
4. Build and measure chunk sizes
5. Verify app functions correctly
