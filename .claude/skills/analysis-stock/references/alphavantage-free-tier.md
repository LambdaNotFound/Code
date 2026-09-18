# Alpha Vantage through the MCP server: what the free key actually does

Everything here was observed in this repo's sessions on 2026-09-18
`[COMPUTED]`; the vendor changes tiers without notice, so a "premium
endpoint" reply is the authority when it disagrees with this page.

## Limits

- 25 requests a day and one a second. The per-second limit bites when
  two Alpha Vantage calls go out in the same turn: the second one comes
  back `rate_limit` and has to be retried, and both count. Fire them
  one per turn.
- A daily-limit reply names the 25-a-day cap in its message; a
  per-second reply says "1 request per second"; a tier reply says
  "premium endpoint" or "premium feature". Read the message before
  retrying: only the per-second case is worth a retry.

## Free on this key

| Endpoint | Notes |
|---|---|
| `TIME_SERIES_DAILY` with `outputsize=compact` | 100 bars, unadjusted |
| `TIME_SERIES_WEEKLY` and `TIME_SERIES_WEEKLY_ADJUSTED` | full history since listing; prefer adjusted (it carries an `adjusted close`; the scripts scale open/high/low by the same factor) |
| Indicator endpoints (`SMA`, `RSI`, `MACD`, ...) | full history; without `return_full_data` the server sends a preview with the two newest values, which is all `--sma200` needs |
| `GLOBAL_QUOTE` | latest price |

## Premium on this key

| Endpoint | Reply |
|---|---|
| `TIME_SERIES_DAILY_ADJUSTED` | premium endpoint |
| `TIME_SERIES_DAILY` with `outputsize=full` | premium feature (the parameter, not the endpoint) |
| `HISTORICAL_OPTIONS`, with or without `expiration` | premium endpoint |
| `REALTIME_OPTIONS` | not tested; assume premium |
| `EARNINGS_CALENDAR` | returned a one-line "Information" message split into CSV cells, not data; treat as unavailable |

## How results arrive

Three shapes, decided by size and by who truncates:

1. **Inline**: the tool result is `{"result": "<csv>"}` in the reply.
   Compact daily (100 rows) always arrives this way, and so does a
   weekly series for a stock with a short listing history (UBER, 384
   rows). It must be written to a file by hand before a script can read
   it; a heredoc of 100 rows is fine, one of 400 rows costs a lot of
   context. Both skills' scripts accept the `{"result": ...}` wrapper
   verbatim, so paste the whole thing rather than unwrapping it.
2. **Saved by the harness**: above roughly 32k tokens the reply is a
   notice with a path under the Claude projects folder's
   `tool-results/`; the file holds `{"result": "<csv>"}` or, for an
   indicator endpoint, the raw Alpha Vantage JSON. Pass the path
   straight to `--weekly` or `--sma200-file`.
3. **Previewed by the server**: an indicator endpoint without
   `return_full_data=true` returns `{"preview": true, "sample_data":
   "<json>", ...}` with the two newest values. Read the value from the
   preview, or save the reply to a file and pass `--sma200-file`.

Which of 2 and 3 you get for the same endpoint varies with the
symbol's history length, so the scripts accept all of them.

## Reuse before refetching

A daily or weekly file fetched earlier the same day is still current
until the next close. Before spending a call, look in the scratchpad
for `<SYMBOL>_daily.csv`, the weekly file, and `<SYMBOL>_ta.json`, and
say in the report which were reused. The analysis-option run on AAPL
cost zero calls for this reason.
