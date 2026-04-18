package world

import (
	"time"

	"github.com/google/uuid"
)

type Region struct {
	ID                  uuid.UUID
	WorldID             uuid.UUID
	Q                   int
	R                   int
	Biome               string
	TerrainType         string
	ClimateType         string
	OwnerCivilizationID *uuid.UUID

	Population        int64
	Stability         int
	DevelopmentLevel  int
	Pollution         int
	DroughtRisk       int
	RevoltRisk        int
	FireRisk          int
	SeismicRisk       int
	EnergyFragility   int
	LogisticFragility int
	HasCapital        bool
	Metadata          map[string]any
	CreatedAt         time.Time
	UpdatedAt         time.Time

	ResourceStocks    map[ResourceType]*RegionResourceStock
	Buildings         []*RegionBuilding
	AdjacentRegionIDs []uuid.UUID
}

type RegionResourceStock struct {
	ID              uuid.UUID
	WorldID         uuid.UUID
	RegionID        uuid.UUID
	ResourceType    ResourceType
	Stock           int64
	ProductionRate  int
	ConsumptionRate int
	Capacity        int
}

type RegionBuilding struct {
	ID             uuid.UUID
	WorldID        uuid.UUID
	RegionID       uuid.UUID
	CivilizationID uuid.UUID
	BuildingType   BuildingType
	Level          int
	State          string
	StartedTick    int64
	CompletedTick  *int64
	Durability     int
	Metadata       map[string]any
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
