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


def today_str(args):
    return args.date or dt.date.today().isoformat()


def load_problems():
    with open(PROBLEMS_FILE) as f:
        return json.load(f)


DEFAULT_CONFIG = {
    # Daily effort budget: Easy costs 1, Medium/Hard cost 2.
    # Budget 2 => two easies per day, or one medium/hard.
    "new_per_day": 1,
    "daily_budget": 2,
    "cost": {"E": 1, "M": 2, "H": 2},
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


def pick_today(state, problems, date):
    if date in state["served"]:
        return state["served"][date], True

    cfg = state["config"]
    order = {p["id"]: i for i, p in enumerate(problems)}
    by_id = {p["id"]: p for p in problems}
    budget = cfg["daily_budget"]

    # Reviews first, in due order. First-fit: an item too big for the
    # remaining budget is skipped, but a later cheaper one may still fit
    # (e.g. budget 1 left -> skip a Medium, take an Easy due later).
    due = [
        pid for pid, c in state["cards"].items()
        if c["due"] <= date and pid in by_id
    ]
    due.sort(key=lambda pid: (state["cards"][pid]["due"], order[pid]))
    reviews = []
    for pid in due:
        c = cost_of(by_id[pid], cfg)
        if c <= budget:
            reviews.append(pid)
            budget -= c
        if budget <= 0:
            break

    # Then at most new_per_day new problems, strictly in curated order
    # (no skipping ahead), only if the next one fits the remaining budget.
    # An unlogged new problem from a previous day is re-served before a
    # fresh one is introduced.
    new = []
    if budget > 0 and cfg["new_per_day"] > 0:
        carried = [
            pid for pid in unlogged_served(state, date)
            if pid not in state["cards"] and pid not in reviews
        ]
        seen = set(state["cards"]) | {
            pid for plan in state["served"].values() for pid in plan["new"]
        }
        fresh = [p["id"] for p in problems if p["id"] not in seen]
        for pid in (carried + fresh):
            if len(new) >= cfg["new_per_day"]:
                break
            c = cost_of(by_id[pid], cfg)
            if c <= budget:
                new.append(pid)
                budget -= c
            else:
                break  # next-in-order doesn't fit; don't skip ahead

    plan = {"review": reviews, "new": new}
    state["served"][date] = plan
    save_state(state)
    return plan, False


def fmt_problem(p, card=None):
    diff = {"E": "Easy", "M": "Medium", "H": "Hard"}[p["difficulty"]]
    url = p["alt_url"] if p.get("paid") else p["url"]
    tags = ", ".join(p["tags"][:3])
    extra = ""
    if card:
        extra = f"  [reps {card['reps']}, last {card['last']}]"
    paid_note = " (premium — free mirror linked)" if p.get("paid") else ""
    return f"#{p['id']} {p['title']} ({diff}) — {tags}{paid_note}\n    {url}{extra}"


def md_problem(p, card=None):
    diff = {"E": "Easy", "M": "Medium", "H": "Hard"}[p["difficulty"]]
    url = p["alt_url"] if p.get("paid") else p["url"]
    line = f"- [ ] **#{p['id']} [{p['title']}]({url})** ({diff}) — {', '.join(p['tags'][:3])}"
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
    backlog = [
        pid for pid in unlogged_served(state, date)
        if pid not in plan["review"] and pid not in plan["new"]
    ]
    if backlog:
        print("\n**Unlogged from previous days:** "
              + ", ".join(f"#{pid} {by_id[pid]['title']}" for pid in backlog))
    print("\n_Log solves with `python3 spaced_repetition/sr.py log <id> "
          "<again|hard|good|easy>` — or just tell Claude._")


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
    pid, grade = args.id, args.grade
    if pid not in by_id:
        sys.exit(f"Unknown problem id {pid}")
    if grade not in GRADES:
        sys.exit(f"Grade must be one of {GRADES}")
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
    p.add_argument("grade", choices=GRADES)
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
