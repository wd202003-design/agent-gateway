# ADR-0003: MVP protocols and cross-model fallback
Status: Accepted, 2026-11-17

## MVP
Both OpenAI-compatible POST /v1/chat/completions and POST /v1/responses are M1 requirements, including streaming and non-streaming, usage, error mapping, cancellation and request tracing. GET /v1/models returns only authorized public models.

Canonical internal representation must preserve protocol-specific semantics and opaque extensions; reject unsupported features explicitly instead of silently dropping reasoning/tool/multimodal fields. Maintain adapter version and upstream API version compatibility tests.

## Fallback
- Same-model across compatible channels is permitted by policy before any irreversible client output, subject to budget, region, capability and quota.
- Cross-model fallback is DISABLED by default.
- Tenant admin may explicitly configure allow_cross_model_fallback and an allowlist of public model substitutions, subject to model entitlement, protocol/feature equivalence, residency and price ceiling.
- Project/key may further restrict but not expand tenant policy.
- Every decision records original public model, actual upstream model, substitution reason, policy version and billing rule.
- Once SSE output has been emitted, do not splice a different upstream stream. End with a structured error/partial state and reconcile actual usage.
- If no approved compatible fallback exists, return model_unavailable instead of changing models.

## Acceptance
Contract tests for both protocols: JSON, SSE, tools, structured output, reasoning, unsupported features, upstream 429/5xx, cancellation, partial output, fallback deny/allow, pricing ceiling, audit and model version drift.
