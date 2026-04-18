package worldrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"project-eonfall/internal/world"
)

type ResearchRepository struct {
	db *pgxpool.Pool
}

func NewResearchRepository(db *pgxpool.Pool) *ResearchRepository {
	return &ResearchRepository{db: db}
}

func (r *ResearchRepository) GetByCivilizationAndType(
	ctx context.Context,
	civilizationID uuid.UUID,
	researchType world.ResearchType,
) (*world.ResearchProgress, error) {
	const q = `
SELECT id, world_id, civilization_id, research_type, state, progress,
       started_tick, completed_tick, created_at, updated_at
FROM civilization_researches
WHERE civilization_id = $1 AND research_type = $2
`
	var rp world.ResearchProgress
	err := r.db.QueryRow(ctx, q, civilizationID, researchType).Scan(
		&rp.ID,
		&rp.WorldID,
		&rp.CivilizationID,
		&rp.ResearchType,
		&rp.State,
		&rp.Progress,
		&rp.StartedTick,
		&rp.CompletedTick,
		&rp.CreatedAt,
		&rp.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("GetByCivilizationAndType scan: %w", err)
	}

	return &rp, nil
}

func (r *ResearchRepository) Insert(ctx context.Context, rp *world.ResearchProgress) error {
	const q = `
INSERT INTO civilization_researches (
    id, world_id, civilization_id, research_type, state, progress,
    started_tick, completed_tick, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10
)
`
	now := time.Now().UTC()

	_, err := r.db.Exec(
		ctx,
		q,
		rp.ID,
		rp.WorldID,
		rp.CivilizationID,
		rp.ResearchType,
		rp.State,
		rp.Progress,
		rp.StartedTick,
		rp.CompletedTick,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("Insert research exec: %w", err)
	}

	return nil
}

func (r *ResearchRepository) UpdateProgress(ctx context.Context, rp *world.ResearchProgress) error {
	const q = `
UPDATE civilization_researches
SET state = $2,
    progress = $3,
    started_tick = $4,
    completed_tick = $5,
    updated_at = $6
WHERE id = $1
`
	_, err := r.db.Exec(
		ctx,
		q,
		rp.ID,
		rp.State,
		rp.Progress,
		rp.StartedTick,
		rp.CompletedTick,
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("UpdateProgress exec: %w", err)
	}

	return nil
}

func (r *ResearchRepository) ListByCivilizationID(ctx context.Context, civilizationID uuid.UUID) ([]*world.ResearchProgress, error) {
	const q = `
SELECT id, world_id, civilization_id, research_type, state, progress,
       started_tick, completed_tick, created_at, updated_at
FROM civilization_researches
WHERE civilization_id = $1
ORDER BY created_at ASC
`
	rows, err := r.db.Query(ctx, q, civilizationID)
	if err != nil {
		return nil, fmt.Errorf("ListByCivilizationID query: %w", err)
	}
	defer rows.Close()

	var out []*world.ResearchProgress
	for rows.Next() {
		var rp world.ResearchProgress
		if err := rows.Scan(
			&rp.ID,
			&rp.WorldID,
			&rp.CivilizationID,
			&rp.ResearchType,
			&rp.State,
			&rp.Progress,
			&rp.StartedTick,
			&rp.CompletedTick,
			&rp.CreatedAt,
			&rp.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("ListByCivilizationID scan: %w", err)
		}
		out = append(out, &rp)
	}

	return out, rows.Err()
}
