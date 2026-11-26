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


class ParseSolves(unittest.TestCase):
    def check(self, text, pairs, skipped=()):
        got_pairs, got_skipped = sr.parse_solves(text)
        self.assertEqual(got_pairs, list(pairs), text)
        self.assertEqual([pid for pid, _ in got_skipped], list(skipped), text)

    def test_documented_forms(self):
        self.check("solved 200 good", [("200", "good")])
        self.check("1 easy, 33 again", [("1", "easy"), ("33", "again")])
        self.check("#5 hard", [("5", "hard")])
        self.check("#33 was hard", [("33", "hard")])
        self.check("solved 200 good, #33 was hard", [("200", "good"), ("33", "hard")])
        self.check("252 good and 9 hard", [("252", "good"), ("9", "hard")])
        self.check("1 good · 20 hard", [("1", "good"), ("20", "hard")])
        self.check("1: good\n20 - hard", [("1", "good"), ("20", "hard")])
        self.check("LC200 good, lc 33 again", [("200", "good"), ("33", "again")])
        self.check("200 (good)", [("200", "good")])
        self.check("**200** GOOD", [("200", "good")])

    def test_every_alias_parses(self):
        for word, grade in sr.GRADE_ALIASES.items():
            self.check(f"#20 {word}", [("20", grade)])

    def test_punctuation_still_separates(self):
        self.check("1. good", [])
        self.check("2 attempts, 1 good", [("1", "good")])

    def test_numbers_with_a_unit_are_not_ids(self):
        self.check("solved 1 in 20 min easy", [("1", "easy")])
        self.check("#70 took 45 minutes hard", [("70", "hard")])
        self.check("#70 took 1.5 hours hard", [("70", "hard")])
        self.check("#70 20-min easy", [("70", "easy")])
        self.check("did 2 problems good", [])
        self.check("beats 95% easy", [])
        self.check("2026-08-04 good", [])

    def test_hash_id_allows_a_title(self):
        self.check("#208 Implement Trie good", [("208", "good")])
        self.check("#15 3Sum again", [("15", "again")])
        self.check("#542 01 Matrix hard", [("542", "hard")])
        self.check("#50 Pow(x, n) good", [("50", "good")])
        self.check("#380 Insert Delete GetRandom O(1) easy", [("380", "easy")])

    def test_urls_and_slashes(self):
        self.check("https://leetcode.ca/all/252.html good", [])
        self.check("#252 https://leetcode.ca/all/252.html good", [("252", "good")])
        self.check("200/33 good", [("33", "good")])

    def test_bare_number_needs_a_nearby_grade(self):
        self.check("started at 9 this morning and it went good", [])

    def test_difficulty_label_is_not_a_grade(self):
        line = ("- [ ] **#208 [Implement Trie (Prefix Tree)]"
                "(https://leetcode.com/problems/implement-trie-prefix-tree/)** (Medium)"
                " — Trie · hash-table, string · reps 2, last 2026-08-04")
        self.check(line, [], ["208"])
        self.check("#9 Palindrome Number (Easy) again", [("9", "again")])

    def test_negated_grade_is_skipped(self):
        self.check("#9 not good", [], ["9"])
        self.check("#9 wasn't easy", [], ["9"])
        self.check("#9 wasn’t easy", [], ["9"])
        self.check("#20 no good", [], ["20"])
        self.check("9 not good, 1 easy", [("1", "easy")], ["9"])

    def test_no_and_yes_only_count_at_the_end(self):
        self.check("#20 no", [("20", "again")])
        self.check("#20 no problem, easy", [], ["20"])
        self.check("#20 no issues good", [("20", "good")])

    def test_conflicting_grades_are_skipped(self):
        self.check("#33 slow but ok", [], ["33"])
        self.check("#33 quick and easy", [("33", "easy")])

    def test_skipped_id_logged_elsewhere_is_not_reported(self):
        self.check("#9 reread the notes\n#9 good", [("9", "good")])

    def test_quoted_lines_are_ignored(self):
        self.check("> Log solves e.g. `solved 1 good`, `#20 hard`\n\n9 easy", [("9", "easy")])

    def test_prose_is_quiet(self):
        self.check("nice plan today", [])


class Scheduling(unittest.TestCase):
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
            self.log("1", "hard", "2026-01-01")
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


if __name__ == "__main__":
    unittest.main()
