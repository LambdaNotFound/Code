#!/usr/bin/env python3
"""Report Claude API token usage over a lookback window, from the Usage and Cost Admin API.

    python3 usage.py --hours 6 [--by model] [--series] [--json]
    python3 usage.py --days 3 --by api_key_id --cost

Stdlib only. Every number printed comes from the API's report; nothing is
estimated or priced locally.

Credentials, first found wins: ANTHROPIC_ADMIN_KEY, ANTHROPIC_API_KEY (must be
admin-scoped: an Admin key or a personal or service-account key not bound to a
workspace), ANTHROPIC_AUTH_TOKEN (OAuth token with the org:admin scope).
ANTHROPIC_BASE_URL overrides the host. Exit 2 with no credential, 1 on an API
or network error, 0 otherwise.

Endpoints and limits are from platform.claude.com/docs/en/manage-claude/usage-cost-api
and the usage-report reference, read 2026-09-18: bucket widths 1m, 1h, 1d with
page limits 1440, 168, 31; pagination by has_more and next_page; the cost report
is daily only, in cents as decimal strings.
"""
from __future__ import annotations

import argparse
import datetime as dt
import json
import os
import sys
import urllib.error
import urllib.parse
import urllib.request
from decimal import Decimal

API_VERSION = "2023-06-01"
DEFAULT_BASE = "https://api.anthropic.com"
USER_AGENT = "token-usage-skill/1.0 (https://github.com/LambdaNotFound/Code)"
PAGE_LIMIT = {"1m": 1440, "1h": 168, "1d": 31}
WIDTH_SECONDS = {"1m": 60, "1h": 3600, "1d": 86400}
GROUPS = ("model", "api_key_id", "workspace_id", "service_tier", "context_window", "inference_geo", "account_id")
TOKEN_FIELDS = ("uncached_input_tokens", "cache_write", "cache_read_input_tokens", "output_tokens")


class ApiError(Exception):
    """The API refused or the network failed; the message is user-facing."""


# ---------------------------------------------------------------- time

def floor_to(t: dt.datetime, width: str) -> dt.datetime:
    t = t.astimezone(dt.timezone.utc)
    if width == "1d":
        return t.replace(hour=0, minute=0, second=0, microsecond=0)
    if width == "1h":
        return t.replace(minute=0, second=0, microsecond=0)
    return t.replace(second=0, microsecond=0)


def ceil_to(t: dt.datetime, width: str) -> dt.datetime:
    f = floor_to(t, width)
    return f if f == t.astimezone(dt.timezone.utc) else f + dt.timedelta(seconds=WIDTH_SECONDS[width])


def window(hours: int | None, days: int | None, now: dt.datetime) -> tuple[dt.datetime, dt.datetime, str]:
    """Start (inclusive), end (exclusive), bucket width. Both ends snap to whole
    buckets so the report never drops the partial bucket the window ends in."""
    if hours is not None:
        width = "1h" if hours <= PAGE_LIMIT["1h"] else "1d"
        back = dt.timedelta(hours=hours)
    else:
        width = "1d"
        back = dt.timedelta(days=days or 7)
    return floor_to(now - back, width), ceil_to(now, width), width


def rfc3339(t: dt.datetime) -> str:
    return t.astimezone(dt.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


# ---------------------------------------------------------------- http

def credential(env: dict) -> tuple[dict, str] | None:
    """Headers for the first credential found, and which variable supplied it."""
    for var in ("ANTHROPIC_ADMIN_KEY", "ANTHROPIC_API_KEY"):
        if env.get(var):
            return {"x-api-key": env[var]}, var
    if env.get("ANTHROPIC_AUTH_TOKEN"):
        return {"Authorization": "Bearer " + env["ANTHROPIC_AUTH_TOKEN"], "anthropic-beta": "oauth-2025-04-20"}, "ANTHROPIC_AUTH_TOKEN"
    return None


def get_json(base: str, path: str, params: list[tuple[str, str]], headers: dict, source: str) -> dict:
    url = base.rstrip("/") + path + "?" + urllib.parse.urlencode(params)
    req = urllib.request.Request(url, headers={"anthropic-version": API_VERSION, "User-Agent": USER_AGENT, **headers})
    try:
        with urllib.request.urlopen(req, timeout=60) as r:
            return json.loads(r.read().decode("utf-8"))
    except urllib.error.HTTPError as e:
        body = e.read().decode("utf-8", errors="replace")
        try:
            msg = json.loads(body)["error"]["message"]
        except (ValueError, KeyError, TypeError):
            msg = body[:300]
        if e.code in (401, 403):
            raise ApiError(f"HTTP {e.code} from {path}: {msg}\n"
                           f"The credential in {source} is not accepted by the Admin API. It needs an Admin API "
                           f"key (sk-ant-admin...), a personal or service-account key not scoped to a workspace, "
                           f"or an OAuth token with org:admin. Workspace keys are rejected, and individual accounts "
                           f"without an organization have no Admin API.") from None
        if e.code == 429:
            raise ApiError(f"HTTP 429 from {path}: {msg} (retry-after: {e.headers.get('retry-after', '?')}s; "
                           f"the API allows sustained polling once per minute)") from None
        raise ApiError(f"HTTP {e.code} from {path}: {msg}") from None
    except urllib.error.URLError as e:
        raise ApiError(f"could not reach {base}: {e.reason}") from None


def paged(base: str, path: str, params: list[tuple[str, str]], headers: dict, source: str) -> list[dict]:
    """All time buckets across pages, oldest first."""
    buckets, page = [], None
    while True:
        q = params + ([("page", page)] if page else [])
        doc = get_json(base, path, q, headers, source)
        buckets.extend(doc.get("data", []))
        if not doc.get("has_more") or not doc.get("next_page"):
            return buckets
        page = doc["next_page"]


def fetch_usage(base, headers, source, start, end, width, group_by) -> list[dict]:
    params = [("starting_at", rfc3339(start)), ("ending_at", rfc3339(end)), ("bucket_width", width),
              ("limit", str(PAGE_LIMIT[width]))]
    if group_by:
        params.append(("group_by[]", group_by))
    return paged(base, "/v1/organizations/usage_report/messages", params, headers, source)


def fetch_cost(base, headers, source, start, end) -> list[dict]:
    params = [("starting_at", rfc3339(floor_to(start, "1d"))), ("ending_at", rfc3339(ceil_to(end, "1d"))),
              ("bucket_width", "1d"), ("limit", str(PAGE_LIMIT["1d"])), ("group_by[]", "description")]
    return paged(base, "/v1/organizations/cost_report", params, headers, source)


# ---------------------------------------------------------------- aggregate

def row_tokens(r: dict) -> dict:
    cc = r.get("cache_creation") or {}
    return {
        "uncached_input_tokens": int(r.get("uncached_input_tokens") or 0),
        "cache_write": int(cc.get("ephemeral_5m_input_tokens") or 0) + int(cc.get("ephemeral_1h_input_tokens") or 0),
        "cache_read_input_tokens": int(r.get("cache_read_input_tokens") or 0),
        "output_tokens": int(r.get("output_tokens") or 0),
        "web_search_requests": int((r.get("server_tool_use") or {}).get("web_search_requests") or 0),
    }


def aggregate(buckets: list[dict], group_by: str | None) -> dict[str, dict]:
    """Group label -> summed token fields across every bucket. The label for a
    null group key is '(none)': the default workspace, Console usage, and so on."""
    out: dict[str, dict] = {}
    for b in buckets:
        for r in b.get("results", []):
            label = "all" if not group_by else str(r.get(group_by) if r.get(group_by) is not None else "(none)")
            acc = out.setdefault(label, dict.fromkeys(TOKEN_FIELDS + ("web_search_requests",), 0))
            for k, v in row_tokens(r).items():
                acc[k] += v
    return out


def series(buckets: list[dict]) -> list[tuple[str, int, int]]:
    """(bucket start, input tokens of every kind, output tokens) per bucket."""
    out = []
    for b in buckets:
        tin = tout = 0
        for r in b.get("results", []):
            t = row_tokens(r)
            tin += t["uncached_input_tokens"] + t["cache_write"] + t["cache_read_input_tokens"]
            tout += t["output_tokens"]
        out.append((b.get("starting_at", "?"), tin, tout))
    return out


def cost_by_description(buckets: list[dict]) -> dict[str, Decimal]:
    """Description -> USD. The API reports cents as decimal strings."""
    out: dict[str, Decimal] = {}
    for b in buckets:
        for r in b.get("results", []):
            label = str(r.get("description") or r.get("cost_type") or "(unlabelled)")
            out[label] = out.get(label, Decimal(0)) + Decimal(str(r.get("amount") or "0")) / 100
    return out


# ---------------------------------------------------------------- print

def fmt(n: int) -> str:
    return f"{n:,}"


def print_usage(start, end, width, group_by, buckets, show_series) -> None:
    n = len(buckets)
    unit = {"1h": "hourly", "1d": "daily", "1m": "minute"}[width]
    print(f"Usage {rfc3339(start)} to {rfc3339(end)} ({n} {unit} bucket{'s' if n != 1 else ''}), "
          f"by {group_by or 'organization'}")
    agg = aggregate(buckets, group_by)
    if not agg:
        print("  the API reports no usage in this window")
        return
    rows = sorted(agg.items(), key=lambda kv: -sum(kv[1][f] for f in TOKEN_FIELDS))
    label_w = max(len(group_by or "scope"), *(len(k) for k, _ in rows))
    head = f"{(group_by or 'scope'):<{label_w}}  {'input':>13} {'cache write':>13} {'cache read':>13} {'output':>13} {'total':>14}"
    print(head)
    print("-" * len(head))
    total = dict.fromkeys(TOKEN_FIELDS + ("web_search_requests",), 0)
    for label, t in rows:
        tot = sum(t[f] for f in TOKEN_FIELDS)
        print(f"{label:<{label_w}}  {fmt(t['uncached_input_tokens']):>13} {fmt(t['cache_write']):>13} "
              f"{fmt(t['cache_read_input_tokens']):>13} {fmt(t['output_tokens']):>13} {fmt(tot):>14}")
        for k in total:
            total[k] += t[k]
    if len(rows) > 1:
        print("-" * len(head))
        print(f"{'total':<{label_w}}  {fmt(total['uncached_input_tokens']):>13} {fmt(total['cache_write']):>13} "
              f"{fmt(total['cache_read_input_tokens']):>13} {fmt(total['output_tokens']):>13} "
              f"{fmt(sum(total[f] for f in TOKEN_FIELDS)):>14}")
    if total["web_search_requests"]:
        print(f"web search requests: {fmt(total['web_search_requests'])}")
    print("input = uncached input tokens; cache write = 5m + 1h cache creation; total = all four columns")
    if show_series:
        print()
        print(f"{'bucket start':<22} {'input (all)':>13} {'output':>13}")
        for when, tin, tout in series(buckets):
            print(f"{when:<22} {fmt(tin):>13} {fmt(tout):>13}")


def print_cost(start, end, buckets) -> None:
    d0, d1 = floor_to(start, "1d"), ceil_to(end, "1d")
    n = len(buckets)
    print()
    print(f"Cost {d0.date()} to {(d1 - dt.timedelta(days=1)).date()} inclusive ({n} daily bucket{'s' if n != 1 else ''}; "
          f"the cost report is daily only, so this covers whole days)")
    by = cost_by_description(buckets)
    if not by:
        print("  the API reports no cost in these days")
        return
    w = max(len("description"), *(len(k) for k in by))
    for label, usd in sorted(by.items(), key=lambda kv: -kv[1]):
        print(f"  {label:<{w}}  ${usd:>12,.2f}")
    print(f"  {'total':<{w}}  ${sum(by.values()):>12,.2f}")


# ---------------------------------------------------------------- cli

def main(argv: list[str] | None = None, env: dict | None = None, now: dt.datetime | None = None) -> int:
    env = os.environ if env is None else env
    now = now or dt.datetime.now(dt.timezone.utc)
    ap = argparse.ArgumentParser(description="Claude API token usage over a lookback window (Admin API).")
    span = ap.add_mutually_exclusive_group()
    span.add_argument("--hours", type=int, help="look back this many hours (hourly buckets up to 168, daily beyond)")
    span.add_argument("--days", type=int, help="look back this many days (daily buckets; default 7)")
    ap.add_argument("--by", choices=GROUPS + ("none",), default="model", help="group rows by this dimension (default model)")
    ap.add_argument("--cost", action="store_true", help="also fetch the cost report (daily granularity, USD)")
    ap.add_argument("--series", action="store_true", help="also print per-bucket totals")
    ap.add_argument("--json", action="store_true", help="print the merged raw buckets as JSON instead of a table")
    args = ap.parse_args(argv)
    for name in ("hours", "days"):
        v = getattr(args, name)
        if v is not None and v <= 0:
            ap.error(f"--{name} must be a positive integer")

    cred = credential(env)
    if not cred:
        print("no credential: set ANTHROPIC_ADMIN_KEY (an Admin API key from Claude Console, Settings, Admin API keys), "
              "or ANTHROPIC_API_KEY if that key is admin-scoped, or ANTHROPIC_AUTH_TOKEN (OAuth, org:admin). "
              "Nothing was sent.", file=sys.stderr)
        return 2
    headers, source = cred
    base = env.get("ANTHROPIC_BASE_URL") or DEFAULT_BASE
    start, end, width = window(args.hours, args.days, now)
    group_by = None if args.by == "none" else args.by

    try:
        buckets = fetch_usage(base, headers, source, start, end, width, group_by)
        cost = fetch_cost(base, headers, source, start, end) if args.cost else None
    except ApiError as e:
        print(f"error: {e}", file=sys.stderr)
        return 1

    if args.json:
        doc = {"starting_at": rfc3339(start), "ending_at": rfc3339(end), "bucket_width": width,
               "group_by": group_by, "usage": buckets}
        if cost is not None:
            doc["cost"] = cost
        print(json.dumps(doc, indent=2))
        return 0
    print_usage(start, end, width, group_by, buckets, args.series)
    if cost is not None:
        print_cost(start, end, cost)
    if end > now:
        print(f"note: the window ends {rfc3339(end)}, after now; the last bucket is still filling, "
              f"and usage lands about five minutes after a request completes")
    return 0


if __name__ == "__main__":
    sys.exit(main())
