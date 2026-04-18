CREATE TABLE world_actions (
                               id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                               world_id UUID NOT NULL REFERENCES worlds(id) ON DELETE CASCADE,
                               civilization_id UUID NOT NULL REFERENCES civilizations(id) ON DELETE CASCADE,
                               user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               action_type TEXT NOT NULL,
                               state TEXT NOT NULL DEFAULT 'pending',
                               target_tick BIGINT NOT NULL,
                               payload_json JSONB NOT NULL,
                               rejection_reason TEXT NULL,
                               created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                               applied_at TIMESTAMPTZ NULL
);