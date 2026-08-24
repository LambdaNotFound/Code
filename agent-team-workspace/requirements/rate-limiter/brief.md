# Brief: API Rate Limiter

Status: signed-off 2026-08-24

## Problem

A public HTTP API has no rate limiting. A single client can exhaust
backend capacity, degrading service for everyone. We need a rate
limiter that bounds per-client request rates across a horizontally
scaled fleet of API servers.

## Goals

1. Enforce a per-client limit (e.g. 100 requests/minute) — verified
   by a test issuing 150 requests in a minute and observing exactly
   50 rejections.
2. The limit holds across the whole fleet, not per server — verified
   by running the same test against 3 servers behind a load balancer.
3. Rejected requests return HTTP 429 with a Retry-After header —
   verified by asserting on the response.
4. Limiter failure must not take down the API — verified by killing
   the limiter's backing store and asserting requests still serve.

## Non-goals

- Per-endpoint or tiered limits. One global per-client limit only.
- Billing, quotas, or long-window (daily/monthly) accounting.
- DDoS protection at the network layer.
- A UI or admin console for managing limits.

## Constraints and invariants

- The fleet is horizontally scaled and stateless; any server may
  receive any client's request.
- Added p99 latency from limiting must stay under 5ms (imposed).
- Clients are identified by API key present on every request.

## Decomposition

| piece | goal it serves | route | depends on |
|---|---|---|---|
| rate limiter design | 1,2,3,4 | design-loop | — |

## Open questions

- Exact limit value per client — owner: user; default 100/min.
- Whether Redis is already in the stack — owner: user; default is to
  assume it is available.

## Context

Greenfield design. No existing rate-limiting code in this repo; this
is a system-design exercise, so the design is the deliverable and no
implementation follows unless separately requested.
