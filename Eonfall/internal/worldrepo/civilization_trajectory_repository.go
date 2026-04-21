package worldrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"project-eonfall/internal/world"
)

type CivilizationTrajectoryRepository struct {
	db *pgxpool.Pool
}

func NewCivilizationTrajectoryRepository(db *pgxpool.Pool) *CivilizationTrajectoryRepository {
	return &CivilizationTrajectoryRepository{db: db}
}

func (r *CivilizationTrajectoryRepository) GetByCivilizationID(ctx context.Context, civilizationID uuid.UUID) (*world.CivilizationTrajectory, error) {
	const q = `
SELECT civilization_id, world_id, resilience_score, expansion_score, influence_score, science_score
FROM civilization_trajectories
WHERE civilization_id = $1
`
	var t world.CivilizationTrajectory

	err := r.db.QueryRow(ctx, q, civilizationID).Scan(
		&t.CivilizationID,
		&t.WorldID,
		&t.ResilienceScore,
		&t.ExpansionScore,
		&t.InfluenceScore,
		&t.ScienceScore,
	)
	if err != nil {
		return nil, fmt.Errorf("GetByCivilizationID scan: %w", err)
	}

	return &t, nil
}

func (r *CivilizationTrajectoryRepository) Upsert(ctx context.Context, t *world.CivilizationTrajectory) error {
	const q = `
INSERT INTO civilization_trajectories (
    civilization_id, world_id, resilience_score, expansion_score, influence_score, science_score
) VALUES (
    $1, $2, $3, $4, $5, $6
)
ON CONFLICT (civilization_id) DO UPDATE
SET world_id = EXCLUDED.world_id,
    resilience_score = EXCLUDED.resilience_score,
    expansion_score = EXCLUDED.expansion_score,
    influence_score = EXCLUDED.influence_score,
    science_score = EXCLUDED.science_score
`
	_, err := r.db.Exec(
		ctx,
		q,
		t.CivilizationID,
		t.WorldID,
		t.ResilienceScore,
		t.ExpansionScore,
		t.InfluenceScore,
		t.ScienceScore,
	)
	if err != nil {
		return fmt.Errorf("Upsert exec: %w", err)
	}

	return nil
}

func (r *CivilizationTrajectoryRepository) IncrementScores(
	ctx context.Context,
	civilizationID uuid.UUID,
	worldID uuid.UUID,
	resilienceDelta int,
	expansionDelta int,
	influenceDelta int,
	scienceDelta int,
) error {
	const q = `
INSERT INTO civilization_trajectories (
    civilization_id, world_id, resilience_score, expansion_score, influence_score, science_score
) VALUES (
    $1, $2, $3, $4, $5, $6
)
ON CONFLICT (civilization_id) DO UPDATE
SET resilience_score = civilization_trajectories.resilience_score + EXCLUDED.resilience_score,
    expansion_score = civilization_trajectories.expansion_score + EXCLUDED.expansion_score,
    influence_score = civilization_trajectories.influence_score + EXCLUDED.influence_score,
    science_score = civilization_trajectories.science_score + EXCLUDED.science_score
`
	_, err := r.db.Exec(
		ctx,
		q,
		civilizationID,
		worldID,
		resilienceDelta,
		expansionDelta,
		influenceDelta,
		scienceDelta,
	)
	if err != nil {
		return fmt.Errorf("IncrementScores exec: %w", err)
	}

	return nil
}
