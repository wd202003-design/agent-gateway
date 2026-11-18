-- ClickHouse analytics projection, NOT the billing source of truth.
CREATE TABLE IF NOT EXISTS gateway_usage_events
(
 event_id UUID,
 request_id UUID,
 attempt_id UUID,
 tenant_id UUID,
 schema_version UInt16,
 input_tokens UInt64,
 output_tokens UInt64,
 occurred_at DateTime64(3, 'UTC'),
 revision UInt32 DEFAULT 1
)
ENGINE = ReplacingMergeTree(revision)
PARTITION BY toYYYYMM(occurred_at)
ORDER BY (tenant_id, request_id, attempt_id, event_id);
-- ReplacingMergeTree merges asynchronously; query design must account for duplicates
-- and corrections. Do not calculate authoritative invoices directly from this table.
