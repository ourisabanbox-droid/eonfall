package world

import (
	"time"

	"github.com/google/uuid"
)

type Catastrophe struct {
	ID              uuid.UUID
	WorldID         uuid.UUID
	RegionID        uuid.UUID
	CatastropheType CatastropheType
	State           CatastropheState
	Severity        int
	StartsAtTick    int64
	EndsAtTick      *int64
	Payload         map[string]any
	Result          map[string]any
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ResearchProgress struct {
	ID             uuid.UUID
	WorldID        uuid.UUID
	CivilizationID uuid.UUID
	ResearchType   ResearchType
	State          string
	Progress       int
	StartedTick    *int64
	CompletedTick  *int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
