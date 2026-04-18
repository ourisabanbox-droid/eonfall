CREATE TABLE world_alerts (
                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              world_id UUID NOT NULL REFERENCES worlds(id) ON DELETE CASCADE,
                              civilization_id UUID NULL REFERENCES civilizations(id) ON DELETE CASCADE,
                              region_id UUID NULL REFERENCES regions(id) ON DELETE CASCADE,
                              alert_type TEXT NOT NULL,
                              severity TEXT NOT NULL,
                              title TEXT NOT NULL,
                              message TEXT NOT NULL,
                              payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
                              created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                              read_at TIMESTAMPTZ NULL
);