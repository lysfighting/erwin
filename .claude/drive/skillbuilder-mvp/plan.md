# Plan: SkillBuilder MVP
## Contract: .claude/drive/skillbuilder-mvp/contract.md

## Approach Overview
Build a multi-step form component that guides users through creating a Claude Code skill.
The component uses a stepper pattern (5 steps) with navigation buttons.
A collapsible preview panel at the bottom shows the generated skill.md in real-time.
The download button creates a Blob and triggers a file download.

Template generation is the core logic — it produces substantive prompt engineering
content (Cardinal Rule, execution model, phase discipline, etc.) based on user selections,
not just empty placeholders.

## Architecture

```
SkillBuilder.tsx (main component)
├── SkillConfig type (state shape)
├── StepIndicator (visual stepper bar)
├── Step1Meta (name, trigger, description)
├── Step2Model (execution model selector)
├── Step3Phases (phase configuration)
├── Step4Rules (decision rules toggles)
├── Step5Tools (tool permission checkboxes)
├── PreviewPanel (collapsible markdown preview)
└── generateSkillMarkdown() (template engine)

lib/skillTemplate.ts (template generation logic — separate file to keep component clean)
```

## Key Design Decisions
- **Single file for component, separate file for template** — template generation is ~200 lines of string interpolation, keeping it separate avoids bloating the component
- **useState with SkillConfig object** — simple, no reducer needed for this form
- **Step content rendered conditionally** — only current step renders, no hidden DOM
- **Preview uses pre/code block** — styled with bg-bg-dark, font-mono, scrollable

## Implementation Sequence
1. Create SkillConfig type + generateSkillMarkdown in lib/skillTemplate.ts
2. Create SkillBuilder.tsx with stepper UI + all 5 steps
3. Add preview panel + download logic
4. Add i18n translations to translations.ts
5. Integrate into App.tsx (replace placeholder)
6. Verify build + visual check

## Risk Mitigation
- Template quality: pre-write the template content before coding, ensure it's substantive
- Component size: split template logic to separate file
- i18n: add all keys before integration to avoid partial translations
