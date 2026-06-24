┌─────────────────────────────────────────────────┐
│         PAIR 2: FEED GENERATION                 │
│                                                 │
│   News Aggregator  ≈  Social Media Feed         │
│                                                 │
│   Shared: Ingestion → Ranking → Redis feed cache│
│   Differ: News clusters by TOPIC (semantic);    │
│           Social fans out by SOCIAL GRAPH.      │
│           News = fan-out-on-read;               │
│           Social = hybrid fan-out.              │
│           News has no celebrity problem         │
│           (publishers aren't followed by 50M    │
│           individual users in the same way).    │
└─────────────────────────────────────────────────┘