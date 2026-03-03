CREATE TABLE appointments (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    client_id UUID NOT NULL REFERENCES clients(id),
    master_id UUID NOT NULL,
    salon_id UUID NOT NULL,
    service_id UUID NOT NULL,
    service_name VARCHAR(200) NOT NULL,
    service_duration INTERVAL NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'requested',
    price_amount DECIMAL(12,2) NOT NULL,
    price_currency VARCHAR(3) DEFAULT 'RUB',
    source VARCHAR(20) NOT NULL,
    comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE master_schedules (
    id UUID PRIMARY KEY,
    master_id UUID NOT NULL,
    salon_id UUID NOT NULL,
    schedule_date DATE NOT NULL,
    work_start TIME NOT NULL,
    work_end TIME NOT NULL,
    break_start TIME,
    break_end TIME,
    booked_slots JSONB DEFAULT '[]',
    blocked_slots JSONB DEFAULT '[]',
    UNIQUE(master_id, schedule_date)
);
