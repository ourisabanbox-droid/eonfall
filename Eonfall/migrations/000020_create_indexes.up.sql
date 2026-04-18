CREATE INDEX idx_world_memberships_world_id ON world_memberships(world_id);
CREATE INDEX idx_world_memberships_user_id ON world_memberships(user_id);

CREATE INDEX idx_civilizations_world_id ON civilizations(world_id);
CREATE INDEX idx_civilizations_user_id ON civilizations(user_id);

CREATE INDEX idx_regions_world_id ON regions(world_id);
CREATE INDEX idx_regions_world_owner ON regions(world_id, owner_civilization_id);

CREATE INDEX idx_region_resource_stocks_world_id ON region_resource_stocks(world_id);
CREATE INDEX idx_region_resource_stocks_region_id ON region_resource_stocks(region_id);

CREATE INDEX idx_region_buildings_world_id ON region_buildings(world_id);
CREATE INDEX idx_region_buildings_region_id ON region_buildings(region_id);
CREATE INDEX idx_region_buildings_civ_id ON region_buildings(civilization_id);

CREATE INDEX idx_civilization_researches_world_id ON civilization_researches(world_id);
CREATE INDEX idx_civilization_researches_civ_id ON civilization_researches(civilization_id);

CREATE INDEX idx_world_actions_world_tick_state ON world_actions(world_id, target_tick, state);
CREATE INDEX idx_world_actions_civ_id ON world_actions(civilization_id);

CREATE INDEX idx_world_events_world_id ON world_events(world_id);
CREATE INDEX idx_world_events_state_tick ON world_events(world_id, state, starts_at_tick);

CREATE INDEX idx_event_participations_event_id ON world_event_participations(world_event_id);
CREATE INDEX idx_event_participations_world_id ON world_event_participations(world_id);

CREATE INDEX idx_world_catastrophes_world_id ON world_catastrophes(world_id);
CREATE INDEX idx_world_catastrophes_region_id ON world_catastrophes(region_id);
CREATE INDEX idx_world_catastrophes_state ON world_catastrophes(world_id, state);

CREATE INDEX idx_world_alerts_world_civ_created ON world_alerts(world_id, civilization_id, created_at DESC);
CREATE INDEX idx_world_snapshots_world_tick ON world_snapshots(world_id, tick DESC);
CREATE INDEX idx_world_activity_log_world_created ON world_activity_log(world_id, created_at DESC);