package http

import (
	"fmt"
	"strings"
)

func ValidateBuildActionRequest(req BuildActionRequest) error {
	if req.CivilizationID.String() == "" {
		return fmt.Errorf("civilization_id is required")
	}
	if req.UserID.String() == "" {
		return fmt.Errorf("user_id is required")
	}
	if req.RegionID.String() == "" {
		return fmt.Errorf("region_id is required")
	}
	if strings.TrimSpace(req.BuildingType) == "" {
		return fmt.Errorf("building_type is required")
	}
	return nil
}
