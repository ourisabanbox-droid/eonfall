CREATE TABLE civilization_trajectories (
                                           civilization_id uuid PRIMARY KEY REFERENCES civilizations(id) ON DELETE CASCADE,
                                           world_id uuid NOT NULL REFERENCES worlds(id) ON DELETE CASCADE,
                                           resilience_score integer NOT NULL DEFAULT 0,
                                           expansion_score integer NOT NULL DEFAULT 0,
                                           influence_score integer NOT NULL DEFAULT 0,
                                           science_score integer NOT NULL DEFAULT 0
);

CREATE INDEX idx_civilization_trajectories_world_id
    ON civilization_trajectories(world_id);