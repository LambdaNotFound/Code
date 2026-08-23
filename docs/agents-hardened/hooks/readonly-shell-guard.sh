#!/usr/bin/env bash
# readonly-shell-guard.sh
#
# PreToolUse guard for prose agents that hold Bash purely to count things.
# Allowlist: every pipeline segment must start with a read-only binary, and
# the command must contain no redirection or in-place edit.
#
# Wired in via ai-writing-auditor.md frontmatter:
#   hooks:
#     PreToolUse:
#       - matcher: "Bash"
#         hooks:
#           - type: command
#             command: "${CLAUDE_PROJECT_DIR}/.claude/agents/hooks/readonly-shell-guard.sh"
#
# Requires: jq, bash. Fails closed if jq is missing.
# Must be executable: chmod +x this file.

set -uo pipefail

ALLOWED='wc|grep|egrep|fgrep|sort|uniq|comm|diff|head|tail|cat|cut|tr|printf|echo|test|ls|basename|dirname|true|file|nl'

deny() {
  printf '{"hookSpecificOutput":{"hookEventName":"PreToolUse","permissionDecision":"deny","permissionDecisionReason":"%s"}}\n' "$1"
  exit 0
}

input=$(cat)

if ! command -v jq >/dev/null 2>&1; then
  deny "readonly-shell-guard.sh needs jq and jq is not installed. Failing closed."
fi

cmd=$(jq -r '.tool_input.command // ""' <<<"$input")
[ -z "$cmd" ] && exit 0

norm=" $(printf '%s' "$cmd" | tr '\n\t' '  ' | tr -s ' ') "

# Blank out quoted spans before parsing. A pipe, a redirect, or the word "rm"
# inside a grep pattern is data, not shell syntax, and must not trip the guard.
norm=$(printf '%s' "$norm" | sed "s/'[^']*'/QUOTED/g; s/\"[^\"]*\"/QUOTED/g")

# Any redirection writes a file. Deny outright.
if printf '%s' "$norm" | grep -Eq '(^|[^0-9<>])>{1,2}[^&]'; then
  deny "Shell redirection is blocked. This agent counts with the shell and edits with the Edit and Write tools."
fi

# Explicit mutators and interpreters, even if the leading binary looks harmless.
if printf '%s' "$norm" | grep -Eqi ' (sed +-[a-z]*i|perl +-[a-z]*i|tee|mv|cp|rm|rmdir|ln|chmod|chown|truncate|dd|install|mktemp|touch|python[0-9.]*|perl|ruby|node|osascript|sh|bash|zsh|env|xargs|find) '; then
  deny "Blocked: this agent may only run read-only counting commands. Use Edit or Write to change a file."
fi

# Every pipeline segment must begin with an allowlisted binary.
IFS=$'\n'
for seg in $(printf '%s' "$norm" | tr '|;&' '\n'); do
  first=$(printf '%s' "$seg" | tr -s ' ' | sed 's/^ *//' | cut -d' ' -f1)
  [ -z "$first" ] && continue
  case "$first" in
    *=*) deny "Blocked: variable assignment. Write literal paths into each command." ;;
  esac
  if ! printf '%s' "$first" | grep -Eq "^(${ALLOWED})$"; then
    deny "Blocked: '${first}' is not on the read-only allowlist (${ALLOWED//|/, })."
  fi
done

exit 0
