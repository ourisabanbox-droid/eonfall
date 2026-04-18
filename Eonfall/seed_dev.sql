-- Users
INSERT INTO users (id, email, password_hash, display_name)
VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'solari@test.local', 'devhash', 'Solari Player'),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2', 'azmar@test.local', 'devhash', 'Azmar Player'),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3', 'verdeenne@test.local', 'devhash', 'Verdeenne Player');

-- World
INSERT INTO worlds (
    id, name, season_number, state, phase, tick_rate_ms, current_tick, config_version, seed_value
)
VALUES (
           '11111111-1111-1111-1111-111111111111',
           'Aurelios',
           1,
           'running',
           'foundation',
           1000,
           0,
           'v1',
           424242
       );

-- Civilizations
INSERT INTO civilizations (
    id, world_id, user_id, template_id, name, status,
    population, cohesion, influence, science_stock, credit_stock,
    food_stock, energy_stock, materials_stock, military_score, victory_score
)
VALUES
    (
        '22222222-2222-2222-2222-222222222221',
        '11111111-1111-1111-1111-111111111111',
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1',
        'solari',
        'Dominion Solari',
        'active',
        1000, 100, 10, 50, 100, 120, 100, 100, 10, 0
    ),
    (
        '22222222-2222-2222-2222-222222222222',
        '11111111-1111-1111-1111-111111111111',
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2',
        'azmar',
        'Synode d''Azmar',
        'active',
        1000, 120, 20, 40, 100, 120, 80, 90, 8, 0
    ),
    (
        '22222222-2222-2222-2222-222222222223',
        '11111111-1111-1111-1111-111111111111',
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3',
        'verdeenne',
        'Ligue Verdéenne',
        'active',
        1000, 110, 15, 35, 100, 150, 70, 80, 7, 0
    );

-- Memberships
INSERT INTO world_memberships (world_id, user_id, civilization_id)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', '22222222-2222-2222-2222-222222222221'),
    ('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2', '22222222-2222-2222-2222-222222222222'),
    ('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa3', '22222222-2222-2222-2222-222222222223');

-- Regions
INSERT INTO regions (
    id, world_id, q, r, biome, terrain_type, climate_type, owner_civilization_id,
    population, stability, development_level, pollution,
    drought_risk, revolt_risk, fire_risk, seismic_risk,
    energy_fragility, logistic_fragility, has_capital
)
VALUES
    ('33333333-3333-3333-3333-333333333331','11111111-1111-1111-1111-111111111111',0,0,'plains','flat','temperate','22222222-2222-2222-2222-222222222221',500,100,1,5,10,5,5,5,10,10,true),
    ('33333333-3333-3333-3333-333333333332','11111111-1111-1111-1111-111111111111',1,0,'hills','elevated','temperate','22222222-2222-2222-2222-222222222221',300,95,1,10,15,5,5,10,15,10,false),
    ('33333333-3333-3333-3333-333333333333','11111111-1111-1111-1111-111111111111',2,0,'desert','rough','arid',NULL,100,90,0,0,30,5,5,5,10,20,false),

    ('33333333-3333-3333-3333-333333333334','11111111-1111-1111-1111-111111111111',0,1,'forest','rough','humid','22222222-2222-2222-2222-222222222222',500,100,1,5,5,5,15,5,5,10,true),
    ('33333333-3333-3333-3333-333333333335','11111111-1111-1111-1111-111111111111',1,1,'plains','flat','temperate','22222222-2222-2222-2222-222222222222',300,98,1,3,10,5,5,5,5,10,false),
    ('33333333-3333-3333-3333-333333333336','11111111-1111-1111-1111-111111111111',2,1,'mountains','elevated','cold',NULL,100,85,0,0,10,5,5,20,10,25,false),

    ('33333333-3333-3333-3333-333333333337','11111111-1111-1111-1111-111111111111',0,2,'wetlands','rough','humid','22222222-2222-2222-2222-222222222223',500,100,1,2,5,5,10,5,5,10,true),
    ('33333333-3333-3333-3333-333333333338','11111111-1111-1111-1111-111111111111',1,2,'forest','flat','humid','22222222-2222-2222-2222-222222222223',300,97,1,4,5,5,12,5,5,10,false),
    ('33333333-3333-3333-3333-333333333339','11111111-1111-1111-1111-111111111111',2,2,'plains','flat','temperate',NULL,100,88,0,0,10,5,5,5,10,15,false);

-- Capitals
UPDATE civilizations
SET capital_region_id = '33333333-3333-3333-3333-333333333331'
WHERE id = '22222222-2222-2222-2222-222222222221';

UPDATE civilizations
SET capital_region_id = '33333333-3333-3333-3333-333333333334'
WHERE id = '22222222-2222-2222-2222-222222222222';

UPDATE civilizations
SET capital_region_id = '33333333-3333-3333-3333-333333333337'
WHERE id = '22222222-2222-2222-2222-222222222223';

-- Adjacencies
INSERT INTO region_adjacencies (region_id, adjacent_region_id) VALUES
                                                                   ('33333333-3333-3333-3333-333333333331','33333333-3333-3333-3333-333333333332'),
                                                                   ('33333333-3333-3333-3333-333333333331','33333333-3333-3333-3333-333333333334'),
                                                                   ('33333333-3333-3333-3333-333333333332','33333333-3333-3333-3333-333333333333'),
                                                                   ('33333333-3333-3333-3333-333333333332','33333333-3333-3333-3333-333333333335'),
                                                                   ('33333333-3333-3333-3333-333333333334','33333333-3333-3333-3333-333333333335'),
                                                                   ('33333333-3333-3333-3333-333333333334','33333333-3333-3333-3333-333333333337'),
                                                                   ('33333333-3333-3333-3333-333333333335','33333333-3333-3333-3333-333333333336'),
                                                                   ('33333333-3333-3333-3333-333333333335','33333333-3333-3333-3333-333333333338'),
                                                                   ('33333333-3333-3333-3333-333333333337','33333333-3333-3333-3333-333333333338'),
                                                                   ('33333333-3333-3333-3333-333333333338','33333333-3333-3333-3333-333333333339');

INSERT INTO region_adjacencies (region_id, adjacent_region_id)
SELECT adjacent_region_id, region_id
FROM region_adjacencies
    ON CONFLICT DO NOTHING;

-- Resource stocks
INSERT INTO region_resource_stocks (
    world_id, region_id, resource_type, stock, production_rate, consumption_rate, capacity
)
VALUES
    ('11111111-1111-1111-1111-111111111111','33333333-3333-3333-3333-333333333331','food',100,8,1,500),
    ('11111111-1111-1111-1111-111111111111','33333333-3333-3333-3333-333333333331','energy',80,6,1,500),
    ('11111111-1111-1111-1111-111111111111','33333333-3333-3333-3333-333333333331','knowledge',50,4,0,500),

    ('11111111-1111-1111-1111-111111111111','33333333-3333-3333-3333-333333333334','cohesion',80,5,0,500),
    ('11111111-1111-1111-1111-111111111111','33333333-3333-3333-3333-333333333334','knowledge',40,3,0,500),
    ('11111111-1111-1111-1111-111111111111','33333333-3333-3333-3333-333333333334','food',100,6,1,500),

    ('11111111-1111-1111-1111-111111111111','33333333-3333-3333-3333-333333333337','food',120,9,1,500),
    ('11111111-1111-1111-1111-111111111111','33333333-3333-3333-3333-333333333337','cohesion',70,4,0,500),
    ('11111111-1111-1111-1111-111111111111','33333333-3333-3333-3333-333333333337','materials',60,4,0,500);

-- Buildings
INSERT INTO region_buildings (
    world_id, region_id, civilization_id, building_type, level, state, started_tick, completed_tick, durability
)
VALUES
    ('11111111-1111-1111-1111-111111111111','33333333-3333-3333-3333-333333333331','22222222-2222-2222-2222-222222222221','power_plant',1,'active',0,0,100),
    ('11111111-1111-1111-1111-111111111111','33333333-3333-3333-3333-333333333331','22222222-2222-2222-2222-222222222221','lab',1,'active',0,0,100),
    ('11111111-1111-1111-1111-111111111111','33333333-3333-3333-3333-333333333334','22222222-2222-2222-2222-222222222222','forum',1,'active',0,0,100),
    ('11111111-1111-1111-1111-111111111111','33333333-3333-3333-3333-333333333337','22222222-2222-2222-2222-222222222223','farm',1,'active',0,0,100);

-- Researches
INSERT INTO civilization_researches (
    world_id, civilization_id, research_type, state, progress, started_tick
)
VALUES
    ('11111111-1111-1111-1111-111111111111','22222222-2222-2222-2222-222222222221','method_analytique','in_progress',20,0),
    ('11111111-1111-1111-1111-111111111111','22222222-2222-2222-2222-222222222222','administration_regionale','in_progress',15,0),
    ('11111111-1111-1111-1111-111111111111','22222222-2222-2222-2222-222222222223','agriculture_rationalisee','in_progress',25,0);