CREATE TABLE worlds (
                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        name TEXT NOT NULL,
                        season_number INT NOT NULL DEFAULT 1,
                        state TEXT NOT NULL,
                        phase TEXT NOT NULL,
                        tick_rate_ms INT NOT NULL,
                        current_tick BIGINT NOT NULL DEFAULT 0,
                        config_version TEXT NOT NULL,
                        seed_value BIGINT NOT NULL,
                        started_at TIMESTAMPTZ NULL,
                        ends_at TIMESTAMPTZ NULL,
                        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);