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

type ActionRepository struct {
	db *pgxpool.Pool
}

func NewActionRepository(db *pgxpool.Pool) *ActionRepository {
	return &ActionRepository{db: db}
}

func (r *ActionRepository) Enqueue(ctx context.Context, action *world.WorldAction) error {
	payloadBytes, err := json.Marshal(action.Payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	const q = `
INSERT INTO world_actions (
    id, world_id, civilization_id, user_id, action_type, state,
    target_tick, payload_json, rejection_reason, created_at, applied_at
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8::jsonb, $9, $10, $11
)
`
	_, err = r.db.Exec(
		ctx,
		q,
		action.ID,
		action.WorldID,
		action.CivilizationID,
		action.UserID,
		action.ActionType,
		action.State,
		action.TargetTick,
		string(payloadBytes),
		action.RejectionReason,
		action.CreatedAt,
		action.AppliedAt,
	)
	if err != nil {
		return fmt.Errorf("enqueue action exec: %w", err)
	}

	return nil
}

func (r *ActionRepository) ListPendingForTick(ctx context.Context, worldID uuid.UUID, tick int64) ([]*world.WorldAction, error) {
	const q = `
SELECT id, world_id, civilization_id, user_id, action_type, state,
       target_tick, payload_json, rejection_reason, created_at, applied_at
FROM world_actions
WHERE world_id = $1
  AND target_tick <= $2
  AND state = 'pending'
ORDER BY created_at ASC
`
	rows, err := r.db.Query(ctx, q, worldID, tick)
	if err != nil {
		return nil, fmt.Errorf("ListPendingForTick query: %w", err)
	}
	defer rows.Close()

	var out []*world.WorldAction

	for rows.Next() {
		var a world.WorldAction
		var rawPayload []byte

		err := rows.Scan(
			&a.ID,
			&a.WorldID,
			&a.CivilizationID,
			&a.UserID,
			&a.ActionType,
			&a.State,
			&a.TargetTick,
			&rawPayload,
			&a.RejectionReason,
			&a.CreatedAt,
			&a.AppliedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ListPendingForTick scan: %w", err)
		}

		if len(rawPayload) > 0 {
			if err := json.Unmarshal(rawPayload, &a.Payload); err != nil {
				return nil, fmt.Errorf("unmarshal payload: %w", err)
			}
		} else {
			a.Payload = map[string]any{}
		}

		out = append(out, &a)
	}

	return out, rows.Err()
}

func (r *ActionRepository) MarkApplied(ctx context.Context, actionID uuid.UUID) error {
	const q = `
UPDATE world_actions
SET state = 'applied',
    applied_at = $2
WHERE id = $1
`
	_, err := r.db.Exec(ctx, q, actionID, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("MarkApplied exec: %w", err)
	}
	return nil
}

func (r *ActionRepository) MarkRejected(ctx context.Context, actionID uuid.UUID, reason string) error {
	const q = `
UPDATE world_actions
SET state = 'rejected',
    rejection_reason = $2,
    applied_at = $3
WHERE id = $1
`
	_, err := r.db.Exec(ctx, q, actionID, reason, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("MarkRejected exec: %w", err)
	}
	return nil
}

func (r *ActionRepository) ListByWorldID(ctx context.Context, worldID uuid.UUID, limit int) ([]*world.WorldAction, error) {
	if limit <= 0 {
		limit = 50
	}

	const q = `
SELECT id, world_id, civilization_id, user_id, action_type, state,
       target_tick, payload_json, rejection_reason, created_at, applied_at
FROM world_actions
WHERE world_id = $1
ORDER BY created_at DESC
LIMIT $2
`
	rows, err := r.db.Query(ctx, q, worldID, limit)
	if err != nil {
		return nil, fmt.Errorf("ListByWorldID query: %w", err)
	}
	defer rows.Close()

	var out []*world.WorldAction

	for rows.Next() {
		var a world.WorldAction
		var rawPayload []byte

		err := rows.Scan(
			&a.ID,
			&a.WorldID,
			&a.CivilizationID,
			&a.UserID,
			&a.ActionType,
			&a.State,
			&a.TargetTick,
			&rawPayload,
			&a.RejectionReason,
			&a.CreatedAt,
			&a.AppliedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ListByWorldID scan: %w", err)
		}

		if len(rawPayload) > 0 {
			if err := json.Unmarshal(rawPayload, &a.Payload); err != nil {
				return nil, fmt.Errorf("unmarshal payload: %w", err)
			}
		} else {
			a.Payload = map[string]any{}
		}

		out = append(out, &a)
	}

	return out, rows.Err()
}

func (r *ActionRepository) ListByWorldIDAndRegionID(ctx context.Context, worldID, regionID uuid.UUID, limit int) ([]*world.WorldAction, error) {
	if limit <= 0 {
		limit = 20
	}

	const q = `
SELECT id, world_id, civilization_id, user_id, action_type, state,
       target_tick, payload_json, rejection_reason, created_at, applied_at
FROM world_actions
WHERE world_id = $1
  AND payload_json->>'region_id' = $2
ORDER BY created_at DESC
LIMIT $3
`
	rows, err := r.db.Query(ctx, q, worldID, regionID.String(), limit)
	if err != nil {
		return nil, fmt.Errorf("ListByWorldIDAndRegionID query: %w", err)
	}
	defer rows.Close()

	var out []*world.WorldAction

	for rows.Next() {
		var a world.WorldAction
		var rawPayload []byte

		err := rows.Scan(
			&a.ID,
			&a.WorldID,
			&a.CivilizationID,
			&a.UserID,
			&a.ActionType,
			&a.State,
			&a.TargetTick,
			&rawPayload,
			&a.RejectionReason,
			&a.CreatedAt,
			&a.AppliedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ListByWorldIDAndRegionID scan: %w", err)
		}

		if len(rawPayload) > 0 {
			if err := json.Unmarshal(rawPayload, &a.Payload); err != nil {
				return nil, fmt.Errorf("unmarshal payload: %w", err)
			}
		} else {
			a.Payload = map[string]any{}
		}

		out = append(out, &a)
	}

	return out, rows.Err()
}
