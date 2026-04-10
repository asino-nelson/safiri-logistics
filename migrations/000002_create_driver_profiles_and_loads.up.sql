CREATE TABLE driver_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    national_id TEXT NOT NULL,
    license_number TEXT NOT NULL,
    truck_registration TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending', 'approved', 'rejected')),
    rejection_reason TEXT,
    submitted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    reviewed_at TIMESTAMP,
    verified_at TIMESTAMP
);

CREATE TABLE loads (
    id UUID PRIMARY KEY,
    poster_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    assigned_driver_id UUID REFERENCES users(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    origin TEXT NOT NULL,
    destination TEXT NOT NULL,
    weight_kg NUMERIC(10,2) NOT NULL CHECK (weight_kg > 0),
    status TEXT NOT NULL CHECK (status IN ('posted', 'picked', 'in_transit', 'delivered')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    picked_at TIMESTAMP,
    delivered_at TIMESTAMP
);
