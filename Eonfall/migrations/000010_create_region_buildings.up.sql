CREATE TABLE region_buildings (
                                  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                  world_id UUID NOT NULL REFERENCES worlds(id) ON DELETE CASCADE,
                                  region_id UUID NOT NULL REFERENCES regions(id) ON DELETE CASCADE,
                                  civilization_id UUID NOT NULL REFERENCES civilizations(id) ON DELETE CASCADE,
                                  building_type TEXT NOT NULL,
                                  level INT NOT NULL DEFAULT 1,
                                  state TEXT NOT NULL,
                                  started_tick BIGINT NOT NULL,
                                  completed_tick BIGINT NULL,
                                  durability INT NOT NULL DEFAULT 100,
                                  metadata_json JSONB NOT NULL DEFAULT '{}'::jsonb,
                                  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);