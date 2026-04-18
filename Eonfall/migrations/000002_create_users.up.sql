CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       email TEXT NOT NULL UNIQUE,
                       password_hash TEXT NOT NULL,
                       display_name TEXT NOT NULL,
                       status TEXT NOT NULL DEFAULT 'active',
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);