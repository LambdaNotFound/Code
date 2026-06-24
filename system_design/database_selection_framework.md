The one-liner for interviews: "I always start with Postgres and look for the reason to leave. Read bottleneck? Add replicas before sharding. Write bottleneck beyond what sharded Postgres handles? That's when I reach for Cassandra or DynamoDB, depending on the consistency tier and aggregation shape."

### Read pattern

— Are reads PK lookups, range scans, joins, or ad-hoc aggregations? This is the single biggest driver and eliminates half the options immediately.

1. Point PK reads — "Get me this one row by its key"

=> Essentially everything handles PK lookups.

1. Range scans — "Get me rows between X and Y"

=> Postgres — B-tree index on the range column
=> Cassandra — this is its sweet spot. Partition key + clustering key gives you sorted data within a partition.
=> DynamoDB — partition key + sort key, same model as Cassandra.

1. Joins + aggregations — "Combine data across tables, compute summaries"

=> Postgres / MySQL — this is what relational DBs exist for. Joins, subqueries, GROUP BY, HAVING, window functions. The canonical choice.
=> CockroachD — when you need joins + ACID at horizontal scale.

Eliminate:
Cassandra — no joins at all. No server-side GROUP BY worth using. If you need joins or aggregation, Cassandra is the wrong pick, full stop. You'd have to bolt on Spark, which changes your architecture materially.

DynamoDB — no joins. Aggregation requires scanning the entire table or building a secondary index and doing client-side work.

MongoDB — has an aggregation pipeline that can do lookups (pseudo-joins), but performance degrades fast at scale. If joins are a primary pattern, MongoDB is fighting its own design.

### Consistency tier

— "How stale can your reads be?" Do you need strong consistency (ACID, read-your-writes, linearizability) or is eventual OK? Payments vs. social feeds sit on opposite ends.

```
Strong consistency — after a write completes, every subsequent read sees that write. No stale data, ever. Linearizable reads. This is what you need when reading stale data causes financial or correctness errors.

Read-your-writes — the user who performed the write sees it immediately, but other users might see stale data briefly. Weaker than full linearizability but sufficient for most user-facing flows.

Eventual consistency — writes propagate asynchronously. Reads might return stale data for some window (typically ms to seconds). No staleness guarantee.
```

1. Strong consistency required — "Wrong data = lost money or broken invariants"

Works well:

Postgres / MySQL — single-leader, ACID transactions, SELECT FOR UPDATE for pessimistic locking, serializable isolation if you need it. The default choice. Payment ledger, wallet balance, inventory count.
CockroachDB / Spanner — strong consistency with horizontal sharding. Consensus per shard (Raft/Paxos). Global payments where you need ACID across regions.
DynamoDB (with strongly consistent reads) — you can opt into strong reads per-request, but you pay 2x read cost and lose multi-region. Shopping cart where you can't show stale items.
Eliminate:
Cassandra — eventual by default. 

1. Read-your-writes sufficient — "I need to see my own changes, others can lag"

Works well:

- **Everything from (1)** — strong consistency trivially satisfies read-your-writes.
- **Postgres with read replicas** — route the *writing user's* reads to the primary for a few seconds after their write, everyone else reads from replicas. This is a common production pattern — "sticky reads" after writes.
- **DynamoDB** — use consistent reads only for the session that just wrote, eventual reads for everyone else. Cheaper than strong reads everywhere.
- **MongoDB** — `readConcern: "majority"` with `writeConcern: "majority"` gives you read-your-writes in a replica set. User profile updates.

1. Eventual consistency OK — "Staleness of seconds is fine, show whatever you have"

Works well:

**Cassandra** — this is where it shines. `ONE` consistency reads, blazing fast, linear scale, multi-DC active-active replication

### Scale

— Does it fit on one machine? How many writes/sec? How many shards? This determines whether you even need to leave Postgres.

1. Does the data fit on one machine?
2. If you need to shard, can you shard SQL or do you need native horizontal scaling?

**Path 1 — Shard Postgres (Citus, Vitess)**

Stay on SQL, bolt on a sharding layer. You keep joins (within a shard), ACID (within a shard), and your team's existing Postgres expertise.

**Path 2 — Native horizontal scaling (Cassandra, DynamoDB, CockroachDB)**

The database handles sharding, rebalancing, and replication for you. You give up some SQL capabilities but gain operational simplicity at scale.

1. Scale of reads vs. writes — this changes the strategy

Read-heavy (100:1 read-to-write ratio)
Don't shard the writes. Instead, scale reads horizontally with replicas and caching.  

Write-heavy (1:1 or write-dominant)
Read replicas don't help — the bottleneck is write throughput on the primary. You must either shard writes or pick a write-optimized store.

### Aggregation shape

— This is the step most candidates skip. Are aggregations pre-defined rollups (→ TSDB), ad-hoc analytical scans (→ ClickHouse/Druid), or not a first-class read pattern at all (→ Cassandra is fine)?

1. No aggregation — "Just give me the rows"
2. Pre-defined rollups — "I know exactly what summaries I need at write time"
3. Ad-hoc analytical scans — "Let me slice and dice however I want"
4. Batch aggregation — "Run it overnight, I don't need real-time"

  



|                 | **PostgreSQL**                                               | **DynamoDB**                                                | **Cassandra**                                             | **CockroachDB**                                               |
| --------------- | ------------------------------------------------------------ | ----------------------------------------------------------- | --------------------------------------------------------- | ------------------------------------------------------------- |
| **Type**        | Relational RDBMS                                             | Key-value / document                                        | Wide-column                                               | Distributed SQL                                               |
| **Engine**      | B-tree + heap filesRow-oriented on disk                      | B-tree (managed)AWS proprietary                             | LSM-tree + SSTablesAppend-only, compaction                | LSM-tree (Pebble)RocksDB-derived                              |
| **CAP leaning** | **CP**Chooses consistency over availability during partition | **AP**Available by default; strong reads opt-in per request | **AP**Available + partition-tolerant; consistency tunable | **CP**Raft consensus; brief unavailability on leader election |


