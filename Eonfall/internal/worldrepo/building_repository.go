package worldrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"project-eonfall/internal/world"
)

type BuildingRepository struct {
	db *pgxpool.Pool
}

func NewBuildingRepository(db *pgxpool.Pool) *BuildingRepository {
	return &BuildingRepository{db: db}
}

func (r *BuildingRepository) Insert(ctx context.Context, b *world.RegionBuilding) error {
	metaBytes, err := json.Marshal(b.Metadata)
	if err != nil {
		return fmt.Errorf("marshal building metadata: %w", err)
	}

	const q = `
INSERT INTO region_buildings (
    id, world_id, region_id, civilization_id, building_type,
    level, state, started_tick, completed_tick, durability,
    metadata_json, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10,
    $11::jsonb, $12, $13
)
`

	now := time.Now().UTC()

	_, err = r.db.Exec(
		ctx,
		q,
		b.ID,
		b.WorldID,
		b.RegionID,
		b.CivilizationID,
		b.BuildingType,
		b.Level,
		b.State,
		b.StartedTick,
		b.CompletedTick,
		b.Durability,
		string(metaBytes),
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("insert building exec: %w", err)
	}

	return nil
}
