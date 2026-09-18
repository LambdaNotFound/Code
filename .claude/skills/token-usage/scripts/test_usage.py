#!/usr/bin/env python3
"""Tests for usage.py against a local fake of the Admin API. Stdlib only.

    python3 .claude/skills/token-usage/scripts/test_usage.py
"""
import contextlib
import datetime as dt
import io
import json
import os
import sys
import threading
import unittest
from decimal import Decimal
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from pathlib import Path
from urllib.parse import parse_qs, urlsplit

sys.path.insert(0, str(Path(__file__).resolve().parent))
import usage  # noqa: E402

NOW = dt.datetime(2026, 9, 18, 3, 31, 7, tzinfo=dt.timezone.utc)


def result(model="claude-opus-5", uncached=1500, c5=100, c1=20, read=200, out=500, ws=3, **extra):
    r = {"model": model, "uncached_input_tokens": uncached,
         "cache_creation": {"ephemeral_5m_input_tokens": c5, "ephemeral_1h_input_tokens": c1},
         "cache_read_input_tokens": read, "output_tokens": out, "server_tool_use": {"web_search_requests": ws},
         "api_key_id": None, "workspace_id": None, "service_tier": None, "context_window": None}
    r.update(extra)
    return r


class FakeAdminAPI:
    """Serves programmed pages for the two report endpoints and records every request."""

    def __init__(self):
        self.requests: list[dict] = []
        self.usage_pages: list[dict] = [{"data": [], "has_more": False, "next_page": None}]
        self.cost_pages: list[dict] = [{"data": [], "has_more": False, "next_page": None}]
        self.status = 200
        self.error_body = None
        api = self

        class H(BaseHTTPRequestHandler):
            def log_message(self, *a):
                pass

            def do_GET(self):
                u = urlsplit(self.path)
                q = parse_qs(u.query)
                api.requests.append({"path": u.path, "query": q, "headers": {k.lower(): v for k, v in self.headers.items()}})
                if api.status != 200:
                    body = json.dumps(api.error_body or {"type": "error", "error": {"type": "x", "message": "nope"}}).encode()
                    self.send_response(api.status)
                    if api.status == 429:
                        self.send_header("retry-after", "30")
                    self.send_header("Content-Type", "application/json")
                    self.send_header("Content-Length", str(len(body)))
                    self.end_headers()
                    self.wfile.write(body)
                    return
                pages = api.usage_pages if u.path.endswith("/usage_report/messages") else api.cost_pages
                page = q.get("page", [None])[0]
                idx = int(page.split("_")[1]) if page else 0
                body = json.dumps(pages[idx]).encode()
                self.send_response(200)
                self.send_header("Content-Type", "application/json")
                self.send_header("Content-Length", str(len(body)))
                self.end_headers()
                self.wfile.write(body)

        self.server = ThreadingHTTPServer(("127.0.0.1", 0), H)
        self.thread = threading.Thread(target=self.server.serve_forever, daemon=True)
        self.thread.start()
        self.base = f"http://127.0.0.1:{self.server.server_address[1]}"

    def close(self):
        self.server.shutdown()
        self.server.server_close()


def run(api, argv, env_extra=None):
    env = {"ANTHROPIC_BASE_URL": api.base, "ANTHROPIC_ADMIN_KEY": "sk-ant-admin01-test"}
    env.update(env_extra or {})
    out, err = io.StringIO(), io.StringIO()
    with contextlib.redirect_stdout(out), contextlib.redirect_stderr(err):
        rc = usage.main(argv, env=env, now=NOW)
    return rc, out.getvalue(), err.getvalue()


class Window(unittest.TestCase):
    def test_hours_snap_to_whole_hourly_buckets_including_the_current_one(self):
        start, end, width = usage.window(6, None, NOW)
        self.assertEqual(width, "1h")
        self.assertEqual(usage.rfc3339(start), "2026-09-17T21:00:00Z")
        self.assertEqual(usage.rfc3339(end), "2026-09-18T04:00:00Z")

    def test_hours_beyond_one_page_of_hourly_buckets_switch_to_daily(self):
        _, _, width = usage.window(169, None, NOW)
        self.assertEqual(width, "1d")
        _, _, width = usage.window(168, None, NOW)
        self.assertEqual(width, "1h")

    def test_days_snap_to_whole_days(self):
        start, end, width = usage.window(None, 3, NOW)
        self.assertEqual(width, "1d")
        self.assertEqual(usage.rfc3339(start), "2026-09-15T00:00:00Z")
        self.assertEqual(usage.rfc3339(end), "2026-09-19T00:00:00Z")

    def test_default_window_is_seven_days(self):
        start, _, _ = usage.window(None, None, NOW)
        self.assertEqual(usage.rfc3339(start), "2026-09-11T00:00:00Z")

    def test_ceil_on_an_exact_boundary_does_not_add_a_bucket(self):
        t = dt.datetime(2026, 9, 18, 4, 0, tzinfo=dt.timezone.utc)
        self.assertEqual(usage.ceil_to(t, "1h"), t)


class Requests(unittest.TestCase):
    def setUp(self):
        self.api = FakeAdminAPI()
        os.environ["no_proxy"] = "127.0.0.1," + os.environ.get("no_proxy", "")

    def tearDown(self):
        self.api.close()

    def test_query_shape_and_headers(self):
        rc, out, err = run(self.api, ["--hours", "6"])
        self.assertEqual(rc, 0, err)
        req = self.api.requests[0]
        self.assertEqual(req["path"], "/v1/organizations/usage_report/messages")
        q = req["query"]
        self.assertEqual(q["starting_at"], ["2026-09-17T21:00:00Z"])
        self.assertEqual(q["ending_at"], ["2026-09-18T04:00:00Z"])
        self.assertEqual(q["bucket_width"], ["1h"])
        self.assertEqual(q["limit"], ["168"])
        self.assertEqual(q["group_by[]"], ["model"])
        self.assertEqual(req["headers"]["anthropic-version"], "2023-06-01")
        self.assertEqual(req["headers"]["x-api-key"], "sk-ant-admin01-test")
        self.assertIn("token-usage-skill", req["headers"]["user-agent"])

    def test_by_none_sends_no_group_by(self):
        run(self.api, ["--days", "1", "--by", "none"])
        self.assertNotIn("group_by[]", self.api.requests[0]["query"])

    def test_pagination_follows_next_page_and_merges_oldest_first(self):
        self.api.usage_pages = [
            {"data": [{"starting_at": "2026-09-17T21:00:00Z", "ending_at": "2026-09-17T22:00:00Z", "results": [result()]}],
             "has_more": True, "next_page": "page_1"},
            {"data": [{"starting_at": "2026-09-17T22:00:00Z", "ending_at": "2026-09-17T23:00:00Z", "results": [result(out=7)]}],
             "has_more": False, "next_page": None},
        ]
        rc, out, err = run(self.api, ["--hours", "2", "--json"])
        self.assertEqual(rc, 0, err)
        self.assertEqual(len(self.api.requests), 2)
        self.assertEqual(self.api.requests[1]["query"]["page"], ["page_1"])
        doc = json.loads(out)
        self.assertEqual([b["starting_at"] for b in doc["usage"]], ["2026-09-17T21:00:00Z", "2026-09-17T22:00:00Z"])

    def test_oauth_token_goes_in_bearer_with_the_oauth_beta(self):
        run(self.api, ["--hours", "1"], {"ANTHROPIC_ADMIN_KEY": "", "ANTHROPIC_AUTH_TOKEN": "tok"})
        h = self.api.requests[0]["headers"]
        self.assertEqual(h["authorization"], "Bearer tok")
        self.assertEqual(h["anthropic-beta"], "oauth-2025-04-20")
        self.assertNotIn("x-api-key", h)

    def test_no_credential_sends_nothing_and_exits_2(self):
        rc, out, err = run(self.api, ["--hours", "1"], {"ANTHROPIC_ADMIN_KEY": ""})
        self.assertEqual(rc, 2)
        self.assertEqual(self.api.requests, [])
        self.assertIn("ANTHROPIC_ADMIN_KEY", err)

    def test_401_names_the_admin_requirement_and_the_variable(self):
        self.api.status = 401
        self.api.error_body = {"type": "error", "error": {"type": "authentication_error", "message": "invalid x-api-key"}}
        rc, out, err = run(self.api, ["--hours", "1"], {"ANTHROPIC_ADMIN_KEY": "", "ANTHROPIC_API_KEY": "sk-ant-api03-x"})
        self.assertEqual(rc, 1)
        self.assertIn("HTTP 401", err)
        self.assertIn("ANTHROPIC_API_KEY", err)
        self.assertIn("Admin API key", err)
        self.assertNotIn("sk-ant-api03-x", err + out, "the credential must never be printed")

    def test_429_reports_retry_after(self):
        self.api.status = 429
        rc, out, err = run(self.api, ["--hours", "1"])
        self.assertEqual(rc, 1)
        self.assertIn("retry-after: 30", err)

    def test_cost_flag_queries_the_cost_report_for_whole_days(self):
        self.api.cost_pages = [{"data": [{"starting_at": "2026-09-17T00:00:00Z", "ending_at": "2026-09-18T00:00:00Z",
                                          "results": [{"description": "Claude Opus 5 Usage - Input Tokens", "amount": "12345.5",
                                                       "currency": "USD", "cost_type": "tokens"}]}],
                                "has_more": False, "next_page": None}]
        rc, out, err = run(self.api, ["--hours", "6", "--cost"])
        self.assertEqual(rc, 0, err)
        cost_req = [r for r in self.api.requests if r["path"].endswith("/cost_report")][0]
        self.assertEqual(cost_req["query"]["bucket_width"], ["1d"])
        self.assertEqual(cost_req["query"]["starting_at"], ["2026-09-17T00:00:00Z"])
        self.assertEqual(cost_req["query"]["ending_at"], ["2026-09-19T00:00:00Z"])
        self.assertEqual(cost_req["query"]["group_by[]"], ["description"])
        self.assertIn("$      123.46", out)   # 12345.5 cents


class Aggregation(unittest.TestCase):
    def test_sums_across_buckets_and_groups_and_labels_null_keys(self):
        buckets = [
            {"results": [result(model="claude-opus-5"), result(model="claude-sonnet-5", uncached=10, c5=0, c1=0, read=0, out=1, ws=0)]},
            {"results": [result(model="claude-opus-5", out=5), result(model=None, uncached=1, c5=0, c1=0, read=0, out=0, ws=0)]},
        ]
        agg = usage.aggregate(buckets, "model")
        self.assertEqual(agg["claude-opus-5"]["uncached_input_tokens"], 3000)
        self.assertEqual(agg["claude-opus-5"]["cache_write"], 240)      # (100 + 20) * 2
        self.assertEqual(agg["claude-opus-5"]["output_tokens"], 505)
        self.assertEqual(agg["claude-opus-5"]["web_search_requests"], 6)
        self.assertEqual(agg["claude-sonnet-5"]["output_tokens"], 1)
        self.assertEqual(agg["(none)"]["uncached_input_tokens"], 1)
        self.assertEqual(usage.aggregate(buckets, None).keys(), {"all"})

    def test_series_reports_all_input_kinds_per_bucket(self):
        buckets = [{"starting_at": "t0", "results": [result()]}]
        self.assertEqual(usage.series(buckets), [("t0", 1500 + 120 + 200, 500)])

    def test_cost_amounts_are_cents_summed_as_decimals(self):
        buckets = [{"results": [{"description": "a", "amount": "0.1"}, {"description": "a", "amount": "0.2"},
                                {"description": "b", "amount": "100"}]}]
        by = usage.cost_by_description(buckets)
        self.assertEqual(by["a"], Decimal("0.003"))
        self.assertEqual(by["b"], Decimal("1"))

    def test_table_prints_every_group_a_total_row_and_the_legend(self):
        buckets = [{"starting_at": "2026-09-17T21:00:00Z", "results": [result(), result(model="claude-sonnet-5")]}]
        out = io.StringIO()
        with contextlib.redirect_stdout(out):
            usage.print_usage(NOW, NOW, "1h", "model", buckets, show_series=True)
        text = out.getvalue()
        for needle in ("claude-opus-5", "claude-sonnet-5", "total", "web search requests: 6", "cache write = 5m + 1h", "bucket start"):
            self.assertIn(needle, text)

    def test_empty_report_says_so(self):
        out = io.StringIO()
        with contextlib.redirect_stdout(out):
            usage.print_usage(NOW, NOW, "1d", "model", [], show_series=False)
        self.assertIn("no usage in this window", out.getvalue())


class Cli(unittest.TestCase):
    def test_hours_and_days_are_exclusive_and_positive(self):
        with self.assertRaises(SystemExit):
            with contextlib.redirect_stderr(io.StringIO()):
                usage.main(["--hours", "1", "--days", "1"], env={}, now=NOW)
        with self.assertRaises(SystemExit):
            with contextlib.redirect_stderr(io.StringIO()):
                usage.main(["--hours", "0"], env={}, now=NOW)


if __name__ == "__main__":
    unittest.main(verbosity=1)
