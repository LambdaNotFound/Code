Design Autopay System

Functional Requirements
FR1 — Schedule configuration. Payroll admin configures payment schedules per payee with support for monthly, bi-weekly, or ad-hoc one-time payments.
FR2 — Exactly-once execution. Every scheduled payment is executed once and only once — no double payments, no missed payments.
FR3 — Audit trail. Payroll team can query the full payment history per payee with pagination and date range filtering.
FR4 — Multi-region tax. System calculates tax deductions per payee's jurisdiction, breaking down gross → tax → net before executing payment.

Non-Functional Requirements
NFR1 — Throughput. ~1M employees processed within a 5-hour window → ~100 TPS sustained.
NFR2 — Strong consistency on writes. Payment creation is guarded by UNIQUE(schedule_id, pay_period_date) to prevent double-scheduling at the DB level.
NFR3 — Eventual consistency on reads. Aggregate views, dashboards, and audit queries can lag by minutes. This is acceptable for reporting.

Data Model

Payee 
— payee_id (UUID PK), name (VARCHAR), country_code (CHAR(2)), tax_profile_id (UUID FK → external tax profile).

Schedule 
— schedule_id (UUID PK), payee_id (UUID FK), cadence (ENUM: MONTHLY | BIWEEKLY | ONETIME), next_pay_date (DATE, indexed), amount_minor (BIGINT), currency_code (CHAR(3) ISO 4217), status (ENUM: ACTIVE | INACTIVE), created_at (TIMESTAMPTZ).

Payment 
— payment_id (UUID PK), schedule_id (UUID FK), payee_id (UUID FK), pay_period_date (DATE), gross_minor (BIGINT), tax_minor (BIGINT), net_minor (BIGINT), currency_code (CHAR(3)), status (ENUM: PENDING → COMPLETED | FAILED), created_at (TIMESTAMPTZ), updated_at (TIMESTAMPTZ). Idempotency guard: UNIQUE(schedule_id, pay_period_date).

Outbox 
— outbox_id (UUID PK), payment_id (UUID FK), payload (JSONB), published (BOOLEAN), created_at (TIMESTAMPTZ). Written atomically with Payment in a single transaction. Drained by outbox poller to Kafka.

API 

POST /schedules — Create a schedule. Body: payee_id, cadence, amount_minor, currency_code. Returns schedule_id.
PUT /schedules/{id} — Update amount or cadence. Does not affect in-flight payments.
PATCH /schedules/{id}/deactivate — Soft-delete. Sets status to INACTIVE. Past payment rows preserved for audit.
GET /payees/{id}/payments — Paginated audit ledger. Query params: ?from=&to=&cursor=&limit=
POST /payroll-runs — Optional manual trigger. Primary path is cron-driven.

┌─────────────────────────────────────────────────────────────────────────────────┐
│                              PAYROLL SYSTEM                                     │
│                                                                                 │
│  ┌──────────┐         ┌─────────────────────────────────────────────┐           │
│  │ Payroll  │  POST   │              API Service                    │           │
│  │ Admin UI ├───────► │  /schedules, /payees/{id}/payments          │           │
│  └──────────┘         └──────────────┬──────────────────────────────┘           │
│                                      │ CRUD                                     │
│                                      ▼                                          │
│                       ┌──────────────────────────────┐                          │
│                       │         Postgres (Leader)    │                          │
│                       │                              │                          │
│                       │  ┌─────────┐  ┌───────────┐  │                          │
│                       │  │Schedule │  │  Payment  │  │                          │
│                       │  │         │  │           │  │                          │
│                       │  └─────────┘  └───────────┘  │                          │
│                       │  ┌─────────┐  ┌───────────┐  │                          │
│                       │  │ Payee   │  │  Outbox   │  │                          │
│                       │  └─────────┘  └───────────┘  │                          │
│                       └──────┬───────────┬───────────┘                          │
│                              │           │                                      │
│                     Read     │           │ WAL stream                           │
│                              │           ▼                                      │
│  ┌───────────┐        ┌──────┴───┐   ┌──────────┐      ┌─────────────────────┐  │
│  │  Cron     │ SELECT │ Schedule │   │ Outbox   │      │                     │  │
│  │  Trigger  ├───────►│ WHERE    │   │ Poller / │─────►│   Kafka             │  │
│  │ (every N  │  due   │ next_pay │   │ Debezium │      │                     │  │
│  │  minutes) │  rows  │ _date <  │   │ (CDC)    │      │   payment.pending   │  │
│  └───────────┘        │ now+5h   │   └──────────┘      │   topic             │  │
│                       │ AND      │                     │                     │  │
│                       │ status=  │                     │   Partitioned by    │  │
│                       │ ACTIVE   │                     │   schedule_id       │  │
│                       └──────────┘                     └──────────┬──────────┘  │
│                                                                   │             │
│       ┌─────────────────────────────────────────────┐             │             │
│       │         Cron Transaction (atomic)           │             │             │
│       │                                             │             │             │
│       │  1. SELECT schedules due in next 5h         │             │             │
│       │  2. INSERT Payment (status=PENDING)         │             │             │
│       │  3. INSERT Outbox row (same txn)            │             │             │
│       │                                             │             │             │
│       │  Idempotency: UNIQUE(schedule_id,           │             │             │
│       │               pay_period_date)              │             │             │
│       └─────────────────────────────────────────────┘             │             │
│                                                                   │             │
│                              ┌────────────────────────────────────┘             │
│                              │                                                  │
│                              ▼                                                  │
│               ┌──────────────────────────────┐                                  │
│               │   Consumer Group (N workers) │                                  │
│               │                              │                                  │
│               │  Per message:                │                                  │
│               │  1. Look up Payee            │        ┌──────────────┐          │
│               │  2. Call Tax Engine ─────────────────►│  Tax Engine  │          │
│               │     (gross → tax → net)      │◄───────│  (per-country│          │
│               │  3. Call Payment Gateway ────────────►│   rules)     │          │
│               │  4. UPDATE Payment status    │        └──────────────┘          │
│               │     → COMPLETED | FAILED     │                                  │
│               │  5. Commit Kafka offset      │        ┌──────────────┐          │
│               │                              │───────►│  Payment     │          │
│               └──────────────┬───────────────┘        │  Gateway     │          │
│                              │                        │  (bank API)  │          │
│                              │ on persistent          └──────────────┘          │
│                              │ failure                                          │
│                              ▼                                                  │
│                       ┌──────────────┐                                          │
│                       │     DLQ      │                                          │
│                       │  + Alert     │                                          │
│                       └──────────────┘                                          │
│                                                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘