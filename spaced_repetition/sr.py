#!/usr/bin/env python3
"""LeetCode spaced-repetition scheduler (SM-2 variant).

State lives in state.json next to this file. Problem set in problems.json.

Commands:
    python3 sr.py today [--date YYYY-MM-DD]   print today's plan (idempotent per day)
    python3 sr.py log ID GRADE [--date D]     record a solve; GRADE in {again,hard,good,easy}
    python3 sr.py stats                       progress overview
    python3 sr.py due                         list all cards with due dates
    python3 sr.py status [--date D]           "active" or "paused YYYY-MM-DD"
    python3 sr.py resume [--date D]           end a pause (logging a solve does too)
    python3 sr.py undo ID                     take back the latest log of a problem
    python3 sr.py handle --body-file F ...    apply a GitHub comment (log-solve.yml)
    python3 sr.py sweep ...                   log comments whose run never finished

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
import re
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


# --- reading log comments ---
#
# A comment is a list of items, separated by new lines, commas, semicolons,
# "\u00b7" or "and". An item is one problem and one grade, "53 good", optionally
# written "#53" and optionally followed by a note ("54 hard (needed a hint)").
# Two commands share the format: "undo 53" and "resume". Anything else that
# names a problem is reported back, never guessed: a wrong grade silently
# moves a review date.

_WORDS = "|".join(sorted(GRADE_ALIASES, key=len, reverse=True))
_ID = r"(?:\#|lc[\s-]?)?(?P<id>[1-9]\d*)"
_ITEM = re.compile(rf"""
    (?:[-+]\s+)?(?:\[[ x]\]\s*)?         # a list bullet, a checkbox
    (?:(?:solved|did|redid)\s+)?
    {_ID}
    \s*[^\w\s]*\s*                        # ":", "-", "=>", an emoji
    (?:(?:was|felt)\s+)?
    \(?(?P<grade>{_WORDS})\)?
    [^\w\s?]*(?:\s+[^\w\s?]+)*            # "!", ".", a trailing emoji
    (?P<note>\s*(?:[(\-\u2013\u2014:]|(?:but|though|after|with|because)\b).*)?
    """, re.I | re.X)
_UNDO = re.compile(rf"undo\s+{_ID}\W*", re.I)
_RESUME = re.compile(r"resume\W*", re.I)
_GRADE_WORD = re.compile(rf"\b(?:{_WORDS})\b", re.I)
_NUMBER = re.compile(r"(?<![\w./=])(#)?([1-9]\d*)(?![\w.])")
# A comment that matched nothing gets the format explained only if it looks
# like it meant to log something.
_ATTEMPT = re.compile(r"\b(?:again|hard|good|easy|solved|logged|undo)\b|#\d", re.I)


def parse_comment(text, deck=None):
    """Read a log comment.

    Returns a dict: "solves" [(id, grade)] to log, "undo" [id] whose latest
    log to take back, "resume" (bool), and "skipped" [(id, why)] for problems
    the comment names but that are not logged. With deck, a set of problem
    ids, a number outside the set is reported instead of logged.
    """
    text = re.sub(r"^\s*>.*$", "", text, flags=re.M)   # a quoted reply
    text = re.sub(r"https?://\S+|[*`]", " ", text)      # links, bold, code marks
    graded = bool(_GRADE_WORD.search(text))
    found, undo, skipped, resume = [], [], [], False
    for item in re.split(r"[\n,;\u00b7]|\band\b", text, flags=re.I):
        item = item.strip()
        m = _ITEM.fullmatch(item)
        if not item:
            continue
        if _RESUME.fullmatch(item):
            resume = True
        elif _UNDO.fullmatch(item):
            undo.append(_UNDO.fullmatch(item)["id"])
        elif m and _GRADE_WORD.search(m["note"] or ""):
            skipped.append((m["id"], "more than one grade"))
        elif m and deck is not None and m["id"] not in deck:
            skipped.append((m["id"], "no such problem in the set"))
        elif m:
            found.append((m["id"], normalize_grade(m["grade"])))
        else:
            # Only report numbers that could be problems, in a comment that
            # grades something or marks a number with "#".
            for mark, n in _NUMBER.findall(item):
                if mark or (graded and (deck is None or n in deck)):
                    skipped.append((n, "couldn't read that line"))
    grades = {}
    for pid, grade in found:
        grades.setdefault(pid, set()).add(grade)
    solves = []
    for pid, grade in found:
        if len(grades[pid]) > 1:
            skipped.append((pid, "more than one grade"))
        elif (pid, grade) not in solves:
            solves.append((pid, grade))
    logged = {pid for pid, _ in solves} | set(undo)
    reported = []
    for pid, why in skipped:
        if pid not in logged and pid not in {p for p, _ in reported}:
            reported.append((pid, why))
    return {"solves": solves, "undo": undo, "resume": resume, "skipped": reported}


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
    # A problem served this many times (not necessarily on consecutive days)
    # without being logged is assumed skipped: it gets pushed back by
    # defer_days so the plan keeps rotating instead of showing the same item
    # forever. Repeat deferrals double the window (capped at defer_max_days)
    # so a problem you keep ignoring fades out instead of reappearing every
    # few days.
    "carry_limit": 2,
    "defer_days": 4,
    "defer_max_days": 24,
    # When reviews pile up, add extra review capacity on top of the normal
    # budget (1 extra slot per backlog_step overdue cards, up to
    # backlog_bonus_max). New problems are never blocked - starving them
    # just makes every day look the same, which is the opposite of the goal.
    "backlog_step": 5,
    "backlog_bonus_max": 2,
    # This many days without a log pauses everything: no plans and no
    # deferrals until a solve is logged or `resume` is run. 0 disables.
    "pause_after": 7,
    # At most this many Hard problems in one day's plan, reviews and new
    # problems together. Reviews claim them first. A Hard review over the cap
    # waits for tomorrow; a Hard new problem gives way to the next candidate
    # (in curated order, to none, as when it does not fit the budget).
    "max_hard": 2,
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


def add_days(date, n):
    return (dt.date.fromisoformat(date) + dt.timedelta(days=n)).isoformat()


def last_touch(state, pid):
    """Most recent date this problem was logged, deferred or reset by a resume."""
    dates = [e["date"] for e in state["log"] if e["id"] == pid]
    d = state.get("deferred", {}).get(pid, {}).get("at")
    if d:
        dates.append(d)
    start = state.get("comeback", {}).get("start")
    if start:
        # The day before, so a showing on the day of the resume itself counts.
        dates.append(add_days(start, -1))
    return max(dates) if dates else ""


def paused_since(state, date):
    """First day plans stopped, or None while active.

    Counts from the last log or resume; someone who has never logged counts
    from their first plan.
    """
    after = state["config"].get("pause_after", 7)
    starts = [e["date"] for e in state["log"]]
    if state.get("comeback"):
        starts.append(state["comeback"]["start"])
    starts = starts or sorted(state["served"])[:1]
    if after <= 0 or not starts:
        return None
    since = add_days(max(starts), after + 1)
    return since if since <= date else None


def resume(state, date):
    """End a pause, undoing what being away did to the schedule.

    Clearing `deferred` alone would backfire: every showing that was ignored
    would count again and push the old cards to the back of the review queue.
    The comeback date makes last_touch treat those showings as settled.
    """
    state["comeback"] = {"start": date}
    state["deferred"] = {}
    state.pop("paused", None)
    for card in state["cards"].values():
        card["due"] = add_days(card["last"], card["interval"])  # undo deferrals


def apply_deferrals(state, date):
    """Push back anything served carry_limit times without a log.

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

    cap = cfg.get("defer_max_days", 24)
    deferred_now = []
    for pid, n in sorted(counts.items()):
        if n >= limit:
            prev = state["deferred"].get(pid, {})
            times = prev.get("count", 0) + 1
            # escalate: 4d, 8d, 16d, ... capped
            window = min(days * (2 ** (times - 1)), cap)
            until = add_days(date, window)
            state["deferred"][pid] = {"at": date, "until": until,
                                      "count": times}
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

    # Reviews first, in due order, but normally only up to review_budget so
    # new problems can't be starved. First-fit: an item too big for the
    # remaining budget is skipped, but a later cheaper one may still fit
    # (e.g. budget 1 left -> skip a Medium, take an Easy due later).
    due = [
        pid for pid, c in state["cards"].items()
        if c["due"] <= date and pid in by_id and not is_deferred(state, pid, date)
    ]

    # A card already shown recently without being graded goes to the BACK of
    # the queue. Without this it keeps winning on due-date order and is
    # re-served day after day, starving cards that are genuinely waiting.
    shown_recently = set()
    for d, plan in state["served"].items():
        if d < date:
            for pid in plan["review"]:
                if d > last_touch(state, pid):
                    shown_recently.add(pid)
    due.sort(key=lambda pid: (pid in shown_recently,
                              state["cards"][pid]["due"], order[pid]))

    # Overdue reviews earn extra capacity, but never at the cost of new
    # problems - variety is what keeps the queue from looking identical.
    bonus = min(len(due) // cfg.get("backlog_step", 5),
                cfg.get("backlog_bonus_max", 2))
    review_budget = cfg["daily_budget"] - cfg.get("new_budget", 0) + bonus
    hard_left = cfg.get("max_hard", 2)

    def fits(pid, budget):
        hard = by_id[pid]["difficulty"] == "H"
        return cost_of(by_id[pid], cfg) <= budget and (hard_left > 0 or not hard)

    def take(pid):
        nonlocal hard_left
        hard_left -= by_id[pid]["difficulty"] == "H"
        return cost_of(by_id[pid], cfg)

    reviews = []
    for pid in due:
        if fits(pid, review_budget):
            reviews.append(pid)
            review_budget -= take(pid)
        if review_budget <= 0:
            break

    # New problems get the reserved slice plus whatever reviews left over.
    budget = cfg.get("new_budget", 0) + max(0, review_budget)

    # Then at most new_per_day new problems, if they fit the remaining
    # budget. An unlogged new problem from a previous day is re-served once
    # before a fresh one is introduced. Never-seen problems come before ones
    # that were deferred. Within each group the order depends on config:
    #   category_rotation - pick from the category whose last introduction
    #     is oldest (never-served categories first), curated order within
    #     a category; guarantees every topic gets fresh coverage regularly
    #   curated - strict problems.json order
    new = []
    if budget > 0 and cfg["new_per_day"] > 0:
        # Only problems shown since they were last deferred carry over. Once
        # a deferral expires the problem is `revived` below, not carried.
        shown = {pid for d, plan in state["served"].items() if d < date
                 for pid in plan["new"] if d > last_touch(state, pid)}
        carried = [
            pid for pid in unlogged_served(state, date)
            if pid in shown and pid not in state["cards"] and pid not in reviews
            and not is_deferred(state, pid, date)
        ]
        seen = set(state["cards"]) | {
            pid for plan in state["served"].values() for pid in plan["new"]
        }
        fresh = [p["id"] for p in problems if p["id"] not in seen]
        # A problem shown and ignored re-enters the pool once any defer window
        # passes, behind every never-seen one, so problems that keep getting
        # ignored can't crowd out the rest of the deck.
        revived = [p["id"] for p in problems
                   if p["id"] in seen and p["id"] not in state["cards"]
                   and not is_deferred(state, p["id"], date)
                   and p["id"] not in carried]
        fresh = fresh + revived
        if cfg.get("new_order") == "category_rotation":
            # last date each category had a new problem introduced
            last_intro = {}
            for d, plan in sorted(state["served"].items()):
                for pid in plan["new"]:
                    if pid in by_id:
                        last_intro[by_id[pid].get("category", "?")] = d
            fresh.sort(key=lambda pid: (
                pid in revived,
                last_intro.get(by_id[pid].get("category", "?"), ""),
                order[pid],
            ))
        candidates = carried + fresh
        while candidates and len(new) < cfg["new_per_day"] and budget > 0:
            pid = candidates.pop(0)
            if fits(pid, budget):
                new.append(pid)
                budget -= take(pid)
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


def today_md(state, problems, date, plan):
    """The plan as Markdown, for the daily issue and the reply to `resume`."""
    by_id = {p["id"]: p for p in problems}
    out = [f"### Practice plan — {date}"]
    if not plan["review"] and not plan["new"]:
        out.append("\nNothing due and no new problems left. Done!")
    if plan["review"]:
        out.append(f"\n**Reviews due ({len(plan['review'])}):**\n")
        out += [md_problem(by_id[pid], state["cards"].get(pid)) for pid in plan["review"]]
    if plan["new"]:
        out.append(f"\n**New ({len(plan['new'])}):**\n")
        out += [md_problem(by_id[pid]) for pid in plan["new"]]
    # Logging a problem clears its deferral, so on a reprint it may be gone.
    pushed = [pid for pid in plan.get("deferred", [])
              if pid in by_id and pid in state.get("deferred", {})]
    if pushed:
        out.append("\n_Pushed back (served repeatedly without being logged): "
                   + ", ".join(f"#{pid} {by_id[pid]['title']} (not before {state['deferred'][pid]['until']})"
                               for pid in pushed) + "._")
    out.append("\n_Log solves by commenting here, one per line: `1 good`, `20 hard` "
               "(grades: again / hard / good / easy). Or use the "
               "[Daily Practice page](https://lambdanotfound.github.io/practice/)._")
    return "\n".join(out)


def cmd_today_md(state, problems, date, plan):
    print(today_md(state, problems, date, plan))


def cmd_paused(state, date, since):
    # The month makes a pause change state.json, and so get committed, once a
    # month: GitHub disables scheduled workflows after 60 days without one.
    note = {"since": since, "month": date[:7]}
    if state.get("paused") != note:
        state["paused"] = note
        save_state(state)
    days = state["config"]["pause_after"]
    print(f"### Practice plans paused since {since}\n\n"
          f"Nothing was logged for {days} days, so plans stopped and nothing is "
          "piling up.\n\n"
          "To come back, comment `resume` here for a plan with your old problems "
          "first, or just log any solve, e.g. `1 good`.")


def cmd_today(args):
    state = load_state()
    problems = load_problems()
    date = today_str(args)
    by_id = {p["id"]: p for p in problems}
    since = paused_since(state, date)
    if since and date not in state["served"]:
        if not args.resume:
            cmd_paused(state, date, since)
            return
        resume(state, date)
    plan, repeated = pick_today(state, problems, date)
    if not repeated:
        save_state(state)

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


def latest_log(state, pid):
    for entry in reversed(state["log"]):
        if entry["id"] == pid:
            return entry
    return None


def restore(state, entry):
    """Take back one log entry: its card and deferral go back as they were."""
    pid, prev = entry["id"], entry["prev"]
    if prev["card"] is None:
        state["cards"].pop(pid, None)
    else:
        state["cards"][pid] = dict(prev["card"])
    if prev["deferred"] is None:
        state.get("deferred", {}).pop(pid, None)
    else:
        state.setdefault("deferred", {})[pid] = dict(prev["deferred"])
    state["log"] = [e for e in state["log"] if e is not entry]


def log_solve(state, by_id, pid, grade, date):
    """Record one solve and reschedule its card. Returns the lines to print.

    Each entry keeps the card and deferral as they were before it ("prev"),
    which is what makes undo exact: a log can't be worked backwards, since
    "again" sets any interval to 1 and ease stops at EASE_MIN. A second log
    of the same problem on the same day replaces the first, so a redo or a
    corrected grade doesn't grow the interval twice.
    """
    out = []
    if paused_since(state, date):
        resume(state, date)
        out.append("Welcome back - plans resume today, with your old problems first.")
    replaced = None
    last = latest_log(state, pid)
    if last and last["date"] == date and "prev" in last:
        replaced = last["grade"]
        restore(state, last)
    old_card = state["cards"].get(pid)
    old_deferral = state.get("deferred", {}).get(pid)
    prev = {"card": dict(old_card) if old_card else None,
            "deferred": dict(old_deferral) if old_deferral else None}

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
        iv = max(iv + 1, round(iv * 1.2))
        ease = max(EASE_MIN, ease - 0.15)
    elif grade == "good":
        iv = max(iv + 1, round(iv * ease))
    else:  # easy
        iv = max(iv + 2, round(iv * ease * 1.3))
        ease += 0.15

    iv = min(iv, MAX_INTERVAL)
    card.update({
        "interval": iv, "ease": round(ease, 2),
        "due": add_days(date, iv), "reps": card["reps"] + 1, "last": date,
    })
    state["cards"][pid] = card
    # A logged problem is no longer being ignored: drop the deferral, which
    # would otherwise hide the review just scheduled, and its escalation count.
    state.get("deferred", {}).pop(pid, None)
    state["log"].append({"date": date, "id": pid, "grade": grade, "prev": prev})
    note = f" (replaces {replaced} from earlier today)" if replaced else ""
    out.append(f"Logged #{pid} {by_id[pid]['title']}: {grade}{note}. "
               f"Next review {card['due']} (interval {iv}d, ease {card['ease']}).")
    return out


def undo_solve(state, by_id, pid):
    """Take back the latest log of pid. Returns (message, changed)."""
    name = f"#{pid} {by_id[pid]['title']}" if pid in by_id else f"#{pid}"
    entry = latest_log(state, pid)
    if entry is None:
        return f"Nothing to undo: {name} has never been logged.", False
    if "prev" not in entry:
        return (f"Can't undo {name}: its latest log, on {entry['date']}, "
                "was made before undo existed."), False
    restore(state, entry)
    card = state["cards"].get(pid)
    after = f"next review back to {card['due']}" if card else "it counts as a new problem again"
    return f"Undid {name}: {entry['grade']} on {entry['date']}; {after}.", True


def cmd_log(args):
    state = load_state()
    by_id = {p["id"]: p for p in load_problems()}
    pid, grade = args.id, normalize_grade(args.grade)
    if pid not in by_id:
        sys.exit(f"Unknown problem id {pid}")
    if grade is None:
        sys.exit(f"Unknown grade {args.grade!r}; use one of "
                 f"{sorted(set(GRADE_ALIASES))}")
    lines = log_solve(state, by_id, pid, grade, today_str(args))
    save_state(state)
    print("\n".join(lines))


def cmd_undo(args):
    state = load_state()
    msg, changed = undo_solve(state, {p["id"]: p for p in load_problems()}, args.id)
    if changed:
        save_state(state)
    print(msg)


# --- comments from GitHub ---

HELP = (
    "Nothing logged: I couldn't find a line like `53 good` in that comment.\n\n"
    "Write one problem per line, or several separated by commas: `53 good, 1 hard`.\n"
    "Grades: `again` (needed the solution) · `hard` (slow or messy) · `good` · `easy`.\n"
    "`undo 53` takes back the latest log of #53, and `resume` ends a pause.")
EVENTS_KEPT_DAYS = 30


def now_utc():
    return dt.datetime.now(dt.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def handle_comment(state, problems, text, date, event=None, late=False):
    """Apply one log comment or log issue to state.

    Returns (reply lines, done): done is True when something was logged,
    undone or resumed. With event (e.g. "comment-123"), a comment already
    applied is not applied twice, so a retried push or a re-run is safe.
    late is the missed-log check replaying a comment whose own run never
    finished: it skips undo, which by then could hit a newer log, and any
    problem that was logged again after the comment.
    """
    by_id = {p["id"]: p for p in problems}
    if event and event in state.get("events", {}).get("done", {}):
        return ([] if late else ["Already logged from this comment; nothing changed."]), False
    parsed = parse_comment(text, set(by_id))
    lines, done = [], False

    if parsed["resume"]:
        if paused_since(state, date):
            resume(state, date)
            done = True
            lines.append("Welcome back - plans resume today, with your old problems first.")
        elif not late:
            lines.append("Plans are not paused - nothing to resume.")
        if not late:
            plan, _ = pick_today(state, problems, date)
            lines += ["", today_md(state, problems, date, plan)]
    for pid in [] if late else parsed["undo"]:
        msg, changed = undo_solve(state, by_id, pid)
        lines.append("- " + msg)
        done = done or changed
    for pid, grade in parsed["solves"]:
        if late and any(e["id"] == pid and e["date"] >= date for e in state["log"]):
            continue  # logged again since this comment; the newer log stands
        lines += ["- " + line for line in log_solve(state, by_id, pid, grade, date)]
        done = True
    if not late and parsed["skipped"]:
        lines += [f"- NOT LOGGED: #{pid} ({why})." for pid, why in parsed["skipped"]]
        lines.append("\nWrite one problem per line, like `53 good`.")
    if not lines and not late and _ATTEMPT.search(text):
        lines.append(HELP)

    if done and event:
        events = state.setdefault("events", {"since": now_utc(), "done": {}})
        events["done"][event] = date
        cutoff = add_days(date, -EVENTS_KEPT_DAYS)
        events["done"] = {k: v for k, v in events["done"].items() if v >= cutoff}
    return lines, done


def cmd_handle(args):
    state, problems = load_state(), load_problems()
    before = json.dumps(state, sort_keys=True)
    text = pathlib.Path(args.body_file).read_text() if args.body_file else sys.stdin.read()
    lines, done = handle_comment(state, problems, text, today_str(args), args.event)
    if json.dumps(state, sort_keys=True) != before:
        save_state(state)
    reply = "\n".join(lines).strip()
    if args.reply_file:
        pathlib.Path(args.reply_file).write_text(reply + "\n" if reply else "")
    else:
        print(reply)
    if args.status_file:
        pathlib.Path(args.status_file).write_text(f"done={'yes' if done else 'no'}\n")


def is_log_issue(title):
    return title.lower().startswith("practice log")


def takes_logs(title):
    """Issues whose comments log-solve.yml reads: plan issues and log issues."""
    t = title.lower()
    return is_log_issue(title) or "daily practice" in t or "daily leetcode" in t


def read_jsonl(path):
    try:
        with open(path) as f:
            return [json.loads(line) for line in f if line.strip()]
    except FileNotFoundError:
        return []


def pacific_date(timestamp):
    """The Pacific date of a GitHub timestamp such as 2026-11-26T01:30:00Z."""
    from zoneinfo import ZoneInfo
    t = dt.datetime.fromisoformat(timestamp.replace("Z", "+00:00"))
    return t.astimezone(ZoneInfo("America/Los_Angeles")).date().isoformat()


def missed_events(issues, comments, owner, since, before):
    """Log issues and log comments from owner, created in [since, before).

    issues and comments are GitHub API records reduced to the fields the
    daily workflow keeps: issues {number, id, title, user, body, created_at},
    comments {id, issue, user, body, created_at}.
    """
    titles = {i["number"]: i["title"] for i in issues}
    found = [{"event": f"issue-{i['id']}", "issue": i["number"], "title": i["title"],
              "body": i.get("body") or "", "created_at": i["created_at"]}
             for i in issues if i["user"] == owner and is_log_issue(i["title"])]
    found += [{"event": f"comment-{c['id']}", "issue": c["issue"], "title": titles[c["issue"]],
               "body": c.get("body") or "", "created_at": c["created_at"]}
              for c in comments
              if c["user"] == owner and c["issue"] in titles and takes_logs(titles[c["issue"]])]
    return sorted((e for e in found if since <= e["created_at"] < before),
                  key=lambda e: e["created_at"])


def sweep(state, problems, issues, comments, owner, now):
    """Log what log-solve.yml missed. Returns the replies to post.

    A run can fail to finish, or GitHub may never start it. Every comment
    log-solve applies is recorded in state["events"], so one from the last
    few days that isn't there was missed. Only comments made after the record
    began count, and none from the last hour, whose own run may still be
    going.
    """
    since = state.get("events", {}).get("since")
    if not since:
        return []
    before = (dt.datetime.fromisoformat(now.replace("Z", "+00:00"))
              - dt.timedelta(hours=1)).strftime("%Y-%m-%dT%H:%M:%SZ")
    replies = []
    for e in missed_events(issues, comments, owner, since, before):
        lines, done = handle_comment(state, problems, e["body"], pacific_date(e["created_at"]),
                                     e["event"], late=True)
        if done:
            replies.append({
                "issue": e["issue"], "close": is_log_issue(e["title"]),
                "reply": "Logged late by the daily check, because the run for "
                         "this comment didn't finish:\n\n" + "\n".join(lines)})
    return replies


def cmd_sweep(args):
    state, problems = load_state(), load_problems()
    before = json.dumps(state, sort_keys=True)
    replies = sweep(state, problems, read_jsonl(args.issues), read_jsonl(args.comments),
                    args.owner, args.now or now_utc())
    if json.dumps(state, sort_keys=True) != before:
        save_state(state)
    with open(args.replies, "w") as f:
        json.dump(replies, f)
    print(f"Missed-log check: {len(replies)} comment(s) logged late.")


def cmd_status(args):
    state, date = load_state(), today_str(args)
    # A day that already has a plan is active, as it is for `today`.
    since = None if date in state["served"] else paused_since(state, date)
    print(f"paused {since}" if since else "active")


def cmd_resume(args):
    state = load_state()
    date = today_str(args)
    if not paused_since(state, date):
        print("Plans are not paused - nothing to resume.")
        return
    resume(state, date)
    save_state(state)
    print("Welcome back - plans resume today, with your old problems first.")


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
    p.add_argument("--resume", action="store_true",
                   help="end a pause first, if there is one")
    p.set_defaults(fn=cmd_today)

    p = sub.add_parser("status")
    p.add_argument("--date")
    p.set_defaults(fn=cmd_status)

    p = sub.add_parser("resume")
    p.add_argument("--date")
    p.set_defaults(fn=cmd_resume)

    p = sub.add_parser("log")
    p.add_argument("id")
    p.add_argument("grade", choices=sorted(set(GRADE_ALIASES)),
                   metavar="again|hard|good|easy")
    p.add_argument("--date")
    p.set_defaults(fn=cmd_log)

    p = sub.add_parser("undo", help="take back the latest log of a problem")
    p.add_argument("id")
    p.set_defaults(fn=cmd_undo)

    p = sub.add_parser("handle", help="apply a GitHub comment (log-solve.yml)")
    p.add_argument("--body-file", help="the comment text (default: stdin)")
    p.add_argument("--event", help="comment-ID or issue-ID, so it is applied once")
    p.add_argument("--date")
    p.add_argument("--reply-file", help="write the reply here instead of stdout")
    p.add_argument("--status-file", help="write done=yes|no here")
    p.set_defaults(fn=cmd_handle)

    p = sub.add_parser("sweep", help="log comments whose run never finished")
    p.add_argument("--issues", required=True, help="JSON lines of recent issues")
    p.add_argument("--comments", required=True, help="JSON lines of recent comments")
    p.add_argument("--owner", required=True)
    p.add_argument("--replies", required=True, help="write the replies to post here")
    p.add_argument("--now", help="UTC time, YYYY-MM-DDTHH:MM:SSZ (default: now)")
    p.set_defaults(fn=cmd_sweep)

    p = sub.add_parser("stats")
    p.set_defaults(fn=cmd_stats)

    p = sub.add_parser("due")
    p.set_defaults(fn=cmd_due)

    args = ap.parse_args()
    args.fn(args)


if __name__ == "__main__":
    main()
