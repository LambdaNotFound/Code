#!/usr/bin/env bash
# kubectl-guard.sh
#
# PreToolUse guard for the kubernetes-specialist subagent.
# Turns that agent's prose safety rules into an actual control:
# blocks mutating cluster commands and Secret value reads before the
# Bash tool runs them.
#
# Wired in via kubernetes-specialist.md frontmatter:
#   hooks:
#     PreToolUse:
#       - matcher: "Bash"
#         hooks:
#           - type: command
#             command: "${CLAUDE_PROJECT_DIR}/.claude/agents/hooks/kubectl-guard.sh"
#
# Requires: jq, bash. Fails closed if jq is missing.
# Must be executable: chmod +x this file.

set -uo pipefail

deny() {
  printf '{"hookSpecificOutput":{"hookEventName":"PreToolUse","permissionDecision":"deny","permissionDecisionReason":"%s"}}\n' "$1"
  exit 0
}

input=$(cat)

if ! command -v jq >/dev/null 2>&1; then
  deny "kubectl-guard.sh needs jq and jq is not installed. Failing closed."
fi

cmd=$(jq -r '.tool_input.command // ""' <<<"$input")
[ -z "$cmd" ] && exit 0

# Flatten to a single spaced line so the word-boundary patterns below hold
# across multi-line and pipelined commands.
norm=" $(printf '%s' "$cmd" | tr '\n\t' '  ' | tr -s ' ') "

# Only inspect commands that actually invoke a cluster client.
if ! printf '%s' "$norm" | grep -Eqi '(^| |[;&|(]) *(kubectl|oc|helm) '; then
  exit 0
fi

allow_dry_run() {
  printf '%s' "$norm" | grep -Eqi -- '--dry-run(=| |$)'
}

# --- Mutating verbs -------------------------------------------------------
MUTATE='apply|create|replace|patch|edit|scale|autoscale|expose|run|delete|cordon|uncordon|drain|taint|annotate|label|set|attach|exec|cp|debug'
if printf '%s' "$norm" | grep -Eqi " (${MUTATE}) "; then
  if ! allow_dry_run; then
    deny "Mutating cluster command blocked by kubectl-guard. Write the manifest or print the command and let the user run it."
  fi
fi

# rollout status/history are reads; restart/undo/pause/resume are not.
if printf '%s' "$norm" | grep -Eqi ' rollout +(restart|undo|pause|resume) '; then
  deny "kubectl rollout restart/undo/pause/resume blocked by kubectl-guard. Print the command and let the user run it."
fi

# helm state changes.
if printf '%s' "$norm" | grep -Eqi ' helm +(install|upgrade|uninstall|delete|rollback) '; then
  if ! allow_dry_run; then
    deny "helm state change blocked by kubectl-guard. Render with 'helm template' instead."
  fi
fi

# --- Secret value reads ---------------------------------------------------
# Plain 'get secrets' (table output) is fine. Any -o/--output on a secret
# can print base64 values, so block it.
if printf '%s' "$norm" | grep -Eqi ' (secret|secrets)(/[^ ]+)? '; then
  if printf '%s' "$norm" | grep -Eqi -- '(-o|--output)(=| )'; then
    deny "Secret value read blocked by kubectl-guard. Use 'kubectl describe secret NAME' to see keys without values."
  fi
fi

exit 0
