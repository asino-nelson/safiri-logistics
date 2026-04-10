CREATE TABLE payments (
    id UUID PRIMARY KEY,
    load_id UUID NOT NULL REFERENCES loads(id) ON DELETE CASCADE,
    customer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    provider_reference TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    amount_kes NUMERIC(12,2) NOT NULL CHECK (amount_kes > 0),
    currency TEXT NOT NULL DEFAULT 'KES',
    status TEXT NOT NULL CHECK (status IN ('initiated', 'paid', 'failed')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_load_id ON payments (load_id);
