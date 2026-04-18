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

type AlertRepository struct {
	db *pgxpool.Pool
}

func NewAlertRepository(db *pgxpool.Pool) *AlertRepository {
	return &AlertRepository{db: db}
}

func (r *AlertRepository) Insert(ctx context.Context, alert *world.Alert) error {
	payloadBytes, err := json.Marshal(alert.Payload)
	if err != nil {
		return fmt.Errorf("marshal alert payload: %w", err)
	}

	const q = `
INSERT INTO world_alerts (
    id, world_id, civilization_id, region_id, alert_type, severity,
    title, message, payload_json, created_at, read_at
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9::jsonb, $10, $11
)
`

	_, err = r.db.Exec(
		ctx,
		q,
		alert.ID,
		alert.WorldID,
		alert.CivilizationID,
		alert.RegionID,
		alert.AlertType,
		alert.Severity,
		alert.Title,
		alert.Message,
		string(payloadBytes),
		alert.CreatedAt,
		alert.ReadAt,
	)
	if err != nil {
		return fmt.Errorf("insert alert exec: %w", err)
	}

	return nil
}

func (r *AlertRepository) ListByWorldID(ctx context.Context, worldID uuid.UUID, limit int) ([]*world.Alert, error) {
	if limit <= 0 {
		limit = 50
	}

	const q = `
SELECT id, world_id, civilization_id, region_id, alert_type, severity,
       title, message, payload_json, created_at, read_at
FROM world_alerts
WHERE world_id = $1
ORDER BY created_at DESC
LIMIT $2
`

	rows, err := r.db.Query(ctx, q, worldID, limit)
	if err != nil {
		return nil, fmt.Errorf("ListByWorldID query: %w", err)
	}
	defer rows.Close()

	var out []*world.Alert
	for rows.Next() {
		var a world.Alert
		var rawPayload []byte

		err := rows.Scan(
			&a.ID,
			&a.WorldID,
			&a.CivilizationID,
			&a.RegionID,
			&a.AlertType,
			&a.Severity,
			&a.Title,
			&a.Message,
			&rawPayload,
			&a.CreatedAt,
			&a.ReadAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ListByWorldID scan: %w", err)
		}

		if len(rawPayload) > 0 {
			if err := json.Unmarshal(rawPayload, &a.Payload); err != nil {
				return nil, fmt.Errorf("unmarshal alert payload: %w", err)
			}
		} else {
			a.Payload = map[string]any{}
		}

		out = append(out, &a)
	}

	return out, rows.Err()
}

func NewWorldAlert(
	worldID uuid.UUID,
	civID *uuid.UUID,
	regionID *uuid.UUID,
	alertType, severity, title, message string,
	payload map[string]any,
) *world.Alert {
	now := time.Now().UTC()
	return &world.Alert{
		ID:             uuid.New(),
		WorldID:        worldID,
		CivilizationID: civID,
		RegionID:       regionID,
		AlertType:      alertType,
		Severity:       severity,
		Title:          title,
		Message:        message,
		Payload:        payload,
		CreatedAt:      now,
	}
}

func (r *AlertRepository) ExistsRecentSimilar(
	ctx context.Context,
	worldID uuid.UUID,
	civID *uuid.UUID,
	regionID *uuid.UUID,
	alertType string,
	resourceType *string,
	within time.Duration,
) (bool, error) {
	const q = `
SELECT EXISTS (
    SELECT 1
    FROM world_alerts
    WHERE world_id = $1
      AND alert_type = $2
      AND (
            ($3::uuid IS NULL AND civilization_id IS NULL)
         OR civilization_id = $3
      )
      AND (
            ($4::uuid IS NULL AND region_id IS NULL)
         OR region_id = $4
      )
      AND (
            $5 IS NULL
         OR payload_json->>'resource_type' = $5
      )
      AND created_at >= NOW() - ($6 * INTERVAL '1 second')
)
`
	var exists bool
	seconds := int(within.Seconds())

	err := r.db.QueryRow(
		ctx,
		q,
		worldID,
		alertType,
		civID,
		regionID,
		resourceType,
		seconds,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
