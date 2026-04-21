package http

import (
	"context"

	"github.com/google/uuid"

	"project-eonfall/internal/world"
)

func (h *Handler) buildRegionAvailableActions(
	ctx context.Context,
	worldID uuid.UUID,
	region *world.Region,
) ([]AvailableActionResponse, *CivilizationContextResponse) {
	actions := make([]AvailableActionResponse, 0, 2)

	if region.OwnerCivilizationID == nil {
		return actions, nil
	}

	var civilizationContext *CivilizationContextResponse
	var profile *world.CivilizationProfile

	trajectory, err := h.trajectoryRepo.GetByCivilizationID(ctx, *region.OwnerCivilizationID)
	if err == nil && trajectory.WorldID == worldID {
		computedProfile := world.ComputeCivilizationProfile(trajectory)
		profile = &computedProfile

		civilizationContext = &CivilizationContextResponse{
			DominantAxis:  computedProfile.DominantAxis,
			SecondaryAxis: computedProfile.SecondaryAxis,
			ProfileLabel:  computedProfile.ProfileLabel,
		}
	}

	stabilizeReasons := []string{}
	stabilizeRecommended := false

	if region.RevoltRisk >= 60 {
		stabilizeRecommended = true
		stabilizeReasons = append(stabilizeReasons, "revolt_risk_high")
	}
	if region.Stability <= 30 {
		stabilizeRecommended = true
		stabilizeReasons = append(stabilizeReasons, "stability_low")
	}
	if profile != nil && (profile.DominantAxis == world.CivilizationAxisResilience || profile.SecondaryAxis == world.CivilizationAxisResilience) {
		stabilizeReasons = append(stabilizeReasons, "civilization_resilience_affinity")
	}

	actions = append(actions, AvailableActionResponse{
		ActionType:  "stabilize_region",
		Label:       "Stabiliser la région",
		Description: "Réduit le risque de révolte et améliore temporairement la stabilité.",
		Recommended: stabilizeRecommended,
		Reasons:     stabilizeReasons,
	})

	droughtReasons := []string{}
	droughtRecommended := false

	if region.DroughtRisk >= 60 {
		droughtRecommended = true
		droughtReasons = append(droughtReasons, "drought_risk_high")
	}
	if profile != nil && (profile.DominantAxis == world.CivilizationAxisResilience || profile.SecondaryAxis == world.CivilizationAxisResilience) {
		droughtReasons = append(droughtReasons, "civilization_resilience_affinity")
	}

	actions = append(actions, AvailableActionResponse{
		ActionType:  "drought_relief",
		Label:       "Secours hydrique",
		Description: "Réduit le risque de sécheresse et peut réinjecter du stock alimentaire.",
		Recommended: droughtRecommended,
		Reasons:     droughtReasons,
	})

	return actions, civilizationContext
}
