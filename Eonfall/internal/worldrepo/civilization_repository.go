package worldrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"project-eonfall/internal/world"
)

type CivilizationRepository struct {
	db *pgxpool.Pool
}

func NewCivilizationRepository(db *pgxpool.Pool) *CivilizationRepository {
	return &CivilizationRepository{db: db}
}

func (r *CivilizationRepository) ListByWorldID(ctx context.Context, worldID uuid.UUID) ([]*world.Civilization, error) {
	const q = `
SELECT id, world_id, user_id, template_id, name, status,
       population, cohesion, influence, science_stock, credit_stock,
       food_stock, energy_stock, materials_stock,
       military_score, victory_score, capital_region_id,
       created_at, updated_at
FROM civilizations
WHERE world_id = $1
`
	rows, err := r.db.Query(ctx, q, worldID)
	if err != nil {
		return nil, fmt.Errorf("ListByWorldID query: %w", err)
	}
	defer rows.Close()

	var out []*world.Civilization
	for rows.Next() {
		var c world.Civilization
		err := rows.Scan(
			&c.ID, &c.WorldID, &c.UserID, &c.TemplateID, &c.Name, &c.Status,
			&c.Population, &c.Cohesion, &c.Influence, &c.ScienceStock, &c.CreditStock,
			&c.FoodStock, &c.EnergyStock, &c.MaterialsStock,
			&c.MilitaryScore, &c.VictoryScore, &c.CapitalRegionID,
			&c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ListByWorldID scan: %w", err)
		}
		c.Researches = map[world.ResearchType]*world.ResearchProgress{}
		out = append(out, &c)
	}
	return out, rows.Err()
}

func (r *CivilizationRepository) UpdateRuntimeState(ctx context.Context, civ *world.Civilization) error {
	const q = `
UPDATE civilizations
SET population = $2,
    cohesion = $3,
    influence = $4,
    science_stock = $5,
    credit_stock = $6,
    food_stock = $7,
    energy_stock = $8,
    materials_stock = $9,
    military_score = $10,
    victory_score = $11,
    updated_at = NOW()
WHERE id = $1
`
	_, err := r.db.Exec(
		ctx,
		q,
		civ.ID,
		civ.Population,
		civ.Cohesion,
		civ.Influence,
		civ.ScienceStock,
		civ.CreditStock,
		civ.FoodStock,
		civ.EnergyStock,
		civ.MaterialsStock,
		civ.MilitaryScore,
		civ.VictoryScore,
	)
	if err != nil {
		return fmt.Errorf("UpdateRuntimeState exec: %w", err)
	}

	return nil
}
