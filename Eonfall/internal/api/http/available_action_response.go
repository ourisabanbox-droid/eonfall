package http

type AvailableActionResponse struct {
	ActionType  string   `json:"action_type"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Recommended bool     `json:"recommended"`
	Reasons     []string `json:"reasons,omitempty"`
}

type AvailableActionSignals struct {
	Stability   int `json:"stability"`
	RevoltRisk  int `json:"revolt_risk"`
	DroughtRisk int `json:"drought_risk"`
}
