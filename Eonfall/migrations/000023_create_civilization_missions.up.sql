CREATE TABLE civilization_missions (
                                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                       world_id UUID NOT NULL REFERENCES worlds(id) ON DELETE CASCADE,
                                       civilization_id UUID NOT NULL REFERENCES civilizations(id) ON DELETE CASCADE,
                                       mission_type TEXT NOT NULL,
                                       target_region_id UUID REFERENCES regions(id) ON DELETE CASCADE,
                                       status TEXT NOT NULL,
                                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_civilization_missions_active
    ON civilization_missions (world_id, civilization_id, mission_type, target_region_id);