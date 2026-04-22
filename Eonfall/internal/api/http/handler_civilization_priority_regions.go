package http

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"project-eonfall/internal/world"
)

func (h *Handler) GetCivilizationPriorityRegions(w http.ResponseWriter, r *http.Request) {
	worldID, err := uuid.Parse(chi.URLParam(r, "worldID"))
	if err != nil {
		http.Error(w, "invalid worldID", http.StatusBadRequest)
		return
	}

	civilizationID, err := uuid.Parse(chi.URLParam(r, "civilizationID"))
	if err != nil {
		http.Error(w, "invalid civilizationID", http.StatusBadRequest)
		return
	}

	civilizations, err := h.civilizationRepo.ListByWorldID(r.Context(), worldID)
	if err != nil {
		http.Error(w, "failed to load civilizations", http.StatusInternalServerError)
		return
	}

	civilizationFound := false
	for _, civilization := range civilizations {
		if civilization.ID == civilizationID {
			civilizationFound = true
			break
		}
	}

	if !civilizationFound {
		http.Error(w, "civilization not found", http.StatusNotFound)
		return
	}

	limit := 10
	if raw := r.URL.Query().Get("limit"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v <= 0 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = v
	}

	recommendedOnly := false
	if raw := r.URL.Query().Get("recommended_only"); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			http.Error(w, "invalid recommended_only", http.StatusBadRequest)
			return
		}
		recommendedOnly = v
	}

	regions, err := h.regionRepo.ListByWorldID(r.Context(), worldID)
	if err != nil {
		http.Error(w, "failed to load regions", http.StatusInternalServerError)
		return
	}

	var civilizationContext *CivilizationContextResponse

	trajectory, err := h.trajectoryRepo.GetByCivilizationID(r.Context(), civilizationID)
	if err == nil && trajectory.WorldID == worldID {
		profile := world.ComputeCivilizationProfile(trajectory)
		civilizationContext = &CivilizationContextResponse{
			DominantAxis:  profile.DominantAxis,
			SecondaryAxis: profile.SecondaryAxis,
			ProfileLabel:  profile.ProfileLabel,
		}
	}

	out := make([]CivilizationPriorityRegionResponse, 0, len(regions))

	for _, region := range regions {
		if region.OwnerCivilizationID == nil || *region.OwnerCivilizationID != civilizationID {
			continue
		}

		score := region.RevoltRisk + region.DroughtRisk + maxIntHTTP(0, 100-region.Stability)

		recommendedActions, _ := h.buildRegionAvailableActions(r.Context(), worldID, region)

		if recommendedOnly {
			filtered := make([]AvailableActionResponse, 0, len(recommendedActions))
			for _, action := range recommendedActions {
				if action.Recommended {
					filtered = append(filtered, action)
				}
			}
			recommendedActions = filtered

			if len(recommendedActions) == 0 {
				continue
			}
		}

		trajectoryAffinities := collectTrajectoryAffinities(recommendedActions)

		var topRecommendedAction *string
		hasStabilize := false
		hasDroughtRelief := false

		for _, action := range recommendedActions {
			if !action.Recommended {
				continue
			}
			switch action.ActionType {
			case "stabilize_region":
				hasStabilize = true
			case "drought_relief":
				hasDroughtRelief = true
			}
		}

		if hasStabilize && (!hasDroughtRelief || region.RevoltRisk >= region.DroughtRisk) {
			v := "stabilize_region"
			topRecommendedAction = &v
		} else if hasDroughtRelief {
			v := "drought_relief"
			topRecommendedAction = &v
		}

		out = append(out, CivilizationPriorityRegionResponse{
			RegionID:             region.ID.String(),
			Stability:            region.Stability,
			RevoltRisk:           region.RevoltRisk,
			DroughtRisk:          region.DroughtRisk,
			PriorityScore:        score,
			TrajectoryAffinities: trajectoryAffinities,
			RecommendedActions:   recommendedActions,
			TopRecommendedAction: topRecommendedAction,
			Actionable:           true,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].PriorityScore > out[j].PriorityScore
	})

	if limit < len(out) {
		out = out[:limit]
	}

	var summary *CivilizationPriorityRegionsSummary

	if civilizationContext != nil {
		summary = &CivilizationPriorityRegionsSummary{
			CriticalRegionCount: len(out),
			DominantAxis:        civilizationContext.DominantAxis,
			SecondaryAxis:       civilizationContext.SecondaryAxis,
			ProfileLabel:        civilizationContext.ProfileLabel,
		}

		if len(out) > 0 {
			topRegionID := out[0].RegionID
			summary.TopUrgentRegionID = &topRegionID
			summary.TopUrgentAction = out[0].TopRecommendedAction
		}

		summary.StrategicNote = buildCivilizationStrategicNote(summary)
		summary.SuggestedObjective = buildSuggestedObjective(summary)
		if len(out) > 0 {
			summary.Mission = buildMissionFromRegion(out[0], "high")
		}
		secondary := make([]MissionResponse, 0, 3)

		for i := 1; i < len(out); i++ {
			mission := buildMissionFromRegion(out[i], "medium")
			if mission == nil {
				continue
			}

			secondary = append(secondary, *mission)

			if len(secondary) >= 3 {
				break
			}
		}

		if len(secondary) > 0 {
			summary.SecondaryMissions = secondary
		}
	}

	writeJSON(w, http.StatusOK, CivilizationPriorityRegionsResponse{
		CivilizationID:      civilizationID.String(),
		WorldID:             worldID.String(),
		CivilizationContext: civilizationContext,
		Regions:             out,
		Summary:             summary,
	})
}

func collectTrajectoryAffinities(actions []AvailableActionResponse) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, 4)

	for _, action := range actions {
		for _, reason := range action.Reasons {
			if len(reason) >= len("civilization_") && reason[:len("civilization_")] == "civilization_" {
				if _, ok := seen[reason]; ok {
					continue
				}
				seen[reason] = struct{}{}
				out = append(out, reason)
			}
		}
	}

	return out
}

func buildCivilizationStrategicNote(summary *CivilizationPriorityRegionsSummary) string {
	if summary == nil {
		return ""
	}

	var base string
	switch {
	case summary.TopUrgentAction != nil && *summary.TopUrgentAction == "stabilize_region":
		base = "Vos régions prioritaires subissent une forte pression de révolte."
	case summary.TopUrgentAction != nil && *summary.TopUrgentAction == "drought_relief":
		base = "Vos régions prioritaires subissent une forte pression de sécheresse."
	default:
		base = "Votre civilisation doit surveiller plusieurs tensions régionales."
	}

	var profileNote string
	switch summary.DominantAxis {
	case world.CivilizationAxisResilience:
		profileNote = "Votre civilisation excelle dans les réponses de survie et de stabilisation."
	case world.CivilizationAxisScience:
		if summary.SecondaryAxis == world.CivilizationAxisResilience {
			profileNote = "Votre trajectoire résiliente favorise les réponses de stabilisation."
		} else {
			profileNote = "Votre trajectoire scientifique favorise les réponses structurées et rationnelles."
		}
	case world.CivilizationAxisInfluence:
		profileNote = "Votre trajectoire d'influence favorise les réponses de contrôle et de coordination."
	case world.CivilizationAxisExpansion:
		profileNote = "Votre trajectoire expansionniste favorise les réponses de projection et de consolidation."
	default:
		profileNote = "Votre trajectoire civilisationnelle influence déjà vos priorités stratégiques."
	}

	return base + " " + profileNote
}

func buildSuggestedObjective(summary *CivilizationPriorityRegionsSummary) *SuggestedObjectiveResponse {
	if summary == nil || summary.TopUrgentRegionID == nil || summary.TopUrgentAction == nil {
		return nil
	}

	switch *summary.TopUrgentAction {
	case "stabilize_region":
		return &SuggestedObjectiveResponse{
			ObjectiveType: "stabilize_priority_region",
			RegionID:      *summary.TopUrgentRegionID,
			Label:         "Stabiliser la région prioritaire",
			Reason:        "revolt_pressure_high",
		}
	case "drought_relief":
		return &SuggestedObjectiveResponse{
			ObjectiveType: "relieve_priority_region_drought",
			RegionID:      *summary.TopUrgentRegionID,
			Label:         "Déployer un secours hydrique prioritaire",
			Reason:        "drought_pressure_high",
		}
	default:
		return nil
	}
}

func buildMissionFromRegion(region CivilizationPriorityRegionResponse, priority string) *MissionResponse {
	if region.TopRecommendedAction == nil {
		return nil
	}

	switch *region.TopRecommendedAction {
	case "stabilize_region":
		description := "Réduisez la pression de révolte dans une région critique de votre civilisation."
		if priority == "high" {
			description = "Réduisez la pression de révolte dans la région la plus urgente de votre civilisation."
		}

		return &MissionResponse{
			MissionType:       "regional_stabilization",
			MissionFamily:     "crisis",
			Scope:             "region",
			Status:            "available",
			TargetRegionID:    region.RegionID,
			Title:             "Stabiliser la région prioritaire",
			Description:       description,
			Priority:          priority,
			RiskLevel:         "high",
			Reason:            "revolt_pressure_high",
			RecommendedAction: "stabilize_region",
			SuggestedApproaches: []string{
				"governance",
				"civil",
				"military",
			},
		}

	case "drought_relief":
		description := "Réduisez la pression de sécheresse dans une région critique de votre civilisation."
		riskLevel := "medium"
		if priority == "high" {
			description = "Réduisez la pression de sécheresse dans la région la plus urgente de votre civilisation."
			riskLevel = "high"
		}

		return &MissionResponse{
			MissionType:       "regional_drought_relief",
			MissionFamily:     "crisis",
			Scope:             "region",
			Status:            "available",
			TargetRegionID:    region.RegionID,
			Title:             "Déployer un secours hydrique prioritaire",
			Description:       description,
			Priority:          priority,
			RiskLevel:         riskLevel,
			Reason:            "drought_pressure_high",
			RecommendedAction: "drought_relief",
			SuggestedApproaches: []string{
				"civil",
				"logistics",
				"science",
			},
		}

	default:
		return nil
	}
}
