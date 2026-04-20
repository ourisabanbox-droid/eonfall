package worldengine

import (
	"context"

	"project-eonfall/internal/world"
	"project-eonfall/internal/worldrepo"
)

type BasicPersistenceService struct {
	worldRepo        *worldrepo.WorldRepository
	civilizationRepo *worldrepo.CivilizationRepository
	regionRepo       *worldrepo.RegionRepository
	flushEveryTicks  int64
}

func NewBasicPersistenceService(
	worldRepo *worldrepo.WorldRepository,
	civilizationRepo *worldrepo.CivilizationRepository,
	regionRepo *worldrepo.RegionRepository,
	flushEveryTicks int64,
) *BasicPersistenceService {
	return &BasicPersistenceService{
		worldRepo:        worldRepo,
		civilizationRepo: civilizationRepo,
		regionRepo:       regionRepo,
		flushEveryTicks:  flushEveryTicks,
	}
}

func (s *BasicPersistenceService) FlushIfNeeded(ctx context.Context, w *world.World) error {
	if s.flushEveryTicks <= 0 {
		s.flushEveryTicks = 1
	}

	if w.CurrentTick%s.flushEveryTicks != 0 {
		return nil
	}

	if err := s.worldRepo.UpdateTick(ctx, w.ID, w.CurrentTick); err != nil {
		return err
	}

	for _, civ := range w.Civilizations {
		if err := s.civilizationRepo.UpdateRuntimeState(ctx, civ); err != nil {
			return err
		}
	}

	for _, region := range w.Regions {
		if err := s.regionRepo.UpdateRuntimeState(ctx, region); err != nil {
			return err
		}
		for _, rs := range region.ResourceStocks {
			if err := s.regionRepo.UpsertResourceStock(ctx, rs); err != nil {
				return err
			}
		}
	}

	return nil
}
