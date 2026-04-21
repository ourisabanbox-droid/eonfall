package worldengine

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"project-eonfall/internal/world"
	"project-eonfall/internal/worldrepo"
)

const (
	stabilizeStabilityGain             = 20
	stabilizeRevoltRiskReduction       = 25
	stabilizePressureReduction         = 15
	stabilizeDurationTicks       int64 = 5
	stabilizeCohesionCost        int64 = 50

	droughtReliefRiskReduction       = 25
	droughtReliefFoodGrant     int64 = 50
	droughtReliefFoodCapacity        = 200
)

type QueuedActionService struct {
	actionRepo     *worldrepo.ActionRepository
	buildingRepo   *worldrepo.BuildingRepository
	researchRepo   *worldrepo.ResearchRepository
	alertRepo      *worldrepo.AlertRepository
	regionRepo     *worldrepo.RegionRepository
	trajectoryRepo *worldrepo.CivilizationTrajectoryRepository
}

func NewQueuedActionService(
	actionRepo *worldrepo.ActionRepository,
	buildingRepo *worldrepo.BuildingRepository,
	researchRepo *worldrepo.ResearchRepository,
	alertRepo *worldrepo.AlertRepository,
	regionRepo *worldrepo.RegionRepository,
	trajectoryRepo *worldrepo.CivilizationTrajectoryRepository,
) *QueuedActionService {
	return &QueuedActionService{
		actionRepo:     actionRepo,
		buildingRepo:   buildingRepo,
		researchRepo:   researchRepo,
		alertRepo:      alertRepo,
		regionRepo:     regionRepo,
		trajectoryRepo: trajectoryRepo,
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
	case world.ActionStabilizeRegion:
		return s.applyStabilizeRegion(ctx, w, action)
	case world.ActionDroughtRelief:
		return s.applyDroughtRelief(ctx, w, action)
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

	var resilienceDelta, expansionDelta, influenceDelta, scienceDelta int

	switch rawBuildingType {
	case "lab":
		scienceDelta = 2
	case "power_plant":
		expansionDelta = 1
		scienceDelta = 1
	}

	if resilienceDelta != 0 || expansionDelta != 0 || influenceDelta != 0 || scienceDelta != 0 {
		if err := s.trajectoryRepo.IncrementScores(
			ctx,
			action.CivilizationID,
			action.WorldID,
			resilienceDelta,
			expansionDelta,
			influenceDelta,
			scienceDelta,
		); err != nil {
			return fmt.Errorf("increment build trajectory: %w", err)
		}
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
	if err := s.trajectoryRepo.IncrementScores(
		ctx,
		action.CivilizationID,
		action.WorldID,
		0, // resilience
		0, // expansion
		0, // influence
		2, // science
	); err != nil {
		return fmt.Errorf("increment research trajectory: %w", err)
	}

	return nil
}

func (s *QueuedActionService) applyStabilizeRegion(ctx context.Context, w *world.World, action *world.WorldAction) error {
	rawRegionID, ok := action.Payload["region_id"].(string)
	if !ok || rawRegionID == "" {
		return fmt.Errorf("missing region_id")
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

	cohesion := region.ResourceStocks[world.ResourceCohesion]
	if cohesion == nil {
		return fmt.Errorf("missing cohesion stock")
	}
	if cohesion.Stock < stabilizeCohesionCost {
		return fmt.Errorf("insufficient cohesion")
	}

	cohesion.Stock -= stabilizeCohesionCost

	region.Stability = clampInt(0, 100, region.Stability+stabilizeStabilityGain)
	region.RevoltRisk = clampInt(0, 100, region.RevoltRisk-stabilizeRevoltRiskReduction)

	if region.Metadata == nil {
		region.Metadata = map[string]any{}
	}
	region.Metadata["stabilize_until_tick"] = w.CurrentTick + stabilizeDurationTicks
	region.Metadata["stabilize_revolt_pressure_reduction"] = stabilizePressureReduction

	alert := worldrepo.NewWorldAlert(
		w.ID,
		&action.CivilizationID,
		&region.ID,
		"region_stabilized",
		"info",
		"Région stabilisée",
		"Des mesures d'urgence ont temporairement renforcé la stabilité régionale.",
		map[string]any{
			"region_id":   region.ID.String(),
			"stability":   region.Stability,
			"revolt_risk": region.RevoltRisk,
			"tick":        w.CurrentTick,
		},
	)

	if err := s.alertRepo.Insert(ctx, alert); err != nil {
		return fmt.Errorf("insert stabilize alert: %w", err)
	}
	if err := s.trajectoryRepo.IncrementScores(
		ctx,
		action.CivilizationID,
		action.WorldID,
		2, // resilience
		0, // expansion
		1, // influence
		0, // science
	); err != nil {
		return fmt.Errorf("increment stabilize trajectory: %w", err)
	}
	return nil
}

func (s *QueuedActionService) applyDroughtRelief(ctx context.Context, w *world.World, action *world.WorldAction) error {
	rawRegionID, ok := action.Payload["region_id"].(string)
	if !ok || rawRegionID == "" {
		return fmt.Errorf("missing region_id")
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

	region.DroughtRisk = clampInt(0, 100, region.DroughtRisk-droughtReliefRiskReduction)

	food := region.ResourceStocks[world.ResourceFood]
	if food == nil {
		if region.ResourceStocks == nil {
			region.ResourceStocks = map[world.ResourceType]*world.RegionResourceStock{}
		}

		food = &world.RegionResourceStock{
			ID:              uuid.New(),
			WorldID:         region.WorldID,
			RegionID:        region.ID,
			ResourceType:    world.ResourceFood,
			Stock:           0,
			ProductionRate:  0,
			ConsumptionRate: 0,
			Capacity:        droughtReliefFoodCapacity,
		}
		if err := s.regionRepo.InsertResourceStock(ctx, food); err != nil {
			return fmt.Errorf("insert drought relief food stock: %w", err)
		}
		region.ResourceStocks[world.ResourceFood] = food
	}

	food.Stock += droughtReliefFoodGrant
	if food.Capacity > 0 && food.Stock > int64(food.Capacity) {
		food.Stock = int64(food.Capacity)
	}
	if err := s.regionRepo.UpdateResourceStock(ctx, food); err != nil {
		return fmt.Errorf("update drought relief food stock: %w", err)
	}
	alert := worldrepo.NewWorldAlert(
		w.ID,
		&action.CivilizationID,
		&region.ID,
		"region_drought_relieved",
		"info",
		"Secours hydrique",
		"Des mesures d'urgence ont réduit temporairement la pression liée à la sécheresse.",
		map[string]any{
			"region_id":    region.ID.String(),
			"drought_risk": region.DroughtRisk,
			"food_stock": func() int64 {
				if food != nil {
					return food.Stock
				}
				return 0
			}(),
			"tick": w.CurrentTick,
		},
	)

	if err := s.alertRepo.Insert(ctx, alert); err != nil {
		return fmt.Errorf("insert drought relief alert: %w", err)
	}
	if err := s.trajectoryRepo.IncrementScores(
		ctx,
		action.CivilizationID,
		action.WorldID,
		2, // resilience
		0, // expansion
		0, // influence
		0, // science
	); err != nil {
		return fmt.Errorf("increment drought relief trajectory: %w", err)
	}
	return nil
}
