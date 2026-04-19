package worldengine

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"project-eonfall/internal/worldrepo"
)

func (e *Engine) StabilizeRegion(ctx context.Context, worldID, regionID uuid.UUID) error {
	if e == nil || e.world == nil {
		return fmt.Errorf("engine world is not initialized")
	}
	if e.world.ID != worldID {
		return fmt.Errorf("world mismatch")
	}

	region := e.world.Regions[regionID]
	if region == nil {
		return fmt.Errorf("region not found")
	}

	region.Stability = clampInt(0, 100, region.Stability+20)
	region.RevoltRisk = clampInt(0, 100, region.RevoltRisk-25)

	if e.normalizerService != nil {
		if err := e.normalizerService.Apply(ctx, e.world); err != nil {
			return fmt.Errorf("normalize stabilized region: %w", err)
		}
	}

	if repo, ok := e.normalizerService.(*SimulationNormalizer); ok && repo.alertRepo != nil {
		_ = repo.alertRepo.Insert(ctx, worldrepo.NewWorldAlert(
			e.world.ID,
			region.OwnerCivilizationID,
			&regionID,
			"region_stabilized",
			"info",
			"Région stabilisée",
			"Des mesures d'urgence ont temporairement renforcé la stabilité régionale.",
			map[string]any{
				"region_id":   regionID.String(),
				"stability":   region.Stability,
				"revolt_risk": region.RevoltRisk,
				"tick":        e.world.CurrentTick,
			},
		))
	}

	return nil
}
