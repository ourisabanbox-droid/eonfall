package worldengine

import (
	"context"

	"github.com/google/uuid"

	"project-eonfall/internal/world"
)

type RegionPressure struct {
	DroughtPressure int `json:"drought_pressure"`
	RevoltPressure  int `json:"revolt_pressure"`
	ResourceStress  int `json:"resource_stress"`
	SocialStress    int `json:"social_stress"`
	StabilityStress int `json:"stability_stress"`
}

type ActionService interface {
	ApplyPending(context.Context, *world.World) error
}

type ProductionService interface {
	Apply(context.Context, *world.World) error
}

type ConsumptionService interface {
	Apply(context.Context, *world.World) error
}

type ResearchService interface {
	Apply(context.Context, *world.World) error
}

type RiskService interface {
	Apply(context.Context, *world.World) error
}

type PressureService interface {
	Compute(context.Context, *world.World) (map[uuid.UUID]RegionPressure, error)
}

type CatastropheService interface {
	Apply(context.Context, *world.World, map[uuid.UUID]RegionPressure) error
}

type PersistenceService interface {
	FlushIfNeeded(context.Context, *world.World) error
}

type NormalizerService interface {
	Apply(context.Context, *world.World) error
}
