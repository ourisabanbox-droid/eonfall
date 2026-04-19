INSERT INTO world_catastrophes (
    id,
    world_id,
    region_id,
    catastrophe_type,
    state,
    severity,
    starts_at_tick,
    ends_at_tick,
    payload_json,
    result_json,
    created_at,
    updated_at
) VALUES (
             gen_random_uuid(),
             '11111111-1111-1111-1111-111111111111',
             '33333333-3333-3333-3333-333333333334',
             'drought',
             'active',
             2,
             15340,
             15400,
             '{"food_loss_per_tick":6,"stability_loss_per_tick":2,"risk_gain_per_tick":2}'::jsonb,
             '{}'::jsonb,
             NOW(),
             NOW()
         );