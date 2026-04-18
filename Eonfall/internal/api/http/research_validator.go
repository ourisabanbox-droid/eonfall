package http

import (
	"fmt"
	"strings"
)

func ValidateResearchActionRequest(req ResearchActionRequest) error {
	if req.CivilizationID.String() == "" {
		return fmt.Errorf("civilization_id is required")
	}
	if req.UserID.String() == "" {
		return fmt.Errorf("user_id is required")
	}
	if strings.TrimSpace(req.ResearchType) == "" {
		return fmt.Errorf("research_type is required")
	}
	return nil
}
