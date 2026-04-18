CREATE TABLE diplomatic_pacts (
                                  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                  world_id UUID NOT NULL REFERENCES worlds(id) ON DELETE CASCADE,
                                  civilization_a_id UUID NOT NULL REFERENCES civilizations(id) ON DELETE CASCADE,
                                  civilization_b_id UUID NOT NULL REFERENCES civilizations(id) ON DELETE CASCADE,
                                  pact_type TEXT NOT NULL,
                                  state TEXT NOT NULL DEFAULT 'active',
                                  trust_score INT NOT NULL DEFAULT 0,
                                  started_tick BIGINT NULL,
                                  ended_tick BIGINT NULL,
                                  metadata_json JSONB NOT NULL DEFAULT '{}'::jsonb,
                                  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                  CHECK (civilization_a_id <> civilization_b_id)
);