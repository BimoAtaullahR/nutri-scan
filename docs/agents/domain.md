# Domain Docs

How the engineering skills should consume this repo's domain documentation when exploring the codebase.

## Layout

This repo uses a multi-context layout.

Start at `CONTEXT-MAP.md` in the repo root. It should point to the relevant context docs, expected initially to cover:

- mobile/web app context
- backend context
- AI/ML inference context
- shared product/domain context, if needed

System-wide architectural decisions live in `docs/adr/`.

Context-specific decisions should live near their context, for example:

- `apps/mobile/docs/adr/`
- `apps/backend/docs/adr/`
- `apps/ai/docs/adr/`

## Before exploring, read these

- `CONTEXT-MAP.md` at the repo root.
- The context-specific `CONTEXT.md` files relevant to the task.
- `docs/adr/` for system-wide decisions.
- Context-specific `docs/adr/` folders for the area being changed.

If any of these files don't exist, proceed silently. Don't flag their absence; don't suggest creating them upfront. The producer skill (`/grill-with-docs`) creates them lazily when terms or decisions actually get resolved.

## Use the glossary's vocabulary

When your output names a domain concept, use the term as defined in the relevant `CONTEXT.md`. Don't drift to synonyms the glossary explicitly avoids.

If the concept you need isn't in the glossary yet, either reconsider whether you're inventing language the project doesn't use, or note the gap for `/grill-with-docs`.

## Flag ADR conflicts

If your output contradicts an existing ADR, surface it explicitly rather than silently overriding.
