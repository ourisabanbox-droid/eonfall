package world

import (
	"time"

	"github.com/google/uuid"
)

type WorldAction struct {
	ID              uuid.UUID
	WorldID         uuid.UUID
	CivilizationID  uuid.UUID
	UserID          uuid.UUID
	ActionType      ActionType
	State           ActionState
	TargetTick      int64
	Payload         map[string]any
	RejectionReason *string
	CreatedAt       time.Time
	AppliedAt       *time.Time
}
