package worldengine

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"project-eonfall/internal/world"
)

type WorldMetadataReader interface {
	GetByID(ctx context.Context, id uuid.UUID) (*world.World, error)
}

type Engine struct {
	world *world.World

	worldRepo WorldMetadataReader

	tickInterval time.Duration

	actionService      ActionService
	productionService  ProductionService
	consumptionService ConsumptionService
	researchService    ResearchService
	riskService        RiskService
	pressureService    PressureService
	catastropheService CatastropheService
	normalizerService  NormalizerService
	persistenceService PersistenceService
}

func NewEngine(
	w *world.World,
	worldRepo WorldMetadataReader,
	tickInterval time.Duration,
	actionService ActionService,
	productionService ProductionService,
	consumptionService ConsumptionService,
	researchService ResearchService,
	riskService RiskService,
	pressureService PressureService,
	catastropheService CatastropheService,
	normalizerService NormalizerService,
	persistenceService PersistenceService,
) *Engine {
	return &Engine{
		world:              w,
		worldRepo:          worldRepo,
		tickInterval:       tickInterval,
		actionService:      actionService,
		productionService:  productionService,
		consumptionService: consumptionService,
		researchService:    researchService,
		riskService:        riskService,
		pressureService:    pressureService,
		catastropheService: catastropheService,
		normalizerService:  normalizerService,
		persistenceService: persistenceService,
	}
}

func (e *Engine) Run(ctx context.Context) error {
	ticker := time.NewTicker(e.tickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			meta, err := e.worldRepo.GetByID(ctx, e.world.ID)
			if err != nil {
				return fmt.Errorf("engine load world metadata: %w", err)
			}

			e.world.IsFrozen = meta.IsFrozen

			if e.world.IsFrozen {
				continue
			}

			if err := e.Tick(ctx); err != nil {
				return err
			}
		}
	}
}

func (e *Engine) Tick(ctx context.Context) error {
	e.world.CurrentTick++

	if err := e.actionService.ApplyPending(ctx, e.world); err != nil {
		return err
	}
	if err := e.productionService.Apply(ctx, e.world); err != nil {
		return err
	}
	if err := e.consumptionService.Apply(ctx, e.world); err != nil {
		return err
	}
	if err := e.researchService.Apply(ctx, e.world); err != nil {
		return err
	}
	if err := e.riskService.Apply(ctx, e.world); err != nil {
		return err
	}

	pressures, err := e.pressureService.Compute(ctx, e.world)
	if err != nil {
		return err
	}

	if err := e.catastropheService.Apply(ctx, e.world, pressures); err != nil {
		return err
	}

	if err := e.normalizerService.Apply(ctx, e.world); err != nil {
		return err
	}
	if err := e.persistenceService.FlushIfNeeded(ctx, e.world); err != nil {
		return err
	}

	log.Printf("world=%s tick=%d", e.world.ID.String(), e.world.CurrentTick)
	return nil
}
