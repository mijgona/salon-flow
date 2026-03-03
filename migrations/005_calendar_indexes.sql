-- Calendar query optimization indexes
CREATE INDEX idx_appointments_tenant_start ON appointments(tenant_id, start_time);
CREATE INDEX idx_appointments_master_start ON appointments(master_id, start_time);
CREATE INDEX idx_appointments_salon_start ON appointments(salon_id, start_time);
CREATE INDEX idx_appointments_status ON appointments(status);
