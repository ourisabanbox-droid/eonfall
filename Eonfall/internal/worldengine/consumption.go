package worldengine

import (
	"context"

	"project-eonfall/internal/world"
)

type BasicConsumptionService struct{}

func NewBasicConsumptionService() *BasicConsumptionService {
	return &BasicConsumptionService{}
}

func (s *BasicConsumptionService) Apply(ctx context.Context, w *world.World) error {
	for _, civ := range w.Civilizations {
		if civ.FoodStock > 0 {
			civ.FoodStock--
		} else {
			civ.Cohesion -= 1
		}

		if civ.EnergyStock > 0 {
			civ.EnergyStock--
		}

		if civ.FoodStock > 20 && civ.EnergyStock > 20 && civ.Cohesion < 100 {
			civ.Cohesion += 1
		}

		if civ.FoodStock < 0 {
			civ.FoodStock = 0
		}
		if civ.EnergyStock < 0 {
			civ.EnergyStock = 0
		}
		if civ.Cohesion < 0 {
			civ.Cohesion = 0
		}
		if civ.Cohesion > 100 {
			civ.Cohesion = 100
		}
	}

	return nil
}
