# P-002 JetStream event contract draft

Stream: AGENT_GATEWAY_USAGE (subject: gateway.usage.v1), durable, explicit ACK, retention and replay policies to be defined before production.

Envelope:
```json
{
 "event_id":"uuid",
 "event_type":"usage.finalized",
 "schema_version":1,
 "occurred_at":"2026-11-17T00:00:00Z",
 "tenant_id":"uuid",
 "request_id":"uuid",
 "attempt_id":"uuid",
 "protocol_family":"openai.chat_completions",
 "usage":{"input_tokens":120,"output_tokens":30,"source":"provider","status":"final"}
}
```

Producer writes a PostgreSQL outbox row in the same transaction as its authoritative state change, then publishes with stable event_id. Consumer ACK only after durable processing, and deduplicates by event_id + revision. ClickHouse is a projection and may be replayed. Never send API keys, prompts, or response bodies through the billing subject. Contract not wired to runtime in current slice.
