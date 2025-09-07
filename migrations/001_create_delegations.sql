-- 001_create_delegations.sql

BEGIN;

CREATE TABLE IF NOT EXISTS delegations (
    id          UUID PRIMARY KEY,
    user_id     TEXT NOT NULL,
    client_id   TEXT NOT NULL,
    scopes      TEXT NOT NULL, -- space-separated list of scopes
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at  TIMESTAMPTZ,   -- optional: null means no expiry
    remember    BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT delegations_user_client UNIQUE (user_id, client_id)
);

-- Index to speed up lookups by user/client
CREATE INDEX IF NOT EXISTS idx_delegations_user_client
    ON delegations (user_id, client_id);

COMMIT;
