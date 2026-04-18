package http

import "github.com/google/uuid"

type BuildActionRequest struct {
	CivilizationID uuid.UUID `json:"civilization_id"`
	UserID         uuid.UUID `json:"user_id"`
	RegionID       uuid.UUID `json:"region_id"`
	BuildingType   string    `json:"building_type"`
}
