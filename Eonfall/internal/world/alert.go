package world

import (
	"time"

	"github.com/google/uuid"
)

type Alert struct {
	ID             uuid.UUID
	WorldID        uuid.UUID
	CivilizationID *uuid.UUID
	RegionID       *uuid.UUID
	AlertType      string
	Severity       string
	Title          string
	Message        string
	Payload        map[string]any
	CreatedAt      time.Time
	ReadAt         *time.Time
}
