CREATE TABLE world_snapshots (
                                 id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                 world_id UUID NOT NULL REFERENCES worlds(id) ON DELETE CASCADE,
                                 tick BIGINT NOT NULL,
                                 snapshot_path TEXT NOT NULL,
                                 checksum TEXT NULL,
                                 created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                 UNIQUE (world_id, tick)
);