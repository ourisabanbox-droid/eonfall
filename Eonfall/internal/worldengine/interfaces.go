package worldengine

import (
	"context"

	"project-eonfall/internal/world"
)

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

type PersistenceService interface {
	FlushIfNeeded(context.Context, *world.World) error
}

type NormalizerService interface {
	Apply(context.Context, *world.World) error
}
