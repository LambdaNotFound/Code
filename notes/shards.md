One Cluster has many shards
Each shard is a single leader w/ multiple replicas group, they are distributed multiple AZs

REGION  us-east-1
┌─ AZ-a  (DOWN) ───┐ ┌─ AZ-b ───────────┐ ┌─ AZ-c ───────────┐
│  S1  leader   ✗  │ │  S1  replica ⇒ ★ │ │  S1  replica     │
│  S2  replica  ✗  │ │  S2  leader   ✓  │ │  S2  replica     │
│  S3  replica  ✗  │ │  S3  replica     │ │  S3  leader   ✓  │
└──────────────────┘ └──────────────────┘ └──────────────────┘

  S1  lost its LEADER → election (~sub-second), ★ new leader in AZ-b, writes resume
  S2  lost a replica  → leader alive, writes never paused
  S3  lost a replica  → leader alive, writes never paused
  each shard still has 2 of 3 replicas → quorum holds → zero committed-data loss


LEADER-BASED  ·  one leader per shard   (Spanner, CockroachDB, MongoDB, HBase) These system leans CP

  client ──write/read──▶ ┌ router ┐  (maps key → shard)
                         └───┬────┘
            ┌───────────────┼───────────────┐
            ▼               ▼               ▼
       ┌ shard A ┐     ┌ shard B ┐     ┌ shard C ┐
       │ LEADER  │     │ LEADER  │     │ LEADER  │   ◀ every write lands here
       │ replica │     │ replica │     │ replica │   ◀ read / standby
       │ replica │     │ replica │     │ replica │
       └─────────┘     └─────────┘     └─────────┘

  WRITE : router → that key's leader → leader commits on a majority → ack
  READ  : from the leader (strong)  or  a follower (cheaper, maybe stale)
  FAIL  : leader down → followers elect a new one (~sub-sec), no committed loss

What actually delivers strong consistency is three things working together:
  Single leader → total order on writes (no conflicts)
  Consensus quorum (Raft) → the leader doesn't ack until a majority has the write durably, so committed writes survive failover
  Read from the leader → the read sees the latest committed state (read from a follower = possible stale read = not linearizable)

  Brief write unavailability during election — sub-second, per-shard, automatic. Not a human-intervention outage; a momentary pause.
  Minority-side unavailability during a network partition — if a partition splits the cluster, the side without a majority can't elect a leader and refuses writes. The majority side stays fully available. This is the CP lean: it sacrifices availability on the minority side of a partition to guarantee consistency (no split-brain, no divergent writes).

LEADERLESS  ·  no leader, quorum reads/writes   (Dynamo, Cassandra, Riak, Scylla) These system leans AP.

  client ──write/read──▶ ┌ coordinator ┐  (any node; consistent-hash key → N replicas)
                         └──────┬───────┘
              send to all N  ┌──┴──┬──────┐
              wait W or R    ▼     ▼      ▼
                       ┌ rep 1 ┐ ┌ rep 2 ┐ ┌ rep 3 ┐    N = 3
                       │ node  │ │ node  │ │ node  │
                       └───────┘ └───────┘ └───────┘

  WRITE : send to N, succeed once W ack          (e.g. W = 2)
  READ  : query N, return once R respond         (e.g. R = 2)
  RULE  : W + R > N  ⇒  read & write sets overlap ⇒ a read sees the newest write
  CONFLICTS : concurrent writes → last-write-wins / version vectors / CRDTs
  REPAIR    : read-repair · hinted handoff · Merkle-tree anti-entropy