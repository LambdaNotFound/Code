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
              │   ┌──────────────────┴───────────────────-┐
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

  ┌────────────────────────────────────────────────────────────────────---──┐
  │  CONCURRENCY  —  DB uniqueness constraint on short_code IS the lock     │
  │                                                                         │
  │   Two requests, SAME url    → both compute same code → one wins insert, │
  │                               loser reads it back → both return same.   │
  │   Two requests, DIFF urls   → hash-collide on 7-char prefix → loser     │
  │   (rare)                      sees stored != mine → rehash with seed.   │
  │                                                                         │
  │   No distributed lock needed: conflict is rare + cheap to resolve after.│
  └────────────────────────────────────────────────────────────────────---──┘

  ┌────────────────────────────────────────────────────────────────────---──-┐
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