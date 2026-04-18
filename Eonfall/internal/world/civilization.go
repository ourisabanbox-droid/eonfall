package world

import (
	"time"

	"github.com/google/uuid"
)

type Civilization struct {
	ID              uuid.UUID
	WorldID         uuid.UUID
	UserID          uuid.UUID
	TemplateID      string
	Name            string
	Status          string
	Population      int64
	Cohesion        int
	Influence       int
	ScienceStock    int
	CreditStock     int
	FoodStock       int
	EnergyStock     int
	MaterialsStock  int
	MilitaryScore   int
	VictoryScore    int
	CapitalRegionID *uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Researches map[ResearchType]*ResearchProgress
}
