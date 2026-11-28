#!/usr/bin/env python3
"""Tests for sr.py. Run: python3 spaced_repetition/test_sr.py"""
import argparse
import collections
import contextlib
import datetime as dt
import io
import pathlib
import sys
import tempfile
import unittest

sys.dont_write_bytecode = True
sys.path.insert(0, str(pathlib.Path(__file__).resolve().parent))
import sr  # noqa: E402


DECK = {p["id"] for p in sr.load_problems()}


class ParseComment(unittest.TestCase):
    def check(self, text, solves, skipped=(), undo=(), resume=False):
        got = sr.parse_comment(text, DECK)
        self.assertEqual(got["solves"], list(solves), text)
        self.assertEqual([pid for pid, _ in got["skipped"]], list(skipped), text)
        self.assertEqual(got["undo"], list(undo), text)
        self.assertEqual(got["resume"], resume, text)

    def test_one_problem_per_line(self):
        self.check("20 good\n54 hard\n217 easy",
                   [("20", "good"), ("54", "hard"), ("217", "easy")])
        self.check("20 good\r\n54 hard", [("20", "good"), ("54", "hard")])
        # What the Daily Practice page puts in a log issue.
        self.check("53 good\n1 hard\n21 easy", [("53", "good"), ("1", "hard"), ("21", "easy")])

    def test_separators(self):
        self.check("20 good, 54 hard; 217 easy",
                   [("20", "good"), ("54", "hard"), ("217", "easy")])
        self.check("1 good · 20 hard", [("1", "good"), ("20", "hard")])
        self.check("252 good and 9 hard", [("252", "good"), ("9", "hard")])

    def test_written_forms(self):
        for text in ("53 good", "#53 good", "solved 53 good", "#53 was good", "53 felt good",
                     "LC53 good", "lc-53 good", "lc 53 good", "53 (good)", "**53** GOOD",
                     "53: good", "53 - good", "53 => good", "53 ✅ good", "- [x] 53 good",
                     "- 53 good", "53 good!", "53. good", "53 good 👍", "`53 good`"):
            self.check(text, [("53", "good")])

    def test_every_alias(self):
        for word, grade in sr.GRADE_ALIASES.items():
            self.check(f"#20 {word}", [("20", grade)])

    def test_a_note_after_the_grade(self):
        for text in ("54 hard (needed a hint)", "54 hard - off by one", "54 hard: messy loop",
                     "54 hard but got there", "54 hard (12 min)", "54 hard — off by one"):
            self.check(text, [("54", "hard")])

    def test_a_note_with_another_grade_is_reported(self):
        self.check("#33 slow but ok", [], ["33"])
        self.check("54 easy - quick one", [], ["54"])

    def test_lines_it_cannot_read_are_reported(self):
        for text, pid in (("53 was pretty good", "53"), ("#53 Maximum Subarray good", "53"),
                          ("53, good", "53"), ("good: 53", "53"), ("#9 not good", "9"),
                          ("#9 easy? not really", "9"), ("#20 quick review, hard", "20"),
                          ("#9 wasn't easy", "9"), ("53 honestly felt good", "53")):
            got = sr.parse_comment(text, DECK)
            self.assertEqual(got["solves"], [], text)
            self.assertIn(pid, [p for p, _ in got["skipped"]], text)

    def test_one_grade_per_number(self):
        self.check("1 and 20 good", [("20", "good")], ["1"])
        self.check("1, 20 good", [("20", "good")], ["1"])

    def test_two_grades_for_one_problem(self):
        self.check("20 good\n20 hard", [], ["20"])
        self.check("20 good, 20 good", [("20", "good")])

    def test_numbers_outside_the_set(self):
        self.check("2026 good", [], ["2026"])
        self.check("#9999 easy", [], ["9999"])
        self.check("2026 good year", [])

    def test_remarks_log_nothing(self):
        for text in ("did 3 easy ones today", "only 2 good ones", "3 hard problems left",
                     "Day 3 good streak", "top 1 good", "4 good days in a row",
                     "5 hard ones on the list", "1 more easy one tomorrow", "part 2 was easy",
                     "n=5 easy", "k = 3 was easy", "#1 priority is easy stuff",
                     "took me 2 tries, good", "ok", "looks fine", "no", "thanks!"):
            self.assertEqual(sr.parse_comment(text, DECK)["solves"], [], text)

    def test_remarks_without_a_grade_report_nothing(self):
        for text in ("started at 9 this morning", "nice plan today", "back on 12-01"):
            self.check(text, [])

    def test_quoted_lines_are_ignored(self):
        self.check("> 20 good\n> 1 easy\n\n9 easy", [("9", "easy")])

    def test_links_are_ignored(self):
        self.check("https://leetcode.com/problems/two-sum/ 1 good", [("1", "good")])
        self.check("1 good https://leetcode.com/problems/two-sum/", [("1", "good")])

    def test_commands(self):
        self.check("undo 54", [], undo=["54"])
        self.check("undo #54", [], undo=["54"])
        self.check("Resume.", [], resume=True)
        self.check("resume\n20 good", [("20", "good")], resume=True)
        self.check("undo 54\n54 good", [("54", "good")], undo=["54"])

    def test_a_pasted_plan_logs_nothing(self):
        state = {"cards": {}, "deferred": {}}
        plan = {"review": ["208"], "new": ["53"]}
        text = sr.today_md(state, sr.load_problems(), "2026-11-26", plan)
        self.assertEqual(sr.parse_comment(text, DECK)["solves"], [])


class StateTest(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        old = sr.STATE_FILE
        self.addCleanup(lambda: setattr(sr, "STATE_FILE", old))
        sr.STATE_FILE = pathlib.Path(self.tmp.name) / "state.json"

    def use_deck(self, size):
        deck, old = sr.load_problems()[:size], sr.load_problems
        self.addCleanup(lambda: setattr(sr, "load_problems", old))
        sr.load_problems = lambda: deck

    def serve_days(self, n):
        """Serve n plans without logging; return each day's new problems."""
        self.set_config(pause_after=0)
        start = dt.date(2026, 1, 1)
        for i in range(n):
            self.today((start + dt.timedelta(i)).isoformat())
        served = sr.load_state()["served"]
        return [served[d]["new"] for d in sorted(served)]

    def run_cmd(self, fn, **kw):
        out = io.StringIO()
        with contextlib.redirect_stdout(out):
            fn(argparse.Namespace(**kw))
        return out.getvalue()

    def today(self, date, resume=False):
        return self.run_cmd(sr.cmd_today, date=date, md=True, json=False, resume=resume)

    def status(self, date):
        return self.run_cmd(sr.cmd_status, date=date).strip()

    def resume(self, date):
        return self.run_cmd(sr.cmd_resume, date=date)

    def set_config(self, **kw):
        old = dict(sr.DEFAULT_CONFIG)
        self.addCleanup(lambda: (sr.DEFAULT_CONFIG.clear(), sr.DEFAULT_CONFIG.update(old)))
        sr.DEFAULT_CONFIG.update(kw)

    def log(self, pid, grade, date):
        return self.run_cmd(sr.cmd_log, id=pid, grade=grade, date=date)

    def defer_one(self):
        """Serve plans without logging until something is deferred."""
        for day in range(1, 4):
            self.today(f"2026-01-0{day}")
        state = sr.load_state()
        return sorted(state["deferred"])[0], "2026-01-03"


class Scheduling(StateTest):
    def test_logging_clears_the_deferral(self):
        pid, date = self.defer_one()
        self.assertTrue(sr.is_deferred(sr.load_state(), pid, "2026-01-04"))
        self.log(pid, "again", date)
        state = sr.load_state()
        self.assertEqual(state["cards"][pid]["due"], "2026-01-04")
        self.assertFalse(sr.is_deferred(state, pid, "2026-01-04"))
        self.assertIn(f"#{pid} ", self.today("2026-01-04"))

    def test_logging_resets_the_escalation(self):
        pid, date = self.defer_one()
        self.log(pid, "again", date)
        for day in range(4, 10):
            self.today(f"2026-01-{day:02d}")
        self.assertEqual(sr.load_state()["deferred"][pid]["count"], 1)

    def test_reprint_after_logging_a_deferred_problem(self):
        pid, date = self.defer_one()
        self.log(pid, "good", date)
        self.assertIn(f"plan — {date}", self.today(date))

    def test_each_deferral_prints_its_own_date(self):
        self.defer_one()
        state = sr.load_state()
        a, b = state["served"]["2026-01-03"]["deferred"][:2]
        state["deferred"][b]["until"] = "2026-02-01"
        sr.save_state(state)
        out = self.today("2026-01-03")
        self.assertIn(state["deferred"][a]["until"], out)
        self.assertIn("2026-02-01", out)

    def test_unseen_problems_come_before_ignored_ones(self):
        for order in ("category_rotation", "curated"):
            with self.subTest(order):
                self.set_config(new_order=order)
                sr.STATE_FILE.unlink(missing_ok=True)
                days = self.serve_days(60)
                shown = collections.Counter(pid for day in days for pid in day)
                # Offered once, re-offered once, then nothing until the unseen
                # problems run out - which takes longer than 60 days.
                self.assertLessEqual(max(shown.values()), 2)
                self.assertGreaterEqual(len(shown), 40)

    def test_ignored_problems_return_after_the_unseen_run_out(self):
        self.use_deck(6)
        days = self.serve_days(40)
        shown = collections.Counter(pid for day in days for pid in day)
        self.assertEqual(len(shown), 6)
        self.assertGreater(min(shown.values()), 2)

    def test_hard_grows_every_time_but_never_past_good(self):
        last, ease = 0, sr.EASE_START
        for i in range(60):
            # One log a day: a second log on the same day replaces the first.
            self.log("1", "hard", sr.add_days("2026-01-01", i))
            card = sr.load_state()["cards"]["1"]
            if last < sr.MAX_INTERVAL:
                self.assertGreater(card["interval"], last)
            good = max(last + 1, round(last * ease)) if i else 2
            self.assertLessEqual(card["interval"], min(good, sr.MAX_INTERVAL))
            last, ease = card["interval"], card["ease"]
        self.assertEqual(last, sr.MAX_INTERVAL)

    def test_plan_is_idempotent(self):
        first = self.today("2026-01-01")
        self.assertEqual(first, self.today("2026-01-01"))

    # --- at most two Hard problems a day ---

    def hard_in(self, plan):
        by_id = {p["id"]: p for p in sr.load_problems()}
        return [pid for pid in plan["review"] + plan["new"]
                if by_id[pid]["difficulty"] == "H"]

    def four_hard_reviews(self, **config):
        """Four Hard cards due together, budget for all; return their due day."""
        self.set_config(daily_budget=20, pause_after=0, **config)
        for pid in ("23", "42", "76", "84"):
            self.log(pid, "good", "2026-01-01")
        return sr.load_state()["cards"]["23"]["due"]

    def test_hard_reviews_over_the_cap_wait_a_day(self):
        due = self.four_hard_reviews()
        self.today(due)
        first = sr.load_state()["served"][due]
        self.assertEqual(len(self.hard_in(first)), 2)
        self.assertEqual(len(first["new"]), 2)  # filled with non-Hard ones
        nxt = sr.add_days(due, 1)
        self.today(nxt)
        second = sr.load_state()["served"][nxt]
        self.assertEqual(set(self.hard_in(first) + self.hard_in(second)),
                         {"23", "42", "76", "84"})

    def test_without_the_cap_all_four_come_at_once(self):
        due = self.four_hard_reviews(max_hard=99)
        self.today(due)
        self.assertEqual(len(self.hard_in(sr.load_state()["served"][due])), 4)

    def test_a_hard_new_problem_gives_way(self):
        self.set_config(max_hard=0)
        days = self.serve_days(30)
        self.assertTrue(all(days))  # a new problem every day all the same
        self.assertEqual(self.hard_in({"review": [], "new": sum(days, [])}), [])

    def test_hard_new_problems_appear_by_default(self):
        days = self.serve_days(30)
        self.assertTrue(self.hard_in({"review": [], "new": sum(days, [])}))

    # --- pausing and coming back ---

    def break_off(self):
        """Log two problems on Jan 1, keep serving, stop logging: paused Jan 9."""
        self.today("2026-01-01")
        plan = sr.load_state()["served"]["2026-01-01"]["new"]
        for pid in plan:
            self.log(pid, "good", "2026-01-01")
        for day in range(2, 9):
            self.today(f"2026-01-{day:02d}")
        return plan

    def test_pauses_after_a_week_without_a_log(self):
        self.break_off()
        self.assertEqual(self.status("2026-01-08"), "active")
        self.assertEqual(self.status("2026-01-09"), "paused 2026-01-09")
        self.assertEqual(self.status("2026-03-01"), "paused 2026-01-09")

    def test_paused_day_serves_and_defers_nothing(self):
        self.break_off()
        before = sr.load_state()
        out = self.today("2026-01-20")
        after = sr.load_state()
        self.assertIn("paused", out.lower())
        self.assertIn("resume", out)
        for key in ("served", "deferred", "cards", "log"):
            self.assertEqual(before[key], after[key], key)

    def test_paused_state_changes_once_a_month(self):
        # Each change is a commit. GitHub disables schedules after 60 days
        # without one, so a pause must keep committing, but not daily.
        self.break_off()
        seen = set()
        day = dt.date(2026, 1, 9)
        last_change = day
        for i in range(120):
            d = day + dt.timedelta(i)
            before = sr.STATE_FILE.read_text()
            self.today(d.isoformat())
            if sr.STATE_FILE.read_text() != before:
                self.assertLess((d - last_change).days, 32)
                last_change = d
                seen.add(d)
        self.assertEqual(len(seen), 5)  # Jan 9, then Feb 1, Mar 1, Apr 1, May 1

    def test_status_agrees_with_a_plan_already_made(self):
        self.break_off()
        state = sr.load_state()
        state["served"]["2026-01-10"] = {"review": [], "new": ["1"]}
        sr.save_state(state)
        self.assertEqual(self.status("2026-01-10"), "active")
        self.assertIn("plan — 2026-01-10", self.today("2026-01-10"))

    def test_disabled_and_empty_never_pause(self):
        self.assertEqual(self.status("2026-01-01"), "active")
        self.break_off()
        state = sr.load_state()
        state["config"]["pause_after"] = 0
        sr.save_state(state)
        self.assertEqual(self.status("2027-01-01"), "active")

    def test_never_logged_pauses_from_the_first_plan(self):
        self.today("2026-01-01")
        self.assertEqual(self.status("2026-01-08"), "active")
        self.assertEqual(self.status("2026-01-09"), "paused 2026-01-09")

    def test_logging_resumes(self):
        plan = self.break_off()
        self.today("2026-01-20")
        self.assertIn("paused", sr.load_state())
        self.log("20", "good", "2026-02-01")
        self.assertEqual(self.status("2026-02-01"), "active")
        state = sr.load_state()
        self.assertEqual(state["deferred"], {})
        self.assertNotIn("paused", state)
        self.assertEqual(state["cards"]["20"]["due"], "2026-02-03")
        for pid in plan:
            self.assertEqual(state["cards"][pid]["due"], "2026-01-03")

    def test_resume_brings_back_every_old_card(self):
        self.set_config(pause_after=0)
        for pid in ("1", "20", "33", "70", "121", "206", "217", "704"):
            self.log(pid, "good", "2026-01-01")
        self.serve_days(90)  # ignored: shown over and over, deferred repeatedly
        state = sr.load_state()
        self.assertTrue(state["deferred"])
        state["config"]["pause_after"] = 7
        sr.save_state(state)
        self.resume("2026-04-01")
        back, start = set(), dt.date(2026, 4, 1)
        for i in range(7):
            d = (start + dt.timedelta(i)).isoformat()
            self.today(d)
            served = sr.load_state()["served"][d]
            back |= set(served["review"])
            for pid in served["review"] + served["new"]:
                self.log(pid, "good", d)
        self.assertLessEqual({"1", "20", "33", "70", "121", "206", "217", "704"}, back)

    def test_resume_then_ignored_new_problem_is_offered_again(self):
        self.break_off()
        self.resume("2026-02-01")
        self.today("2026-02-01")
        new = sr.load_state()["served"]["2026-02-01"]["new"][0]
        self.today("2026-02-02")
        self.assertIn(new, sr.load_state()["served"]["2026-02-02"]["new"])

    def test_resume_loses_no_problem(self):
        # Shown but never logged before the break, so in neither cards nor
        # (after the reset) deferred: it must still come back eventually.
        self.use_deck(6)
        self.serve_days(10)
        state = sr.load_state()
        state["config"]["pause_after"] = 7
        sr.save_state(state)
        self.resume("2026-03-01")
        state = sr.load_state()
        state["config"]["pause_after"] = 0  # keep serving without logs
        sr.save_state(state)
        shown = set()
        for day in range(1, 15):
            d = f"2026-03-{day:02d}"
            self.today(d)
            shown |= set(sr.load_state()["served"][d]["new"])
        self.assertEqual(len(shown), 6)

    def test_resume_when_active_changes_nothing(self):
        self.break_off()
        before = sr.STATE_FILE.read_text()
        self.assertIn("not paused", self.resume("2026-01-05"))
        self.today("2026-01-05", resume=True)
        self.assertEqual(before, sr.STATE_FILE.read_text())

    def test_today_can_resume(self):
        self.break_off()
        self.assertIn("plan — 2026-02-01", self.today("2026-02-01", resume=True))
        self.assertEqual(self.status("2026-02-01"), "active")

    def test_pauses_again_after_an_unused_resume(self):
        self.break_off()
        self.resume("2026-02-01")
        self.assertEqual(self.status("2026-02-08"), "active")
        self.assertEqual(self.status("2026-02-09"), "paused 2026-02-09")


class Comments(StateTest):
    """handle_comment: what log-solve.yml runs for each comment or log issue."""

    def handle(self, text, date="2026-01-01", event=None, late=False):
        state = sr.load_state()
        lines, done = sr.handle_comment(state, sr.load_problems(), text, date, event, late)
        sr.save_state(state)
        return "\n".join(lines), done

    def test_reply_lists_each_log(self):
        reply, done = self.handle("20 good\n54 hard")
        self.assertTrue(done)
        self.assertIn("Logged #20 Valid Parentheses: good", reply)
        self.assertIn("Logged #54 Spiral Matrix: hard", reply)
        self.assertEqual(len(sr.load_state()["log"]), 2)

    def test_a_comment_is_applied_once(self):
        self.handle("20 good", event="comment-1")
        reply, done = self.handle("20 good", date="2026-01-02", event="comment-1")
        self.assertFalse(done)
        self.assertIn("Already logged", reply)
        self.assertEqual(len(sr.load_state()["log"]), 1)

    def test_old_event_ids_are_dropped(self):
        self.handle("20 good", event="comment-1")
        self.handle("54 good", date="2026-03-01", event="comment-2")
        self.assertEqual(list(sr.load_state()["events"]["done"]), ["comment-2"])

    def test_nothing_done_records_nothing(self):
        self.handle("nice plan", event="comment-1")
        self.assertNotIn("events", sr.load_state())

    def test_reports_what_it_could_not_read(self):
        reply, done = self.handle("20 good\n53 was pretty good")
        self.assertTrue(done)
        self.assertIn("NOT LOGGED: #53", reply)
        self.assertIn("one problem per line", reply)

    def test_explains_the_format_for_an_attempt(self):
        reply, done = self.handle("solved it, felt good")
        self.assertFalse(done)
        self.assertIn("Nothing logged", reply)

    def test_stays_quiet_for_chatter(self):
        self.assertEqual(self.handle("thanks, nice plan"), ("", False))

    def paused(self):
        """Log on Jan 1 and nothing after: paused from Jan 9."""
        self.log("20", "good", "2026-01-01")
        self.assertTrue(sr.paused_since(sr.load_state(), "2026-01-09"))

    def test_resume_replies_with_the_plan(self):
        self.paused()
        reply, done = self.handle("resume", date="2026-01-20")
        self.assertTrue(done)
        self.assertIn("Welcome back", reply)
        self.assertIn("Practice plan — 2026-01-20", reply)
        self.assertEqual(self.status("2026-01-20"), "active")

    def test_undo_by_comment(self):
        self.handle("20 good")
        reply, done = self.handle("undo 20", date="2026-01-02")
        self.assertTrue(done)
        self.assertIn("Undid #20 Valid Parentheses: good on 2026-01-01", reply)
        self.assertNotIn("20", sr.load_state()["cards"])

    def test_late_skips_undo_and_problems_logged_since(self):
        self.handle("20 good", date="2026-01-03")
        reply, done = self.handle("20 hard\nundo 54\n54 easy", date="2026-01-02", late=True)
        self.assertTrue(done)
        state = sr.load_state()
        self.assertEqual([(e["id"], e["grade"]) for e in state["log"]],
                         [("20", "good"), ("54", "easy")])
        self.assertNotIn("NOT LOGGED", reply)


class Undo(StateTest):
    def snapshot(self):
        s = sr.load_state()
        return {k: s.get(k) for k in ("cards", "deferred", "log")}

    def undo(self, pid):
        return self.run_cmd(sr.cmd_undo, id=pid)

    def test_puts_the_card_back_exactly(self):
        self.log("20", "good", "2026-01-01")
        before = self.snapshot()
        self.log("20", "again", "2026-01-05")
        self.assertIn("Undid #20 Valid Parentheses: again on 2026-01-05", self.undo("20"))
        self.assertEqual(self.snapshot(), before)

    def test_a_first_solve_becomes_new_again(self):
        before = self.snapshot()
        self.log("20", "easy", "2026-01-01")
        self.assertIn("counts as a new problem again", self.undo("20"))
        after = self.snapshot()
        self.assertEqual(after["cards"], before["cards"])
        self.assertEqual(after["log"], [])

    def test_puts_the_deferral_back(self):
        pid, date = self.defer_one()
        deferral = sr.load_state()["deferred"][pid]
        self.log(pid, "good", date)
        self.assertNotIn(pid, sr.load_state()["deferred"])
        self.undo(pid)
        self.assertEqual(sr.load_state()["deferred"][pid], deferral)

    def test_walks_back_one_log_at_a_time(self):
        self.log("20", "good", "2026-01-01")
        first = self.snapshot()
        self.log("20", "good", "2026-01-04")
        self.log("20", "good", "2026-01-10")
        self.undo("20")
        self.undo("20")
        self.assertEqual(self.snapshot(), first)

    def test_refuses_a_log_from_before_undo_existed(self):
        self.log("20", "good", "2026-01-01")
        state = sr.load_state()
        del state["log"][-1]["prev"]
        sr.save_state(state)
        before = sr.STATE_FILE.read_text()
        self.assertIn("before undo existed", self.undo("20"))
        self.assertEqual(sr.STATE_FILE.read_text(), before)

    def test_nothing_to_undo(self):
        self.assertIn("never been logged", self.undo("20"))
        self.assertFalse(sr.STATE_FILE.exists())

    def test_a_second_log_the_same_day_replaces_the_first(self):
        self.log("54", "good", "2026-01-01")
        self.log("54", "hard", "2026-01-05")
        once = self.snapshot()
        self.undo("54")
        self.log("54", "good", "2026-01-05")
        out = self.log("54", "hard", "2026-01-05")
        self.assertIn("replaces good from earlier today", out)
        self.assertEqual(self.snapshot(), once)

    def test_logs_on_different_days_both_count(self):
        self.log("54", "good", "2026-01-01")
        self.log("54", "good", "2026-01-03")
        self.assertEqual(sr.load_state()["cards"]["54"]["reps"], 2)


class Sweep(StateTest):
    """The daily run's check for comments whose log-solve run never finished."""
    OWNER = "me"
    NOW = "2026-01-05T23:00:00Z"

    def setUp(self):
        super().setUp()
        self.issues = [
            {"number": 7, "id": 70, "title": "Daily Practice 2026-01-04", "user": "github-actions[bot]",
             "body": "plan", "created_at": "2026-01-04T23:30:00Z"},
            {"number": 8, "id": 80, "title": "Practice log 2026-01-05", "user": self.OWNER,
             "body": "54 hard", "created_at": "2026-01-05T18:00:00Z"},
            {"number": 9, "id": 90, "title": "Something else", "user": self.OWNER,
             "body": "", "created_at": "2026-01-01T00:00:00Z"},
        ]
        self.comments = []
        # The record starts with the first comment log-solve applies.
        state = sr.load_state()
        sr.handle_comment(state, sr.load_problems(), "1 good", "2026-01-04", "comment-1")
        state["events"]["since"] = "2026-01-04T00:00:00Z"
        sr.save_state(state)

    def comment(self, cid, issue, body, created_at, user=OWNER):
        self.comments.append({"id": cid, "issue": issue, "user": user, "body": body,
                              "created_at": created_at})

    def sweep(self, now=NOW):
        state = sr.load_state()
        replies = sr.sweep(state, sr.load_problems(), self.issues, self.comments, self.OWNER, now)
        sr.save_state(state)
        return replies

    def logged(self):
        return [(e["id"], e["grade"], e["date"]) for e in sr.load_state()["log"]]

    def test_logs_a_missed_comment_and_a_missed_log_issue(self):
        # 03:00 UTC on Jan 5 is Jan 4 in Los Angeles: logged under Jan 4.
        self.comment(2, 7, "20 good", "2026-01-05T03:00:00Z")
        replies = self.sweep()
        self.assertEqual(self.logged(), [("1", "good", "2026-01-04"), ("20", "good", "2026-01-04"),
                                         ("54", "hard", "2026-01-05")])
        self.assertEqual([(r["issue"], r["close"]) for r in replies], [(7, False), (8, True)])
        self.assertIn("Logged late", replies[0]["reply"])
        self.assertEqual(self.sweep(), [])  # recorded, so never twice

    def test_ignores_what_it_should(self):
        self.comment(1, 7, "1 good", "2026-01-04T20:00:00Z")    # already applied
        self.comment(3, 7, "20 good", "2026-01-03T20:00:00Z")   # before the record began
        self.comment(4, 7, "21 good", "2026-01-05T22:30:00Z")   # within the hour: may be running
        self.comment(5, 7, "33 good", "2026-01-05T10:00:00Z", user="someone")
        self.comment(6, 9, "70 good", "2026-01-05T10:00:00Z")   # not a plan or log issue
        self.issues[1]["user"] = "someone"
        self.assertEqual(self.sweep(), [])
        self.assertEqual(self.logged(), [("1", "good", "2026-01-04")])

    def test_a_problem_logged_again_since_keeps_the_newer_log(self):
        self.comment(2, 7, "20 hard", "2026-01-04T20:00:00Z")
        state = sr.load_state()
        sr.handle_comment(state, sr.load_problems(), "20 easy", "2026-01-05", "comment-3")
        sr.save_state(state)
        self.issues.pop(1)
        self.assertEqual(self.sweep(), [])
        self.assertEqual(self.logged()[-1], ("20", "easy", "2026-01-05"))

    def test_does_not_replay_undo(self):
        self.comment(2, 7, "undo 1", "2026-01-05T10:00:00Z")
        self.issues.pop(1)
        self.assertEqual(self.sweep(), [])
        self.assertIn("1", sr.load_state()["cards"])

    def test_does_nothing_before_any_comment_was_recorded(self):
        sr.STATE_FILE.unlink()
        self.comment(2, 7, "20 good", "2026-01-05T03:00:00Z")
        self.assertEqual(self.sweep(), [])
        self.assertEqual(self.logged(), [])

    def test_a_late_log_does_not_top_up(self):
        self.set_config(pause_after=0)
        self.run_cmd(sr.cmd_today, date="2026-01-04", md=True, json=False, resume=False)
        self.run_cmd(sr.cmd_today, date="2026-01-05", md=True, json=False, resume=False)
        state = sr.load_state()
        before = {k: list(v) for k, v in state["served"]["2026-01-05"].items()}
        carried = before["carried"][0]
        self.comment(2, 7, f"{carried} good", "2026-01-05T10:00:00Z")
        self.issues.pop(1)
        self.sweep()
        self.assertEqual(sr.load_state()["served"]["2026-01-05"], before)


class TopUp(StateTest):
    """A plan made overnight still holds yesterday's leftovers; logging one
    of them later that day gives its budget back."""
    DAY1, DAY2 = "2026-01-01", "2026-01-02"

    def setUp(self):
        super().setUp()
        self.set_config(pause_after=0)
        self.today(self.DAY1)
        self.today(self.DAY2)  # nothing logged in between
        self.by_id = {p["id"]: p for p in sr.load_problems()}

    def plan(self, date=DAY2):
        return sr.load_state()["served"][date]

    def handle(self, text, date=DAY2, event=None):
        state = sr.load_state()
        lines, _ = sr.handle_comment(state, sr.load_problems(), text, date, event)
        sr.save_state(state)
        return "\n".join(lines)

    def test_plan_records_what_it_carried(self):
        day1 = self.plan(self.DAY1)
        self.assertNotIn("carried", day1)
        self.assertEqual(day1["new"], self.plan()["new"])
        self.assertEqual(self.plan()["carried"], day1["new"])

    def test_logging_a_carried_problem_tops_up(self):
        first = self.plan()["new"][0]
        reply = self.handle(f"{first} good")
        plan = self.plan()
        self.assertEqual(plan["new"][:2], self.plan(self.DAY1)["new"])
        self.assertEqual(len(plan["new"]), 3)
        added = plan["new"][2]
        self.assertIn(f"today's plan gains #{added}", reply)
        self.assertNotIn(added, sr.load_state()["cards"])
        # still one category per new problem, counting the kept one
        kept = plan["new"][1]
        self.assertNotEqual(self.by_id[added]["category"], self.by_id[kept]["category"])

    def test_tops_up_once(self):
        first = self.plan()["new"][0]
        self.handle(f"{first} good")
        after = self.plan()
        self.handle(f"{first} easy")  # same-day re-log replaces the first
        self.assertEqual(self.plan(), after)
        state = sr.load_state()
        self.assertEqual(sr.top_up(state, sr.load_problems(), self.DAY2), [])

    def test_a_problem_first_shown_today_does_not_top_up(self):
        added = self.handle(f"{self.plan()['new'][0]} good")  # tops up once
        new_one = self.plan()["new"][-1]
        before = self.plan()
        reply = self.handle(f"{new_one} good")
        self.assertEqual(self.plan(), before)
        self.assertNotIn("gains", reply)
        self.assertIn("gains", added)

    def test_a_log_on_another_day_does_not_top_up(self):
        before = self.plan()
        self.handle(f"{before['new'][0]} good", date=self.DAY1)
        self.assertEqual(self.plan(), before)

    def test_a_carried_review_tops_up(self):
        sr.STATE_FILE.unlink()  # start over, with the config below
        self.set_config(new_per_day=0)  # reviews only: 2 points a day
        for pid in ("1", "20", "121"):  # all Easy, all due 01-03
            self.log(pid, "good", "2026-01-01")
        self.today("2026-01-03")  # 1 and 20; 121 waits
        self.today("2026-01-04")  # 121 first, then 1 again (shown, not logged)
        plan = self.plan("2026-01-04")
        self.assertEqual(plan["review"], ["121", "1"])
        self.assertEqual(plan["carried"], ["1"])
        self.handle("1 good", date="2026-01-04")
        self.assertEqual(self.plan("2026-01-04")["review"], ["121", "1", "20"])

    def test_handle_writes_the_plan_only_when_it_changed(self):
        body = pathlib.Path(self.tmp.name) / "body.txt"
        out = pathlib.Path(self.tmp.name) / "plan.md"

        def handle(text):
            body.write_text(text)
            self.run_cmd(sr.cmd_handle, body_file=str(body), event=None, date=self.DAY2,
                         reply_file=str(pathlib.Path(self.tmp.name) / "r.md"),
                         status_file=None, plan_file=str(out))
            return out.read_text()

        first, second = self.plan()["new"]
        md = handle(f"{first} good")
        self.assertIn("### Practice plan — 2026-01-02", md)
        self.assertIn(f"#{self.plan()['new'][2]} ", md)
        self.assertEqual(handle(f"{self.plan()['new'][2]} good"), "")


class StalePause(StateTest):
    def test_today_clears_a_pause_that_ended_without_resume(self):
        self.log("1", "good", "2026-01-01")
        self.today("2026-01-10")  # paused since 01-09
        self.assertIn("paused", sr.load_state())
        # The missed-log check logs a comment from 01-08 under its own date,
        # which ends the pause without calling resume().
        state = sr.load_state()
        sr.log_solve(state, {p["id"]: p for p in sr.load_problems()}, "20", "good", "2026-01-08")
        sr.save_state(state)
        self.assertEqual(self.status("2026-01-10"), "active")
        self.today("2026-01-10")
        self.assertNotIn("paused", sr.load_state())
        self.assertIn("2026-01-10", sr.load_state()["served"])


if __name__ == "__main__":
    unittest.main()
