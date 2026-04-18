package http

import "github.com/google/uuid"

type ResearchActionRequest struct {
	CivilizationID uuid.UUID `json:"civilization_id"`
	UserID         uuid.UUID `json:"user_id"`
	ResearchType   string    `json:"research_type"`
}
