-- 002_seed_delegations.sql

INSERT INTO delegations (id, user_id, client_id, scopes, remember)
VALUES
    (gen_random_uuid(), 'user123', 'client_abc', 'openid profile email', true),
    (gen_random_uuid(), 'user456', 'client_xyz', 'openid', false);
