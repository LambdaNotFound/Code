#!/usr/bin/env bash
#
# Tests for validate-definitions.sh, covering the three ways a definition
# set changes: add, remove, update.
#
# Everything runs against a throwaway copy of .claude/ and
# agent-team-workspace/ under /tmp, so the real tree is never mutated. Each
# case restores the sandbox from a pristine copy before the next one, so a
# failure cannot cascade.
#
#   bash .claude/hooks/test-validate-definitions.sh

set -uo pipefail
REPO="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
python3 -c 'import yaml' 2>/dev/null || { echo "needs PyYAML: python3 -m pip install pyyaml"; exit 1; }
# An explicit template: plain `mktemp -d` is not under /tmp on macOS, where
# cleanup would then refuse to remove it.
SRC=$(mktemp -d /tmp/validate-definitions.XXXXXX); BOX=$(mktemp -d /tmp/validate-definitions.XXXXXX)
cleanup() { case "$SRC$BOX" in /tmp/*/tmp/*) rm -rf "$SRC" "$BOX";; esac; }
trap cleanup EXIT

# Pristine copy, plus stubs for the source trees CLAUDE.md cites but that a
# definitions-only sandbox does not contain.
cp -r "$REPO/.claude" "$REPO/agent-team-workspace" "$SRC"/; cp "$REPO/CLAUDE.md" "$SRC"/
for p in $(cd "$SRC" && python3 agent-team-workspace/validate-definitions.py 2>&1 \
           | sed -n 's/.*FAIL dangling path: .* -> //p' | sort -u); do
  if [ -d "$REPO/$p" ]; then mkdir -p "$SRC/$p"; else mkdir -p "$SRC/$(dirname "$p")"; touch "$SRC/$p"; fi
done

reset() { rm -rf "$BOX"; mkdir -p "$BOX"; cp -r "$SRC"/. "$BOX"/; }
HOOK() { CLAUDE_PROJECT_DIR="$BOX" bash "$BOX/.claude/hooks/validate-definitions.sh"; }
pass=0; fail=0

# run <desc> <expected-exit> <payload-json> [setup shell run inside $BOX]
#     [text the hook's output must contain]
run() {
  local desc=$1 want=$2 payload=$3 setup=${4:-} text=${5:-} rc out
  reset
  if [ -n "$setup" ] && ! ( cd "$BOX" && eval "$setup" ) >/dev/null 2>&1; then
    fail=$((fail+1)); printf "  FAIL  %-52s setup failed\n" "$desc"; return
  fi
  out=$(printf '%s' "$payload" | HOOK 2>&1); rc=$?
  if [ "$rc" = "$want" ] && [[ "$out" == *"$text"* ]]; then pass=$((pass+1)); printf "  ok    %-52s exit %s\n" "$desc" "$rc"
  else fail=$((fail+1)); printf "  FAIL  %-52s exit %s (want %s)%s\n" "$desc" "$rc" "$want" "${text:+ with: $text}"; fi
}
# subst <file> <old> <new>: replace literal text, failing if it is not there.
# Stands in for `sed -i`, whose syntax differs between GNU and BSD/macOS.
subst() {
  python3 -c 'import io,sys
p,old,new=sys.argv[1:]; s=io.open(p,encoding="utf-8").read()
assert old in s
io.open(p,"w",encoding="utf-8").write(s.replace(old,new))' "$@"
}
edit() { printf '{"tool_input":{"file_path":"%s/%s"}}' "$BOX" "$1"; }
bash_() { python3 -c 'import json,sys;print(json.dumps({"tool_input":{"command":sys.argv[1]}}))' "$1"; }

echo "=== baseline ==="
run "untouched sandbox is valid" 0 "$(edit .claude/agents/rust-pro.md)"

echo "=== ADD ==="
run "new agent, not yet in CLAUDE.md"            2 "$(edit .claude/agents/demo.md)" \
    "printf -- '---\nname: demo\ndescription: D. Use never. Not for anything (use rust-pro).\ntools: Read\nmodel: inherit\n---\nBody.\n' > .claude/agents/demo.md"
run "new agent, registered in CLAUDE.md"         0 "$(edit .claude/agents/demo.md)" \
    "printf -- '---\nname: demo\ndescription: D. Use never. Not for anything (use rust-pro).\ntools: Read\nmodel: inherit\n---\nBody.\n' > .claude/agents/demo.md; printf -- '\n- \`.claude/agents/demo.md\` — demo\n' >> CLAUDE.md"
run "new skill, not yet in CLAUDE.md"            2 "$(edit .claude/skills/demo/SKILL.md)" \
    "mkdir -p .claude/skills/demo; printf -- '---\ndescription: D. Use never. Not for anything (use review-code).\n---\nBody.\n' > .claude/skills/demo/SKILL.md"
run "new orphaned reference file"                2 "$(edit .claude/skills/review-pr/references/orphan.md)" \
    "cp .claude/skills/review-pr/references/go.md .claude/skills/review-pr/references/orphan.md"
run "new agent with a duplicate name"            2 "$(edit .claude/agents/dupe.md)" \
    "cp .claude/agents/rust-pro.md .claude/agents/dupe.md"

echo "=== REMOVE ==="
run "rm an agent other definitions reference"    2 "$(bash_ 'rm .claude/agents/rust-pro.md')" \
    "rm .claude/agents/rust-pro.md"
run "git rm a skill"                             2 "$(bash_ 'git rm -r .claude/skills/scope-problem')" \
    "rm -rf .claude/skills/scope-problem"
run "rm a reference file SKILL.md still links"   2 "$(bash_ 'rm .claude/skills/review-pr/references/go.md')" \
    "rm .claude/skills/review-pr/references/go.md"
run "rm the guard hook an agent points at"       2 "$(bash_ 'rm .claude/agents/hooks/kubectl-guard.sh')" \
    "rm .claude/agents/hooks/kubectl-guard.sh"
run "rm an unrelated file"                       0 "$(bash_ 'rm /tmp/scratch.txt')"

echo "=== UPDATE ==="
run "break an agent's YAML frontmatter"          2 "$(edit .claude/agents/rust-pro.md)" \
    "subst .claude/agents/rust-pro.md 'name: rust-pro' 'bad: \"unterminated
name: rust-pro'"
run "strip a memory grant's scope rule"          2 "$(edit .claude/agents/coding-bar-raiser.md)" \
    "python3 -c \"import io,re;p='.claude/agents/coding-bar-raiser.md';s=io.open(p).read();io.open(p,'w').write(re.sub(r'## Memory.*?(?=\n## )','',s,flags=re.S))\""
run "strip a capped agent's turn budget"         2 "$(edit .claude/agents/rust-pro.md)" \
    "python3 -c \"import io,re;p='.claude/agents/rust-pro.md';s=io.open(p).read();io.open(p,'w').write(re.sub(r'## Your turn budget.*','',s,flags=re.S))\""
run "point a boundary at a nonexistent agent"    2 "$(edit .claude/skills/write-rust/SKILL.md)" \
    "subst .claude/skills/write-rust/SKILL.md '(use rust-pro)' '(use ghost-agent)'"
run "introduce a dangling repo path"             2 "$(edit .claude/skills/run-pr-loop/SKILL.md)" \
    "subst .claude/skills/run-pr-loop/SKILL.md agent-specs/pr-spec.md agent-specs/nope.md"
run "drop an agent from the CLAUDE.md roster"    2 "$(edit CLAUDE.md)" \
    "subst CLAUDE.md '\`multi-agent-coordinator\`' 'the coordinator'"
run "add a stale name to the CLAUDE.md roster"   2 "$(edit CLAUDE.md)" \
    "subst CLAUDE.md '\`review-code\` (Go/Python' '\`ghost-skill\`, \`review-code\` (Go/Python'"
run "duplicate a resume state number"            2 "$(edit .claude/skills/scope-problem/SKILL.md)" \
    "subst .claude/skills/scope-problem/SKILL.md '6. Questions logged' '5. Questions logged'"
run "benign body edit"                           0 "$(edit .claude/agents/rust-pro.md)" \
    "printf -- '\nA clarifying sentence.\n' >> .claude/agents/rust-pro.md"

echo "=== THE VALIDATOR CANNOT RUN ==="
PYTHONPATH="$BOX/noyaml" run "PyYAML missing: says so, does not block" 0 "$(edit .claude/agents/rust-pro.md)" \
    "mkdir noyaml; echo 'raise ImportError' > noyaml/yaml.py" '"additionalContext":"definition validator skipped: PyYAML'
run "validator crashes: reported as a crash"     2 "$(edit agent-team-workspace/validate-definitions.py)" \
    "subst agent-team-workspace/validate-definitions.py 'HAVE_YAML=True' 'raise RuntimeError'" "validator crashed"
run "invalid definitions: reported as invalid"   2 "$(bash_ 'rm .claude/agents/rust-pro.md')" \
    "rm .claude/agents/rust-pro.md" "definitions are now invalid"

echo "=== PATHS THE HOOK MUST IGNORE ==="
run "a Go solution"        0 "$(edit golang/tree/foo.go)"
run "the README"           0 "$(edit README.md)"
run "go test via Bash"     0 "$(bash_ 'go test ./...')"
run "git status via Bash"  0 "$(bash_ 'git status --short')"

echo
echo "$pass passed, $fail failed"
[ "$fail" -eq 0 ]
