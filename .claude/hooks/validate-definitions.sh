#!/usr/bin/env bash
#
# PostToolUse hook: re-validate agent and skill definitions after any edit
# that touches them.
#
# Wired in .claude/settings.json under hooks.PostToolUse, matcher
# "Edit|Write|MultiEdit". Reads the hook payload on stdin and does its own
# path filtering rather than relying on an `if` condition, so one entry
# covers editing tools and Bash (where deletes happen) and the logic stays
# testable in isolation:
#
#   echo '{"tool_input":{"file_path":"/abs/path"}}' | .claude/hooks/validate-definitions.sh
#
# PostToolUse cannot block — the edit has already happened. Exit 2 shows
# stderr to Claude so it fixes what it just broke; exit 0 stays silent.
# This is a lint, not a safety guard, so it fails open: a missing
# interpreter reports itself and gets out of the way rather than crying
# wolf on every edit.

set -uo pipefail

root="${CLAUDE_PROJECT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)}"
validator="$root/agent-team-workspace/validate-definitions.py"

payload=$(cat)

note() {  # non-blocking message to Claude, exit 0
  printf '{"hookSpecificOutput":{"hookEventName":"PostToolUse","systemMessage":"%s"}}\n' "$1"
  exit 0
}

command -v jq >/dev/null 2>&1 || note "definition validator skipped: jq not installed"
path=$(jq -r '.tool_input.file_path // ""' <<<"$payload")

if [ -n "$path" ]; then
  # Edit/Write/MultiEdit: decide on the file that was written.
  case "$path" in
    */.claude/agents/*|*/.claude/skills/*|*/CLAUDE.md|*/validate-definitions.py) ;;
    *) exit 0 ;;
  esac
else
  # Bash has no file_path. A delete or move goes through it, and deleting a
  # definition leaves stale references that the write-side check never sees,
  # so inspect the command instead. Anything naming .claude or the validator
  # is worth one 174ms run; everything else leaves immediately.
  cmd=$(jq -r '.tool_input.command // ""' <<<"$payload")
  case "$cmd" in
    *.claude/agents/*|*.claude/skills/*|*CLAUDE.md*|*validate-definitions*) ;;
    *) exit 0 ;;
  esac
  path="$cmd"
fi

[ -f "$validator" ] || note "definition validator missing at $validator"
command -v python3 >/dev/null 2>&1 || note "definition validator skipped: python3 not installed"

if output=$(python3 "$validator" 2>&1); then
  exit 0
fi

{
  echo "Agent/skill definitions are now invalid — you edited ${path#"$root"/} and the validator fails:"
  echo
  grep -E '^  FAIL|hard checks passed' <<<"$output" || echo "$output" | tail -20
  echo
  echo "Fix these before continuing. Re-run: python3 agent-team-workspace/validate-definitions.py"
} >&2
exit 2
