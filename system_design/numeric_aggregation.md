┌─────────────────────────────────────────────────┐
│         PAIR 1: NUMERIC AGGREGATION             │
│                                                 │
│   Metrics Aggregator  ≈  Ad Click Aggregator    │
│                                                 │
│   Shared: Kafka → Flink pre-agg → ClickHouse    │
│   Differ: Metrics tolerates approximation;      │
│           Ad Click requires billing-grade       │
│           exactly-once + fraud layer            │
└─────────────────────────────────────────────────┘

                    ┌──────────────────────────────┐
                    │       Metric Sources         │
                    │  (app agents, infra agents,  │
                    │ StatsD, Prometheus exporters)│
                    └──────────────┬───────────────┘
                                   │ UDP/HTTP push
                                   ▼
                          ┌────────────────────┐
                          │  Ingestion Gateway │
                          │  (stateless, LB'd) │
                          │  • validate schema │
                          │  • tag enrichment  │
                          │  • batch → Kafka   │
                          └────────┬───────────┘
                                   │
                                   ▼
                          ┌─────────────────┐
                          │      Kafka      │
                          │  (raw-metrics)  │
                          │  partition by   │
                          │  metric_name +  │
                          │  source_id      │
                          └───┬─────────┬───┘
                              │         │
              ┌───────────────┘         └───────────────┐
              │                                         │
              ▼                                         ▼
   ┌─────────────────────┐                ┌─────────────────────┐
   │   Hot Path Consumer │                │  Warm Path Consumer │
   │   (bare consumer)   │                │  (Flink streaming)  │
   │                     │                │                     │
   │  • sliding window   │                │  • tumbling windows │
   │    counters         │                │    1min, 5min, 1hr  │
   │  • real-time alerts │                │  • pre-aggregate    │
   │  • latest gauge     │                │    SUM/AVG/P99 by   │
   │    values           │                │    (metric, tags,   │
   │                     │                │     time_bucket)    │
   └────────┬────────────┘                └──┬──────────────┬───┘
            │                                │              │
            ▼                                ▼              ▼
   ┌─────────────────┐            ┌────────────────┐ ┌───────────┐
   │     Redis       │            │   ClickHouse   │ │ S3/Iceberg│
   │                 │            │                │ │ (cold)    │
   │ • sorted sets   │            │ • 1min/5min/   │ │           │
   │   (sliding      │            │   1hr rollup   │ │ • raw     │
   │   window ctrs)  │            │   tables       │ │   events  │
   │ • latest gauge  │            │ • partitioned  │ │ • 90d+    │
   │   hash maps     │            │   by time +    │ │   rollups │
   │ • alert state   │            │   metric_name  │ │           │
   │                 │            │ • TTL: 30d raw,│ │           │
   │ noeviction      │            │   90d rollups  │ │           │
   └────────┬────────┘            └───────┬────────┘ └───────────┘
            │                             │
            │         ┌───────────────────┘
            │         │
            ▼         ▼
   ┌──────────────────────────┐
   │       Query Service      │
   │                          │
   │  • route by time range:  │
   │    last 1hr  → Redis     │
   │    last 30d  → ClickHouse│
   │    older     → S3/Iceberg│
   │                          │
   │  • merge across tiers    │
   │  • downsample on the fly │
   │    if range too wide     │
   └─────────┬────────────────┘
             │
             ▼
   ┌──────────────────────────┐
   │     API Gateway / LB     │
   └──┬──────────────┬────────┘
      │              │
      ▼              ▼
┌───────────┐  ┌──────────────┐
│ Dashboard │  │ Alert Engine │
│   UI      │  │              │
│           │  │ • threshold  │
│ • Grafana │  │   rules      │
│   style   │  │ • anomaly    │
│ • custom  │  │   detection  │
│   charts  │  │ • PagerDuty/ │
│           │  │   webhook    │
└───────────┘  └──────────────┘