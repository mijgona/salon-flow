CREATE TABLE certificates (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    purchased_by UUID REFERENCES clients(id),
    activated_by UUID REFERENCES clients(id),
    balance_amount DECIMAL(12,2) NOT NULL,
    balance_currency VARCHAR(3) DEFAULT 'RUB',
    status VARCHAR(20) NOT NULL DEFAULT 'created',
    activated_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE outbox (
    id UUID PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);
CREATE INDEX idx_outbox_pending ON outbox(created_at) WHERE processed_at IS NULL;
