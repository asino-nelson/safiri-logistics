ALTER TABLE loads
    DROP CONSTRAINT IF EXISTS loads_assignment_source_check;

ALTER TABLE loads
    DROP CONSTRAINT IF EXISTS loads_priority_check;

ALTER TABLE loads
    DROP CONSTRAINT IF EXISTS loads_status_check;

ALTER TABLE loads
    ADD CONSTRAINT loads_status_check
    CHECK (status IN ('posted', 'picked', 'in_transit', 'delivered'));

ALTER TABLE loads
    DROP COLUMN IF EXISTS assignment_source,
    DROP COLUMN IF EXISTS matching_score,
    DROP COLUMN IF EXISTS matched_at,
    DROP COLUMN IF EXISTS quoted_price_kes,
    DROP COLUMN IF EXISTS equipment_type,
    DROP COLUMN IF EXISTS cargo_type,
    DROP COLUMN IF EXISTS priority,
    DROP COLUMN IF EXISTS pickup_at,
    DROP COLUMN IF EXISTS dropoff_longitude,
    DROP COLUMN IF EXISTS dropoff_latitude,
    DROP COLUMN IF EXISTS pickup_longitude,
    DROP COLUMN IF EXISTS pickup_latitude;

ALTER TABLE driver_profiles
    DROP COLUMN IF EXISTS last_location_at,
    DROP COLUMN IF EXISTS current_longitude,
    DROP COLUMN IF EXISTS current_latitude,
    DROP COLUMN IF EXISTS is_available,
    DROP COLUMN IF EXISTS is_online,
    DROP COLUMN IF EXISTS equipment_type,
    DROP COLUMN IF EXISTS max_load_kg,
    DROP COLUMN IF EXISTS years_experience;
