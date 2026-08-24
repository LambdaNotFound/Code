# Design: Fleet-Wide API Rate Limiter

Round 0 draft. Requirements: `docs/research/rate-limiter/brief.md`.

## Answer up front

Enforce the limit in **HTTP middleware on every API server**, backed by a
**single shared Redis** executing **one atomic Lua script per request** that
maintains a **sliding-window log** (a sorted set of admitted-request
timestamps) per client. The Redis call carries a hard deadline below the
latency budget and **fails open** into a per-server fallback limiter, so a
dead store degrades protection instead of the API.

Sliding-window log is not a default choice: it is the only common algorithm
that produces goal 1's *exactly 50 rejections* under every arrival pattern
the goal's wording permits. Token bucket and fixed window both produce 0
rejections for the same 150 requests under a legal reading of the same test
(measured, [Algorithm choice](#algorithm-choice-derived-low-level)).

## Requirements traced to mechanism

Every element below names the requirement that forces it. Anything with no
requirement in this table is not in the design.

| Requirement (brief) | Mechanism | Section |
|---|---|---|
| G1 exactly 50 rejections of 150/min at 100/min | sliding-window log; deny path consumes no quota | [Algorithm](#algorithm-choice-derived-low-level) |
| G2 limit holds fleet-wide | one shared Redis key per client; read-modify-write inside one atomic Lua script | [Consistency](#consistency-and-atomicity-boundary-high-level) |
| G3 429 + `Retry-After` | script returns time until the oldest in-window entry expires | [Script](#the-script-low-level) |
| G4 store failure must not take down the API | per-call deadline, fail-open, circuit breaker, local fallback limiter | [Degradation](#failure-domains-and-degradation-high-level) |
| C1 stateless fleet, any server any request | no per-server affinity; all shared state in Redis, keyed by client | [Components](#components-and-responsibilities-high-level) |
| C2 added p99 < 5 ms | exactly one round trip; deadline `< 5 ms` covering pool acquisition; breaker removes the round trip when the store is dead | [Latency budget](#latency-budget-high-level) |
| C3 clients identified by API key | limiter keyed by the authenticated client ID, placed after authentication | [Placement](#placement-in-the-request-path-high-level) |

---

## Placement in the request path [high-level]

```
client -> LB -> [ API server ]
                  TLS -> routing -> auth -> RATE LIMITER -> handler -> backend
                                                 |
                                                 +--(1 EVALSHA, deadline 3ms)--> Redis
```

Two placement decisions, each with the requirement that forces it and the
kill reason for the simpler alternative.

**In the API server's middleware chain, not in a separate limiter service.**
Forced by C2: a dedicated limiter service adds a network hop *in addition to*
its own store access, spending latency budget to buy language-independence
that a single-language Go fleet does not need. Kill reason for the service:
it adds a hop and a deployment without satisfying any requirement the
middleware misses.

**After authentication, not before.** Forced by a bypass that otherwise
defeats the whole design: if the limiter keys on the raw API-key string
before that key is validated, an attacker rotates fabricated keys and each
one receives its own fresh 100/min bucket — unlimited aggregate throughput,
plus one new Redis key per fabricated key (TTL-bounded, but unbounded in
rate). Keying on the *authenticated* client ID means unauthenticated floods
are rejected by auth (cheaper than a Redis round trip) and the limiter's key
space is bounded by the customer count. (**inferred** from C3 plus the key
construction; no auth code exists in this repo to cite.)

Consequence, stated as a cost rather than hidden: this limiter protects
*backend capacity* — which is what the brief's Problem names — but not the
API servers' own front door. A flood still pays for TLS, routing, auth, and a
goroutine per request before being rejected. Protecting the front door is
gateway or network-layer work, and network-layer DDoS is an explicit non-goal.

## Components and responsibilities [high-level]

| Component | Responsibility | State | Failure behavior |
|---|---|---|---|
| `Middleware` | extract client ID, call limiter, write 429 + `Retry-After`, emit metrics | none | cannot fail; limiter returns a decision, never an error |
| `RedisLimiter` | one atomic decision per request within a deadline | none locally | on error/timeout: fail open, record, consult breaker |
| Redis | the only fleet-wide state: per-client sorted set of admitted timestamps | authoritative, ephemeral (TTL = window) | unavailable -> limiter degrades |
| `denyCache` | per-server memo of "this client is denied until T" | per-server, bounded LRU | miss costs one Redis call; never over-admits |
| `breaker` | stop calling a dead store | per-server, atomic | open -> local fallback |
| `localLimiter` | per-server token bucket used only while the breaker is open | per-server, bounded LRU | none |

The last three exist only because a requirement kills the version without
them; each kill reason is stated where the component is specified. The
minimal design that satisfies G1, G2, G3 is `Middleware` + `RedisLimiter` +
Redis, and everything else is degradation machinery forced by G4 and C2.

## Data model and ownership [high-level]

One Redis key per client, owned exclusively by the limiter:

```
key:    rl:<clientID>           (string key, one per active client)
type:   sorted set
member: <requestID>             unique per request (ULID/UUIDv7 from the API server)
score:  <admission time in microseconds since epoch, from Redis TIME>
TTL:    window (60_000 ms), refreshed on each admission
```

- The set contains **admitted requests only**. Denied requests write nothing;
  this is what makes 150 requests at limit 100 produce exactly 50 denials
  rather than a feedback loop.
- Cardinality is bounded by `limit` (100) per key, because an admission
  happens only when cardinality is below the limit.
- No durability requirement. State is reconstructible by waiting one window,
  so AOF/RDB persistence is unnecessary (Redis 7.2 default is `appendonly no`
  — **observed**, redis.conf). A Redis restart grants every client a fresh
  window; bounded and acceptable because long-window accounting and billing
  are explicit non-goals.
- Keys carry the client ID, never the API key itself. API keys in key names
  leak through `MONITOR`, `SLOWLOG`, and keyspace dumps. If no stable client
  ID exists post-auth, use `sha256(apiKey)[:16]` hex.

## Consistency and atomicity boundary [high-level]

The decision "is this client under the limit, and if so record this request"
is a read-modify-write that must be serializable across the fleet, or two
servers both observe count 99 and both admit (G2 fails; every extra admission
costs one rejection, so G1's count comes out at 49 or lower).

The whole read-modify-write is one Lua script:

> "Redis guarantees the script's atomic execution. While executing the
> script, all server activities are blocked during its entire runtime."
> — Redis docs, *Scripting with Lua* (**observed**, accessed 2026-08-24)

That single sentence is the entire consistency argument: the sequence
evict -> count -> admit is indivisible, so at most `limit` admissions exist
in any trailing window on that key. No `WATCH`/`MULTI` retry loop, no
distributed lock, no read-then-write from the application. It also fixes the
transaction boundary at exactly one key, which is what keeps the design
Redis-Cluster-ready with no change (see [Scaling](#capacity-and-scaling-shape-high-level)).

Consistency is *not* maintained across a Redis failover: replication is
asynchronous ("Redis uses by default asynchronous replication" — **observed**,
Redis docs) and "acknowledged writes can still be lost during a failover".
A promoted replica may be missing the most recent admissions, letting an
affected client briefly exceed the limit by up to the number of lost writes.
Accepted: the brief asks for bounding, not accounting, and quotas/billing are
non-goals. `WAIT` is rejected as a fix — it "does not turn a set of Redis
instances into a CP system" (**observed**) and would put a replication
round trip inside the 5 ms budget.

## End-to-end data flow [high-level]

Common path (breaker closed, no cached deny) — one round trip:

1. Middleware reads the authenticated client ID via `ClientIDFunc`. Absent
   (`ok=false`) -> pass through unlimited; the limiter never invents an
   identity, and auth is what rejects unauthenticated traffic.
2. `denyCache` lookup. Hit and not expired -> deny locally, `Retry-After` =
   `deadline - now`. No Redis call.
3. `breaker` check (one atomic load). Open -> local fallback limiter decides.
4. `EVALSHA` the script with `KEYS=[rl:<clientID>]`,
   `ARGV=[windowMicros, limit, requestID]`, on a context whose deadline is
   `now + 3 ms` and which covers connection-pool acquisition.
5. Script returns `{allowed, remaining, retryAfterMicros}`.
   - allowed -> `next.ServeHTTP`.
   - denied -> record in `denyCache`, respond `429` + `Retry-After`.
6. Error or deadline exceeded -> record a failure with the breaker, count a
   fail-open, **allow** the request.

Degraded path (breaker open): steps 4-5 are skipped entirely, so the store's
death costs microseconds, not the timeout, per request.

## Failure domains and degradation [high-level]

| Failure | Detected by | Behavior | Protection during |
|---|---|---|---|
| Redis unreachable / dead | dial error, then breaker after N consecutive failures | fail open; breaker opens; local fallback limiter takes over | per-server limit (worst case `limit x servers` fleet-wide) |
| Redis slow (GC, big script, saturated) | 3 ms deadline | fail open per request, breaker opens if sustained | same as above |
| Redis at `maxmemory` with `noeviction` | script returns OOM error | fail open (error path) | same as above |
| Redis failover | NOSCRIPT after cache flush; lost writes | go-redis `Script.Run` falls back to `EVAL` automatically (**observed**, go-redis `script.go`); brief over-admission window | full, with a bounded overshoot |
| One cluster shard down (if sharded) | errors on that shard only | clients hashed to that shard degrade; others unaffected | partial |
| Bad limit config (limit set to 0) | 429 rate metric spikes | operator flips mode to `shadow` | none, by choice |

G4's requirement — "limiter failure must not take down the API" — is met by
construction at three levels: the limiter's Go API cannot return an error to
the middleware (see [`Limiter`](#interfaces-low-level)), every store call has a
deadline strictly below the latency budget, and the breaker removes the store
from the path entirely once it is proven dead.

**Why a local fallback limiter rather than plain fail-open.** Plain fail-open
satisfies G4 as literally written. It is rejected because it leaves the API in
exactly the pre-project state — one client able to exhaust backend capacity —
precisely during a store outage, and store outages correlate with the overload
episodes the limiter exists for. Kill reason for the simpler version: it
restores the brief's Problem statement in full at the worst moment.
(**inferred**; this is the one component justified by the Problem rather than
by a numbered goal, and the honest reading is that G4 alone does not demand it.)

**Why the fallback uses the full limit per server, not `limit / fleetSize`.**
Dividing requires each server to know the live fleet size — a service-discovery
dependency, and a stale value under-limits during scale-out and over-limits
during scale-in (rejecting legitimate traffic while the store is already
broken). Kill reason for the divided version: it buys tighter bounds during a
degraded window by adding a dependency that can itself be degraded. Worst case
accepted: `limit x servers` (300/min at 3 servers) while the store is down.

## Capacity and scaling shape [high-level]

- **Redis ops/sec = admitted request rate + one probe per freed slot per
  client per server.** The `denyCache` is what makes the second term small:
  without it, a client flooding at F requests/sec generates F ops/sec against
  a *single key on a single shard* — the limiter becomes the cheapest thing in
  the system to DoS. With it, a flooding client generates at most one Redis
  call per server per freed slot (~1 per 600 ms at 100/min). Kill reason for
  omitting the deny cache: it converts an attack on the API into an attack on
  the limiter's hottest key. (**inferred** from the algorithm; the load
  reduction is arithmetic, not measured.)
- **Per-op cost** is bounded by the limit: at most 100 members per key, under
  the Redis 7.2 default `zset-max-listpack-entries 128` (**observed**,
  redis.conf), so every key stays in the compact listpack encoding. Script
  work per call is `ZREMRANGEBYSCORE` + `ZCARD` + `ZADD` + `PEXPIRE`, all on a
  <=100-element structure.
- **Memory** = active clients x (key overhead + up to 100 packed entries).
  I have no byte figure and will not invent one; the measurement that produces
  it is `MEMORY USAGE rl:<client>` on a key filled to 100 entries, multiplied
  by peak concurrent active clients.
- **Scale-out path**: Redis Cluster. The script touches exactly one key and
  passes it in `KEYS`, which is what the EVAL docs require — "all names of
  keys that a script accesses must be explicitly provided as input key
  arguments" (**observed**) — so keys shard by CRC16 with no hash tags and no
  code change. Refinement needed at that point: breaker state must become
  per-shard, or one dead shard opens the breaker for all clients.
- **Single instance until measured otherwise.** Start with one primary plus a
  replica. The measurement that decides whether that suffices is
  `redis-benchmark` at the fleet's peak request rate against a key of 100
  entries, compared with the peak rate.

## Latency budget [high-level]

C2 (added p99 < 5 ms) is met by construction rather than by hope:

```
added latency per request <= localWork + redisDeadline
                          =  O(microseconds) + 3 ms   < 5 ms
```

This holds only if the deadline covers **connection-pool acquisition**, not
just the round trip. With go-redis, an exhausted pool makes the caller wait;
if that wait sits outside the deadline the bound is void. The context passed
to `Script.Run` therefore carries the 3 ms deadline and the pool is sized so
that acquisition is not the common blocking point.

Three numbers here are unmeasured and load-bearing: typical round-trip time
(decides whether 3 ms is generous or marginal), the fail-open rate at a 3 ms
deadline (a deadline below the store's natural p99 silently disables
limiting), and the actual added p99. The measurement: run the fleet's load
test twice at identical offered load, limiter enabled and disabled, and take
`p99_with - p99_without`; then histogram the store call and set the deadline
above its p999. 3 ms is a starting value, not a finding.

## Operational surface [high-level]

- **Config** (process flags/env, no config service — nothing in the brief
  needs dynamic limits, and a Redis-hosted config would make the kill switch
  depend on the component it must survive): `limit` (100), `window` (60 s),
  `redisAddr`, `deadline` (3 ms), `mode`.
- **Mode** is the rollout and rollback control: `off` (no store call, no
  enforcement) | `shadow` (call the store, record the decision, always allow)
  | `enforce`. Deploy in `shadow` first; the denial rate observed there is the
  real-traffic estimate of how many customers `enforce` will start rejecting,
  which is the number that decides whether 100/min is the right limit at all.
  Rollback is a mode flip.
- **Metrics**: `ratelimit_decisions_total{result=allow|deny, source=redis|deny_cache|fail_open|local}`,
  `ratelimit_backend_seconds` histogram, `ratelimit_breaker_open` gauge.
- **Alerts**: any sustained `source=fail_open` (limiting is silently off);
  breaker open > 1 minute; denial rate step-change after a deploy.
- **Runbook line that must exist**: fail-open means unlimited. An alert on
  `fail_open > 0` for a minute is the difference between "we have a limiter"
  and "we had one".

---

## Algorithm choice, derived [low-level]

G1 is a test, not a preference: 150 requests in a minute against 100/min must
yield **exactly 50** rejections. The brief does not say how those 150 requests
are spaced, so the design must satisfy the test under every legal spacing.

I simulated four candidate algorithms across arrival patterns — a burst, and
150 requests evenly spread over the minute — each at several phases relative
to an aligned window boundary (probe:
`/tmp/claude-0/-home-user-Code/f70f1ff3-2f85-578d-9b8f-cf95aa35610e/scratchpad/sim.py`,
run 2026-08-24). Rejection counts (**observed**):

| algorithm | burst@0 | burst@37.3 | spread@0 | spread@37.3 | spread@59.9 |
|---|---|---|---|---|---|
| fixed window (aligned, 100/min) | 50 | 50 | 50 | **0** | **49** |
| token bucket (cap 100, refill 100/min) | 50 | 50 | **0** | **0** | **0** |
| **sliding window log** | **50** | **50** | **50** | **50** | **50** |
| sliding window counter (weighted) | 50 | 50 | 50 | **14** | **49** |

The two "obvious" answers fail:

- **Token bucket fails the spread test outright.** 150 requests over 60 s is
  2.5 req/s; the bucket refills at 100/60 = 1.67 tokens/s, so a bucket
  starting full never empties within the minute and rejects *nothing*. Kill
  reason: it enforces a rate with burst allowance, and G1's test measures a
  count in a window. Only sizing the bucket smaller than the limit would help,
  and that breaks the "100 requests per minute" promise in the other
  direction.
- **Fixed window fails on phase.** When the test's minute straddles a window
  boundary, the counter resets mid-test.
- **Sliding window counter is an approximation** and its error is exactly what
  G1 forbids.
- **Sliding window log is exact by construction** — it answers "how many
  admissions in the last 60 s" directly — at a memory cost of `limit` entries
  per active client, which the [capacity](#capacity-and-scaling-shape-high-level)
  section bounds and the 128-entry listpack threshold keeps cheap.

Leaky bucket is absent from the table deliberately: this repo's prior
investigation established that leaky-bucket-as-a-meter is "mathematically
equivalent to token bucket: same admit/reject decisions, described from
opposite ends" (**observed**, `docs/research/leaky-bucket.md:20-24`), and its
shaping variant delays traffic rather than rejecting it, which G3's 429
forbids.

**A caveat the verification plan must respect** (**observed**, same probe):
sliding-window log yields exactly 50 only if the 150 requests all land inside
one 60 s span. Stretch the test to 60.5 s and it yields 49; to 62 s, 46; to
65 s, 39. That is correct limiter behavior — the earliest entries have aged
out — but it makes a wall-clock test flaky. The test must therefore drive an
injected clock, or issue its 150 requests well inside the window (see
[Verification](#verification-plan)).

## The script [low-level]

```lua
-- KEYS[1] = rl:<clientID>
-- ARGV[1] = window in microseconds   (60000000)
-- ARGV[2] = limit                    (100)
-- ARGV[3] = unique request id
-- returns  {allowed(0|1), remaining, retry_after_micros}
local window = tonumber(ARGV[1])
local limit  = tonumber(ARGV[2])

local t   = redis.call('TIME')
local now = tonumber(t[1]) * 1000000 + tonumber(t[2])

redis.call('ZREMRANGEBYSCORE', KEYS[1], '-inf', now - window)
local count = redis.call('ZCARD', KEYS[1])

if count < limit then
  redis.call('ZADD', KEYS[1], 'NX', now, ARGV[3])
  redis.call('PEXPIRE', KEYS[1], math.ceil(window / 1000))
  return {1, limit - count - 1, 0}
end

local oldest = redis.call('ZRANGE', KEYS[1], 0, 0, 'WITHSCORES')
local retry  = (tonumber(oldest[2]) + window) - now
if retry < 0 then retry = 0 end
return {0, 0, retry}
```

Line-by-line justification of the parts that are not obvious:

- **`redis.call('TIME')` rather than a timestamp from the caller.** Removes
  clock skew across the fleet as a correctness concern: all scores come from
  one clock. Legal since Redis 5.0 — with effects replication, "the
  restrictions on non-deterministic functions are removed. You can, for
  example, use the `TIME` or `SRANDMEMBER` commands inside your scripts freely
  at any place" (**observed**, Redis scripting docs, accessed 2026-08-24). On
  Redis 4 or older this script is invalid; that is the design's minimum
  version. Residual risk: a failover changes the clock source, and an entry
  scored by a fast clock can linger up to one window.
- **Microsecond scores.** Sorted-set scores are IEEE-754 doubles, exact for
  integers up to 2^53 = 9007199254740992 (**observed**, ZADD docs); epoch
  microseconds are ~1.8e15, comfortably inside. Millisecond precision would
  also fit but coarsens `Retry-After` for no gain.
- **`ZREMRANGEBYSCORE ... '-inf' (now - window)`** is inclusive of the bound,
  so survivors satisfy `score > now - window` — a request exactly one window
  old does *not* count. This matches the exclusive-boundary convention already
  documented in this repo for windowed counting
  (`golang_containers/hit_counter.go:15-18`, **observed**) and is what makes
  a client's retry at exactly `Retry-After` succeed rather than miss by a tick.
- **`ZADD ... NX`**: "Only add new elements. Don't update already existing
  elements" (**observed**, ZADD docs). Combined with a caller-generated unique
  request ID this makes the call **idempotent**: a client-side retry of a
  request whose response was lost re-presents the same member and consumes no
  second slot, and does not shift its own timestamp forward. Without a unique
  member, two requests sharing a member would collapse into one entry and the
  limiter would over-admit.
- **`PEXPIRE` on the admit path only.** TTL = one window from the most recent
  admission, which is exactly when the last entry leaves the window and the key
  would be empty anyway. Idle clients cost nothing; memory is bounded by
  *active* clients, not by customer count.
- **Deny path reads the oldest entry** to compute when a slot frees:
  `retry = oldest.score + window - now`. This is the exact earliest instant the
  request would succeed, which is what `Retry-After` is supposed to mean, and
  it costs one extra command only when denying.
- **Complexity**: `ZREMRANGEBYSCORE` O(log N + M), `ZCARD` O(1), `ZADD`
  O(log N), with N <= 100 — under listpack encoding these are short linear
  scans over a packed buffer. Bounded work matters more than usual here
  because the script blocks the whole server while it runs (**observed**,
  quoted above).

## Interfaces [low-level]

New package `ratelimit/` at repo root, alongside the existing non-LeetCode
top-level packages `golang_containers/` and `golang_concurrency/`
(**observed**). All files below are `new file`.

```go
// ratelimit/limiter.go  (new file)

// Source records which tier produced a Decision. It exists so that
// "we are limiting" and "we are failing open and limiting nothing"
// are distinguishable in metrics rather than both looking like 200s.
type Source uint8

const (
    SourceRedis     Source = iota // authoritative, fleet-wide
    SourceDenyCache               // local memo of a previous authoritative deny
    SourceLocal                   // breaker open: per-server fallback
    SourceFailOpen                // store call failed; request allowed unlimited
)

type Decision struct {
    Allowed    bool
    Remaining  int           // best-effort; 0 when !Allowed
    RetryAfter time.Duration // > 0 only when !Allowed
    Source     Source
}

// Limiter returns a Decision and never an error: fail-open is the
// limiter's own responsibility, not the caller's. This is goal 4
// expressed in the type system — a caller cannot forget to handle a
// store outage, because it is never told about one.
type Limiter interface {
    Allow(ctx context.Context, clientID string) Decision
}
```

```go
// ratelimit/middleware.go  (new file)

// ClientIDFunc extracts the authenticated client identity. Returning
// ok=false means "not this middleware's business" and the request
// passes through untouched; the limiter never fabricates an identity.
type ClientIDFunc func(*http.Request) (id string, ok bool)

func Middleware(l Limiter, id ClientIDFunc, m Metrics) func(http.Handler) http.Handler
```

On deny the middleware writes exactly:

```
HTTP/1.1 429 Too Many Requests
Retry-After: <ceil(RetryAfter in seconds), minimum 1>
Content-Type: application/json
{"error":"rate_limit_exceeded","retry_after":<same integer>}
```

`Retry-After` is emitted in `delay-seconds` form; the minimum of 1 exists
because a 0 is ambiguous and sub-second precision is not expressible in that
form. (Status 429 is defined by RFC 6585 §4 and `Retry-After` by RFC 9110
§10.2.3 — **assumed**, not verified: rfc-editor.org, datatracker.ietf.org,
httpwg.org, iana.org and MDN are all blocked by this environment's egress
proxy, checked 2026-08-24. What would close it: fetch either RFC from any
reachable mirror and confirm the `delay-seconds` ABNF and 429's definition.)

```go
// ratelimit/redis_limiter.go  (new file)

type RedisLimiter struct {
    rdb      redis.Scripter   // *redis.Client or *redis.ClusterClient
    script   *redis.Script    // go-redis Script; Run() = EVALSHA with EVAL fallback
    limit    int
    window   time.Duration
    deadline time.Duration    // 3ms, must satisfy deadline + localWork < 5ms
    deny     *denyCache
    breaker  *breaker
    local    *localLimiter
    newID    func() string    // ULID/UUIDv7; unique per request
    metrics  Metrics
}

func (l *RedisLimiter) Allow(ctx context.Context, clientID string) Decision
```

`Allow`'s control flow is the [data flow](#end-to-end-data-flow-high-level)
above, with one detail that matters: the context handed to `script.Run` is
`context.WithTimeout(ctx, l.deadline)` derived from the *request* context, so
client disconnects abort the call too, and pool-acquisition time is inside the
deadline.

Dependencies this adds to `go.mod` (**observed**: the module currently
requires only `github.com/stretchr/testify v1.8.4`, `go.mod:5`):
`github.com/redis/go-redis/v9`, and `github.com/alicebob/miniredis/v2` for
tests.

## denyCache [low-level]

```go
// ratelimit/denycache.go  (new file)
type denyCache struct {
    mu    sync.Mutex
    lru   *list.List                    // most-recent first
    index map[string]*list.Element      // clientID -> element{clientID, until time.Time}
    cap   int                           // fixed, e.g. 4096
}

func (c *denyCache) until(clientID string, now time.Time) (time.Time, bool)
func (c *denyCache) set(clientID string, until time.Time)
```

**Invariant**: the cache may deny a request the store would have allowed
(bounded by the cached deadline, which came from the store), but it can never
allow one the store would have denied — it only ever produces denials. G1's
count is therefore unchanged by its presence; only the *location* of the
denial moves. Entries are dropped lazily on lookup when expired and by LRU
when at capacity.

Fixed capacity rather than an unbounded map: memory becomes O(1) per server
instead of O(distinct denied clients), and the only penalty for an eviction is
one extra Redis call. O(1) per operation; one mutex, held for pointer moves
only — if profiling shows contention, shard by `fnv(clientID) % 16`.

## breaker [low-level]

State machine — the point is that an *open* breaker costs a single atomic
load, so a dead store adds microseconds instead of 3 ms per request:

```
        consecutive failures >= F (=5)
CLOSED ------------------------------> OPEN(until = now + cooldown, 1s)
  ^                                       |
  |                                       | now >= until: first goroutine to
  |  probe succeeds                       | win a CAS becomes the probe
  +------------------- HALF_OPEN <--------+
              probe fails -> OPEN(now + cooldown)
```

```go
// ratelimit/breaker.go  (new file)
type breaker struct {
    failures  atomic.Int64  // consecutive
    openUntil atomic.Int64  // unix micros; 0 = closed
    cooldown  time.Duration
    threshold int64
}

func (b *breaker) allowCall(nowMicros int64) bool // one atomic load on the hot path
func (b *breaker) onSuccess()                     // failures.Store(0); openUntil.Store(0)
func (b *breaker) onFailure(nowMicros int64)      // increment; trip at threshold
```

Concurrency model: no locks and no background goroutine. `allowCall` is a
single `openUntil.Load()` compared against now; the half-open probe is claimed
by `openUntil.CompareAndSwap(old, now+cooldown)`, so exactly one in-flight
request probes and the rest take the fallback path. Deliberately *not* a
sliding error-rate breaker: a consecutive-failure counter needs one atomic and
no windowing, and the failure it must catch (store gone) is not subtle.

## localLimiter [low-level]

Per-server token bucket per client, consulted only when the breaker is open.

```go
// ratelimit/local.go  (new file)
type localLimiter struct {
    mu      sync.Mutex
    buckets map[string]*bucket   // bounded LRU, same shape as denyCache
    limit   float64              // tokens = limit
    refill  float64              // tokens per second = limit / window.Seconds()
}
type bucket struct { tokens float64; last time.Time }
```

Refill must be `float64`: 100/min is 1.667 tokens/s. This is why the repo's
existing `TokenBucket` cannot be reused as written — its `refillRate int //
tokens per second` (**observed**, `golang/design/rate_limiter.go:57`) cannot
express any rate below 60/min. That type also lives in package `design`, a
LeetCode answers package, which production middleware should not import.

Chosen deliberately over reusing the sliding-window log locally: during
degradation, exactness has no meaning (the fleet is not coordinating anyway),
and a token bucket is O(1) time and O(1) memory per client instead of O(limit).

## Error handling, retries, idempotency [low-level]

| Condition | Handling |
|---|---|
| `context.DeadlineExceeded` | fail open, `SourceFailOpen`, `breaker.onFailure` |
| dial / connection error | same |
| `NOSCRIPT` (post-failover cache flush) | handled inside go-redis: `Run` "optimistically uses EVALSHA... If script does not exist it is retried using EVAL" (**observed**, go-redis `script.go`) |
| `OOM command not allowed` (maxmemory) | fail open + failure; alert fires on fail-open rate |
| script returns a malformed reply | fail open + failure; a parse error is a bug, not a limit |
| timeout after the script actually ran | request allowed although a slot was consumed (or a denial lost). Bounded by the timeout rate; **no retry** — a retry would double the latency the deadline exists to bound |

Retries are deliberately absent from the hot path. The unique request ID
exists so that a *client-level* retry is idempotent (`ZADD NX`), not so that
the limiter can retry itself.

Redis configuration this design depends on: a limiter-dedicated instance or
logical DB, `maxmemory` sized from the measured per-key cost, and
`maxmemory-policy volatile-ttl` — every key here carries a TTL, so evicting
the nearest-to-expiry key is the least damaging choice. The 7.2 default is
`noeviction` (**observed**, redis.conf), under which hitting `maxmemory`
turns *all* limiting off at once via the error path instead of shedding the
coldest buckets.

---

## Verification plan

Each goal maps to a test, and the third column names what makes the test
deterministic rather than time-flaky.

| Goal | Test | Determinism |
|---|---|---|
| G1 exactly 50 rejections | 150 requests, limit 100/min, against `miniredis` | `miniredis.SetTime()` drives `TIME`, so the window is controlled, not raced. miniredis v2 supports `EVAL`/`EVALSHA`/`SCRIPT LOAD` and "TIME -- returns time.Now() or value set by SetTime()" (**observed**, miniredis README, accessed 2026-08-24) |
| G1 (real store) | same count against a real Redis | issue all 150 within a fraction of the window; a run stretching past 60 s legitimately yields 49 or fewer (see the [caveat](#algorithm-choice-derived-low-level)) |
| G2 fleet-wide | 3 limiter instances sharing one store, interleaved round-robin | same injected clock; assert 50 rejections in aggregate and that no instance's count matters |
| G2 concurrency | N goroutines racing at count = limit-1 | assert exactly one admission; this is the test that would catch a non-atomic implementation |
| G3 429 + Retry-After | assert status, integer `Retry-After` >= 1, and that a retry at exactly that offset succeeds | the exclusive eviction boundary is what makes the retry succeed rather than miss |
| G4 store failure | kill the store mid-test | assert 200s continue, `SourceFailOpen` then `SourceLocal` counters rise, breaker opens after F failures, and added latency drops back to microseconds once open |
| C2 latency | load test with `mode=off` vs `mode=enforce` at identical offered load | report `p99_with - p99_without`; also histogram the store call and confirm the deadline sits above its p999 |

If implementation is requested, the risk order is: (1) spike the C2
measurement against a real store — every latency claim here rests on it and a
bad result changes the topology, not the code; (2) script + `RedisLimiter` +
G1/G2 tests; (3) middleware + G3; (4) breaker, deny cache, local fallback +
G4; (5) metrics, shadow mode, runbook.

## Alternatives rejected

| Alternative | Real advantage | Kill reason |
|---|---|---|
| Token bucket in Redis | O(1) memory per client, smooth pacing, industry default | Measured 0 rejections on G1's spread test (table above) |
| Fixed-window counter | Cheapest possible: one `INCR` + `EXPIRE` | Phase-dependent; 0 rejections when the test straddles a boundary; allows 2x limit across a boundary |
| Sliding-window counter | O(1) memory, close to exact | Approximate by construction; G1 forbids approximation |
| Leaky bucket (queueing) | Smooths bursts into a constant output rate | Delays instead of rejecting; G3 requires a 429. As a meter it is token bucket in other words (`docs/research/leaky-bucket.md:20-24`) |
| Per-server local limits at `limit/servers` | Zero added latency, no shared store, no new failure domain | Violates G2 under uneven LB distribution and during scale events; needs live fleet size |
| Gossip / CRDT approximate counting | No central store; survives partitions | Approximate (fails G1) and O(n^2) chatter |
| Dedicated rate-limit service (gRPC) | Language-independent, central policy, own scaling | Adds a hop inside a 5 ms budget while still needing the same shared store; no requirement needs it |
| Enforce at the LB / API gateway | Strictly better placement — rejects before TLS/auth/goroutine cost, protects the API servers themselves | Not chosen because the brief names no gateway. **If a gateway with shared-state limiting already exists, revisit this before building** (see open questions) |
| Async / fire-and-forget accounting | Removes the round trip from the critical path entirely | Races admit more than the limit; G1's exactness dies |
| `WATCH`/`MULTI` instead of Lua | No scripting dependency | Optimistic-retry loop under contention adds round trips and unbounded latency; the hot key is precisely where contention lives |

## Tradeoffs this design accepts

- **Memory scales with the limit.** 100 entries per active client versus 2
  numbers for a token bucket. Bounded and cheap at 100/min; a 10,000/min limit
  would push past the 128-entry listpack threshold and this choice should be
  revisited then.
- **Over-admission across failover.** Async replication can lose recent
  admissions; a client can briefly exceed 100/min.
- **Over-admission on timeout.** A request whose script ran but whose reply
  was late is allowed. Bounded by the timeout rate, visible in metrics.
- **Under-limiting while the store is down.** Up to `limit x servers`
  fleet-wide. The alternative (dividing the limit) needs fleet-size discovery.
- **The API servers' own front door is unprotected.** A flood still costs TLS,
  auth, and a goroutine per rejected request.
- **Slight over-denial from the deny cache.** A client denied at T is denied
  locally until the store-computed free time even if capacity frees earlier by
  another path. There is no other path today, so the current error is zero,
  but it is a real coupling to the algorithm.
- **A Redis restart grants everyone a fresh window.** Accepted: persistence
  would trade a 60-second accounting inaccuracy for fsync latency in the hot
  path.

## What would invalidate this design

1. **The store's p99 is not comfortably under 3 ms** (unverified). Then the
   deadline either eats the latency budget or silently disables limiting via
   fail-open. Response: move the store closer (same AZ), or accept a
   larger budget, or reconsider gateway placement.
2. **A gateway with shared-state rate limiting already fronts the fleet**
   (unknown). Then middleware is the second-best placement.
3. **The real constraint is concurrency, not rate.** This repo's prior
   research flags the distinction: "is the actual constraint a *rate* ... or a
   *concurrency* limit ... The two only coincide when request duration is
   constant" (**observed**, `docs/research/leaky-bucket.md:37-42`). The brief
   fixes the rate form (100 requests/minute), so this design answers what was
   asked; if backend exhaustion is really about in-flight work, a rate limiter
   bounds the wrong quantity.
4. **Redis is not actually available**, or is version 4 or older. The `TIME`
   call inside the script requires Redis >= 5.
5. **Clients are not uniquely identified after auth** — if many customers
   share one API key, the limit lands on the wrong subject.

## Open questions

| # | Question | What closes it |
|---|---|---|
| 1 | Is there already an API gateway / LB with rate-limiting support in front of the fleet? | Ask the owner; inspect the LB config. Changes the placement decision, not the algorithm |
| 2 | What is the actual Redis round-trip p99 from an API server? | `redis-benchmark -n 100000` plus an in-app histogram; sets the deadline |
| 3 | Bytes per 100-entry key, and peak concurrent active clients? | `MEMORY USAGE rl:<client>` on a full key x peak actives; sizes `maxmemory` |
| 4 | Which header carries the API key, and what is the post-auth client ID? | Ask the owner. Default assumed: `Authorization: Bearer <key>`, stable customer ID from auth |
| 5 | Exact `Retry-After` / 429 wording per RFC | Fetch RFC 6585 §4 and RFC 9110 §10.2.3 from a reachable mirror; all standards hosts were egress-blocked on 2026-08-24 |
| 6 | Is the fleet single-region? | Ask the owner. A cross-region store makes the 5 ms budget unreachable and forces per-region limits |
| 7 | Package placement: `ratelimit/` at repo root vs. under `golang/`? | Owner's call; `golang/` is documented as LeetCode categories (`CLAUDE.md`), which this is not |

**Gaps in the brief** (it is `signed-off` and conforms to `docs/brief-spec.md`;
these are ambiguities rather than missing sections): G1 does not specify the
arrival pattern of its 150 requests, which is the single most consequential
unstated requirement here — it decides the algorithm (see the table). G4 does
not say what protection level must survive a store outage, only that requests
must serve; the local fallback tier answers the Problem statement rather than
the goal, and could be cut if the owner prefers pure fail-open.

## Finding: the brief's "no existing rate-limiting code" is not accurate

The brief states "No existing rate-limiting code in this repo"
(`docs/research/rate-limiter/brief.md:52-53`). Observed, there is:

- `golang/design/rate_limiter.go:17` `FixedWindowLimiter` and `:54`
  `TokenBucket` — single-process, mutex-guarded.
- `golang_concurrency/leaky_bucket.go:64` `LeakyBucket` — ticker-drained
  bounded queue with `Allow`/`Wait` policies.
- `golang_containers/hit_counter.go` — 300-second sliding-window hit counter.
- `golang/interview/rippling_rate_limiter.go:32` — LC 359 logger limiter,
  whose closing comment names "Follow-up 4: Distributed rate limiting".

The brief's operative conclusion still holds: none is distributed, none is
HTTP-facing, and none satisfies G2. The correction matters in two places —
`golang/design/rate_limiter.go:54` looks like a reusable local fallback but is
not (integer refill rate, LeetCode package), and
`docs/research/leaky-bucket.md` is prior art this design cites rather than
repeats.

## Sources

- Redis, *Scripting with Lua* (eval-intro), `redis/docs` repo `main`, accessed
  2026-08-24 — script atomicity; `TIME` legal in scripts since Redis 5.0;
  script cache cleared on failover.
- Redis, `EVAL` command reference, `redis/redis-doc` `master`, accessed
  2026-08-24 — all accessed keys must be passed in `KEYS` (cluster safety).
- Redis, `ZADD` command reference, same repo, accessed 2026-08-24 — `NX`
  semantics; scores are doubles exact to 2^53.
- Redis, *Replication*, `redis/docs` `main`, accessed 2026-08-24 — async by
  default; acknowledged writes can be lost in failover; `WAIT` is not CP.
- Redis 7.2 `redis.conf`, `redis/redis` tag 7.2, accessed 2026-08-24 —
  `zset-max-listpack-entries 128`, `maxmemory-policy noeviction`,
  `appendonly no`.
- go-redis `script.go`, `redis/go-redis` `master`, accessed 2026-08-24 —
  `Run` = `EVALSHA` with automatic `EVAL` fallback on `NOSCRIPT`.
- miniredis v2 README, `alicebob/miniredis` `master`, accessed 2026-08-24 —
  `EVAL`/`EVALSHA`/`SCRIPT LOAD` supported; `TIME` honors `SetTime()`.
- Local probe `sim.py` (scratchpad, run 2026-08-24) — the algorithm comparison
  table and the test-span sensitivity figures.
- This repo: `go.mod:5`, `golang/design/rate_limiter.go:17,54,57`,
  `golang_concurrency/leaky_bucket.go:64`,
  `golang_containers/hit_counter.go:15-18`,
  `golang/interview/rippling_rate_limiter.go:32`,
  `docs/research/leaky-bucket.md:20-24,37-42`.
- RFC 6585 §4 (429) and RFC 9110 §10.2.3 (`Retry-After`) — **not fetched**;
  every standards host was egress-blocked on 2026-08-24. Claims about them are
  labeled assumed.

## Revision log

- R0: initial design. Algorithm selected by measurement rather than
  convention; degradation tier derived from G4 plus the Problem statement;
  deny cache added after the hot-key load analysis.

## Objection responses

None yet — no review rounds have run.
