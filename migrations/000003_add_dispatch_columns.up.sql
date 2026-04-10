ALTER TABLE driver_profiles
    ADD COLUMN years_experience INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN max_load_kg NUMERIC(10,2) NOT NULL DEFAULT 0,
    ADD COLUMN equipment_type TEXT NOT NULL DEFAULT 'flatbed',
    ADD COLUMN is_online BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN is_available BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN current_latitude DOUBLE PRECISION,
    ADD COLUMN current_longitude DOUBLE PRECISION,
    ADD COLUMN last_location_at TIMESTAMP;

ALTER TABLE loads
    ADD COLUMN pickup_latitude DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN pickup_longitude DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN dropoff_latitude DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN dropoff_longitude DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN pickup_at TIMESTAMP NOT NULL DEFAULT NOW(),
    ADD COLUMN priority INTEGER NOT NULL DEFAULT 3,
    ADD COLUMN cargo_type TEXT NOT NULL DEFAULT 'general',
    ADD COLUMN equipment_type TEXT NOT NULL DEFAULT 'flatbed',
    ADD COLUMN quoted_price_kes NUMERIC(12,2) NOT NULL DEFAULT 0,
    ADD COLUMN matched_at TIMESTAMP,
    ADD COLUMN matching_score DOUBLE PRECISION,
    ADD COLUMN assignment_source TEXT NOT NULL DEFAULT 'manual';

ALTER TABLE loads
    DROP CONSTRAINT IF EXISTS loads_status_check;

ALTER TABLE loads
    ADD CONSTRAINT loads_status_check
    CHECK (status IN ('posted', 'matched', 'picked', 'in_transit', 'delivered'));

ALTER TABLE loads
    ADD CONSTRAINT loads_priority_check
    CHECK (priority BETWEEN 1 AND 5);

ALTER TABLE loads
    ADD CONSTRAINT loads_assignment_source_check
    CHECK (assignment_source IN ('manual', 'auto', 'admin_dispatch'));
