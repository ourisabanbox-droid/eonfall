CREATE TABLE world_event_participations (
                                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                            world_event_id UUID NOT NULL REFERENCES world_events(id) ON DELETE CASCADE,
                                            world_id UUID NOT NULL REFERENCES worlds(id) ON DELETE CASCADE,
                                            civilization_id UUID NOT NULL REFERENCES civilizations(id) ON DELETE CASCADE,
                                            participation_state TEXT NOT NULL DEFAULT 'registered',
                                            contribution_json JSONB NOT NULL DEFAULT '{}'::jsonb,
                                            score INT NOT NULL DEFAULT 0,
                                            reward_json JSONB NOT NULL DEFAULT '{}'::jsonb,
                                            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                            updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                            UNIQUE (world_event_id, civilization_id)
);