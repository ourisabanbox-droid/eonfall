package world

import (
	"time"

	"github.com/google/uuid"
)

type WorldEvent struct {
	ID           uuid.UUID
	WorldID      uuid.UUID
	EventType    EventType
	State        EventState
	StartsAtTick int64
	EndsAtTick   *int64
	Severity     int
	Payload      map[string]any
	Result       map[string]any
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
