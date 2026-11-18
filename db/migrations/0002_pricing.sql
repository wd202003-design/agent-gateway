-- P-002.5 draft commercial pricing and invoice schema.
CREATE TABLE price_plans (
 plan_id UUID NOT NULL,
 version INTEGER NOT NULL CHECK(version>0),
 currency CHAR(3) NOT NULL,
 effective_from TIMESTAMPTZ NOT NULL,
 effective_until TIMESTAMPTZ,
 status TEXT NOT NULL,
 PRIMARY KEY(plan_id,version)
);
CREATE TABLE price_rules (
 rule_id UUID PRIMARY KEY,
 plan_id UUID NOT NULL,
 plan_version INTEGER NOT NULL,
 public_model_id UUID REFERENCES public_models(public_model_id),
 metric TEXT NOT NULL,
 unit TEXT NOT NULL,
 unit_price NUMERIC(28,12) NOT NULL CHECK(unit_price>=0),
 conditions JSONB NOT NULL DEFAULT '{}'::jsonb,
 FOREIGN KEY(plan_id,plan_version) REFERENCES price_plans(plan_id,version)
);
CREATE TABLE price_assignments (
 assignment_id UUID PRIMARY KEY,
 tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
 plan_id UUID NOT NULL,
 plan_version INTEGER NOT NULL,
 priority INTEGER NOT NULL DEFAULT 100,
 effective_from TIMESTAMPTZ NOT NULL,
 effective_until TIMESTAMPTZ,
 FOREIGN KEY(plan_id,plan_version) REFERENCES price_plans(plan_id,version)
);
CREATE TABLE charge_records (
 charge_id UUID PRIMARY KEY,
 request_id UUID NOT NULL REFERENCES invocations(request_id),
 plan_id UUID NOT NULL,
 plan_version INTEGER NOT NULL,
 amount NUMERIC(28,12) NOT NULL CHECK(amount>=0),
 currency CHAR(3) NOT NULL,
 status TEXT NOT NULL,
 idempotency_key TEXT NOT NULL UNIQUE,
 FOREIGN KEY(plan_id,plan_version) REFERENCES price_plans(plan_id,version)
);
CREATE TABLE invoices (
 invoice_id UUID PRIMARY KEY,
 tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
 period_start TIMESTAMPTZ NOT NULL,
 period_end TIMESTAMPTZ NOT NULL,
 currency CHAR(3) NOT NULL,
 total NUMERIC(28,12) NOT NULL,
 status TEXT NOT NULL,
 CHECK(period_end>period_start)
);
-- Immutable ledger and reconciled adjustments must be implemented and audited
-- before any real customer charge or invoice is issued.
