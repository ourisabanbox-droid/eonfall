package worldengine

import (
	"context"
	"time"

	"project-eonfall/internal/world"
	"project-eonfall/internal/worldrepo"
)

type SimulationNormalizer struct {
	alertRepo *worldrepo.AlertRepository
}

func NewSimulationNormalizer(alertRepo *worldrepo.AlertRepository) *SimulationNormalizer {
	return &SimulationNormalizer{
		alertRepo: alertRepo,
	}
}

func (s *SimulationNormalizer) Apply(ctx context.Context, w *world.World) error {
	// --- Civilizations ---
	for _, civ := range w.Civilizations {
		// Clamp cohesion
		if civ.Cohesion < 0 {
			civ.Cohesion = 0
		}
		if civ.Cohesion > 100 {
			civ.Cohesion = 100
		}

		// Clamp civilization stocks
		if civ.FoodStock < 0 {
			civ.FoodStock = 0
		}
		if civ.EnergyStock < 0 {
			civ.EnergyStock = 0
		}
		if civ.ScienceStock < 0 {
			civ.ScienceStock = 0
		}
		if civ.CreditStock < 0 {
			civ.CreditStock = 0
		}
		if civ.MaterialsStock < 0 {
			civ.MaterialsStock = 0
		}

		// Low cohesion alert with cooldown
		if civ.Cohesion <= 20 {
			exists, err := s.alertRepo.ExistsRecentSimilar(
				ctx,
				w.ID,
				&civ.ID,
				nil,
				"low_cohesion",
				nil,
				30*time.Second,
			)
			if err == nil && !exists {
				alert := worldrepo.NewWorldAlert(
					w.ID,
					&civ.ID,
					nil,
					"low_cohesion",
					"warning",
					"Cohésion faible",
					"La cohésion de la civilisation est dangereusement basse.",
					map[string]any{
						"civilization_id": civ.ID.String(),
						"cohesion":        civ.Cohesion,
						"tick":            w.CurrentTick,
					},
				)
				_ = s.alertRepo.Insert(ctx, alert)
			}
		}
	}

	// --- Regions ---
	for _, region := range w.Regions {
		// Clamp stability
		if region.Stability < 0 {
			region.Stability = 0
		}
		if region.Stability > 100 {
			region.Stability = 100
		}

		for _, rs := range region.ResourceStocks {
			// Clamp minimum
			if rs.Stock < 0 {
				rs.Stock = 0
			}

			reachedCap := false

			// Clamp to capacity
			if rs.Capacity > 0 && rs.Stock > int64(rs.Capacity) {
				rs.Stock = int64(rs.Capacity)
				reachedCap = true
			}

			// Resource capacity alert with cooldown + resource_type dedupe
			if reachedCap {
				regionID := region.ID
				rt := string(rs.ResourceType)

				exists, err := s.alertRepo.ExistsRecentSimilar(
					ctx,
					w.ID,
					region.OwnerCivilizationID,
					&regionID,
					"resource_capacity_reached",
					&rt,
					30*time.Second,
				)
				if err == nil && !exists {
					alert := worldrepo.NewWorldAlert(
						w.ID,
						region.OwnerCivilizationID,
						&regionID,
						"resource_capacity_reached",
						"info",
						"Capacité atteinte",
						"Une ressource régionale a atteint sa capacité maximale.",
						map[string]any{
							"region_id":     region.ID.String(),
							"resource_type": rs.ResourceType,
							"capacity":      rs.Capacity,
							"tick":          w.CurrentTick,
						},
					)
					_ = s.alertRepo.Insert(ctx, alert)
				}
			}
		}
	}

	return nil
}
