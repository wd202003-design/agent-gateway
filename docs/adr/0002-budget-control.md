# ADR-0002: Strict budget reservation and settlement
Status: Accepted policy, implementation details pending validation

## State machine
REQUESTED -> RESERVED -> CAPTURED or RELEASED; EXPIRED is a recovery state requiring reconciliation. Reservation has reservation_id, tenant_id, account_id, request_id, currency, amount, expiry, state, version and fencing token.

## Admission
1. Authenticate and authorize tenant/key/model.
2. Resolve immutable price snapshot, request limits, and bounded maximum estimated charge.
3. Atomically reserve available balance (budget - posted charges - active reservations) using Redis Lua/Function with a single authoritative coordination slot per budget account.
4. Persist reservation intent and lifecycle in PostgreSQL; admission must not proceed if durable reservation confirmation fails. Reconcile Redis and PostgreSQL after partial failure.
5. Stream request; track usage; settle in PostgreSQL using idempotent charge/ledger transaction, then release or adjust Redis hold. Never double capture on retry.

## Strictness boundary
Redis and PostgreSQL cannot provide an atomic cross-system transaction. Therefore strict no-overspend under partitions/crashes requires a defined authority and recovery protocol. MVP uses PostgreSQL as final authority with atomic conditional reservation/ledger updates for hard-budget accounts; Redis accelerates distributed admission and rate limits but cannot overrule a PostgreSQL denial. Benchmark the extra database hop. Future escrow allocations may reduce hot-path writes if proven safe.

## Degradation
- Unknown maximum liability: enforce max_output_tokens and maximum charge, or reject strict-budget request.
- Redis failure: reject new strict-budget admissions.
- PostgreSQL failure: reject new strict-budget admissions; reconcile in-flight.
- Upstream disconnect or missing usage: provisional charge with timeout/reconciliation; never assume zero usage.
- Concurrent over-limit: return structured budget_exceeded / quota_exceeded.
- Reservation expiry must not automatically free liability for an attempt still running; use lease/fencing and attempt-state reconciliation.

## Test gates
1000 concurrent reservations on same account, duplicate idempotency keys, retry after timeout, worker crash, Redis failover, PostgreSQL failover, partial stream, reconciliation and ledger balance invariants.
