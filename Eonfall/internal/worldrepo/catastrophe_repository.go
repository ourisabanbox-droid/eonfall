package worldrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"project-eonfall/internal/world"
)

type CatastropheRepository struct {
	db *pgxpool.Pool
}

func NewCatastropheRepository(db *pgxpool.Pool) *CatastropheRepository {
	return &CatastropheRepository{db: db}
}

func (r *CatastropheRepository) Insert(ctx context.Context, c *world.Catastrophe) error {
	payloadBytes, err := json.Marshal(c.Payload)
	if err != nil {
		return fmt.Errorf("marshal catastrophe payload: %w", err)
	}

	resultBytes, err := json.Marshal(c.Result)
	if err != nil {
		return fmt.Errorf("marshal catastrophe result: %w", err)
	}

	const q = `
INSERT INTO world_catastrophes (
	id, world_id, region_id, catastrophe_type, state, severity,
	starts_at_tick, ends_at_tick, payload_json, result_json, created_at, updated_at
) VALUES (
	$1, $2, $3, $4, $5, $6,
	$7, $8, $9::jsonb, $10::jsonb, $11, $12
)
`

	_, err = r.db.Exec(
		ctx,
		q,
		c.ID,
		c.WorldID,
		c.RegionID,
		c.CatastropheType,
		c.State,
		c.Severity,
		c.StartsAtTick,
		c.EndsAtTick,
		string(payloadBytes),
		string(resultBytes),
		c.CreatedAt,
		c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert catastrophe exec: %w", err)
	}

	return nil
}

func (r *CatastropheRepository) ListActiveByWorldID(ctx context.Context, worldID uuid.UUID) ([]*world.Catastrophe, error) {
	const q = `
SELECT id, world_id, region_id, catastrophe_type, state, severity,
       starts_at_tick, ends_at_tick, payload_json, result_json, created_at, updated_at
FROM world_catastrophes
WHERE world_id = $1
  AND state = 'active'
ORDER BY created_at ASC
`
	rows, err := r.db.Query(ctx, q, worldID)
	if err != nil {
		return nil, fmt.Errorf("ListActiveByWorldID query: %w", err)
	}
	defer rows.Close()

	var out []*world.Catastrophe
	for rows.Next() {
		var c world.Catastrophe
		var rawPayload []byte
		var rawResult []byte

		if err := rows.Scan(
			&c.ID, &c.WorldID, &c.RegionID, &c.CatastropheType, &c.State, &c.Severity,
			&c.StartsAtTick, &c.EndsAtTick, &rawPayload, &rawResult, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("ListActiveByWorldID scan: %w", err)
		}

		c.Payload = map[string]any{}
		c.Result = map[string]any{}

		if len(rawPayload) > 0 {
			if err := json.Unmarshal(rawPayload, &c.Payload); err != nil {
				return nil, fmt.Errorf("unmarshal catastrophe payload: %w", err)
			}
		}
		if len(rawResult) > 0 {
			if err := json.Unmarshal(rawResult, &c.Result); err != nil {
				return nil, fmt.Errorf("unmarshal catastrophe result: %w", err)
			}
		}

		out = append(out, &c)
	}

	return out, rows.Err()
}

func (r *CatastropheRepository) ListActiveByRegionID(ctx context.Context, regionID uuid.UUID) ([]*world.Catastrophe, error) {
	const q = `
SELECT id, world_id, region_id, catastrophe_type, state, severity,
       starts_at_tick, ends_at_tick, payload_json, result_json, created_at, updated_at
FROM world_catastrophes
WHERE region_id = $1
  AND state = 'active'
ORDER BY created_at ASC
`
	rows, err := r.db.Query(ctx, q, regionID)
	if err != nil {
		return nil, fmt.Errorf("ListActiveByRegionID query: %w", err)
	}
	defer rows.Close()

	var out []*world.Catastrophe
	for rows.Next() {
		var c world.Catastrophe
		var rawPayload []byte
		var rawResult []byte

		if err := rows.Scan(
			&c.ID, &c.WorldID, &c.RegionID, &c.CatastropheType, &c.State, &c.Severity,
			&c.StartsAtTick, &c.EndsAtTick, &rawPayload, &rawResult, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("ListActiveByRegionID scan: %w", err)
		}

		c.Payload = map[string]any{}
		c.Result = map[string]any{}

		if len(rawPayload) > 0 {
			_ = json.Unmarshal(rawPayload, &c.Payload)
		}
		if len(rawResult) > 0 {
			_ = json.Unmarshal(rawResult, &c.Result)
		}

		out = append(out, &c)
	}

	return out, rows.Err()
}

func (r *CatastropheRepository) MarkResolved(ctx context.Context, id uuid.UUID) error {
	const q = `
UPDATE world_catastrophes
SET state = 'resolved',
    updated_at = $2
WHERE id = $1
`
	_, err := r.db.Exec(ctx, q, id, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("MarkResolved exec: %w", err)
	}
	return nil
}
