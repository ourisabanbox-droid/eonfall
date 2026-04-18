CREATE TABLE world_catastrophes (
                                    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                    world_id UUID NOT NULL REFERENCES worlds(id) ON DELETE CASCADE,
                                    region_id UUID NOT NULL REFERENCES regions(id) ON DELETE CASCADE,
                                    catastrophe_type TEXT NOT NULL,
                                    state TEXT NOT NULL,
                                    severity INT NOT NULL DEFAULT 1,
                                    starts_at_tick BIGINT NOT NULL,
                                    ends_at_tick BIGINT NULL,
                                    payload_json JSONB NOT NULL DEFAULT '{}'::jsonb,
                                    result_json JSONB NOT NULL DEFAULT '{}'::jsonb,
                                    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);