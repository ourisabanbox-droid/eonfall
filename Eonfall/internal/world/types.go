package world

type WorldState string
type WorldPhase string
type ResourceType string
type BuildingType string
type ResearchType string
type ActionType string
type ActionState string
type EventType string
type EventState string
type CatastropheType string
type CatastropheState string

const (
	WorldStatePending WorldState = "pending"
	WorldStateRunning WorldState = "running"
	WorldStateEnded   WorldState = "ended"

	PhaseFoundation   WorldPhase = "foundation"
	PhaseExpansion    WorldPhase = "expansion"
	PhasePolarization WorldPhase = "polarization"
	PhaseCrisis       WorldPhase = "crisis"
	PhaseEndgame      WorldPhase = "endgame"

	ResourceFood      ResourceType = "food"
	ResourceMaterials ResourceType = "materials"
	ResourceEnergy    ResourceType = "energy"
	ResourceCredit    ResourceType = "credit"
	ResourceKnowledge ResourceType = "knowledge"
	ResourceCohesion  ResourceType = "cohesion"

	ActionBuild           ActionType = "build"
	ActionResearchStart   ActionType = "research_start"
	ActionExpand          ActionType = "expand"
	ActionStabilizeRegion ActionType = "stabilize_region"
	ActionDroughtRelief   ActionType = "drought_relief"

	ActionPending  ActionState = "pending"
	ActionApplied  ActionState = "applied"
	ActionRejected ActionState = "rejected"

	CatastropheDrought CatastropheType = "drought"
	CatastropheRevolt  CatastropheType = "revolt"

	CatastropheStateActive   CatastropheState = "active"
	CatastropheStateResolved CatastropheState = "resolved"
)
