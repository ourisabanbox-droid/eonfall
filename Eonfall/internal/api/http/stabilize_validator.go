package http

import (
	"fmt"

	"github.com/google/uuid"
)

func ValidateStabilizeRegionActionRequest(req StabilizeRegionActionRequest) error {
	if req.CivilizationID == uuid.Nil {
		return fmt.Errorf("civilization_id is required")
	}
	if req.UserID == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}
	if req.RegionID == uuid.Nil {
		return fmt.Errorf("region_id is required")
	}
	return nil
}
