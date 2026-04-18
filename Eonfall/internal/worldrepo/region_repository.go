package worldrepo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"project-eonfall/internal/world"
)

type RegionRepository struct {
	db *pgxpool.Pool
}

func NewRegionRepository(db *pgxpool.Pool) *RegionRepository {
	return &RegionRepository{db: db}
}

func (r *RegionRepository) ListByWorldID(ctx context.Context, worldID uuid.UUID) ([]*world.Region, error) {
	const q = `
SELECT id, world_id, q, r, biome, terrain_type, climate_type, owner_civilization_id,
       population, stability, development_level, pollution,
       drought_risk, revolt_risk, fire_risk, seismic_risk,
       energy_fragility, logistic_fragility, has_capital,
       metadata_json, created_at, updated_at
FROM regions
WHERE world_id = $1
`
	rows, err := r.db.Query(ctx, q, worldID)
	if err != nil {
		return nil, fmt.Errorf("ListByWorldID query: %w", err)
	}
	defer rows.Close()

	var out []*world.Region
	for rows.Next() {
		var rg world.Region
		err := rows.Scan(
			&rg.ID, &rg.WorldID, &rg.Q, &rg.R, &rg.Biome, &rg.TerrainType, &rg.ClimateType, &rg.OwnerCivilizationID,
			&rg.Population, &rg.Stability, &rg.DevelopmentLevel, &rg.Pollution,
			&rg.DroughtRisk, &rg.RevoltRisk, &rg.FireRisk, &rg.SeismicRisk,
			&rg.EnergyFragility, &rg.LogisticFragility, &rg.HasCapital,
			&rg.Metadata, &rg.CreatedAt, &rg.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ListByWorldID scan: %w", err)
		}
		rg.ResourceStocks = map[world.ResourceType]*world.RegionResourceStock{}
		out = append(out, &rg)
	}
	return out, rows.Err()
}

func (r *RegionRepository) GetByID(ctx context.Context, worldID, regionID uuid.UUID) (*world.Region, error) {
	const q = `
SELECT id, world_id, q, r, biome, terrain_type, climate_type, owner_civilization_id,
       population, stability, development_level, pollution,
       drought_risk, revolt_risk, fire_risk, seismic_risk,
       energy_fragility, logistic_fragility, has_capital,
       metadata_json, created_at, updated_at
FROM regions
WHERE world_id = $1 AND id = $2
`
	var rg world.Region
	err := r.db.QueryRow(ctx, q, worldID, regionID).Scan(
		&rg.ID, &rg.WorldID, &rg.Q, &rg.R, &rg.Biome, &rg.TerrainType, &rg.ClimateType, &rg.OwnerCivilizationID,
		&rg.Population, &rg.Stability, &rg.DevelopmentLevel, &rg.Pollution,
		&rg.DroughtRisk, &rg.RevoltRisk, &rg.FireRisk, &rg.SeismicRisk,
		&rg.EnergyFragility, &rg.LogisticFragility, &rg.HasCapital,
		&rg.Metadata, &rg.CreatedAt, &rg.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("GetByID scan: %w", err)
	}

	rg.ResourceStocks = map[world.ResourceType]*world.RegionResourceStock{}
	return &rg, nil
}

func (r *RegionRepository) ListBuildingsByRegionID(ctx context.Context, regionID uuid.UUID) ([]*world.RegionBuilding, error) {
	const q = `
SELECT id, world_id, region_id, civilization_id, building_type,
       level, state, started_tick, completed_tick, durability,
       metadata_json, created_at, updated_at
FROM region_buildings
WHERE region_id = $1
ORDER BY created_at ASC
`
	rows, err := r.db.Query(ctx, q, regionID)
	if err != nil {
		return nil, fmt.Errorf("ListBuildingsByRegionID query: %w", err)
	}
	defer rows.Close()

	var out []*world.RegionBuilding
	for rows.Next() {
		var b world.RegionBuilding
		var rawMeta []byte

		err := rows.Scan(
			&b.ID, &b.WorldID, &b.RegionID, &b.CivilizationID, &b.BuildingType,
			&b.Level, &b.State, &b.StartedTick, &b.CompletedTick, &b.Durability,
			&rawMeta, &b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ListBuildingsByRegionID scan: %w", err)
		}

		if len(rawMeta) > 0 {
			if err := json.Unmarshal(rawMeta, &b.Metadata); err != nil {
				return nil, fmt.Errorf("unmarshal building metadata: %w", err)
			}
		} else {
			b.Metadata = map[string]any{}
		}

		out = append(out, &b)
	}

	return out, rows.Err()
}

func (r *RegionRepository) ListResourceStocksByRegionID(ctx context.Context, regionID uuid.UUID) ([]*world.RegionResourceStock, error) {
	const q = `
SELECT id, world_id, region_id, resource_type, stock, production_rate, consumption_rate, capacity
FROM region_resource_stocks
WHERE region_id = $1
ORDER BY resource_type ASC
`
	rows, err := r.db.Query(ctx, q, regionID)
	if err != nil {
		return nil, fmt.Errorf("ListResourceStocksByRegionID query: %w", err)
	}
	defer rows.Close()

	var out []*world.RegionResourceStock
	for rows.Next() {
		var rs world.RegionResourceStock
		err := rows.Scan(
			&rs.ID,
			&rs.WorldID,
			&rs.RegionID,
			&rs.ResourceType,
			&rs.Stock,
			&rs.ProductionRate,
			&rs.ConsumptionRate,
			&rs.Capacity,
		)
		if err != nil {
			return nil, fmt.Errorf("ListResourceStocksByRegionID scan: %w", err)
		}
		out = append(out, &rs)
	}

	return out, rows.Err()
}

func (r *RegionRepository) ListResourceStocksByWorldID(ctx context.Context, worldID uuid.UUID) ([]*world.RegionResourceStock, error) {
	const q = `
SELECT id, world_id, region_id, resource_type, stock, production_rate, consumption_rate, capacity
FROM region_resource_stocks
WHERE world_id = $1
ORDER BY region_id, resource_type
`
	rows, err := r.db.Query(ctx, q, worldID)
	if err != nil {
		return nil, fmt.Errorf("ListResourceStocksByWorldID query: %w", err)
	}
	defer rows.Close()

	var out []*world.RegionResourceStock
	for rows.Next() {
		var rs world.RegionResourceStock
		if err := rows.Scan(
			&rs.ID,
			&rs.WorldID,
			&rs.RegionID,
			&rs.ResourceType,
			&rs.Stock,
			&rs.ProductionRate,
			&rs.ConsumptionRate,
			&rs.Capacity,
		); err != nil {
			return nil, fmt.Errorf("ListResourceStocksByWorldID scan: %w", err)
		}
		out = append(out, &rs)
	}

	return out, rows.Err()
}

func (r *RegionRepository) UpdateResourceStock(ctx context.Context, rs *world.RegionResourceStock) error {
	const q = `
UPDATE region_resource_stocks
SET stock = $2,
    production_rate = $3,
    consumption_rate = $4,
    capacity = $5,
    updated_at = NOW()
WHERE id = $1
`
	_, err := r.db.Exec(
		ctx,
		q,
		rs.ID,
		rs.Stock,
		rs.ProductionRate,
		rs.ConsumptionRate,
		rs.Capacity,
	)
	if err != nil {
		return fmt.Errorf("UpdateResourceStock exec: %w", err)
	}

	return nil
}
