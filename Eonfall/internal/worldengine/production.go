package worldengine

import (
	"context"

	"project-eonfall/internal/world"
)

type BasicProductionService struct{}

func NewBasicProductionService() *BasicProductionService {
	return &BasicProductionService{}
}

func (s *BasicProductionService) Apply(ctx context.Context, w *world.World) error {
	for _, region := range w.Regions {
		if region.OwnerCivilizationID == nil {
			continue
		}

		civ := w.Civilizations[*region.OwnerCivilizationID]
		if civ == nil {
			continue
		}

		for _, stock := range region.ResourceStocks {
			stock.Stock += int64(stock.ProductionRate)
			switch stock.ResourceType {
			case world.ResourceFood:
				civ.FoodStock += stock.ProductionRate
			case world.ResourceMaterials:
				civ.MaterialsStock += stock.ProductionRate
			case world.ResourceEnergy:
				civ.EnergyStock += stock.ProductionRate
			case world.ResourceCredit:
				civ.CreditStock += stock.ProductionRate
			case world.ResourceKnowledge:
				civ.ScienceStock += stock.ProductionRate
			case world.ResourceCohesion:
				civ.Cohesion += stock.ProductionRate
			}
		}
	}
	return nil
}
