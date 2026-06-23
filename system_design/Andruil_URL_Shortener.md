Design URL Shortener

FR:
FR-1 — Shorten: Given a long URL, return a short code. POST /shorten { long_url } → { short_url }
FR-2 — Redirect: Given a short code, redirect to the original long URL. GET /{short_code} → 301 → long_url
FR-3 — Deterministic mapping: The same long URL must always map to the same short code. This is the defining constraint — it makes the system content-addressable and rules out random/sequential ID generation.
FR-4 — Idempotent creation: Re-submitting an already-shortened URL returns the existing code, not a new one (falls directly out of FR-3).
FR-5 — Collision correctness: Two different long URLs must never resolve to the same code, even when their hash prefixes collide.

NFR:
NFR-1 — Latency: Redirect (read path) should be fast, p99 under ~10ms on a cache hit. The redirect is in the critical path of every embedded link, so users feel it directly. Shorten (write path) can tolerate more latency since it's off the hot path.
NFR-2 — Availability: Read path should target 99.9% or higher, ideally 99.99%. Links are embedded permanently across the internet, so a dead redirect service breaks third-party content you don't control. Write path can tolerate slightly lower availability.
NFR-3 — Scalability: Must scale horizontally on both storage and throughput. Storage is often the real driver here (billions of rows over time), not just QPS.
NFR-4 — Read-heavy skew: Roughly 100:1 read-to-write ratio. This is the single assumption that drives the cache-first architecture; writes are comparatively rare.
NFR-5 — Durability: No lost mappings, ever. Losing a single short_code to long_url row permanently breaks every link pointing at it.
NFR-6 — Consistency: Eventual consistency is acceptable on reads; creation must be atomic. Because mappings are immutable, a stale read is never wrong, only briefly missing a brand-new code. The write path still needs strongly-consistent atomic insert to enforce idempotency and collision-correctness. (Different consistency requirements on read versus write is the subtle, senior-level point.)
NFR-7 — Security and abuse resistance: Rate limiting on writes plus malicious-URL filtering. Open shorteners are routinely abused for phishing and spam, so you throttle write volume and scan target URLs.

Table: url_mappings

  short_code      string   PRIMARY KEY      7-char Base62, the lookup key
  long_url        string                    the original URL to redirect to
  url_hash        string   SECONDARY INDEX  SHA-256(long_url), for "seen this?" checks
  collision_seed  int                       0 normally; >0 if bumped by a collision
  created_at      timestamp                 when the mapping was created
  ttl             timestamp (nullable)      optional expiry; null = permanent

API
POST /shorten
  Request body:   { "long_url": "https://example.com/very/long/path" }
  Response 200:   { "short_url": "https://sho.rt/abc1234",
                    "short_code": "abc1234" }
  Response 400:   invalid or malformed URL
  Response 429:   rate limit exceeded

  Semantics: idempotent. Submitting the same long_url returns the same
  short_code every time — never creates a duplicate.
GET /{short_code}
  Response 301:   Location: https://example.com/very/long/path
                  (permanent redirect — immutable mapping makes this safe
                   and lets CDN/browser cache aggressively)
  Response 404:   unknown or expired short_code

  Semantics: read-only, cache-first. Served from CDN/Redis on hit;
  falls through to the store on miss.

                              ┌─────────────────────────────────────┐
                              │              CLIENTS                │
                              │   POST /shorten      GET /{code}    │
                              └──────┬───────────────────┬──────────┘
                                     │                   │
              ═══ WRITE PATH ════════╪═══════════════════╪════ READ PATH ═══
                                     │                   │
                                     ▼                   ▼
                          ┌───────────────────=─┐   ┌-──────────────────┐
                          │    API Gateway      │   │     CDN / Edge    │
                          │  • authn / API key  │   │  • caches 301s    │
                          │  • rate limit       │   │  • TTL on entries │
                          │  • URL validation   │   └─-────────┬────────┘
                          └──────────┬─────────=┘              │ miss
                                     │                         ▼
                                     ▼                ┌──=────────────────┐
                          ┌────────────────────┐      │   Redis Cache     │
                          │     Hash Engine    │      │  key: short_code  │
                          │  SHA-256(long_url) │      │  val: long_url    │
                          │  → Base62 → take 7 │      │  • LRU eviction   │
                          └──────────┬─────────┘      └─-────────┬────────┘
                                     │                           │ miss
                                     ▼                           │
                          ┌────────────────────┐                 │
                          │  Collision Handler │                 │
                          │  insert-if-absent  │                 │
                          └──────────┬─────────┘                 │
                                     │                           │
                                     ▼                           │
                    ┌─────────────────────────────-───┐          │
              ┌────▶│       URL Mapping Store         │◀────────-┘
              │     │  (DynamoDB / Cassandra / PG)    │
              │     │ ┌─────────────────────────────┐ │
              │     │ │ PK:  short_code  (7 chars)  │ │
              │     │ │ SI:  hash(long_url)         │ │  ◀── idempotency:
              │     │ │ val: long_url, seed, ts,ttl │ │      "seen this URL?"
              │     │ └─────────────────────────────┘ │
              │     └────────────────┬────────────────┘
              │                      │
              │   ┌──────────────────┴────────────────────┐
              │   │   INSERT ON CONFLICT DO NOTHING       │
              │   │   then read-after-write, branch:      │
              │   │                                       │
              │   │   ┌──────────────┬──────────────────┐ │
              │   │   ▼              ▼                  ▼ │
              │   │ won insert   stored==mine      stored!=mine
              │   │  (new URL)   (dup/race)        (COLLISION)
              │   │   │              │                  │  │
              │   │   ▼              ▼                  ▼  │
              │   │ return       return same      rehash with
              │   │ new code     code (idem-      seed+1 ──────┐
              │   │              potent)          retry insert │
              │   └─────────────────────────────────────────┘  │
              │                                                │
              └────────────────────────────────────────────────┘
                            (seeded rehash loop, bounded ~3-5x)

  ┌────────────────────────────────────────────────────────────────────────┐
  │  CONCURRENCY  —  DB uniqueness constraint on short_code IS the lock     │
  │                                                                         │
  │   Two requests, SAME url    → both compute same code → one wins insert, │
  │                               loser reads it back → both return same.   │
  │   Two requests, DIFF urls   → hash-collide on 7-char prefix → loser     │
  │   (rare)                      sees stored != mine → rehash with seed.   │
  │                                                                         │
  │   No distributed lock needed: conflict is rare + cheap to resolve after.│
  └─────────────────────────────────────────────────────────────────────────┘

  ┌──────────────────────────────────────────────────────────────────────────┐
  │  KEY DESIGN DECISIONS                                                    │
  │                                                                          │
  │  • Content-addressable: hash(url) = deterministic key, no counter,       │
  │    no central ID sequence. Same input → same output, for free.           │
  │  • 301 (permanent) not 302: safe BECAUSE mapping is immutable. Lets      │ 
  │    CDN + browser cache hard. Use 302 only if you need per-click analytics│
  │  • Secondary index on hash(long_url): makes "have I seen this?" O(1) —   │
  │    the piece that makes idempotent creates fast.                         │
  │  • 62^7 ≈ 3.5T keys before exhaustion. Reads ≫ writes (100:1+),          │
  │    so the whole system is cache-first, sub-10ms redirect on hit.         │
  └─────────────────────────────────────────────────────────────────────----- d┘