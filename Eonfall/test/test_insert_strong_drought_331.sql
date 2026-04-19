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
             '33333333-3333-3333-3333-333333333331',
             'drought',
             'active',
             3,
             15340,
             25000,
             '{"food_loss_per_tick":50,"stability_loss_per_tick":5,"risk_gain_per_tick":5}'::jsonb,
             '{}'::jsonb,
             NOW(),
             NOW()
         );