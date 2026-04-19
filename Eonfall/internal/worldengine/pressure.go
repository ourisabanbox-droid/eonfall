package worldengine

import (
	"context"

	"github.com/google/uuid"

	"project-eonfall/internal/world"
)

type BasicPressureService struct{}

func NewBasicPressureService() *BasicPressureService {
	return &BasicPressureService{}
}

func (s *BasicPressureService) Compute(ctx context.Context, w *world.World) (map[uuid.UUID]RegionPressure, error) {
	out := make(map[uuid.UUID]RegionPressure, len(w.Regions))

	for regionID, region := range w.Regions {
		resourceStress := s.resourceStress(region)
		stabilityStress := clampInt(0, 100, 70-region.Stability)

		socialStress := 0
		if region.OwnerCivilizationID != nil {
			if civ := w.Civilizations[*region.OwnerCivilizationID]; civ != nil {
				socialStress = clampInt(0, 100, 65-civ.Cohesion)
			}
		}

		droughtPressure := clampInt(
			0,
			100,
			region.DroughtRisk+
				resourceStress+
				maxInt(0, region.Pollution-30)/2+
				maxInt(0, 70-region.Stability)/4,
		)

		revoltPressure := clampInt(
			0,
			100,
			region.RevoltRisk+
				stabilityStress+
				socialStress+
				maxInt(0, resourceStress/2),
		)

		out[regionID] = RegionPressure{
			DroughtPressure: droughtPressure,
			RevoltPressure:  revoltPressure,
			ResourceStress:  resourceStress,
			SocialStress:    socialStress,
			StabilityStress: stabilityStress,
		}
	}

	return out, nil
}

func (s *BasicPressureService) resourceStress(region *world.Region) int {
	total := 0
	samples := 0

	for _, rt := range []world.ResourceType{
		world.ResourceFood,
		world.ResourceEnergy,
		world.ResourceMaterials,
	} {
		rs := region.ResourceStocks[rt]
		if rs == nil {
			// Missing strategic resource should create real structural pressure.
			total += 35
			samples++
			continue
		}

		stress := 0

		// Level-based stress
		if rs.Capacity > 0 {
			fillRatio := float64(rs.Stock) / float64(rs.Capacity)

			switch {
			case fillRatio < 0.10:
				stress += 45
			case fillRatio < 0.25:
				stress += 35
			case fillRatio < 0.40:
				stress += 25
			case fillRatio < 0.60:
				stress += 10
			}
		}

		// Flow imbalance stress
		imbalance := rs.ConsumptionRate - rs.ProductionRate
		if imbalance > 0 {
			stress += minInt(35, imbalance*4)
		} else if imbalance == 0 {
			// Flat equilibrium should still generate mild pressure,
			// because no buffer is being built.
			stress += 5
		}

		// Scarcity on key resources should matter more
		if rt == world.ResourceFood && rs.Stock < 100 {
			stress += 15
		}
		if rt == world.ResourceEnergy && rs.Stock < 100 {
			stress += 10
		}
		if rt == world.ResourceMaterials && rs.Stock < 100 {
			stress += 10
		}

		total += clampInt(0, 100, stress)
		samples++
	}

	if samples == 0 {
		return 0
	}

	return clampInt(0, 100, total/samples)
}

func clampInt(minV, maxV, v int) int {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
