# Review ledger: API Rate Limiter

Reviewer: design-bar-raiser. Append-only; one section per round.

## Round 1

**Independent derivation** (from `brief.md` alone, before reading `design.md`):

1. Invariant: the check-and-record for one client must be indivisible fleet-wide,
   or two servers both observe count = limit-1 and both admit; denials must
   consume no quota; per-client state must self-expire so store memory is
   O(active clients), with the expiry set in the same atomic step as the first
   write (a lost expiry permanently limits a client).
2. G1's acceptance test is the binding constraint, not the "100/min" phrasing:
   "exactly 50 rejections of 150" is a count-in-window test, so any algorithm
   whose admitted count depends on window phase (fixed window) or on refill
   during the test (token bucket) can fail it; conversely no algorithm is exact
   unless the test's span is pinned to at most one window.
3. Hard constraints: at most one store round trip per request, under a timeout
   strictly below the 5 ms p99 budget (a timeout equal to the budget spends all
   of it); any store error must degrade to allow (G4); `Retry-After` must be
   derived from the same state that produced the denial, or it lies.
4. Minimum moving parts: middleware placed after authentication and keyed by the
   authenticated client ID (keying on an unvalidated API key gives an attacker a
   fresh bucket per fabricated key); one shared store executing one atomic
   script; bounded timeout with fail-open; static config for limit/window/mode;
   metrics that distinguish "denied" from "failed open". Anything beyond this —
   local fallback, deny cache, breaker — is optional and owes a kill reason for
   the version without it.
5. Capacity: the brief states no QPS and no client count, so every sizing claim
   must be either labeled unmeasured with a named measurement, or absent.

Diff against the proposal: it matches 1-4 by a different route on the algorithm
(it derives sliding-window log from the same test I used to pin the algorithm
class) and adds three degradation tiers, each of which does carry a written kill
reason — so they are not per se unjustified, but their attribution to the
numbered requirements is wrong (R1-6). Objections below come from that diff and
from verification of the low-level claims.

**Objections**

R1-1 | blocking | Idempotency claim is contradicted by the design's own interface: `ZADD NX` "makes the call idempotent: a client-side retry ... consumes no second slot" (design.md:381-384) and "The unique request ID exists so that a *client-level* retry is idempotent" (design.md:593-595), but `newID func() string // ULID/UUIDv7; unique per request` (design.md:476) mints the member inside `RedisLimiter`, so a client retry presents a different member and `NX` can never fire; client retries consume quota. | Source the member from a client-supplied idempotency key and state the new client-side requirement plus its abuse surface (a client reusing one key gets free requests), or delete the idempotency claim from both places and state that `NX` is defensive against ID collision only.

R1-2 | blocking | G3 has no mechanism on the deny path the design itself introduces: `localLimiter` is the enforcer while the breaker is open, G4's own test asserts `SourceLocal` denials occur (design.md:619), and `Decision.RetryAfter` must be "> 0 only when !Allowed" (design.md:423) — but the `localLimiter` section (design.md:557-580) specifies a struct and no method, and nowhere derives a retry-after from a token bucket. The one component that enforces during the failure G4 exists to test is still a box. | Give `localLimiter` its decision signature and the retry-after derivation (time until one token accrues) so every 429 the design can emit has a defined `Retry-After`, or state explicitly that the local tier admits-only and never denies.

R1-3 | should-fix | The failure table's misconfiguration row — "Bad limit config (limit set to 0) | detected by 429 rate metric spikes" (design.md:175) — is contradicted by the script as written: with limit 0 and an empty key, `count < limit` is false, the deny branch runs `ZRANGE` on a nonexistent key, `oldest[2]` is nil, and `(tonumber(nil) + window)` raises a Lua error (Lua 5.1 `tonumber(nil)` returns nil), which reaches Go as an error and takes the fail-open path. Limit 0 therefore produces zero 429s and silently disables limiting; the operator watches the wrong signal. | Validate `limit >= 1` at construction or guard the deny branch against a short reply, and correct the row's detection column to the fail-open metric.

R1-4 | should-fix | The deny cache's load claim is stated without its regime: "a flooding client generates at most one Redis call per server per freed slot" (design.md:206-210) holds only while denials are authoritative. The cache is populated exclusively by store denials, so in the fail-open regime (store slow, breaker still closed because failures are not consecutive) nothing is cached and the full flood lands on the hot key — the protection lapses exactly under store stress, which is the correlated case the design names elsewhere. | Show the arithmetic that bounds the window (probability of 5 consecutive failures at the failure rates and per-server QPS that matter, hence time-to-trip), or state the residual with its trigger, or populate a short local throttle on fail-open.

R1-5 | should-fix | The headline selection claim overreaches: sliding-window log is "the only common algorithm that produces goal 1's *exactly 50 rejections* under every arrival pattern the goal's wording permits" (design.md:15-18), yet the design's own caveat reports 49/46/39 for 60.5/62/65 s spans (design.md:320-326) — spans that "150 requests in a minute" arguably permits. I reproduced both the table and the caveat figures; the substantive conclusion survives, the absolute phrasing does not. | Restate the claim with its true domain ("exact for any arrival pattern whose span is at most one window") and carry the span condition into the G1 acceptance criterion listed under gaps in the brief.

R1-6 | should-fix | Necessity is mis-attributed: "The minimal design that satisfies G1, G2, G3 is Middleware + RedisLimiter + Redis, and everything else is degradation machinery forced by G4 and C2" (design.md:85-86) is contradicted by the design's own text — plain fail-open "satisfies G4 as literally written" (design.md:183-184) and a 3 ms deadline already meets C2 without a breaker (design.md:236-238). Breaker, deny cache and local limiter are Problem-derived, not requirement-forced; only the local limiter is admitted as such (design.md:697-703). | Label the three tiers as one optional layer with a stated cut line, each with what the owner loses by cutting it, so the reviewer of record can approve the two-component core independently of the tiers.

R1-7 | should-fix | The determinism column of the verification plan rests on `miniredis.SetTime()` controlling the `TIME` that the *Lua script* reads (design.md:614), but the cited README line documents the `TIME` command, not `redis.call('TIME')` inside `EVAL`. I fetched the README: it lists EVAL/EVALSHA/SCRIPT LOAD and says "SetTime() also sets the value returned by TIME", and says nothing about the script path. The citation is one step short of the claim it carries, and G1/G2/G3 determinism all hang on it. | An observed run of `EVAL "return redis.call('TIME')"` against miniredis after `SetTime`, or a stated test-only fallback (pass `now` in ARGV) that leaves the production script unchanged.

R1-8 | nit | The scratchpad artifacts behind two **observed** Redis-doc quotes are 404 stubs: `eval-intro.md` and `replication.md` are 14 bytes of "404: Not Found". I independently confirmed both quoted sentences (script atomicity; TIME/SRANDMEMBER free under effects replication since Redis 5), so the claims hold and the label stands — but the trail does not reproduce. | Cite the URL actually fetched (redis.io/docs/latest/develop/programmability/eval-intro/ and the replication page) so a later reader can re-verify.

R1-9 | nit | `Decision.Remaining` (design.md:422) is specified and never consumed: the middleware emits only `Retry-After`, no `X-RateLimit-*` headers exist, and no requirement asks for remaining. | Emit it as a header or drop the field.

R1-10 | nit | "There is no other path today, so the current error is zero" for deny-cache over-denial (design.md:658-660) overlooks the design's own Redis-restart case (design.md:107-109): after a restart the store would admit, while a cached deny still rejects for up to one window. | Name the restart as the second path and bound the error at one window.

**Spot-checks** (16; failures widen the sample, none failed)

- `sim.py` at the cited scratchpad path, re-run by me — reproduces the algorithm table exactly (50/50/50/0/49 fixed, 50/50/0/0/0 token bucket, 50 across sliding log, 50/50/50/14/49 counter) | held
- Stretch caveat 49/46/39 at 60.5/62/65 s — reproduced with 150 arrivals spaced span/150 | held
- `go.mod:5` testify-only direct requirement | held
- `golang/design/rate_limiter.go:17,54,57` (FixedWindowLimiter, TokenBucket, `refillRate int // tokens per second`) | held
- `golang_containers/hit_counter.go:15-18` exclusive window boundary | held
- `golang_concurrency/leaky_bucket.go:64` LeakyBucket type | held
- `agent-team-workspace/research/leaky-bucket/leaky-bucket.md:20-24` and `:37-42` quotes (meter-equivalence; rate vs. concurrency) | held
- `golang/interview/rippling_rate_limiter.go:32` + "Follow-up 4: Distributed rate limiting" at :97 | held
- `brief.md:52-53` "No existing rate-limiting code in this repo" | held (the design's correction stands)
- Redis 7.2 `redis.conf`: `zset-max-listpack-entries 128`, `maxmemory-policy noeviction`, `appendonly no` | held (fetched from the 7.2 tag)
- ZADD docs: "Only add new elements..." and 9007199254740992 | held
- EVAL docs: all accessed key names must be passed as key arguments | held
- ZREMRANGEBYSCORE inclusive bound (basis for the exclusive-survivor claim) | held
- Redis scripting atomicity sentence | held (confirmed independently of the 404 artifact)
- TIME/SRANDMEMBER free in scripts under effects replication since Redis 5 | held (same)
- go-redis `Run` doc comment "optimistically uses EVALSHA ... retried using EVAL" | held (verbatim)

Verdict: revise
