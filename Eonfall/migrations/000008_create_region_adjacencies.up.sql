CREATE TABLE region_adjacencies (
                                    region_id UUID NOT NULL REFERENCES regions(id) ON DELETE CASCADE,
                                    adjacent_region_id UUID NOT NULL REFERENCES regions(id) ON DELETE CASCADE,
                                    PRIMARY KEY (region_id, adjacent_region_id),
                                    CHECK (region_id <> adjacent_region_id)
);