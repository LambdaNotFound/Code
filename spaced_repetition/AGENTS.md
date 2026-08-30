# spaced_repetition/

Python, stdlib-only, independent of the Go module. Read [README.md](README.md) first — it documents the daily GitHub Actions flow, the SM-2-style scheduling algorithm, and the `sr.py` CLI in full; this file only orients an agent among the four files.

- `sr.py` — the scheduler: builds the daily plan, logs solves, computes review dates. Entry point for all commands (`today`, `log`, `stats`, `due`).
- `problems.json` — the fixed 165-problem set (Grind 75 ∪ Grind 169 ∪ Blind 75, deduplicated, bit-manipulation excluded): id, title, slug, difficulty, tags, category, URLs.
- `state.json` — generated, not hand-edited except via the `config` block (tuning knobs) — everything else (review cards, solve log, served plans) is written by `sr.py`. Two GitHub Actions workflows (`.github/workflows/daily-leetcode.yml`, `log-solve.yml`) also commit changes to this file directly to `main`, so `git pull` before editing it locally and `git push` after, per the root CLAUDE.md.

Root CLAUDE.md has the trigger for when to run `sr.py log` from a Claude session ("solved 200 good" style requests) — don't relitigate that here.
