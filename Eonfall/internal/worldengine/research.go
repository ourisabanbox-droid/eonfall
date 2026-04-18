package worldengine

import (
	"context"
	"fmt"

	"project-eonfall/internal/world"
	"project-eonfall/internal/worldrepo"
)

type BasicResearchService struct {
	researchRepo *worldrepo.ResearchRepository
	alertRepo    *worldrepo.AlertRepository
}

func NewBasicResearchService(
	researchRepo *worldrepo.ResearchRepository,
	alertRepo *worldrepo.AlertRepository,
) *BasicResearchService {
	return &BasicResearchService{
		researchRepo: researchRepo,
		alertRepo:    alertRepo,
	}
}

func (s *BasicResearchService) Apply(ctx context.Context, w *world.World) error {
	for _, civ := range w.Civilizations {
		for _, rp := range civ.Researches {
			if rp.State != "in_progress" {
				continue
			}

			if civ.ScienceStock <= 0 {
				continue
			}

			civ.ScienceStock--
			rp.Progress += 1

			completedNow := false
			if rp.Progress >= 100 {
				rp.Progress = 100
				rp.State = "completed"
				tick := w.CurrentTick
				rp.CompletedTick = &tick
				completedNow = true
			}

			if err := s.researchRepo.UpdateProgress(ctx, rp); err != nil {
				return err
			}

			if completedNow {
				title := "Recherche terminée"
				message := fmt.Sprintf("La recherche %s est terminée pour %s.", rp.ResearchType, civ.Name)

				alert := worldrepo.NewWorldAlert(
					w.ID,
					&civ.ID,
					nil,
					"research_completed",
					"info",
					title,
					message,
					map[string]any{
						"research_type": rp.ResearchType,
						"civilization":  civ.Name,
						"tick":          w.CurrentTick,
					},
				)

				if err := s.alertRepo.Insert(ctx, alert); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
