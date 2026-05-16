#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
AGENT="${RALPH_AGENT:-codex}"
MAX_ITERS="${1:-${RALPH_MAX_ITERS:-10}}"
SLEEP_SECONDS="${RALPH_SLEEP_SECONDS:-0}"
PRD="docs/prd/backend-core-scan-loop-50.md"
PROGRESS="docs/ralph/backend-core-scan-loop-progress.md"

if ! [[ "${MAX_ITERS}" =~ ^[0-9]+$ ]] || [[ "${MAX_ITERS}" -lt 1 ]]; then
  echo "Invalid iteration count: ${MAX_ITERS}. Usage: $0 <positive-integer>"
  exit 2
fi

if ! [[ "${SLEEP_SECONDS}" =~ ^[0-9]+$ ]]; then
  echo "Invalid RALPH_SLEEP_SECONDS: ${SLEEP_SECONDS}. Must be a non-negative integer."
  exit 2
fi

PROMPT="$(cat <<'PROMPT_EOF'
You are running one AFK Ralph iteration for the NutriScan backend.

Read:
- AGENTS.md
- CONTEXT-MAP.md
- services/backend/CONTEXT.md
- docs/adr/*.md
- docs/prd/backend-core-scan-loop-50.md
- docs/ralph/backend-core-scan-loop-progress.md

Task:
1. Pick exactly one unchecked backend issue from docs/ralph/backend-core-scan-loop-progress.md.
2. Prefer the lowest-numbered unblocked issue.
3. Implement only that issue.
4. Run focused tests for touched backend/contracts code.
5. Update docs/ralph/backend-core-scan-loop-progress.md:
   - mark the issue done only if done
   - append a short run log entry
6. Commit changes with a conventional commit message.

Stop condition:
- If all scoped backend issues are complete, output exactly: <promise>COMPLETE</promise>

Constraints:
- ONLY ONE ISSUE PER RUN.
- Do not implement out-of-scope PRD items.
- Do not store scan images by default.
- Preserve Backend API, AI/ML Inference, and shared contract boundaries.
- Use project domain language: Scan, Scan Lifecycle, Nudge Decision, Anonymous User, User Profile, Core Scan Loop.
- If blocked, update progress with blocker and commit only documentation if useful.
PROMPT_EOF
)"

cd "$ROOT"

for ((i = 1; i <= MAX_ITERS; i++)); do
  echo "=== Ralph AFK iteration ${i}/${MAX_ITERS} (${AGENT}) ==="

  case "$AGENT" in
    claude)
      if ! command -v claude >/dev/null 2>&1; then
        echo "claude not found. Install Claude Code or run with RALPH_AGENT=codex."
        exit 127
      fi
      RESULT="$(claude --permission-mode acceptEdits "@${PRD} @${PROGRESS}

${PROMPT}

ONLY DO ONE TASK AT A TIME.")"
      ;;
    codex)
      if ! command -v codex >/dev/null 2>&1; then
        echo "codex not found. Install Codex CLI or run with RALPH_AGENT=claude."
        exit 127
      fi
      RESULT="$(codex exec -C "$ROOT" -s workspace-write "${PROMPT}")"
      ;;
    *)
      echo "Unknown RALPH_AGENT: ${AGENT}. Use codex or claude."
      exit 2
      ;;
  esac

  echo "${RESULT}"

  if [[ "${RESULT}" == *"<promise>COMPLETE</promise>"* ]]; then
    echo "Ralph reported completion. Exiting early."
    exit 0
  fi

  if [[ "${SLEEP_SECONDS}" -gt 0 ]] && [[ "${i}" -lt "${MAX_ITERS}" ]]; then
    sleep "${SLEEP_SECONDS}"
  fi
done

echo "Reached max iterations (${MAX_ITERS}) without completion marker."
