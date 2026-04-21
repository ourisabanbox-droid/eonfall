package world

import "github.com/google/uuid"

type CivilizationTrajectory struct {
	CivilizationID  uuid.UUID
	WorldID         uuid.UUID
	ResilienceScore int
	ExpansionScore  int
	InfluenceScore  int
	ScienceScore    int
}
