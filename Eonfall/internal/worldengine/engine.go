package worldengine

import (
	"context"
	"log"
	"time"

	"project-eonfall/internal/world"
)

type Engine struct {
	world *world.World

	tickInterval time.Duration

	actionService      ActionService
	productionService  ProductionService
	consumptionService ConsumptionService
	researchService    ResearchService
	riskService        RiskService
	normalizerService  NormalizerService
	persistenceService PersistenceService
}

func NewEngine(
	w *world.World,
	tickInterval time.Duration,
	actionService ActionService,
	productionService ProductionService,
	consumptionService ConsumptionService,
	researchService ResearchService,
	riskService RiskService,
	normalizerService NormalizerService,
	persistenceService PersistenceService,
) *Engine {
	return &Engine{
		world:              w,
		tickInterval:       tickInterval,
		actionService:      actionService,
		productionService:  productionService,
		consumptionService: consumptionService,
		researchService:    researchService,
		riskService:        riskService,
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
	if err := e.normalizerService.Apply(ctx, e.world); err != nil {
		return err
	}
	if err := e.persistenceService.FlushIfNeeded(ctx, e.world); err != nil {
		return err
	}

	log.Printf("world=%s tick=%d", e.world.ID.String(), e.world.CurrentTick)
	return nil
}
