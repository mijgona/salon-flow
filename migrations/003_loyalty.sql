CREATE TABLE loyalty_accounts (
    id UUID PRIMARY KEY,
    client_id UUID NOT NULL UNIQUE REFERENCES clients(id),
    tenant_id UUID NOT NULL,
    tier VARCHAR(10) NOT NULL DEFAULT 'Bronze',
    balance INT NOT NULL DEFAULT 0,
    lifetime_points INT NOT NULL DEFAULT 0
);

CREATE TABLE points_transactions (
    id UUID PRIMARY KEY,
    loyalty_account_id UUID NOT NULL REFERENCES loyalty_accounts(id),
    amount INT NOT NULL,
    type VARCHAR(20) NOT NULL,
    reason VARCHAR(100),
    related_entity_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE referrals (
    id UUID PRIMARY KEY,
    loyalty_account_id UUID NOT NULL REFERENCES loyalty_accounts(id),
    referred_client_id UUID NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    bonus_earned INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
