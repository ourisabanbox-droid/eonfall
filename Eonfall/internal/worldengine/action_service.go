package worldengine

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"project-eonfall/internal/world"
	"project-eonfall/internal/worldrepo"
)

type QueuedActionService struct {
	actionRepo   *worldrepo.ActionRepository
	buildingRepo *worldrepo.BuildingRepository
	researchRepo *worldrepo.ResearchRepository
	alertRepo    *worldrepo.AlertRepository
}

func NewQueuedActionService(
	actionRepo *worldrepo.ActionRepository,
	buildingRepo *worldrepo.BuildingRepository,
	researchRepo *worldrepo.ResearchRepository,
	alertRepo *worldrepo.AlertRepository,
) *QueuedActionService {
	return &QueuedActionService{
		actionRepo:   actionRepo,
		buildingRepo: buildingRepo,
		researchRepo: researchRepo,
		alertRepo:    alertRepo,
	}
}

func (s *QueuedActionService) ApplyPending(ctx context.Context, w *world.World) error {
	actions, err := s.actionRepo.ListPendingForTick(ctx, w.ID, w.CurrentTick)
	if err != nil {
		return err
	}

	for _, action := range actions {
		if err := s.applyOne(ctx, w, action); err != nil {
			_ = s.actionRepo.MarkRejected(ctx, action.ID, err.Error())
			continue
		}
		if err := s.actionRepo.MarkApplied(ctx, action.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *QueuedActionService) applyOne(ctx context.Context, w *world.World, action *world.WorldAction) error {
	switch action.ActionType {
	case world.ActionBuild:
		return s.applyBuild(ctx, w, action)
	case world.ActionResearchStart:
		return s.applyResearchStart(ctx, w, action)
	default:
		return fmt.Errorf("unsupported action type: %s", action.ActionType)
	}
}

func (s *QueuedActionService) applyBuild(ctx context.Context, w *world.World, action *world.WorldAction) error {
	rawRegionID, ok := action.Payload["region_id"].(string)
	if !ok {
		return fmt.Errorf("missing region_id")
	}

	rawBuildingType, ok := action.Payload["building_type"].(string)
	if !ok || rawBuildingType == "" {
		return fmt.Errorf("missing building_type")
	}

	regionID, err := uuid.Parse(rawRegionID)
	if err != nil {
		return fmt.Errorf("invalid region_id")
	}

	region := w.Regions[regionID]
	if region == nil {
		return fmt.Errorf("region not found")
	}

	if region.OwnerCivilizationID == nil || *region.OwnerCivilizationID != action.CivilizationID {
		return fmt.Errorf("region is not owned by civilization")
	}

	building := &world.RegionBuilding{
		ID:             uuid.New(),
		WorldID:        w.ID,
		RegionID:       region.ID,
		CivilizationID: action.CivilizationID,
		BuildingType:   world.BuildingType(rawBuildingType),
		Level:          1,
		State:          "active",
		StartedTick:    w.CurrentTick,
		CompletedTick:  ptrInt64(w.CurrentTick),
		Durability:     100,
		Metadata:       map[string]any{},
	}

	region.Buildings = append(region.Buildings, building)

	switch building.BuildingType {
	case "farm":
		increaseRegionProduction(region, world.ResourceFood, 5)
	case "quarry":
		increaseRegionProduction(region, world.ResourceMaterials, 4)
	case "power_plant":
		increaseRegionProduction(region, world.ResourceEnergy, 4)
	case "lab":
		increaseRegionProduction(region, world.ResourceKnowledge, 3)
	case "forum":
		increaseRegionProduction(region, world.ResourceCohesion, 2)
	case "logistics_hub":
		region.LogisticFragility -= 10
		if region.LogisticFragility < 0 {
			region.LogisticFragility = 0
		}
	default:
		return fmt.Errorf("unsupported building_type: %s", rawBuildingType)
	}

	if err := s.buildingRepo.Insert(ctx, building); err != nil {
		return fmt.Errorf("persist building: %w", err)
	}

	title := "Construction terminée"
	message := fmt.Sprintf("Le bâtiment %s a été construit dans la région %s.", rawBuildingType, region.ID.String())

	alert := worldrepo.NewWorldAlert(
		w.ID,
		&action.CivilizationID,
		&region.ID,
		"building_constructed",
		"info",
		title,
		message,
		map[string]any{
			"building_type": rawBuildingType,
			"region_id":     region.ID.String(),
			"tick":          w.CurrentTick,
		},
	)

	if err := s.alertRepo.Insert(ctx, alert); err != nil {
		return fmt.Errorf("insert building alert: %w", err)
	}
	return nil
}

func increaseRegionProduction(region *world.Region, resourceType world.ResourceType, amount int) {
	stock := region.ResourceStocks[resourceType]
	if stock == nil {
		stock = &world.RegionResourceStock{
			ID:              uuid.New(),
			WorldID:         region.WorldID,
			RegionID:        region.ID,
			ResourceType:    resourceType,
			Stock:           0,
			ProductionRate:  0,
			ConsumptionRate: 0,
			Capacity:        1000,
		}
		region.ResourceStocks[resourceType] = stock
	}

	stock.ProductionRate += amount
}

func ptrInt64(v int64) *int64 {
	return &v
}

func (s *QueuedActionService) applyResearchStart(ctx context.Context, w *world.World, action *world.WorldAction) error {
	rawResearchType, ok := action.Payload["research_type"].(string)
	if !ok || rawResearchType == "" {
		return fmt.Errorf("missing research_type")
	}

	civ := w.Civilizations[action.CivilizationID]
	if civ == nil {
		return fmt.Errorf("civilization not found")
	}

	researchType := world.ResearchType(rawResearchType)

	if existing, ok := civ.Researches[researchType]; ok {
		if existing.State == "in_progress" || existing.State == "completed" {
			return fmt.Errorf("research already exists with state %s", existing.State)
		}
	}

	startTick := w.CurrentTick

	rp := &world.ResearchProgress{
		ID:             uuid.New(),
		WorldID:        w.ID,
		CivilizationID: civ.ID,
		ResearchType:   researchType,
		State:          "in_progress",
		Progress:       0,
		StartedTick:    &startTick,
		CompletedTick:  nil,
	}

	if civ.Researches == nil {
		civ.Researches = map[world.ResearchType]*world.ResearchProgress{}
	}
	civ.Researches[researchType] = rp

	if err := s.researchRepo.Insert(ctx, rp); err != nil {
		return fmt.Errorf("persist research: %w", err)
	}

	return nil
}
