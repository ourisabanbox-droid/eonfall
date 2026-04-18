package worldrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"project-eonfall/internal/world"
)

type WorldRepository struct {
	db *pgxpool.Pool
}

func NewWorldRepository(db *pgxpool.Pool) *WorldRepository {
	return &WorldRepository{db: db}
}

func (r *WorldRepository) GetByID(ctx context.Context, id uuid.UUID) (*world.World, error) {
	const q = `
SELECT id, name, season_number, state, phase, tick_rate_ms, current_tick,
       config_version, seed_value, started_at, ends_at, created_at, updated_at
FROM worlds
WHERE id = $1
`
	var w world.World
	err := r.db.QueryRow(ctx, q, id).Scan(
		&w.ID,
		&w.Name,
		&w.SeasonNumber,
		&w.State,
		&w.Phase,
		&w.TickRateMs,
		&w.CurrentTick,
		&w.ConfigVersion,
		&w.SeedValue,
		&w.StartedAt,
		&w.EndsAt,
		&w.CreatedAt,
		&w.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("GetByID scan: %w", err)
	}

	w.Civilizations = map[uuid.UUID]*world.Civilization{}
	w.Regions = map[uuid.UUID]*world.Region{}
	w.Events = map[uuid.UUID]*world.WorldEvent{}
	w.Catastrophes = map[uuid.UUID]*world.Catastrophe{}

	return &w, nil
}

func (r *WorldRepository) UpdateTick(ctx context.Context, worldID uuid.UUID, tick int64) error {
	const q = `
UPDATE worlds
SET current_tick = $2, updated_at = $3
WHERE id = $1
`
	_, err := r.db.Exec(ctx, q, worldID, tick, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("UpdateTick exec: %w", err)
	}
	return nil
}

func (r *WorldRepository) GetCurrentTick(ctx context.Context, worldID uuid.UUID) (int64, error) {
	const q = `SELECT current_tick FROM worlds WHERE id = $1`

	var tick int64
	err := r.db.QueryRow(ctx, q, worldID).Scan(&tick)
	if err != nil {
		return 0, fmt.Errorf("GetCurrentTick scan: %w", err)
	}
	return tick, nil
}
