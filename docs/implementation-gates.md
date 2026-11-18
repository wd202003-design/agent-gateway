# P-002 implementation gates

## P-002.1 verification
Go tests, race detector, vet, real mock upstream, SSE interruption/cancellation. CI workflow added; no local Go test evidence yet because GitHub checkout failed in this environment.

## P-002.2 persistence
PostgreSQL schema and database/sql repository drafted. A PostgreSQL driver, migrations runner, configuration service and production SK lookup are NOT wired into HTTP gateway. No database integration tests.

## P-002.3 routing
Pure Go candidate filtering/scoring and breaker primitives drafted with unit tests. Not wired into HTTP proxy; no real channel health sampling, retries or route decision persistence.

## P-002.4 budget/metering
PostgreSQL conditional reservation/capture draft and JetStream publisher interface + ClickHouse projection DDL. Redis coordination, NATS JetStream concrete client, outbox worker, ClickHouse ingestion, failure recovery and reconciliation NOT wired. Strict-budget production admission must remain disabled.

## P-002.5 commercial billing
Not yet implemented in runtime. Price/version schemas, charge calculator, immutable accounting and invoices need integration and financial review.

## Non-negotiable release gate
Do not enable external paid traffic or advertise production readiness until all dependencies are wired, concurrency/failover tests pass, tenant isolation is enforced, and accounting invariants are verified.
