---
description: 'Report the organization''s Claude API token usage over a lookback window, straight from the Usage and Cost Admin API. Runs a stdlib script that takes --hours N or --days N, picks the bucket width, follows pagination, and prints tokens by model (or by API key, workspace, service tier, or context window) with cache reads and cache writes separated; --cost adds the billed cost in USD from the cost report. Needs an Admin API key or another admin-scoped credential; workspace keys are rejected by the API. Use when the user asks how many tokens they have used, wants usage or spend for the last N hours or days, or wants a usage breakdown by model or by key. Not for counting the tokens in a prompt before sending it (use claude-api). Not for cutting the bill (use claude-api, whose cost-optimize flow starts from this data).'
argument-hint: '[--hours N | --days N] [--by model|api_key_id|workspace_id|service_tier|context_window|inference_geo|account_id|none] [--cost] [--series] [--json]'
---

You report token usage from the Usage and Cost Admin API. You do not
estimate: every number comes from the API's own report, which counts
what Anthropic metered, not what a client logged. You never turn
tokens into dollars with a price table; `--cost` reads the cost
report, which is the billed figure.

## Arguments

$ARGUMENTS

## Run it

Map the request to flags and run the script. It owns the time math,
the bucket choice, the pagination, and the table:

```
python3 .claude/skills/token-usage/scripts/usage.py --hours 6
python3 .claude/skills/token-usage/scripts/usage.py --days 3 --by api_key_id --cost
```

- "last 6 hours" is `--hours 6`; "this week" is `--days 7`. With no
  window given, run `--days 7` and say that you did.
- Up to 168 hours uses hourly buckets; longer windows, and every
  `--days` window, use daily buckets. The script floors the start to
  a bucket boundary so every bucket is whole, and prints the exact
  range it queried. Windows longer than the API's page size are
  paged, not truncated.
- `--by` groups the table; `model` is the default. `--series` adds a
  per-bucket line of totals. `--json` prints the merged raw buckets
  for the user to keep or pipe elsewhere.

## Credentials

The script reads, in order: `ANTHROPIC_ADMIN_KEY`; `ANTHROPIC_API_KEY`
(works only if it is an Admin key, `sk-ant-admin...`, or a personal or
service-account key not scoped to a workspace); `ANTHROPIC_AUTH_TOKEN`
(an OAuth token with the `org:admin` scope). It never prints a
credential. `ANTHROPIC_BASE_URL` overrides the host.

If none is set, the script says so and stops. Tell the user where an
Admin key comes from (Claude Console, Settings, Admin API keys), that
the Admin API is unavailable to individual accounts without an
organization, and that they export the key in their shell. Never ask
for a key to be pasted into the conversation. A 401 with a key set
means the key is not admin-scoped; say that rather than retrying.

## Report

Lead with the window's total and the largest line, then the table the
script printed. State the window and the bucket width once. Usage
lands about five minutes after a request completes, so when the
window ends now the newest bucket is usually short; say so. If the
API returned no rows, say the API reports no usage in that window.
Add nothing the API did not return: no projections, no per-token
prices, no "roughly".
