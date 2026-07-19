# LeetCode Spaced Repetition

Daily practice scheduler for 165 problems (Grind 75 ∪ Grind 169 ∪ Blind 75, deduplicated, bit-manipulation excluded). Counteracts the forgetting curve with an SM-2-style algorithm: each solved problem comes back for review at growing intervals (2d → 5d → ~2w → ~1mo → ...), shrinking when you struggle.

## Daily flow

1. Every day at 3 PM PT, a GitHub Actions workflow ([`.github/workflows/daily-leetcode.yml`](../.github/workflows/daily-leetcode.yml)) runs `sr.py today --md` on GitHub's servers and opens a **new issue per day** titled "Daily LeetCode YYYY-MM-DD" with the plan as its body (previous days' plan issues are auto-closed; logging by comment still works on closed issues): ~2-3 problems per day — Easy costs 1, Medium/Hard costs 2, against a daily budget of 4. Reviews take up to half the budget (due-order); the other half is reserved for new problems so a heavy review day can't starve them. New problems are introduced by **category rotation**: each problem carries one of 16 Grind-style categories (`category` in problems.json), and the next new problem comes from the category that least recently got one — so every topic gets fresh coverage roughly every 2 weeks, and a full first pass takes ~4-5 months. The comment @-mentions you, so the GitHub mobile app pushes a notification with clickable links — no laptop needed.
2. Solve them on LeetCode (premium problems link to a free mirror on leetcode.ca).
3. Log each solve, any of three ways:
   - **From your phone (no laptop):** comment on the day's "Daily LeetCode" issue — e.g. "solved 200 good, #33 was hard". The [`log-solve` workflow](../.github/workflows/log-solve.yml) parses your comment, updates `state.json`, commits it, and replies with the next review dates. Only your own comments trigger it.
   - Tell Claude in any Cowork chat ("solved 200 good, 33 again").
   - Run it directly:

```bash
python3 sr.py log <id> <again|hard|good|easy>
```

Grades: `again` = needed the solution · `hard` = solved but slow/messy · `good` = minor friction · `easy` = clean and fast.

Unlogged problems don't advance — they carry over to the next day's plan, so skipping a day just pauses the schedule rather than breaking it.

**State sync:** the workflow commits `state.json` to `main` after each run (it records which plan was served). Run `git pull` before logging locally, and `git push` after — GitHub is the source of truth for review state.

**Timing notes:** GitHub cron is UTC, so the workflow schedules both DST variants (22:00/23:00 UTC) and a guard step keeps only the one matching 15:00 America/Los_Angeles. Scheduled runs can start a few minutes late (occasionally more) on GitHub's side. GitHub auto-disables schedules after ~60 days without repo activity — pushing your solutions counts as activity. You can also trigger a run manually from the Actions tab (workflow_dispatch).

## Commands

```bash
python3 sr.py today     # today's plan (idempotent — same output all day)
python3 sr.py log 1 good
python3 sr.py stats     # progress overview
python3 sr.py due       # every card's next review date
```

## Files

- `problems.json` — the 165-problem set: id, title, slug, difficulty, tags, source lists, URL (premium entries carry `alt_url` to leetcode.ca)
- `state.json` — created on first run; review cards, served plans, and full solve log. Delete it to reset all progress.
- `sr.py` — scheduler, stdlib only

## Tuning

Edit `config` inside `state.json`: `daily_budget` (default 4), `new_budget` (slice reserved for new problems, default 2), `new_per_day` (default 2), `new_order` (`category_rotation` or `curated`), and `cost` per difficulty (default E=1, M=2, H=2).

Pace levers, measured by simulation: budget 4 / new 2 ≈ 2.4 problems/day, ~13-day tag rotation, full pass ~4.5 months, but overdue reviews queue up over time. Raise `daily_budget` to 6-7 (keeping `new_budget` 2) to clear reviews on schedule at ~3.5 problems/day; raise both (e.g. 10/5) to finish the deck in ~2 months at ~5 problems/day. Overdue reviews aren't lost — they serve oldest-first as budget allows. To change the delivery time or switch to every-other-day, edit the cron expressions in the workflow file (remember they're UTC) — intervals are date-based, so nothing else needs to change.

## Note for Claude sessions

When the user reports solving a problem, run `python3 spaced_repetition/sr.py log <id> <grade>` from the repo root. Map casual phrasing to grades conservatively ("got it but took forever" → hard). This folder is independent of the Go module — no Go tests apply.
