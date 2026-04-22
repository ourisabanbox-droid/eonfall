package world

import (
	"time"

	"github.com/google/uuid"
)

type CivilizationMission struct {
	ID             uuid.UUID
	WorldID        uuid.UUID
	CivilizationID uuid.UUID
	MissionType    string
	TargetRegionID *uuid.UUID
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	AcceptedAt     *time.Time
}
