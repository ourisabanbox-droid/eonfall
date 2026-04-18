CREATE TABLE world_memberships (
                                   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                   world_id UUID NOT NULL REFERENCES worlds(id) ON DELETE CASCADE,
                                   user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                   civilization_id UUID NULL,
                                   joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                   UNIQUE (world_id, user_id)
);