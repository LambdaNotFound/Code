#!/usr/bin/env bash
#
# Copy GoCode's .claude/ folder into ~/.claude/ so what it holds applies in
# every project on this machine, not just GoCode.
#
# Idempotent. Re-run after any git pull that touches .claude/.
#
#   ./sync-claude.sh                   sync from $GOCODE_REPO (default ~/src/GoCode)
#   ./sync-claude.sh ~/work/GoCode     or name the checkout explicitly
#   ./sync-claude.sh --dry-run         report what would change, write nothing
#   ./sync-claude.sh --exclude PATH    leave PATH alone: a top-level name or a unit
#                                      (see below); repeatable, globs allowed, e.g.
#                                      --exclude settings.json --exclude 'skills/wip-*'
#   ./sync-claude.sh --no-prune        keep things that were deleted from the repo
#   ./sync-claude.sh --strict          refuse to sync if validation reports anything
#
# What is copied: every file under GoCode/.claude that git does not ignore,
# committed or not, except .DS_Store and Python caches. Gitignored files are
# machine-local, so they stay put.
#
# Files are synced in units. A unit is a top-level file (settings.json,
# CLAUDE.md) or one entry inside a top-level folder (skills/<name>/,
# agents/<name>.md). The repo wins: a unit that already exists in ~/.claude is
# deleted and replaced, unless that is unsafe (see "held" below). The top-level
# folders themselves are merged, never replaced, so anything in ~/.claude that
# the repo doesn't have is left alone. Because a unit is replaced whole,
# --exclude can't protect a path inside one, and refuses such a path.
#
# Never written: Claude Code's own data in ~/.claude (projects/, plugins/,
# sessions/, history.jsonl, ...), skills/synced/ from claude.ai, any top-level
# or unit name that starts with a dot (credentials, trash, this script's
# state), and EXCLUDE. Dotfiles inside a unit, like skills/<name>/.env, are
# copied with it.
#
# Pruning removes only units this script installed or replaced, as listed in
# ~/.claude/.gocode-sync. A hand-made unit that the repo once replaced counts
# as installed, so it is deleted if the repo later drops it.

# Needs bash; the 3.2 that ships with macOS is enough. dash and zsh leave
# BASH_VERSION unset, so stop with a clear message instead of a syntax error.
if [ -z "${BASH_VERSION:-}" ]; then
  echo "This is a bash script. Run it as ./sync-claude.sh or bash sync-claude.sh" >&2
  exit 1
fi
# `sh sync-claude.sh` on macOS runs bash in POSIX mode, where bash 3.2 rejects
# the process substitution this script uses.
set +o posix

set -euo pipefail
shopt -s nullglob

# Paths under .claude/ that are never synced. Edit to taste, or add with --exclude.
EXCLUDE=(
  settings.local.json   # project-only: Claude Code never reads it from ~/.claude
  worktrees             # `claude --worktree` checkouts, i.e. whole copies of the repo
)

# Claude Code's own data in ~/.claude: the names on the docs page "Explore the
# .claude directory", plus downloads, ide, teams and telemetry, seen in a live
# install. Not exhaustive. Never written, even if the repo has a path by one of
# these names.
MANAGED=(
  projects plugins sessions file-history plans debug paste-cache image-cache
  uploads dev-mods session-env tasks shell-snapshots backups feedback-bundles
  feedback usage-data todos statsig logs cache jobs daemon
  downloads ide teams telemetry
  history.jsonl stats-cache.json remote-settings.json
  policy-limits.json policy-limits.json.stamp.json
)

DRY_RUN=0
PRUNE=1
STRICT=0
REPO="${GOCODE_REPO:-$HOME/src/GoCode}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run)   DRY_RUN=1 ;;
    --no-prune)  PRUNE=0 ;;
    --strict)    STRICT=1 ;;
    --exclude)   [[ $# -ge 2 ]] || { echo "--exclude needs a path" >&2; exit 2; }
                 EXCLUDE+=("${2%/}"); shift ;;
    --exclude=*) v="${1#--exclude=}"; EXCLUDE+=("${v%/}") ;;
    -h|--help)   awk 'NR>1 { if (!/^#/) exit; sub(/^# ?/, ""); print }' "$0"; exit 0 ;;
    -*)          echo "unknown flag: $1" >&2; exit 2 ;;
    *)           REPO="$1" ;;
  esac
  shift
done

# A unit is copied whole, so excluding a path inside one can't stop that path
# from being overwritten. Refuse the pattern instead.
for pat in ${EXCLUDE[@]+"${EXCLUDE[@]}"}; do
  if [[ $pat == */*/* ]]; then
    rest="${pat#*/}"
    echo "--exclude $pat: that is inside the unit ${pat%%/*}/${rest%%/*}, which is synced whole." >&2
    echo "Exclude the whole unit instead, or remove the path from the repo." >&2
    exit 2
  fi
done

SRC="$REPO/.claude"
DEST="$HOME/.claude"
STATE="$DEST/.gocode-sync"
OLD_STATE="$DEST/.gocode-skills"   # from the skills-only version of this script

if ! git -C "$REPO" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "$REPO is not a git checkout. Pass the GoCode path, or set GOCODE_REPO." >&2
  exit 1
fi
if [[ ! -d "$SRC" ]]; then
  echo "No .claude folder in $REPO" >&2
  exit 1
fi

# matches UNIT PATTERN...: true if UNIT is one of the patterns or sits under one.
matches() {
  local unit="$1" pat
  shift
  for pat in "$@"; do
    # Unquoted on purpose: patterns may be globs.
    [[ $unit == $pat || $unit == $pat/* ]] && return 0
  done
  return 1
}

# Prints why UNIT must never be written to ~/.claude, or nothing if it may be.
reserved_reason() {
  local unit="$1" part
  local -a parts
  IFS=/ read -ra parts <<<"$unit"
  for part in ${parts[@]+"${parts[@]}"}; do
    if [[ $part == .* ]]; then echo "names starting with a dot are never synced"; return; fi
  done
  if matches "${parts[0]}" ${MANAGED[@]+"${MANAGED[@]}"}; then echo "Claude Code manages this in ~/.claude"; return; fi
  if [[ $unit == skills/synced ]]; then echo "reserved for skills synced from claude.ai"; return; fi
}

# True if $1 parses as JSON, or if neither jq nor python3 is installed to check.
json_ok() {
  if command -v jq >/dev/null 2>&1; then
    jq empty "$1" >/dev/null 2>&1
  elif command -v python3 >/dev/null 2>&1; then
    python3 -c 'import json, sys; json.load(open(sys.argv[1]))' "$1" 2>/dev/null
  fi
}

# Bash 3.2 has no associative arrays, so the bookkeeping uses plain arrays:
# units/ufiles/ustatus are parallel (same index = same unit), and tracked is
# a list. These stay small (dozens of entries), so linear scans are fine.
units=()     # unit names, in the order first seen
ufiles=()    # each unit's files, one per line, relative to .claude/
ustatus=()   # each unit's classification (see "classify units")
tracked=()   # units this script installed or replaced

# Sets idx to the index of unit $1 in units, or -1.
find_unit() {
  local i
  idx=-1
  for (( i = ${#units[@]} - 1; i >= 0; i-- )); do
    if [[ ${units[i]} == "$1" ]]; then idx=$i; return; fi
  done
}

# Prints the status of unit $1, or "gone" if the repo no longer has it.
status_of() {
  find_unit "$1"
  if (( idx < 0 )); then echo gone; else echo "${ustatus[idx]}"; fi
}

is_tracked() {
  local t
  for t in ${tracked[@]+"${tracked[@]}"}; do [[ $t == "$1" ]] && return 0; done
  return 1
}
track()   { is_tracked "$1" || tracked+=("$1"); }
untrack() {
  local t
  local -a keep
  keep=()
  for t in ${tracked[@]+"${tracked[@]}"}; do [[ $t == "$1" ]] || keep+=("$t"); done
  tracked=(${keep[@]+"${keep[@]}"})
}

# ------------------------------------------------------------ what we installed
if [[ -f "$STATE" ]]; then
  while IFS= read -r u; do [[ -n $u ]] && track "$u"; done < "$STATE"
elif [[ -f "$OLD_STATE" ]]; then
  while IFS= read -r n; do [[ -n $n ]] && track "skills/$n"; done < "$OLD_STATE"
fi

# ------------------------------------------------------------ list repo files
prev=""
while IFS= read -r -d '' path; do
  # A file with merge conflicts is listed once per conflict stage, back to back.
  [[ $path != "$prev" ]] || continue
  prev="$path"
  rel="${path#.claude/}"
  rel="${rel%/}"
  case "/$rel" in */.DS_Store | */__pycache__/* | *.pyc) continue ;; esac
  # Tracked but deleted in the working tree: treat it as already removed.
  [[ -e "$SRC/$rel" || -L "$SRC/$rel" ]] || continue
  if [[ $rel == */* ]]; then
    rest="${rel#*/}"
    unit="${rel%%/*}/${rest%%/*}"
  else
    unit="$rel"
  fi
  find_unit "$unit"
  if (( idx < 0 )); then
    units+=("$unit"); ufiles+=(""); ustatus+=("")
    idx=$(( ${#units[@]} - 1 ))
  fi
  ufiles[idx]+="$rel"$'\n'
done < <(git -C "$REPO" ls-files -z --cached --others --exclude-standard -- .claude)

if (( ${#units[@]} == 0 )); then
  echo "Nothing to sync: no files under $SRC that git doesn't ignore." >&2
  exit 1
fi

# ------------------------------------------------------------ classify units
# ok:       will be synced
# excluded: matches EXCLUDE                  (not synced; untracked if it was)
# reserved: reserved_reason gives a reason   (not synced; untracked if it was)
# notskill: skills/<x> without SKILL.md      (not synced; pruned like a deletion)
# held:     in the repo, but unsafe to write (not synced; never pruned)
for (( i = 0; i < ${#units[@]}; i++ )); do
  unit="${units[i]}"
  why="$(reserved_reason "$unit")"
  if [[ -n $why ]]; then
    echo "skip    $unit  ($why)" >&2
    ustatus[i]=reserved
  elif matches "$unit" ${EXCLUDE[@]+"${EXCLUDE[@]}"}; then
    echo "exclude $unit"
    ustatus[i]=excluded
  elif [[ $unit == skills/* ]] && ! grep -qxF "$unit/SKILL.md" <<<"${ufiles[i]}"; then
    echo "skip    $unit  (no SKILL.md)" >&2
    ustatus[i]=notskill
  elif [[ $unit == */* && -L "$DEST/${unit%%/*}" ]]; then
    echo "skip    $unit  (~/.claude/${unit%%/*} is a symlink; not writing through it)" >&2
    ustatus[i]=held
  elif [[ $unit != */* && -d "$DEST/$unit" && ! -L "$DEST/$unit" ]]; then
    echo "skip    $unit  (would replace the directory ~/.claude/$unit)" >&2
    ustatus[i]=held
  else
    ustatus[i]=ok
  fi
done

# ------------------------------------------------------------------ validate
issues=0

bad_json=0
for (( i = 0; i < ${#units[@]}; i++ )); do
  [[ ${ustatus[i]} == ok ]] || continue
  while IFS= read -r f; do
    [[ $f == *.json ]] || continue
    if ! json_ok "$SRC/$f"; then
      echo "warning: $f is not valid JSON" >&2
      bad_json=1
    fi
  done <<<"${ufiles[i]}"
done
(( bad_json )) && issues=1

# `claude plugin validate` checks skills, agents, and commands. It only warns
# about bad frontmatter, so a clean exit doesn't mean every file is
# well-formed; --strict turns any reported issue into a refusal.
if command -v claude >/dev/null 2>&1; then
  echo "Validating $SRC ..."
  val_status=0
  val_out="$(claude plugin validate "$SRC" 2>&1)" || val_status=$?
  printf '%s\n\n' "$val_out"
  if (( val_status != 0 )); then
    echo "Validation reported errors. Not syncing; fix the files above." >&2
    exit 1
  fi
  grep -qE 'Found [0-9]+ (warning|error)' <<<"$val_out" && issues=1
else
  echo "note: 'claude' not on PATH, so skills, agents and commands weren't validated." >&2
fi

if (( issues && STRICT )); then
  echo "--strict: validation reported issues. Not syncing." >&2
  exit 1
elif (( issues )); then
  echo "warning: validation reported issues above; syncing anyway." >&2
  echo "         Re-run with --strict to refuse instead." >&2
  echo
fi

# --------------------------------------------------------------------- stage
# Copy every unit to a staging area first. The real run stages inside ~/.claude
# so the final move is a rename; a dry run stages in the system temp folder so
# it writes nothing under your home directory. Either way, copy-in diffs the
# staged files against ~/.claude, so a dry run prints the same ok/copy/replace
# lines a real run would.
if (( DRY_RUN )); then
  stage_base="${TMPDIR:-/tmp}"
else
  mkdir -p "$DEST"
  stage_base="$DEST"
fi
stage="$(mktemp -d "${stage_base%/}/.gocode-sync.XXXXXX")"
trap 'rm -rf "$stage"' EXIT

for (( i = 0; i < ${#units[@]}; i++ )); do
  [[ ${ustatus[i]} == ok ]] || continue
  unit="${units[i]}"
  failed=""
  while IFS= read -r f; do
    [[ -n $f ]] || continue
    if [[ ! -f "$SRC/$f" ]]; then failed="$f is not a regular file"; break; fi
    mkdir -p "$stage/$(dirname "$f")"
    cp -p "$SRC/$f" "$stage/$f" 2>/dev/null || { failed="could not copy $f"; break; }
  done <<<"${ufiles[i]}"
  if [[ -n $failed ]]; then
    echo "skip    $unit  ($failed)" >&2
    rm -rf "${stage:?}/$unit"
    ustatus[i]=held
  fi
done

# --------------------------------------------------------------------- prune
while IFS= read -r unit; do
  [[ -n $unit ]] || continue
  case "$(status_of "$unit")" in
    ok | held)
      ;;
    excluded | reserved)
      echo "untrack $unit  (no longer synced; left in place)"
      untrack "$unit"
      ;;
    *)  # gone from the repo, or no longer a valid skill
      (( PRUNE )) || continue
      if [[ -n "$(reserved_reason "$unit")" ]] || matches "$unit" ${EXCLUDE[@]+"${EXCLUDE[@]}"}; then
        untrack "$unit"
      elif [[ $unit == */* && -L "$DEST/${unit%%/*}" ]]; then
        :   # don't delete through a symlinked top-level folder; keep tracking it
      else
        if [[ -e "$DEST/$unit" || -L "$DEST/$unit" ]]; then
          echo "remove  $unit  (no longer in the repo)"
          (( DRY_RUN )) || rm -rf "${DEST:?}/$unit"
        fi
        untrack "$unit"
      fi
      ;;
  esac
done < <(printf '%s\n' ${tracked[@]+"${tracked[@]}"} | LC_ALL=C sort)

# ------------------------------------------------------------------ copy in
changed_settings=0
changed_claude_md=0
for (( i = 0; i < ${#units[@]}; i++ )); do
  [[ ${ustatus[i]} == ok ]] || continue
  unit="${units[i]}"
  target="$DEST/$unit"
  is_changed=0

  if [[ -e "$target" || -L "$target" ]] && ! is_tracked "$unit"; then
    kind="directory"
    if [[ -L "$target" ]]; then kind="symlink"; elif [[ -f "$target" ]]; then kind="file"; fi
    echo "replace $unit  (existing $kind was not installed by this script)"
    is_changed=1
  elif [[ -e "$target" ]] && diff -rq "$stage/$unit" "$target" >/dev/null 2>&1; then
    echo "ok      $unit  (unchanged)"
  else
    echo "copy    $unit"
    is_changed=1
  fi
  track "$unit"
  if (( is_changed )); then
    case "$unit" in
      settings.json) changed_settings=1 ;;
      CLAUDE.md)     changed_claude_md=1 ;;
    esac
  fi

  (( DRY_RUN )) && continue
  mkdir -p "$(dirname "$target")"
  rm -rf "$target"
  mv "$stage/$unit" "$target"
done

# ------------------------------------------------------------------- notes
# These two files are per-project in GoCode but global in ~/.claude.
if (( changed_settings )); then
  echo
  echo "note: ~/.claude/settings.json now comes from the repo. Its permissions,"
  echo "      hooks and env apply in every project, not just GoCode."
  if grep -q 'CLAUDE_PROJECT_DIR' "$SRC/settings.json"; then
    echo "warning: it uses \$CLAUDE_PROJECT_DIR, which points at whichever project" >&2
    echo "         you're in. Hooks that rely on it will fail outside GoCode." >&2
  fi
fi
if (( changed_claude_md )); then
  echo
  echo "note: ~/.claude/CLAUDE.md now comes from the repo. Its instructions load"
  echo "      in every project, not just GoCode."
fi

# -------------------------------------------------------------- save state
echo
if (( DRY_RUN )); then
  echo "Dry run: nothing written."
else
  printf '%s\n' ${tracked[@]+"${tracked[@]}"} | LC_ALL=C sort | grep -v '^$' > "$STATE" || true
  rm -f "$OLD_STATE"
  echo "${#tracked[@]} item(s) synced into $DEST, tracked in $STATE"
fi
