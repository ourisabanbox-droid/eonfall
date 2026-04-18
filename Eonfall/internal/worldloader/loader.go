package worldloader

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"project-eonfall/internal/world"
	"project-eonfall/internal/worldrepo"
)

type Loader struct {
	worldRepo        *worldrepo.WorldRepository
	civilizationRepo *worldrepo.CivilizationRepository
	regionRepo       *worldrepo.RegionRepository
	researchRepo     *worldrepo.ResearchRepository
}

func New(
	worldRepo *worldrepo.WorldRepository,
	civilizationRepo *worldrepo.CivilizationRepository,
	regionRepo *worldrepo.RegionRepository,
	researchRepo *worldrepo.ResearchRepository,
) *Loader {
	return &Loader{
		worldRepo:        worldRepo,
		civilizationRepo: civilizationRepo,
		regionRepo:       regionRepo,
		researchRepo:     researchRepo,
	}
}

func (l *Loader) LoadWorld(ctx context.Context, worldID uuid.UUID) (*world.World, error) {
	w, err := l.worldRepo.GetByID(ctx, worldID)
	if err != nil {
		return nil, fmt.Errorf("load world base: %w", err)
	}

	civs, err := l.civilizationRepo.ListByWorldID(ctx, worldID)
	if err != nil {
		return nil, fmt.Errorf("load civilizations: %w", err)
	}
	for _, c := range civs {
		if c.Researches == nil {
			c.Researches = map[world.ResearchType]*world.ResearchProgress{}
		}
		w.Civilizations[c.ID] = c
	}

	regions, err := l.regionRepo.ListByWorldID(ctx, worldID)
	if err != nil {
		return nil, fmt.Errorf("load regions: %w", err)
	}
	for _, rg := range regions {
		w.Regions[rg.ID] = rg
	}

	resourceStocks, err := l.regionRepo.ListResourceStocksByWorldID(ctx, worldID)
	if err != nil {
		return nil, fmt.Errorf("load region resource stocks: %w", err)
	}

	for _, rs := range resourceStocks {
		region := w.Regions[rs.RegionID]
		if region == nil {
			continue
		}
		if region.ResourceStocks == nil {
			region.ResourceStocks = map[world.ResourceType]*world.RegionResourceStock{}
		}
		region.ResourceStocks[rs.ResourceType] = rs
	}

	// Load researches for each civilization
	for _, civ := range w.Civilizations {
		researches, err := l.researchRepo.ListByCivilizationID(ctx, civ.ID)
		if err != nil {
			return nil, fmt.Errorf("load researches for civ %s: %w", civ.ID, err)
		}
		for _, rp := range researches {
			civ.Researches[rp.ResearchType] = rp
		}
	}

	return w, nil
}
