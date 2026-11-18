# P-002 MVP development slice

## Run
Requires Go 1.22+.
```sh
export GATEWAY_CLIENT_KEY='local-client-secret'
export UPSTREAM_API_KEY='local-upstream-secret'
export UPSTREAM_BASE_URL='http://127.0.0.1:9000'
export UPSTREAM_MODEL='native-model'
export PUBLIC_MODEL='gateway-default'
go test ./...
go run ./cmd/gateway
```
Mock or real OpenAI-compatible upstream must expose both /v1/chat/completions and /v1/responses. The proxy does not translate between protocols.

```sh
curl -sS http://127.0.0.1:8080/v1/models -H 'Authorization: Bearer local-client-secret'
curl -sS http://127.0.0.1:8080/v1/chat/completions -H 'Authorization: Bearer local-client-secret' -H 'Content-Type: application/json' -d '{"model":"gateway-default","messages":[{"role":"user","content":"hello"}]}'
```

## Security and readiness limitations
This is an explicitly development-only slice: single static SK, single upstream, no tenant store, no billing, no strict budget, no PostgreSQL/Redis/NATS/ClickHouse runtime connection, no model-group pricing, no real protocol transformation or failover. It MUST NOT serve paid or multi-tenant production traffic. /readyz signals development process readiness only, not production readiness. No SK or prompt/response body is logged.

## Planned production gates
PostgreSQL authoritative budget/ledger and SK/tenant models; Redis distributed quota; JetStream durable usage; ClickHouse projections; model capability matrix; multi-channel routing and fallback; TLS, SSRF protection, trusted config snapshots, proper upstream URL path normalization, request IDs, TTFT/usage accounting and integration tests against pinned upstream API versions.
