#!/usr/bin/env bash
#
# Tests for .githooks/commit-msg. Unit cases hand the hook a message file
# directly; the end-to-end cases run real `git commit`s in a throwaway repo
# with core.hooksPath pointed at this directory. Global and system git config
# are shut out so a personal core.commentChar can't change the results.
#
#   bash .githooks/test-commit-msg.sh

set -u

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
HOOK="$DIR/commit-msg"
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT

export HOME="$WORK/home" XDG_CONFIG_HOME="$WORK/home/.config" GIT_CONFIG_NOSYSTEM=1
mkdir -p "$HOME"
pass=0; fail=0

result() {  # result <want> <got> <desc>
  if [ "$1" = "$2" ]; then pass=$((pass + 1)); printf "  ok    %-58s exit %s\n" "$3" "$2"
  else fail=$((fail + 1)); printf "  FAIL  %-58s exit %s (want %s)\n" "$3" "$2" "$1"; fi
}

check() {  # check <want-exit> <desc> <message>
  printf '%s' "$3" >"$WORK/msg"
  (cd "$WORK" && bash "$HOOK" "$WORK/msg") >/dev/null 2>&1
  result "$1" $? "$2"
}

long_desc=$(printf 'x%.0s' $(seq 1 60))   # 12-char prefix + 60 = 72

echo "=== accepted ==="
check 0 "plain subject"                          "11-24-2026: Add a commit-msg hook"
check 0 "subject, body and trailer"              $'11-24-2026: Add hook\n\nWhy it exists.\n\nCo-Authored-By: A <a@example.com>\n'
check 0 "editor template comments before subject" $'# Please enter the commit message\n#\n11-24-2026: Add hook\n'
check 0 "blank lines before subject"             $'\n\n11-24-2026: Add hook\n'
check 0 "CRLF line endings"                      $'11-24-2026: Add hook\r\n\r\nBody\r\n'
check 0 "no trailing newline"                    "01-01-2026: New year"
check 0 "leap day, divisible by 4"               "02-29-2024: Leap"
check 0 "leap day, divisible by 400"             "02-29-2000: Leap"
check 0 "last day of a 30-day month"             "04-30-2026: End of April"
check 0 "exactly 72 characters"                  "11-24-2026: $long_desc"
check 0 "git merge"                              "Merge branch 'feature'"
check 0 "git pull merge"                         "Merge branch 'main' of https://github.com/example/repo"
check 0 "remote-tracking merge"                  "Merge remote-tracking branch 'origin/main'"
check 0 "GitHub merge"                           "Merge pull request #12 from example/branch"
check 0 "git revert"                             $'Revert "11-24-2026: Add hook"\n\nThis reverts commit abc.\n'
check 0 "fixup! marker"                          "fixup! 11-24-2026: Add hook"
check 0 "squash! marker"                         "squash! 11-24-2026: Add hook"

echo "=== rejected ==="
check 1 "no date"                                "Add a commit-msg hook"
check 1 "ISO date order"                         "2026-11-24: Add hook"
check 1 "DD-MM-YYYY that is not a valid month"   "24-11-2026: Add hook"
check 1 "slashes instead of dashes"              "11/24/2026: Add hook"
check 1 "single-digit month and day"             "1-2-2026: Add hook"
check 1 "two-digit year"                         "11-24-26: Add hook"
check 1 "no colon"                               "11-24-2026 Add hook"
check 1 "no space after colon"                   "11-24-2026:Add hook"
check 1 "two spaces after colon"                 "11-24-2026:  Add hook"
check 1 "date only"                              "11-24-2026:"
check 1 "date and colon, description blank"      "11-24-2026:    "
check 1 "month 00"                               "00-10-2026: Add hook"
check 1 "month 13"                               "13-01-2026: Add hook"
check 1 "day 00"                                 "11-00-2026: Add hook"
check 1 "day 31 in a 30-day month"               "04-31-2026: Add hook"
check 1 "Feb 29 in a common year"                "02-29-2026: Add hook"
check 1 "Feb 29 in a century non-leap year"      "02-29-1900: Add hook"
check 1 "73 characters"                          "11-24-2026: ${long_desc}x"
check 1 "text before the date"                   "feat 11-24-2026: Add hook"
check 1 "valid line in the body, bad subject"    $'Add hook\n\n11-24-2026: Add hook\n'
check 1 "only comments and blanks"               $'# comment\n\n# another\n'
check 1 "subject starting with # (git keeps it)" "#12: Fix bug"
check 1 "empty file"                             ""
check 1 "Merge used as an ordinary verb"         "Merge the two heap helpers"
check 1 "Revert without git's quoting"           "Revert the heap change"

echo "=== scissors line (git commit -v) ==="
check 1 "nothing above the scissors, diff below" $'\n# ------------------------ >8 ------------------------\ndiff --git a/x b/x\n'
check 0 "subject above the scissors"             $'11-24-2026: Add hook\n# ------------------------ >8 ------------------------\ndiff --git a/x b/x\n'

echo "=== error output ==="
says() {  # says <desc> <message> <substring the error must contain>
  printf '%s' "$2" >"$WORK/msg"
  local out; out=$( (cd "$WORK" && bash "$HOOK" "$WORK/msg") 2>&1 )
  case "$out" in *"$3"*) result 0 0 "$1" ;; *) result 0 1 "$1" ;; esac
}
# The exit code alone can't tell these apart: an empty subject would fail the
# format regex anyway. The messages are what the extra checks buy.
says "bad format names the expected format"  "Add hook"                 "expected: MM-DD-YYYY: <short description>"
says "bad format echoes the subject"         "Add hook"                 "subject:  Add hook"
says "'#' subject explains comments"         "#12: Fix bug"             "no subject line (lines starting with '#'"
says "scissors-only explains, not the diff"  $'# -------- >8 --------\ndiff --git a/x b/x\n' "no subject line"
says "bad date names the day"                "04-31-2026: Add hook"     "day 31 does not exist in 04/2026"
says "overlong subject gives the length"     "11-24-2026: ${long_desc}x" "subject is 73 characters"

echo "=== end to end ==="
repo="$WORK/repo"
git init -q "$repo"
git -C "$repo" config user.name test
git -C "$repo" config user.email test@example.com
git -C "$repo" config core.hooksPath "$DIR"

echo a >"$repo/a"; git -C "$repo" add a
git -C "$repo" commit -q -m "Add a" >/dev/null 2>&1
result 1 $? "git commit with a bad message is refused"
[ -z "$(git -C "$repo" rev-parse --verify -q HEAD)" ]
result 0 $? "and no commit was created"

git -C "$repo" commit -q -m "11-24-2026: Add a" >/dev/null 2>&1
result 0 $? "git commit with a good message lands"

echo b >"$repo/b"; git -C "$repo" add b
git -C "$repo" commit -q --no-verify -m "Add b" >/dev/null 2>&1
result 0 $? "--no-verify skips the check"

echo c >"$repo/c"; git -C "$repo" add c
printf '; note\n11-24-2026: Add c\n' >"$WORK/cmsg"
GIT_EDITOR=true git -C "$repo" -c core.commentChar=';' commit -q -e -F "$WORK/cmsg" >/dev/null 2>&1
result 0 $? "custom core.commentChar is honoured"
[ "$(git -C "$repo" log -1 --format=%s)" = "11-24-2026: Add c" ]
result 0 $? "and git stripped the comment from the subject"

echo
echo "$pass passed, $fail failed"
[ "$fail" -eq 0 ]
