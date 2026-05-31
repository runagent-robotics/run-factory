-- +migrate Up
CREATE TABLE IF NOT EXISTS factories (
    id         TEXT        NOT NULL PRIMARY KEY,
    name       TEXT        NOT NULL,
    map3d      JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS factories;
