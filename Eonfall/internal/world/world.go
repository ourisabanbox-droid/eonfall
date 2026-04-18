package world

import (
	"time"

	"github.com/google/uuid"
)

type World struct {
	ID            uuid.UUID
	Name          string
	SeasonNumber  int
	State         WorldState
	Phase         WorldPhase
	TickRateMs    int
	CurrentTick   int64
	ConfigVersion string
	SeedValue     int64

	StartedAt *time.Time
	EndsAt    *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	Civilizations map[uuid.UUID]*Civilization
	Regions       map[uuid.UUID]*Region
	Events        map[uuid.UUID]*WorldEvent
	Catastrophes  map[uuid.UUID]*Catastrophe
}
