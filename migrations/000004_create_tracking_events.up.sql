CREATE TABLE tracking_events (
    id UUID PRIMARY KEY,
    load_id UUID NOT NULL REFERENCES loads(id) ON DELETE CASCADE,
    driver_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    speed_kph DOUBLE PRECISION NOT NULL DEFAULT 0,
    heading_degrees DOUBLE PRECISION NOT NULL DEFAULT 0,
    recorded_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tracking_events_load_recorded_at
    ON tracking_events (load_id, recorded_at DESC);
