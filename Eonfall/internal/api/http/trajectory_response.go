package http

import "project-eonfall/internal/world"

type CivilizationTrajectoryScoresResponse struct {
	Resilience int `json:"resilience"`
	Expansion  int `json:"expansion"`
	Influence  int `json:"influence"`
	Science    int `json:"science"`
}

type CivilizationTrajectoryResponse struct {
	CivilizationID string                               `json:"civilization_id"`
	WorldID        string                               `json:"world_id"`
	Scores         CivilizationTrajectoryScoresResponse `json:"scores"`
	DominantAxis   world.CivilizationAxis               `json:"dominant_axis"`
	SecondaryAxis  world.CivilizationAxis               `json:"secondary_axis"`
	ProfileLabel   string                               `json:"profile_label"`
}
