package worldrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"project-eonfall/internal/world"
)

type CivilizationMissionRepository struct {
	db *pgxpool.Pool
}

func NewCivilizationMissionRepository(db *pgxpool.Pool) *CivilizationMissionRepository {
	return &CivilizationMissionRepository{db: db}
}

func (r *CivilizationMissionRepository) Accept(
	ctx context.Context,
	worldID uuid.UUID,
	civilizationID uuid.UUID,
	missionType string,
	targetRegionID *uuid.UUID,
) error {
	const q = `
INSERT INTO civilization_missions (
    world_id,
    civilization_id,
    mission_type,
    target_region_id,
    status,
    accepted_at,
    created_at,
    updated_at
)
VALUES ($1, $2, $3, $4, 'accepted', $5, $5, $5)
ON CONFLICT (world_id, civilization_id, mission_type, target_region_id)
DO UPDATE SET
    status = 'accepted',
    accepted_at = COALESCE(civilization_missions.accepted_at, EXCLUDED.accepted_at),
    updated_at = EXCLUDED.updated_at
`
	now := time.Now().UTC()

	_, err := r.db.Exec(ctx, q, worldID, civilizationID, missionType, targetRegionID, now)
	if err != nil {
		return fmt.Errorf("Accept exec: %w", err)
	}
	return nil
}

func (r *CivilizationMissionRepository) ListByCivilizationID(
	ctx context.Context,
	worldID uuid.UUID,
	civilizationID uuid.UUID,
) ([]*world.CivilizationMission, error) {
	const q = `
SELECT id, world_id, civilization_id, mission_type, target_region_id, status, accepted_at, created_at, updated_at
FROM civilization_missions
WHERE world_id = $1 AND civilization_id = $2
ORDER BY created_at DESC
`

	rows, err := r.db.Query(ctx, q, worldID, civilizationID)
	if err != nil {
		return nil, fmt.Errorf("ListByCivilizationID query: %w", err)
	}
	defer rows.Close()

	out := []*world.CivilizationMission{}
	for rows.Next() {
		var m world.CivilizationMission
		if err := rows.Scan(
			&m.ID,
			&m.WorldID,
			&m.CivilizationID,
			&m.MissionType,
			&m.TargetRegionID,
			&m.Status,
			&m.AcceptedAt,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("ListByCivilizationID scan: %w", err)
		}
		out = append(out, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListByCivilizationID rows: %w", err)
	}

	return out, nil
}
