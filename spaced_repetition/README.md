# LeetCode Spaced Repetition

Daily practice scheduler for 165 problems (Grind 75 ∪ Grind 169 ∪ Blind 75, deduplicated, bit-manipulation excluded). Counteracts the forgetting curve with an SM-2-style algorithm: each solved problem comes back for review at growing intervals (2d → 5d → ~2w → ~1mo → ...), shrinking when you struggle.

## Daily flow

1. Every day from 3 PM PT, a GitHub Actions workflow ([`.github/workflows/daily-leetcode.yml`](../.github/workflows/daily-leetcode.yml)) runs `sr.py today --md` on GitHub's servers and opens a **new issue per day** titled "Daily LeetCode YYYY-MM-DD" with the plan as its body (previous days' issues stay open until you close them; logging by comment works either way): ~2-3 problems per day — Easy costs 1, Medium/Hard costs 2, against a daily budget of 4. Reviews take up to half the budget (due-order); the other half is reserved for new problems so a heavy review day can't starve them. New problems are introduced by **category rotation**: each problem carries one of 16 Grind-style categories (`category` in problems.json), and the next new problem comes from the category that least recently got one — so every topic gets fresh coverage roughly every 2 weeks, and a full first pass takes ~4-5 months. Problems you have never been shown come before ones you were shown and ignored. The comment @-mentions you, so the GitHub mobile app pushes a notification with clickable links — no laptop needed.
2. Solve them on LeetCode (premium problems link to a free mirror on leetcode.ca).
3. Log each solve, any of three ways:
   - **From your phone (no laptop):** comment on the day's "Daily LeetCode" issue — e.g. "solved 200 good, #33 was hard". The [`log-solve` workflow](../.github/workflows/log-solve.yml) parses your comment, updates `state.json`, commits it, and replies with the next review dates. Only your own comments trigger it. When a grade is unclear ("#9 not good", "#33 slow but ok") it logs nothing for that problem and says so in the reply. Put several solves in one comment: GitHub may drop a run when comments arrive in quick succession.
   - Tell Claude in any Cowork chat ("solved 200 good, 33 again").
   - Run it directly:

```bash
python3 sr.py log <id> <again|hard|good|easy>
```

Grades describe **how the attempt went**, not the problem's LeetCode difficulty (that's static and already in `problems.json`):

| Grade | Means | Also accepted |
|---|---|---|
| `again` | needed the solution | `failed`, `fail`, `stuck`, `blanked`, `no` |
| `hard` | solved, but slow or messy (the interval still grows, by at least a day) | `slow`, `rough`, `tough` |
| `good` | solved with minor friction | `medium`, `ok`, `okay`, `fine`, `yes` |
| `easy` | quick and clean | `trivial`, `quick`, `clean` |

`medium` maps to `good` because it reads naturally — but it means "moderate friction on this attempt", not "this is a Medium problem".

Unlogged problems carry over to the next day, then get deferred (see below) so the rotation continues — skipping a day pauses the schedule rather than breaking it.

**State sync:** the workflow commits `state.json` to `main` after each run (it records which plan was served). Run `git pull` before logging locally, and `git push` after — GitHub is the source of truth for review state.

**Timing notes:** GitHub cron is UTC-only and scheduled runs are routinely delayed (about 30 minutes in August 2026, 2-4 hours by November), so the workflow fires hourly at :17 from 19:17 to 06:17 UTC. The guard proceeds only when the local Pacific time is 15:00 or later *and* today's issue doesn't exist yet, so the first attempt that lands wins and the later crons act as automatic retries for a delayed or failed run.

**If you don't log:** a problem served `carry_limit` times (default 2, not necessarily on consecutive days) without a grade is assumed skipped and pushed back — 4 days, then 8, 16, up to `defer_max_days` (24) if you keep ignoring it. Deferred items return automatically and nothing is dropped from the deck: a review comes back when its window ends, and a new problem you never logged waits behind the never-seen ones after that. Logging a problem clears its deferral and resets the count. Separately, a review card already shown recently but never graded is sorted to the *back* of the review queue, so it can't keep winning on due-date order and crowding out cards that are genuinely waiting.

**If you take a break:** after `pause_after` days without a log (default 7) plans stop: no issues, no deferrals, nothing piling up. One "Daily LeetCode paused since YYYY-MM-DD" issue tells you. To come back, comment `resume` on any Daily LeetCode issue, log any solve, or run the workflow by hand. That restores every card's real due date and settles the showings you ignored, so your old problems come back first, with one new problem a day mixed in. While paused, `state.json` changes once a month, because GitHub disables scheduled workflows after about 60 days without activity.

**Why the same problems recur:** three legitimate causes before suspecting a bug — (1) a problem graded `again` is deliberately due the very next day; (2) an ungraded problem is re-offered once before being deferred; (3) a large overdue backlog means the oldest cards keep surfacing. Only the third grows without bound, which is why overdue reviews earn extra daily capacity (`backlog_step`, `backlog_bonus_max`). Grading what you're served is what makes the intervals spread out. GitHub auto-disables schedules after ~60 days without repo activity — pushing your solutions counts as activity. You can also trigger a run manually from the Actions tab (workflow_dispatch).

## Commands

```bash
python3 sr.py today     # today's plan (idempotent — same output all day)
python3 sr.py log 1 good
python3 sr.py stats     # progress overview
python3 sr.py due       # every card's next review date
python3 sr.py status    # "active" or "paused YYYY-MM-DD"
python3 sr.py resume    # end a pause
```

## Files

- `problems.json` — the 165-problem set: id, title, slug, difficulty, tags, source lists, URL (premium entries carry `alt_url` to leetcode.ca)
- `state.json` — created on first run; review cards, served plans, and full solve log. Delete it to reset all progress.
- `sr.py` — scheduler, stdlib only
- `test_sr.py` — its tests, run against a temporary state file: `python3 spaced_repetition/test_sr.py`

## Tuning

Edit `config` inside `state.json`: `daily_budget` (default 4), `new_budget` (slice reserved for new problems, default 2), `new_per_day` (default 2), `new_order` (`category_rotation` or `curated`), `cost` per difficulty (default E=1, M=2, H=2), `carry_limit` (default 2), `defer_days` (default 4), `defer_max_days` (24), `backlog_step` (5), `backlog_bonus_max` (2), `pause_after` (7; 0 turns pausing off) and `max_hard` (2: at most this many Hard problems a day, reviews and new together; a Hard review over it waits a day, a Hard new problem gives way to the next one).

Pace levers, measured by simulation: budget 4 / new 2 ≈ 2.4 problems/day, ~13-day tag rotation, full pass ~4.5 months, but overdue reviews queue up over time. Raise `daily_budget` to 6-7 (keeping `new_budget` 2) to clear reviews on schedule at ~3.5 problems/day; raise both (e.g. 10/5) to finish the deck in ~2 months at ~5 problems/day. Overdue reviews aren't lost — they serve oldest-first as budget allows. To change the delivery time or switch to every-other-day, edit the cron expressions in the workflow file (remember they're UTC) — intervals are date-based, so nothing else needs to change.

## Note for Claude sessions

When the user reports solving a problem, run `python3 spaced_repetition/sr.py log <id> <grade>` from the repo root. Map casual phrasing to grades conservatively ("got it but took forever" → hard). This folder is independent of the Go module — no Go tests apply.
