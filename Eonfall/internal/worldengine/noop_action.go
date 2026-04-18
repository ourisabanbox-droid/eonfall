package worldengine

import (
	"context"

	"project-eonfall/internal/world"
)

type NoopActionService struct{}

func NewNoopActionService() *NoopActionService {
	return &NoopActionService{}
}

func (s *NoopActionService) ApplyPending(ctx context.Context, w *world.World) error {
	return nil
}
