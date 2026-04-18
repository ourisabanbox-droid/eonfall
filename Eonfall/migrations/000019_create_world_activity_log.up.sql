CREATE TABLE world_activity_log (
                                    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                    world_id UUID NOT NULL REFERENCES worlds(id) ON DELETE CASCADE,
                                    civilization_id UUID NULL REFERENCES civilizations(id) ON DELETE CASCADE,
                                    region_id UUID NULL REFERENCES regions(id) ON DELETE CASCADE,
                                    activity_type TEXT NOT NULL,
                                    payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
                                    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);