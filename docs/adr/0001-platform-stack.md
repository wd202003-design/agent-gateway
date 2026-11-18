# ADR-0001: Data and messaging stack
Status: Accepted (design decision), 2026-11-17

## Decision
- PostgreSQL: source of truth for tenants, SK metadata, provider/channel/model catalog, configuration revisions, price and cost snapshots, budget accounts, reservations, invoices and immutable ledger.
- ClickHouse: high-volume invocation/attempt/usage/latency analytics and reporting projections. Not authoritative for account balances or final billing.
- Redis: runtime config cache, distributed rate limits, concurrency quotas and atomic budget reservation coordination. No financial ledger of record.
- NATS JetStream: durable event streams with explicit acknowledgements, bounded retries, dead-letter handling and replay. Core NATS alone is not the billing transport.

## Consistency
PostgreSQL transaction writes state change plus outbox row; publisher retries to JetStream; consumer deduplicates by event_id / aggregate_version and records checkpoints. Exactly-once delivery across systems is not assumed. ClickHouse ingestion is replayable; deduplicate by stable invocation/attempt/usage IDs and version; aggregation queries must account for retries and correction events.

## Failure modes
- PostgreSQL unavailable: reject new paid reservations and configuration mutations; existing in-flight requests use documented completion/reconciliation procedure.
- Redis unavailable: fail closed for strict-budget paid requests and new quota admissions; do not silently switch to per-node approximate budgets. Internal best-effort mode requires explicit policy and cannot bypass tenant authorization.
- JetStream unavailable: durable local/PG outbox buffers events with limits; when durable capacity is exhausted, reject new billable requests instead of dropping usage.
- ClickHouse unavailable: continue if durable events are safely queued, with backlog alert and replay.

## Verification required
Crash between commit/publish, duplicate delivery, consumer replay, Redis failover and script atomicity, budget oversubscription under concurrency, ClickHouse reingestion, outage recovery, and ledger reconciliation.
