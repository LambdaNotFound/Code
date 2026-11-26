# spaced_repetition/

Python, stdlib-only, independent of the Go module. Read [README.md](README.md) first — it documents the daily GitHub Actions flow, the SM-2-style scheduling algorithm, and the `sr.py` CLI in full; this file only orients an agent among the files.

- `sr.py` — the scheduler: builds the daily plan, logs solves, computes review dates. Entry point for all commands (`today`, `log`, `undo`, `stats`, `due`, `status`, `resume`, and `handle`, which `log-solve.yml` runs for each comment), and home of `parse_comment`, the comment parser. Tests: `python3 spaced_repetition/test_sr.py`.
- `problems.json` — the fixed 165-problem set (Grind 75 ∪ Grind 169 ∪ Blind 75, deduplicated, bit-manipulation excluded): id, title, slug, difficulty, tags, category, URLs.
- `state.json` — generated, not hand-edited except via the `config` block (tuning knobs) — everything else (review cards, solve log, served plans) is written by `sr.py`. Two GitHub Actions workflows (`.github/workflows/daily-leetcode.yml`, `log-solve.yml`) also commit changes to this file directly to `main`, so `git pull` before editing it locally and `git push` after, per the root CLAUDE.md.

Three things that have caught edits to `sr.py`:

- **Stored config wins.** `load_state` merges `{**DEFAULT_CONFIG, **state["config"]}` and `save_state` writes the merged result back, so the real `state.json` stores every key. A new key in `DEFAULT_CONFIG` reaches it on the next save, but changing an existing default does nothing to the live schedule: edit the `config` block in `state.json` for that. In a test, set config keys before the first save (before logging), or the stored value overrides them.
- **Which commands write.** `today` (it records `served` and deferrals, or the pause), `log`, `undo`, `resume` and `handle` save `state.json`; `status`, `stats` and `due` only read it. There is no `--state` flag, so experiment on a copy of this directory.
- **Log entries carry `prev`.** Each entry `log_solve` appends keeps the card and deferral as they were before it. `undo`, and a second log of a problem on the same day, restore from it; a log can't be worked backwards from the card. Entries logged before 11-25-2026 have none and can't be undone. `state["events"]` records the comments `handle` applied, so none counts twice.

Root CLAUDE.md has the trigger for when to run `sr.py log` from a Claude session ("solved 200 good" style requests) — don't relitigate that here.
