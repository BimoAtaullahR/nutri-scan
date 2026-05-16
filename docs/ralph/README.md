# Ralph Backend Workflow

This repo is set up for a Ralph-style "one task per run" backend workflow.

Reference article: https://www.aihero.dev/getting-started-with-ralph

## Files

- `docs/prd/backend-core-scan-loop-50.md` — target PRD.
- `docs/ralph/backend-core-scan-loop-progress.md` — queue and run log.
- `scripts/ralph-backend-once.sh` — one Ralph iteration.

## Run Once With Codex

Codex CLI is available on this machine, so this is the default:

```bash
./scripts/ralph-backend-once.sh
```

Each run should:

1. Read repo docs, PRD, and progress file.
2. Pick one unchecked backend issue.
3. Implement only that issue.
4. Run focused tests.
5. Update progress.
6. Commit with a conventional commit message.

## Run Once With Claude Code

Claude Code is not currently installed on this machine. After installing and authenticating it, run:

```bash
RALPH_AGENT=claude ./scripts/ralph-backend-once.sh
```

## Progress Discipline

Only mark an issue complete when implementation, tests, progress update, and commit are done.

If a run is blocked, append the blocker to the run log instead of moving to another issue.
