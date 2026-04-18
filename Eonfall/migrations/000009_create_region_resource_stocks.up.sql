CREATE TABLE region_resource_stocks (
                                        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                        world_id UUID NOT NULL REFERENCES worlds(id) ON DELETE CASCADE,
                                        region_id UUID NOT NULL REFERENCES regions(id) ON DELETE CASCADE,
                                        resource_type TEXT NOT NULL,
                                        stock BIGINT NOT NULL DEFAULT 0,
                                        production_rate INT NOT NULL DEFAULT 0,
                                        consumption_rate INT NOT NULL DEFAULT 0,
                                        capacity INT NOT NULL DEFAULT 0,
                                        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                        UNIQUE (region_id, resource_type)
);