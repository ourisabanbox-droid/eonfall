package http

import "project-eonfall/internal/world"

type CivilizationPriorityRegionResponse struct {
	RegionID             string                    `json:"region_id"`
	Stability            int                       `json:"stability"`
	RevoltRisk           int                       `json:"revolt_risk"`
	DroughtRisk          int                       `json:"drought_risk"`
	PriorityScore        int                       `json:"priority_score"`
	TrajectoryAffinities []string                  `json:"trajectory_affinities,omitempty"`
	RecommendedActions   []AvailableActionResponse `json:"recommended_actions"`
	TopRecommendedAction *string                   `json:"top_recommended_action,omitempty"`
	Actionable           bool                      `json:"actionable"`
	ActionBlockers       []string                  `json:"action_blockers,omitempty"`
}

type CivilizationPriorityRegionsResponse struct {
	CivilizationID      string                               `json:"civilization_id"`
	WorldID             string                               `json:"world_id"`
	CivilizationContext *CivilizationContextResponse         `json:"civilization_context,omitempty"`
	Regions             []CivilizationPriorityRegionResponse `json:"regions"`
	Summary             *CivilizationPriorityRegionsSummary  `json:"summary,omitempty"`
}

type CivilizationPriorityRegionsSummary struct {
	CriticalRegionCount int                         `json:"critical_region_count"`
	TopUrgentRegionID   *string                     `json:"top_urgent_region_id,omitempty"`
	TopUrgentAction     *string                     `json:"top_urgent_action,omitempty"`
	DominantAxis        world.CivilizationAxis      `json:"dominant_axis"`
	SecondaryAxis       world.CivilizationAxis      `json:"secondary_axis"`
	ProfileLabel        string                      `json:"profile_label"`
	StrategicNote       string                      `json:"strategic_note"`
	SuggestedObjective  *SuggestedObjectiveResponse `json:"suggested_objective,omitempty"`
	Mission             *MissionResponse            `json:"mission,omitempty"`
}

type SuggestedObjectiveResponse struct {
	ObjectiveType string `json:"objective_type"`
	RegionID      string `json:"region_id"`
	Label         string `json:"label"`
	Reason        string `json:"reason"`
}

type MissionResponse struct {
	MissionType       string `json:"mission_type"`
	Status            string `json:"status"`
	TargetRegionID    string `json:"target_region_id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Priority          string `json:"priority"`
	Reason            string `json:"reason"`
	RecommendedAction string `json:"recommended_action"`
}
