package worldengine

import (
	"context"

	"project-eonfall/internal/world"
)

type BasicRiskService struct{}

func NewBasicRiskService() *BasicRiskService {
	return &BasicRiskService{}
}

func (s *BasicRiskService) Apply(ctx context.Context, w *world.World) error {
	for _, region := range w.Regions {
		if region.Pollution > 50 {
			region.DroughtRisk += 1
		}
		if region.Stability < 60 {
			region.RevoltRisk += 1
		}
		if region.DroughtRisk > 100 {
			region.DroughtRisk = 100
		}
		if region.RevoltRisk > 100 {
			region.RevoltRisk = 100
		}
	}
	return nil
}
