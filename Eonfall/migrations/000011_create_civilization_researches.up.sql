CREATE TABLE civilization_researches (
                                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                         world_id UUID NOT NULL REFERENCES worlds(id) ON DELETE CASCADE,
                                         civilization_id UUID NOT NULL REFERENCES civilizations(id) ON DELETE CASCADE,
                                         research_type TEXT NOT NULL,
                                         state TEXT NOT NULL,
                                         progress INT NOT NULL DEFAULT 0,
                                         started_tick BIGINT NULL,
                                         completed_tick BIGINT NULL,
                                         created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                         updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                         UNIQUE (civilization_id, research_type)
);