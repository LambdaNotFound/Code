#!/usr/bin/env python3
"""LeetCode spaced-repetition scheduler (SM-2 variant).

State lives in state.json next to this file. Problem set in problems.json.

Commands:
    python3 sr.py today [--date YYYY-MM-DD]   print today's plan (idempotent per day)
    python3 sr.py log ID GRADE [--date D]     record a solve; GRADE in {again,hard,good,easy}
    python3 sr.py stats                       progress overview
    python3 sr.py due                         list all cards with due dates

Grading guide:
    again = couldn't solve / needed the solution
    hard  = solved, but slow or with major hiccups
    good  = solved with minor friction
    easy  = solved quickly, clean
"""
import argparse
import datetime as dt
import json
import pathlib
import sys

DIR = pathlib.Path(__file__).resolve().parent
PROBLEMS_FILE = DIR / "problems.json"
STATE_FILE = DIR / "state.json"

EASE_START = 2.5
EASE_MIN = 1.3
MAX_INTERVAL = 90
GRADES = ("again", "hard", "good", "easy")

# Accepted synonyms -> canonical grade. The grade records how the ATTEMPT
# went, not the problem's LeetCode difficulty (which is static and already
# stored in problems.json). "medium" is accepted as a synonym for "good"
# because it reads naturally, but note it means "moderate friction", not
# "this is a Medium problem".
GRADE_ALIASES = {
    "again": "again", "fail": "again", "failed": "again", "stuck": "again",
    "blanked": "again", "no": "again",
    "hard": "hard", "slow": "hard", "rough": "hard", "tough": "hard",
    "good": "good", "medium": "good", "ok": "good", "okay": "good",
    "fine": "good", "yes": "good",
    "easy": "easy", "trivial": "easy", "quick": "easy", "clean": "easy",
}


def normalize_grade(word):
    return GRADE_ALIASES.get(word.strip().lower())


def today_str(args):
    return args.date or dt.date.today().isoformat()


def load_problems():
    with open(PROBLEMS_FILE) as f:
        return json.load(f)


DEFAULT_CONFIG = {
    # Daily effort budget: Easy costs 1, Medium/Hard cost 2.
    # Budget 4 => ~2 problems/day given the deck's difficulty mix.
    "daily_budget": 4,
    # Reserved slice of daily_budget that reviews may NOT consume, so new
    # problems (and thus category rotation) can't be starved by a heavy
    # review day. Reviews get (daily_budget - new_budget); new problems
    # get new_budget plus whatever reviews left unused.
    "new_budget": 2,
    "new_per_day": 2,
    "cost": {"E": 1, "M": 2, "H": 2},
    # "category_rotation": introduce new problems from the category whose
    # last new-problem introduction is oldest, so every topic gets fresh
    # coverage every ~1-2 weeks. "curated": strict problems.json order.
    "new_order": "category_rotation",
    # A problem served this many days running without being logged is
    # assumed skipped: it gets pushed back by defer_days so the plan keeps
    # rotating instead of showing the same item forever.
    "carry_limit": 2,
    "defer_days": 4,
}


def load_state():
    if STATE_FILE.exists():
        with open(STATE_FILE) as f:
            state = json.load(f)
        state["config"] = {**DEFAULT_CONFIG, **state.get("config", {})}
        state["config"].pop("daily_cap", None)  # pre-budget config key
        return state
    return {
        "config": dict(DEFAULT_CONFIG),
        "cards": {},          # id -> {interval, ease, due, reps, lapses, last}
        "served": {},         # date -> {"review": [...], "new": [...]}
        "log": [],            # {date, id, grade}
    }


def save_state(state):
    with open(STATE_FILE, "w") as f:
        json.dump(state, f, indent=2, sort_keys=True)


def logged_ids(state):
    return {e["id"] for e in state["log"]}


def unlogged_served(state, before_date):
    """Items served on previous days and never logged since."""
    out = []
    seen = logged_ids(state)
    for date, plan in sorted(state["served"].items()):
        if date >= before_date:
            continue
        for pid in plan["review"] + plan["new"]:
            card = state["cards"].get(pid)
            # a review item is 'handled' if logged on/after that date
            handled = any(e["id"] == pid and e["date"] >= date for e in state["log"])
            if not handled and pid not in out:
                out.append(pid)
    return out


def cost_of(problem, cfg):
    return cfg["cost"].get(problem["difficulty"], 2)


def last_touch(state, pid):
    """Most recent date this problem was logged or deferred."""
    dates = [e["date"] for e in state["log"] if e["id"] == pid]
    d = state.get("deferred", {}).get(pid, {}).get("at")
    if d:
        dates.append(d)
    return max(dates) if dates else ""


def apply_deferrals(state, date):
    """Push back anything served carry_limit days running without a log.

    Returns the list of ids deferred on this call. Without this, an item
    the user never logs stays due forever and reappears every single day,
    crowding out the rotation.
    """
    cfg = state["config"]
    limit = cfg.get("carry_limit", 2)
    days = cfg.get("defer_days", 4)
    state.setdefault("deferred", {})
    if limit <= 0:
        return []

    counts = {}
    for d, plan in state["served"].items():
        if d >= date:
            continue
        for pid in plan["review"] + plan["new"]:
            if d > last_touch(state, pid):
                counts[pid] = counts.get(pid, 0) + 1

    until = (dt.date.fromisoformat(date) + dt.timedelta(days=days)).isoformat()
    deferred_now = []
    for pid, n in sorted(counts.items()):
        if n >= limit:
            state["deferred"][pid] = {"at": date, "until": until}
            if pid in state["cards"]:
                state["cards"][pid]["due"] = until
            deferred_now.append(pid)
    return deferred_now


def is_deferred(state, pid, date):
    d = state.get("deferred", {}).get(pid)
    return bool(d) and d["until"] > date


def pick_today(state, problems, date):
    if date in state["served"]:
        return state["served"][date], True

    deferred_now = apply_deferrals(state, date)
    cfg = state["config"]
    order = {p["id"]: i for i, p in enumerate(problems)}
    by_id = {p["id"]: p for p in problems}
    review_budget = cfg["daily_budget"] - cfg.get("new_budget", 0)

    # Reviews first, in due order, but only up to review_budget so new
    # problems can't be starved. First-fit: an item too big for the
    # remaining budget is skipped, but a later cheaper one may still fit
    # (e.g. budget 1 left -> skip a Medium, take an Easy due later).
    due = [
        pid for pid, c in state["cards"].items()
        if c["due"] <= date and pid in by_id and not is_deferred(state, pid, date)
    ]
    due.sort(key=lambda pid: (state["cards"][pid]["due"], order[pid]))
    reviews = []
    for pid in due:
        c = cost_of(by_id[pid], cfg)
        if c <= review_budget:
            reviews.append(pid)
            review_budget -= c
        if review_budget <= 0:
            break

    # New problems get the reserved slice plus whatever reviews left over.
    budget = cfg.get("new_budget", 0) + review_budget

    # Then at most new_per_day new problems, if they fit the remaining
    # budget. An unlogged new problem from a previous day is re-served
    # before a fresh one is introduced. Fresh order depends on config:
    #   category_rotation - pick from the category whose last introduction
    #     is oldest (never-served categories first), curated order within
    #     a category; guarantees every topic gets fresh coverage regularly
    #   curated - strict problems.json order
    new = []
    if budget > 0 and cfg["new_per_day"] > 0:
        carried = [
            pid for pid in unlogged_served(state, date)
            if pid not in state["cards"] and pid not in reviews
            and not is_deferred(state, pid, date)
        ]
        seen = set(state["cards"]) | {
            pid for plan in state["served"].values() for pid in plan["new"]
        }
        fresh = [p["id"] for p in problems if p["id"] not in seen]
        # A deferred new problem is not "fresh" yet, but once its defer
        # window passes it re-enters the pool ahead of never-seen ones.
        revived = [pid for pid in state.get("deferred", {})
                   if pid not in state["cards"] and pid in by_id
                   and not is_deferred(state, pid, date)
                   and pid not in carried]
        fresh = revived + fresh
        if cfg.get("new_order") == "category_rotation":
            # last date each category had a new problem introduced
            last_intro = {}
            for d, plan in sorted(state["served"].items()):
                for pid in plan["new"]:
                    if pid in by_id:
                        last_intro[by_id[pid].get("category", "?")] = d
            fresh.sort(key=lambda pid: (
                last_intro.get(by_id[pid].get("category", "?"), ""),
                order[pid],
            ))
        candidates = carried + fresh
        while candidates and len(new) < cfg["new_per_day"] and budget > 0:
            pid = candidates.pop(0)
            c = cost_of(by_id[pid], cfg)
            if c <= budget:
                new.append(pid)
                budget -= c
                if cfg.get("new_order") == "category_rotation":
                    # don't introduce two problems of the same category today
                    cat = by_id[pid].get("category")
                    candidates = [f for f in candidates
                                  if by_id[f].get("category") != cat]
            elif cfg.get("new_order") != "category_rotation":
                break  # curated mode: don't skip ahead

    plan = {"review": reviews, "new": new}
    if deferred_now:
        plan["deferred"] = deferred_now
    state["served"][date] = plan
    save_state(state)
    return plan, False


def fmt_problem(p, card=None):
    diff = {"E": "Easy", "M": "Medium", "H": "Hard"}[p["difficulty"]]
    url = p["alt_url"] if p.get("paid") else p["url"]
    tags = p.get("category", "") + " · " + ", ".join(p["tags"][:2])
    extra = ""
    if card:
        extra = f"  [reps {card['reps']}, last {card['last']}]"
    paid_note = " (premium — free mirror linked)" if p.get("paid") else ""
    return f"#{p['id']} {p['title']} ({diff}) — {tags}{paid_note}\n    {url}{extra}"


def md_problem(p, card=None):
    diff = {"E": "Easy", "M": "Medium", "H": "Hard"}[p["difficulty"]]
    url = p["alt_url"] if p.get("paid") else p["url"]
    line = (f"- [ ] **#{p['id']} [{p['title']}]({url})** ({diff}) — "
            f"{p.get('category', '')} · {', '.join(p['tags'][:2])}")
    if p.get("paid"):
        line += " · premium, free mirror linked"
    if card:
        line += f" · reps {card['reps']}, last {card['last']}"
    return line


def cmd_today_md(state, problems, date, plan):
    by_id = {p["id"]: p for p in problems}
    print(f"### LeetCode plan — {date}")
    if not plan["review"] and not plan["new"]:
        print("\nNothing due and no new problems left. Done!")
    if plan["review"]:
        print(f"\n**Reviews due ({len(plan['review'])}):**\n")
        for pid in plan["review"]:
            print(md_problem(by_id[pid], state["cards"].get(pid)))
    if plan["new"]:
        print(f"\n**New ({len(plan['new'])}):**\n")
        for pid in plan["new"]:
            print(md_problem(by_id[pid]))
    if plan.get("deferred"):
        until = state["deferred"][plan["deferred"][0]]["until"]
        print(f"\n_Pushed to {until} (served repeatedly without being logged): "
              + ", ".join(f"#{pid} {by_id[pid]['title']}"
                          for pid in plan["deferred"] if pid in by_id) + "._")
    print("\n_Log solves by commenting here — e.g. `solved 1 good`, "
          "`#20 hard` (grades: again / hard / good / easy)._")


def cmd_today(args):
    state = load_state()
    problems = load_problems()
    date = today_str(args)
    by_id = {p["id"]: p for p in problems}
    plan, repeated = pick_today(state, problems, date)

    if args.md:
        cmd_today_md(state, problems, date, plan)
        return

    print(f"=== LeetCode plan for {date} ===")
    if not plan["review"] and not plan["new"]:
        print("Nothing due and no new problems left. Done!")
    if plan["review"]:
        print(f"\nReviews due ({len(plan['review'])}):")
        for pid in plan["review"]:
            print("  " + fmt_problem(by_id[pid], state["cards"].get(pid)))
    if plan["new"]:
        print(f"\nNew ({len(plan['new'])}):")
        for pid in plan["new"]:
            print("  " + fmt_problem(by_id[pid]))

    backlog = [
        pid for pid in unlogged_served(state, date)
        if pid not in plan["review"] and pid not in plan["new"]
    ]
    if backlog:
        print(f"\nUnlogged from previous days ({len(backlog)}): "
              + ", ".join(f"#{pid} {by_id[pid]['title']}" for pid in backlog))
        print("Log them with: python3 sr.py log <id> <again|hard|good|easy>")
    if args.json:
        print("\n" + json.dumps({"date": date, **plan, "backlog": backlog}))


def cmd_log(args):
    state = load_state()
    problems = load_problems()
    by_id = {p["id"]: p for p in problems}
    pid, grade = args.id, normalize_grade(args.grade)
    if pid not in by_id:
        sys.exit(f"Unknown problem id {pid}")
    if grade is None:
        sys.exit(f"Unknown grade {args.grade!r}; use one of "
                 f"{sorted(set(GRADE_ALIASES))}")
    date = today_str(args)

    card = state["cards"].get(pid, {
        "interval": 0, "ease": EASE_START, "due": date,
        "reps": 0, "lapses": 0, "last": None,
    })
    iv, ease = card["interval"], card["ease"]

    if card["reps"] == 0:
        iv = {"again": 1, "hard": 1, "good": 2, "easy": 4}[grade]
        if grade == "again":
            card["lapses"] += 1
    elif grade == "again":
        iv = 1
        ease = max(EASE_MIN, ease - 0.20)
        card["lapses"] += 1
    elif grade == "hard":
        iv = max(2, round(iv * 1.2))
        ease = max(EASE_MIN, ease - 0.15)
    elif grade == "good":
        iv = max(iv + 1, round(iv * ease))
    else:  # easy
        iv = max(iv + 2, round(iv * ease * 1.3))
        ease += 0.15

    iv = min(iv, MAX_INTERVAL)
    d = dt.date.fromisoformat(date) + dt.timedelta(days=iv)
    card.update({
        "interval": iv, "ease": round(ease, 2),
        "due": d.isoformat(), "reps": card["reps"] + 1, "last": date,
    })
    state["cards"][pid] = card
    state["log"].append({"date": date, "id": pid, "grade": grade})
    save_state(state)
    print(f"Logged #{pid} {by_id[pid]['title']}: {grade}. "
          f"Next review {card['due']} (interval {iv}d, ease {card['ease']}).")


def cmd_stats(args):
    state = load_state()
    problems = load_problems()
    cards = state["cards"]
    total = len(problems)
    started = len(cards)
    mature = sum(1 for c in cards.values() if c["interval"] >= 21)
    young = started - mature
    today = dt.date.today().isoformat()
    due_now = sum(1 for c in cards.values() if c["due"] <= today)
    lapses = sum(c["lapses"] for c in cards.values())
    solves = len(state["log"])
    print(f"Problems: {total} total | {started} started | {total - started} unseen")
    print(f"Cards: {mature} mature (interval>=21d) | {young} young | {due_now} due now")
    print(f"Solves logged: {solves} | total lapses: {lapses}")
    if state["log"]:
        days = sorted({e['date'] for e in state['log']})
        print(f"Active days: {len(days)} (first {days[0]}, last {days[-1]})")


def cmd_due(args):
    state = load_state()
    problems = load_problems()
    by_id = {p["id"]: p for p in problems}
    rows = sorted(state["cards"].items(), key=lambda kv: kv[1]["due"])
    for pid, c in rows:
        print(f"{c['due']}  #{pid:>5} {by_id[pid]['title']}  "
              f"(interval {c['interval']}d, reps {c['reps']}, ease {c['ease']})")
    if not rows:
        print("No cards yet.")


def main():
    ap = argparse.ArgumentParser(description=__doc__,
                                 formatter_class=argparse.RawDescriptionHelpFormatter)
    sub = ap.add_subparsers(dest="cmd", required=True)

    p = sub.add_parser("today")
    p.add_argument("--date")
    p.add_argument("--json", action="store_true")
    p.add_argument("--md", action="store_true",
                   help="markdown output (for GitHub issue comments)")
    p.set_defaults(fn=cmd_today)

    p = sub.add_parser("log")
    p.add_argument("id")
    p.add_argument("grade", choices=sorted(set(GRADE_ALIASES)),
                   metavar="again|hard|good|easy")
    p.add_argument("--date")
    p.set_defaults(fn=cmd_log)

    p = sub.add_parser("stats")
    p.set_defaults(fn=cmd_stats)

    p = sub.add_parser("due")
    p.set_defaults(fn=cmd_due)

    args = ap.parse_args()
    args.fn(args)


if __name__ == "__main__":
    main()
