# LeetCode Spaced Repetition

Daily practice scheduler for 165 problems (Grind 75 ∪ Grind 169 ∪ Blind 75, deduplicated, bit-manipulation excluded). Counteracts the forgetting curve with an SM-2-style algorithm: each solved problem comes back for review at growing intervals (2d → 5d → ~2w → ~1mo → ...), shrinking when you struggle.

## Daily flow

1. Every day at 6 PM PT, a GitHub Actions workflow ([`.github/workflows/daily-leetcode.yml`](../.github/workflows/daily-leetcode.yml)) runs `sr.py today --md` on GitHub's servers and posts the plan as a comment on the open **"Daily LeetCode"** issue (created automatically on first run): due reviews first, then 1 new problem, max 3 total. The comment @-mentions you, so the GitHub mobile app pushes a notification with clickable links — no laptop needed.
2. Solve them on LeetCode (premium problems link to a free mirror on leetcode.ca).
3. Log each solve — either tell Claude in any Cowork chat ("solved 200 good, 33 again") or run:

```bash
python3 sr.py log <id> <again|hard|good|easy>
```

Grades: `again` = needed the solution · `hard` = solved but slow/messy · `good` = minor friction · `easy` = clean and fast.

Unlogged problems don't advance — they carry over to the next day's plan, so skipping a day just pauses the schedule rather than breaking it.

**State sync:** the workflow commits `state.json` to `main` after each run (it records which plan was served). Run `git pull` before logging locally, and `git push` after — GitHub is the source of truth for review state.

**Timing notes:** GitHub cron is UTC, so the workflow schedules both DST variants (01:00/02:00 UTC) and a guard step keeps only the one matching 18:00 America/Los_Angeles. Scheduled runs can start a few minutes late (occasionally more) on GitHub's side. GitHub auto-disables schedules after ~60 days without repo activity — pushing your solutions counts as activity. You can also trigger a run manually from the Actions tab (workflow_dispatch).

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

Edit `config` inside `state.json`: `new_per_day` (default 1) and `daily_cap` (default 3, reviews + new combined). To change the delivery time or switch to every-other-day, edit the cron expressions in the workflow file (remember they're UTC) — intervals are date-based, so nothing else needs to change.

## Note for Claude sessions

When the user reports solving a problem, run `python3 spaced_repetition/sr.py log <id> <grade>` from the repo root. Map casual phrasing to grades conservatively ("got it but took forever" → hard). This folder is independent of the Go module — no Go tests apply.
