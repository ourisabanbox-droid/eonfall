package http

type PriorityRegionResponse struct {
	RegionID             string                       `json:"region_id"`
	Stability            int                          `json:"stability"`
	RevoltRisk           int                          `json:"revolt_risk"`
	DroughtRisk          int                          `json:"drought_risk"`
	PriorityScore        int                          `json:"priority_score"`
	RecommendedActions   []AvailableActionResponse    `json:"recommended_actions"`
	TopRecommendedAction *string                      `json:"top_recommended_action,omitempty"`
	Actionable           bool                         `json:"actionable"`
	ActionBlockers       []string                     `json:"action_blockers,omitempty"`
	CivilizationContext  *CivilizationContextResponse `json:"civilization_context,omitempty"`
}
