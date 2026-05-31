-- +migrate Up
CREATE TABLE IF NOT EXISTS factories (
    id         UUID        NOT NULL PRIMARY KEY,
    name       TEXT        NOT NULL,
    map3d      JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER factories_set_updated_at
    BEFORE UPDATE ON factories
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- +migrate Down
DROP TRIGGER IF EXISTS factories_set_updated_at ON factories;
DROP FUNCTION IF EXISTS set_updated_at();
DROP TABLE IF EXISTS factories;
