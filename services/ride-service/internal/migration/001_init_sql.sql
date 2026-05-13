
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE ride_status AS ENUM (
    'REQUESTED',
    'MATCHING',
    'ACCEPTED',
    'DRIVER_ARRIVING',
    'RIDE_STARTED',
    'COMPLETED',
    'CANCELLED'
);

CREATE TABLE IF NOT EXISTS rides (
    id               UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    rider_id         TEXT        NOT NULL,
    driver_id        TEXT,
    pickup_latitude  DOUBLE PRECISION NOT NULL,
    pickup_longitude DOUBLE PRECISION NOT NULL,
    pickup_address   TEXT        NOT NULL,
    drop_latitude    DOUBLE PRECISION NOT NULL,
    drop_longitude   DOUBLE PRECISION NOT NULL,
    drop_address     TEXT        NOT NULL,
    status           ride_status NOT NULL DEFAULT 'REQUESTED',
    estimated_fare   DOUBLE PRECISION NOT NULL DEFAULT 0,
    actual_fare      DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at       TIMESTAMPTZ,
    completed_at     TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_rides_rider_id  ON rides (rider_id);
CREATE INDEX IF NOT EXISTS idx_rides_driver_id ON rides (driver_id);
CREATE INDEX IF NOT EXISTS idx_rides_status    ON rides (status);

-- Auto-update updated_at on every row change
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER rides_set_updated_at
    BEFORE UPDATE ON rides
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
