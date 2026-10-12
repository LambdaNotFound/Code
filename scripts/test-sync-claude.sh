#!/usr/bin/env bash
# Regression suite for sync-claude.sh. Run it from anywhere:
#   ./test-sync-claude.sh
# It runs sync-claude.sh with the same bash that runs this file, against
# throwaway git repos in .testwork/ next to this file and a fake HOME. Your
# real ~/.claude is never touched. Runs on macOS's bash 3.2 with BSD tools, and
# on Linux.
set -u
cd "$(dirname "$0")"
W="$PWD/.testwork"; rm -rf "$W"; mkdir -p "$W"
SYNC="$PWD/sync-claude.sh"

# q() puts a stub `claude` that always passes validation first on PATH, so most
# cases neither pay for nor depend on the real CLI. qv() runs the real one.
FAKEBIN="$W/fakebin"; mkdir -p "$FAKEBIN"
printf '#!/bin/sh\necho "Validation passed"\n' > "$FAKEBIN/claude"; chmod +x "$FAKEBIN/claude"
HAVE_CLAUDE=0; command -v claude >/dev/null 2>&1 && HAVE_CLAUDE=1

P=0; F=0; S=0
chk(){ if [[ "$2" == "$3" ]]; then echo "  PASS  $1"; P=$((P+1)); else echo "  FAIL  $1 (got '$2' want '$3')"; F=$((F+1)); fi; }
skip(){ echo "  SKIP  $1"; S=$((S+1)); }
yn(){ if eval "$1"; then echo y; else echo n; fi; }
count(){ wc -l | tr -d ' '; }
acts(){ grep -E '^(copy|ok|replace|remove|untrack) ' | awk '{print $1,$2}'; }
# Change detection without GNU find/stat: compare the file listing, and look
# for anything modified after a marker file.
listing(){ find "$@" 2>/dev/null | LC_ALL=C sort | cksum; }
mark(){ MARK_LIST="$(listing "$@")"; : > "$W/.mark"; sleep 1.1; }
untouched(){ if [[ "$(listing "$@")" == "$MARK_LIST" && -z "$(find "$@" -newer "$W/.mark" 2>/dev/null | head -1)" ]]; then echo same; else echo changed; fi; }

setup(){
  rm -rf "$W/r"; mkdir -p "$W/r/GoCode/.claude" "$W/r/home/.claude/skills" "$W/r/elsewhere"
  R="$W/r/GoCode"; C="$R/.claude"; H="$W/r/home/.claude"; ST="$H/.gocode-sync"
  git -C "$R" init -q; git -C "$R" config user.email t@t; git -C "$R" config user.name t
  export HOME="$W/r/home"
}
commit(){ git -C "$R" add -A >/dev/null; git -C "$R" commit -qm "${1:-c}" >/dev/null 2>&1 || true; }
mk(){ mkdir -p "$C/skills/$1"; printf -- '---\nname: %s\ndescription: %s\n---\n%s\n' "$1" "Does $1 things." "${2:-body}" > "$C/skills/$1/SKILL.md"; }
q(){  PATH="$FAKEBIN:$PATH" "$BASH" "$SYNC" "$R" "$@" 2>&1; }   # fast
qv(){ "$BASH" "$SYNC" "$R" "$@" 2>&1; }                         # real `claude plugin validate`

echo "Testing with $BASH ($BASH_VERSION)"

# ======================================================================= skills
setup
mk go-review; mk proto-gen; mk deploy
mkdir -p "$C/skills/go-review/scripts"; printf '#!/bin/sh\necho hi\n' > "$C/skills/go-review/scripts/run.sh"; chmod +x "$C/skills/go-review/scripts/run.sh"
mkdir -p "$C/skills/not-a-skill" && echo notes > "$C/skills/not-a-skill/README.md"
mkdir -p "$H/skills/my-own" && echo "HAND MADE" > "$H/skills/my-own/SKILL.md"
mkdir -p "$H/skills/synced/from-account" && echo x > "$H/skills/synced/from-account/SKILL.md"
commit

echo "--- skills ---"
q >/dev/null
chk "3 skills installed"           "$(ls "$H/skills" | grep -vE 'my-own|synced' | count)" "3"
chk "exec bit preserved"           "$(yn "[[ -x $H/skills/go-review/scripts/run.sh ]]")" "y"
chk "dir without SKILL.md skipped" "$(yn "[[ -e $H/skills/not-a-skill ]]")" "n"
chk "other-named hand-made kept"   "$(cat "$H/skills/my-own/SKILL.md")" "HAND MADE"
chk "idempotent re-run"            "$(q | grep -c '^ok')" "3"
printf -- '---\nname: deploy\ndescription: Now with staging.\n---\nv2\n' > "$C/skills/deploy/SKILL.md"; commit
echo stale > "$H/skills/go-review/LEFTOVER"; q >/dev/null
chk "upstream edit lands"          "$(grep -c staging "$H/skills/deploy/SKILL.md")" "1"
chk "stale inner file removed"     "$(yn "[[ -e $H/skills/go-review/LEFTOVER ]]")" "n"
git -C "$R" rm -rq .claude/skills/proto-gen; commit; q >/dev/null
chk "upstream delete pruned"       "$(yn "[[ -e $H/skills/proto-gen ]]")" "n"
chk "prune spares my-own"          "$(cat "$H/skills/my-own/SKILL.md")" "HAND MADE"

echo "--- same name replaces ---"
mkdir -p "$H/skills/lint" && echo "OLD LOCAL" > "$H/skills/lint/SKILL.md" && echo junk > "$H/skills/lint/old-extra.txt"; mk lint "from repo"; commit
chk "dry-run announces replace"    "$(q --dry-run | grep -c '^replace skills/lint')" "1"
chk "dry-run leaves old content"   "$(cat "$H/skills/lint/SKILL.md")" "OLD LOCAL"
out="$(q)"
chk "run prints one line for it"   "$(grep -cE '^(replace|copy|ok) +skills/lint' <<<"$out")" "1"
chk "repo copy replaced local"     "$(grep -c 'from repo' "$H/skills/lint/SKILL.md")" "1"
chk "old extra file gone"          "$(yn "[[ -e $H/skills/lint/old-extra.txt ]]")" "n"
chk "lint now tracked"             "$(grep -cx skills/lint "$ST")" "1"
chk "re-run: no more replace"      "$(q | grep -c '^replace')" "0"
mkdir -p "$W/r/elsewhere/fmt" && echo "PRECIOUS" > "$W/r/elsewhere/fmt/SKILL.md"
ln -s "$W/r/elsewhere/fmt" "$H/skills/fmt"; mk fmt "repo fmt"; commit; q >/dev/null
chk "symlink replaced by real dir" "$(yn "[[ -d $H/skills/fmt && ! -L $H/skills/fmt ]]")" "y"
chk "symlink TARGET left intact"   "$(cat "$W/r/elsewhere/fmt/SKILL.md")" "PRECIOUS"
mkdir -p "$C/skills/synced"; printf -- '---\nname: synced\ndescription: Bad.\n---\nb\n' > "$C/skills/synced/SKILL.md"; commit
chk "'synced' name refused"        "$(q | grep -c 'reserved for skills synced from claude.ai')" "1"
chk "claude.ai synced/ intact"     "$(yn "[[ -f $H/skills/synced/from-account/SKILL.md && ! -e $H/skills/synced/SKILL.md ]]")" "y"
git -C "$R" rm -rq .claude/skills/synced; commit

# ============================================================== whole .claude
echo "--- everything under .claude ---"
mkdir -p "$C/agents" "$C/commands/db" "$C/rules" "$C/output-styles"
printf -- '---\nname: reviewer\ndescription: Reviews.\n---\nYou review.\n' > "$C/agents/reviewer.md"
printf -- '---\ndescription: Migrate.\n---\nRun migrations\n' > "$C/commands/db/migrate.md"
printf 'Prefer table tests.\n' > "$C/rules/go-testing.md"
printf 'Be terse.\n' > "$C/output-styles/terse.md"
printf '# GoCode instructions\n' > "$C/CLAUDE.md"
echo '{"permissions":{"allow":["Bash(go test:*)"]}}' > "$C/settings.json"
echo '{"configurations":[]}' > "$C/launch.json"
echo '{"bindings":[]}' > "$C/keybindings.json"
echo '{"model":"opus"}' > "$H/settings.json"; echo "PERSONAL" > "$H/CLAUDE.md"
mkdir -p "$H/agents" && echo "my agent" > "$H/agents/personal.md"
commit
out="$(q)"
chk "agent file synced"            "$(yn "[[ -f $H/agents/reviewer.md ]]")" "y"
chk "command namespace synced"     "$(yn "[[ -f $H/commands/db/migrate.md ]]")" "y"
chk "rule synced"                  "$(yn "[[ -f $H/rules/go-testing.md ]]")" "y"
chk "output style synced"          "$(yn "[[ -f $H/output-styles/terse.md ]]")" "y"
chk "keybindings synced"           "$(yn "[[ -f $H/keybindings.json ]]")" "y"
chk "launch.json synced"           "$(yn "[[ -f $H/launch.json ]]")" "y"
chk "personal agent kept (merge)"  "$(cat "$H/agents/personal.md")" "my agent"
chk "settings.json replaced"       "$(grep -c 'go test' "$H/settings.json")" "1"
chk "CLAUDE.md replaced"           "$(cat "$H/CLAUDE.md")" "# GoCode instructions"
chk "replace announced (settings)" "$(grep -c '^replace settings.json' <<<"$out")" "1"
chk "global-effect note shown"     "$(grep -c 'note: ~/.claude/settings.json now comes from the repo' <<<"$out")" "1"
chk "note silent when unchanged"   "$(q | grep -c '^note: ~/.claude')" "0"
git -C "$R" rm -q .claude/agents/reviewer.md; commit; q >/dev/null
chk "deleted agent pruned"         "$(yn "[[ -e $H/agents/reviewer.md ]]")" "n"
chk "personal agent still kept"    "$(cat "$H/agents/personal.md")" "my agent"

echo "--- what counts as an existing file ---"
mk wip "uncommitted"
echo '{"permissions":{"allow":["Bash(rm:*)"]}}' > "$C/settings.local.json"
echo "secret.md" > "$C/.gitignore"; echo "LOCAL SECRET" > "$C/secret.md"
git -C "$R" add .claude/.gitignore; commit
q >/dev/null
chk "uncommitted new file synced"  "$(yn "[[ -f $H/skills/wip/SKILL.md ]]")" "y"
chk "gitignored file not synced"   "$(yn "[[ -e $H/secret.md ]]")" "n"
chk "settings.local.json excluded" "$(yn "[[ -e $H/settings.local.json ]]")" "n"
rm "$C/rules/go-testing.md"
mkdir -p "$C/skills/deploy/extra" && echo x > "$C/skills/deploy/extra/a.md"; commit
rm "$C/skills/deploy/extra/a.md"
q >/dev/null
chk "uncommitted delete: pruned"   "$(yn "[[ -e $H/rules/go-testing.md ]]")" "n"
chk "inner uncommitted delete"     "$(yn "[[ -e $H/skills/deploy/extra/a.md ]]")" "n"
chk "...but the skill remains"     "$(yn "[[ -f $H/skills/deploy/SKILL.md ]]")" "y"
git -C "$R" add -A; commit

echo "--- never written ---"
mkdir -p "$C/worktrees/feat" "$C/projects/x" "$C/plugins" "$H/projects/p"
echo "whole checkout" > "$C/worktrees/feat/README"
echo planted > "$C/projects/x/evil.jsonl"; echo planted > "$C/plugins/evil.json"
echo planted > "$C/history.jsonl"; echo planted > "$C/.credentials.json"
echo "transcript" > "$H/projects/p/s.jsonl"; echo "hist" > "$H/history.jsonl"; echo "secret" > "$H/.credentials.json"
commit
out="$(q)"
chk "worktrees excluded"           "$(yn "[[ -e $H/worktrees ]]")" "n"
chk "projects/ not written"        "$(yn "[[ -e $H/projects/x ]]")" "n"
chk "transcripts untouched"        "$(cat "$H/projects/p/s.jsonl")" "transcript"
chk "plugins/ not written"         "$(yn "[[ -e $H/plugins/evil.json ]]")" "n"
chk "history.jsonl untouched"      "$(cat "$H/history.jsonl")" "hist"
chk "credentials untouched"        "$(cat "$H/.credentials.json")" "secret"
chk "managed paths reported"       "$(grep -c 'Claude Code manages this' <<<"$out")" "3"

echo "--- --exclude ---"
echo '{"permissions":{"allow":["Bash(go vet:*)"]}}' > "$C/settings.json"; commit
echo '{"model":"mine"}' > "$H/settings.json"
out="$(q --exclude settings.json)"
chk "--exclude: file untouched"    "$(cat "$H/settings.json")" '{"model":"mine"}'
chk "--exclude: untracked, kept"   "$(grep -c '^untrack settings.json' <<<"$out")" "1"
chk "--exclude: not in state"      "$(grep -cx settings.json "$ST")" "0"
chk "--exclude: survives re-run"   "$(q --exclude settings.json >/dev/null; cat "$H/settings.json")" '{"model":"mine"}'
q --exclude='skills/wip*' >/dev/null
chk "--exclude=glob keeps file"    "$(yn "[[ -f $H/skills/wip/SKILL.md ]]")" "y"
chk "--exclude=glob untracks"      "$(grep -cx skills/wip "$ST")" "0"
echo "LOCAL EDIT" > "$H/skills/deploy/extra.md"
out="$(q --exclude skills/deploy/extra.md)"; code=$?
chk "--exclude inside a unit refused" "$code/$(grep -c 'inside the unit skills/deploy' <<<"$out")" "2/1"
chk "...and nothing overwritten"   "$(cat "$H/skills/deploy/extra.md")" "LOCAL EDIT"
out="$(q --exclude 'skills/*/scripts')"; code=$?
chk "glob inside a unit refused"   "$code" "2"
rm "$H/skills/deploy/extra.md"

echo "--- safety guards ---"
mkdir -p "$W/r/elsewhere/agents" && echo "DOTFILES AGENT" > "$W/r/elsewhere/agents/keep.md"
rm -rf "$H/agents"; ln -s "$W/r/elsewhere/agents" "$H/agents"
mkdir -p "$C/agents"; printf -- '---\nname: a2\ndescription: x\n---\nb\n' > "$C/agents/a2.md"; commit
out="$(q)"
chk "symlinked top dir: skipped"   "$(grep -c '~/.claude/agents is a symlink' <<<"$out")" "1"
chk "...nothing written through"   "$(ls "$W/r/elsewhere/agents")" "keep.md"
rm "$H/agents"
echo '{"bad json' > "$C/keybindings.json"; commit
if command -v jq >/dev/null 2>&1 || command -v python3 >/dev/null 2>&1; then
  chk "bad JSON warns"             "$(q | grep -c 'keybindings.json is not valid JSON')" "1"
  q --strict >/dev/null; chk "bad JSON + --strict blocks" "$?" "1"
else
  skip "bad JSON checks (needs jq or python3)"
fi
echo '{"bindings":[]}' > "$C/keybindings.json"; commit
echo '{"hooks":{"Stop":[{"hooks":[{"type":"command","command":"$CLAUDE_PROJECT_DIR/x.sh"}]}]}}' > "$C/settings.json"; commit
chk "CLAUDE_PROJECT_DIR warning"   "$(q | grep -c 'uses \$CLAUDE_PROJECT_DIR')" "1"
PATH="$FAKEBIN:$PATH" "$BASH" "$SYNC" "$W/r/elsewhere" >/dev/null 2>&1
chk "non-git dir refused"          "$?" "1"

echo "--- migration from the skills-only script ---"
setup; mk go-review; commit
cp -R "$C/skills/go-review" "$H/skills/go-review"; printf 'go-review\n' > "$H/.gocode-skills"
out="$(q)"
chk "old state honoured"           "$(grep -c '^ok *skills/go-review' <<<"$out")" "1"
chk "old state file removed"       "$(yn "[[ -e $H/.gocode-skills ]]")" "n"
chk "new state written"            "$(cat "$ST")" "skills/go-review"

echo "--- dry run ---"
setup; mk go-review; mk deploy; mk old-one
mkdir -p "$C/agents" && printf -- '---\nname: r\ndescription: x\n---\nb\n' > "$C/agents/r.md"
echo '{"a":1}' > "$C/settings.json"; commit
q >/dev/null
printf -- '---\nname: deploy\ndescription: Edit.\n---\nv3\n' > "$C/skills/deploy/SKILL.md"
git -C "$R" rm -rq .claude/skills/old-one; mk brand-new
mkdir -p "$H/skills/clash" && echo LOCAL > "$H/skills/clash/SKILL.md"; mk clash
echo "LOCAL MD" > "$H/CLAUDE.md"; echo "# repo md" > "$C/CLAUDE.md"
touch "$C/skills/go-review/.DS_Store"; mkdir -p "$C/skills/go-review/__pycache__"; touch "$C/skills/go-review/__pycache__/x.pyc"
commit
mark "$HOME"
q --dry-run >/dev/null
chk "zero writes in HOME"          "$(untouched "$HOME")" "same"
if (( HAVE_CLAUDE )); then
  mark "$H/skills" "$H/agents" "$H/settings.json" "$H/CLAUDE.md" "$ST"
  qv --dry-run >/dev/null
  chk "user files untouched w/ validate" "$(untouched "$H/skills" "$H/agents" "$H/settings.json" "$H/CLAUDE.md" "$ST")" "same"
else
  skip "dry run with real validation (claude not installed)"
fi
pred="$(q --dry-run | acts)"; real="$(q | acts)"
chk "prediction == real actions"   "$pred" "$real"
chk "...covers copy/replace/remove/ok" "$(awk '{print $1}' <<<"$real" | sort -u | paste -sd, -)" "copy,ok,remove,replace"
chk "junk files not copied"        "$(find "$H/skills/go-review" -name .DS_Store -o -name __pycache__ | count)" "0"
mkdir -p "$W/fresh/.claude"; mark "$W/fresh"
PATH="$FAKEBIN:$PATH" HOME="$W/fresh" "$BASH" "$SYNC" "$R" --dry-run >/dev/null 2>&1
chk "fresh machine: zero writes"   "$(untouched "$W/fresh")" "same"

echo "--- flags ---"
printf 'no frontmatter\n' > "$C/skills/deploy/SKILL.md"; commit
if (( HAVE_CLAUDE )); then
  qv --strict >/dev/null; chk "--strict blocks on warning" "$?" "1"
  qv          >/dev/null; chk "lenient default"            "$?" "0"
else
  skip "--strict / lenient validation (claude not installed)"
fi
git -C "$R" rm -rq .claude/skills/brand-new; commit; q --no-prune >/dev/null
chk "--no-prune keeps deletion"    "$(yn "[[ -d $H/skills/brand-new ]]")" "y"
q >/dev/null
chk "default prune removes it"     "$(yn "[[ -e $H/skills/brand-new ]]")" "n"
chk "no temp dirs left"            "$(ls -d "$H/.gocode-sync."* "${TMPDIR:-/tmp}"/.gocode-sync.* 2>/dev/null | count)" "0"

echo "--- how it's launched ---"
setup; mk go-review; commit
out="$(PATH="$FAKEBIN:$PATH" "$BASH" --posix "$SYNC" "$R" 2>&1)"; code=$?
chk "bash in POSIX mode (macOS 'sh') works" "$code/$(grep -c '^copy *skills/go-review' <<<"$out")" "0/1"
for other in dash zsh; do
  if command -v "$other" >/dev/null 2>&1; then
    out="$("$other" "$SYNC" "$R" 2>&1)"; code=$?
    chk "$other: clear 'run with bash' message" "$code/$(grep -c 'This is a bash script' <<<"$out")" "1/1"
  else
    skip "$other not installed"
  fi
done

rm -rf "$W"
echo; echo "$P passed, $F failed, $S skipped"; (( F == 0 ))
