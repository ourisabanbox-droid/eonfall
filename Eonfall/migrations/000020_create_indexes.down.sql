DROP INDEX IF EXISTS idx_world_memberships_world_id;
DROP INDEX IF EXISTS idx_world_memberships_user_id;

DROP INDEX IF EXISTS idx_civilizations_world_id;
DROP INDEX IF EXISTS idx_civilizations_user_id;

DROP INDEX IF EXISTS idx_regions_world_id;
DROP INDEX IF EXISTS idx_regions_world_owner;

DROP INDEX IF EXISTS idx_region_resource_stocks_world_id;
DROP INDEX IF EXISTS idx_region_resource_stocks_region_id;

DROP INDEX IF EXISTS idx_region_buildings_world_id;
DROP INDEX IF EXISTS idx_region_buildings_region_id;
DROP INDEX IF EXISTS idx_region_buildings_civ_id;

DROP INDEX IF EXISTS idx_civilization_researches_world_id;
DROP INDEX IF EXISTS idx_civilization_researches_civ_id;

DROP INDEX IF EXISTS idx_world_actions_world_tick_state;
DROP INDEX IF EXISTS idx_world_actions_civ_id;

DROP INDEX IF EXISTS idx_world_events_world_id;
DROP INDEX IF EXISTS idx_world_events_state_tick;

DROP INDEX IF EXISTS idx_event_participations_event_id;
DROP INDEX IF EXISTS idx_event_participations_world_id;

DROP INDEX IF EXISTS idx_world_catastrophes_world_id;
DROP INDEX IF EXISTS idx_world_catastrophes_region_id;
DROP INDEX IF EXISTS idx_world_catastrophes_state;

DROP INDEX IF EXISTS idx_world_alerts_world_civ_created;
DROP INDEX IF EXISTS idx_world_snapshots_world_tick;
DROP INDEX IF EXISTS idx_world_activity_log_world_created;