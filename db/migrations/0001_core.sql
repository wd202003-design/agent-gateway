-- P-002 initial PostgreSQL schema (draft; apply only after review)
CREATE TABLE tenants (
 tenant_id UUID PRIMARY KEY, tenant_type TEXT NOT NULL CHECK (tenant_type IN ('INTERNAL','EXTERNAL')),
 status TEXT NOT NULL DEFAULT 'ACTIVE', created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE api_keys (
 key_id UUID PRIMARY KEY, tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
 key_prefix TEXT NOT NULL, verifier_hash BYTEA NOT NULL UNIQUE,
 status TEXT NOT NULL DEFAULT 'ACTIVE', created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 expires_at TIMESTAMPTZ
);
CREATE TABLE providers (
 provider_id UUID PRIMARY KEY, code TEXT NOT NULL UNIQUE, name TEXT NOT NULL
);
CREATE TABLE channels (
 channel_id UUID PRIMARY KEY, provider_id UUID NOT NULL REFERENCES providers(provider_id),
 endpoint TEXT NOT NULL, credential_ref TEXT NOT NULL, region TEXT, status TEXT NOT NULL
);
CREATE TABLE upstream_models (
 upstream_model_id UUID PRIMARY KEY, provider_id UUID NOT NULL REFERENCES providers(provider_id),
 native_name TEXT NOT NULL, revision TEXT, capabilities JSONB NOT NULL DEFAULT '{}'::jsonb
);
CREATE TABLE channel_deployments (
 deployment_id UUID PRIMARY KEY, channel_id UUID NOT NULL REFERENCES channels(channel_id),
 upstream_model_id UUID NOT NULL REFERENCES upstream_models(upstream_model_id),
 protocol_family TEXT NOT NULL, api_version TEXT, status TEXT NOT NULL,
 UNIQUE(channel_id, upstream_model_id, protocol_family)
);
CREATE TABLE public_models (
 public_model_id UUID PRIMARY KEY, code TEXT NOT NULL UNIQUE, status TEXT NOT NULL
);
CREATE TABLE invocations (
 request_id UUID PRIMARY KEY, tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
 key_id UUID NOT NULL REFERENCES api_keys(key_id), public_model_id UUID NOT NULL REFERENCES public_models(public_model_id),
 protocol_family TEXT NOT NULL, status TEXT NOT NULL, started_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE invocation_attempts (
 attempt_id UUID PRIMARY KEY, request_id UUID NOT NULL REFERENCES invocations(request_id),
 deployment_id UUID NOT NULL REFERENCES channel_deployments(deployment_id),
 status TEXT NOT NULL, started_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE budget_accounts (
 account_id UUID PRIMARY KEY, tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
 currency CHAR(3) NOT NULL, limit_amount NUMERIC(28,12) NOT NULL CHECK(limit_amount >= 0),
 spent_amount NUMERIC(28,12) NOT NULL DEFAULT 0, reserved_amount NUMERIC(28,12) NOT NULL DEFAULT 0,
 CHECK(spent_amount >= 0 AND reserved_amount >= 0)
);
CREATE TABLE budget_reservations (
 reservation_id UUID PRIMARY KEY, account_id UUID NOT NULL REFERENCES budget_accounts(account_id),
 request_id UUID NOT NULL UNIQUE, amount NUMERIC(28,12) NOT NULL CHECK(amount >= 0),
 status TEXT NOT NULL, expires_at TIMESTAMPTZ NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE billing_ledger (
 entry_id UUID PRIMARY KEY, account_id UUID NOT NULL REFERENCES budget_accounts(account_id),
 request_id UUID NOT NULL, amount NUMERIC(28,12) NOT NULL,
 entry_type TEXT NOT NULL, idempotency_key TEXT NOT NULL UNIQUE,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE outbox_events (
 event_id UUID PRIMARY KEY, aggregate_type TEXT NOT NULL, aggregate_id UUID NOT NULL,
 event_type TEXT NOT NULL, payload JSONB NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 published_at TIMESTAMPTZ
);
CREATE INDEX idx_outbox_pending ON outbox_events(created_at) WHERE published_at IS NULL;
-- This is a structural starting point, not a complete financial schema or an executed migration.
