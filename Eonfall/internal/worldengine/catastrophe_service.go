package worldengine

import (
	"context"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/google/uuid"

	"project-eonfall/internal/world"
	"project-eonfall/internal/worldrepo"
)

type BasicCatastropheService struct {
	catastropheRepo *worldrepo.CatastropheRepository
	alertRepo       *worldrepo.AlertRepository
}

func NewBasicCatastropheService(
	catastropheRepo *worldrepo.CatastropheRepository,
	alertRepo *worldrepo.AlertRepository,
) *BasicCatastropheService {
	return &BasicCatastropheService{
		catastropheRepo: catastropheRepo,
		alertRepo:       alertRepo,
	}
}

func (s *BasicCatastropheService) Apply(ctx context.Context, w *world.World, pressures map[uuid.UUID]RegionPressure) error {
	if w.Catastrophes == nil {
		w.Catastrophes = map[uuid.UUID]*world.Catastrophe{}
	}

	if err := s.resolveExpired(ctx, w); err != nil {
		return err
	}

	for _, c := range w.Catastrophes {
		if c.State != world.CatastropheStateActive {
			continue
		}
		if err := s.applyOngoingEffects(w, c); err != nil {
			return err
		}
	}

	for regionID, pressure := range pressures {
		region := w.Regions[regionID]
		if region == nil {
			continue
		}

		if !s.hasActiveTypeInRegion(w, regionID, world.CatastropheDrought) {
			if s.shouldTrigger(w, regionID, world.CatastropheDrought, pressure.DroughtPressure, 72) {
				if err := s.createDrought(ctx, w, region, pressure); err != nil {
					return err
				}
			}
		}

		if !s.hasActiveTypeInRegion(w, regionID, world.CatastropheRevolt) {
			if s.shouldTrigger(w, regionID, world.CatastropheRevolt, pressure.RevoltPressure, 76) {
				if err := s.createRevolt(ctx, w, region, pressure); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *BasicCatastropheService) resolveExpired(ctx context.Context, w *world.World) error {
	for id, c := range w.Catastrophes {
		if c.State != world.CatastropheStateActive {
			continue
		}
		if c.EndsAtTick != nil && w.CurrentTick >= *c.EndsAtTick {
			c.State = world.CatastropheStateResolved
			if err := s.catastropheRepo.MarkResolved(ctx, c.ID); err != nil {
				return err
			}

			regionID := c.RegionID
			exists, err := s.alertRepo.ExistsRecentSimilar(
				ctx,
				w.ID,
				nil,
				&regionID,
				fmt.Sprintf("catastrophe_%s_resolved", c.CatastropheType),
				nil,
				60*time.Second,
			)
			if err == nil && !exists {
				_ = s.alertRepo.Insert(ctx, worldrepo.NewWorldAlert(
					w.ID,
					nil,
					&regionID,
					fmt.Sprintf("catastrophe_%s_resolved", c.CatastropheType),
					"info",
					"Catastrophe terminée",
					"La région commence à se stabiliser après la catastrophe.",
					map[string]any{
						"region_id":        regionID.String(),
						"catastrophe_type": string(c.CatastropheType),
						"tick":             w.CurrentTick,
					},
				))
			}

			delete(w.Catastrophes, id)
		}
	}
	return nil
}

func (s *BasicCatastropheService) createDrought(ctx context.Context, w *world.World, region *world.Region, p RegionPressure) error {
	severity := catastropheSeverity(p.DroughtPressure)
	duration := int64(6 + severity*3)
	endTick := w.CurrentTick + duration

	c := &world.Catastrophe{
		ID:              uuid.New(),
		WorldID:         w.ID,
		RegionID:        region.ID,
		CatastropheType: world.CatastropheDrought,
		State:           world.CatastropheStateActive,
		Severity:        severity,
		StartsAtTick:    w.CurrentTick,
		EndsAtTick:      &endTick,
		Payload: map[string]any{
			"food_loss_per_tick":      3 * severity,
			"stability_loss_per_tick": severity,
			"risk_gain_per_tick":      severity,
			"pressure":                p.DroughtPressure,
		},
		Result:    map[string]any{},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.catastropheRepo.Insert(ctx, c); err != nil {
		return err
	}
	w.Catastrophes[c.ID] = c

	regionID := region.ID
	exists, err := s.alertRepo.ExistsRecentSimilar(
		ctx,
		w.ID,
		region.OwnerCivilizationID,
		&regionID,
		"catastrophe_drought_started",
		nil,
		60*time.Second,
	)
	if err == nil && !exists {
		_ = s.alertRepo.Insert(ctx, worldrepo.NewWorldAlert(
			w.ID,
			region.OwnerCivilizationID,
			&regionID,
			"catastrophe_drought_started",
			"warning",
			"Sécheresse",
			"Une sécheresse frappe la région et affaiblit son économie vivrière.",
			map[string]any{
				"region_id": region.ID.String(),
				"severity":  severity,
				"tick":      w.CurrentTick,
			},
		))
	}

	return nil
}

func (s *BasicCatastropheService) createRevolt(ctx context.Context, w *world.World, region *world.Region, p RegionPressure) error {
	severity := catastropheSeverity(p.RevoltPressure)
	duration := int64(5 + severity*2)
	endTick := w.CurrentTick + duration

	c := &world.Catastrophe{
		ID:              uuid.New(),
		WorldID:         w.ID,
		RegionID:        region.ID,
		CatastropheType: world.CatastropheRevolt,
		State:           world.CatastropheStateActive,
		Severity:        severity,
		StartsAtTick:    w.CurrentTick,
		EndsAtTick:      &endTick,
		Payload: map[string]any{
			"stability_loss_per_tick": 2 * severity,
			"cohesion_loss_per_tick":  severity,
			"materials_loss_per_tick": 2 * severity,
			"energy_loss_per_tick":    2 * severity,
			"risk_gain_per_tick":      severity,
			"pressure":                p.RevoltPressure,
		},
		Result:    map[string]any{},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := s.catastropheRepo.Insert(ctx, c); err != nil {
		return err
	}
	w.Catastrophes[c.ID] = c

	regionID := region.ID
	exists, err := s.alertRepo.ExistsRecentSimilar(
		ctx,
		w.ID,
		region.OwnerCivilizationID,
		&regionID,
		"catastrophe_revolt_started",
		nil,
		60*time.Second,
	)
	if err == nil && !exists {
		_ = s.alertRepo.Insert(ctx, worldrepo.NewWorldAlert(
			w.ID,
			region.OwnerCivilizationID,
			&regionID,
			"catastrophe_revolt_started",
			"critical",
			"Révolte",
			"Une révolte éclate dans la région et fragilise son ordre interne.",
			map[string]any{
				"region_id": region.ID.String(),
				"severity":  severity,
				"tick":      w.CurrentTick,
			},
		))
	}

	return nil
}

func (s *BasicCatastropheService) applyOngoingEffects(w *world.World, c *world.Catastrophe) error {
	region := w.Regions[c.RegionID]
	if region == nil {
		return nil
	}

	switch c.CatastropheType {
	case world.CatastropheDrought:
		foodLoss := payloadInt(c.Payload, "food_loss_per_tick")
		stabilityLoss := payloadInt(c.Payload, "stability_loss_per_tick")
		riskGain := payloadInt(c.Payload, "risk_gain_per_tick")

		if food := region.ResourceStocks[world.ResourceFood]; food != nil {
			food.Stock -= int64(foodLoss)
		}
		region.Stability -= stabilityLoss
		region.DroughtRisk = clampInt(0, 100, region.DroughtRisk+riskGain)

	case world.CatastropheRevolt:
		stabilityLoss := payloadInt(c.Payload, "stability_loss_per_tick")
		cohesionLoss := payloadInt(c.Payload, "cohesion_loss_per_tick")
		materialsLoss := payloadInt(c.Payload, "materials_loss_per_tick")
		energyLoss := payloadInt(c.Payload, "energy_loss_per_tick")
		riskGain := payloadInt(c.Payload, "risk_gain_per_tick")

		region.Stability -= stabilityLoss
		region.RevoltRisk = clampInt(0, 100, region.RevoltRisk+riskGain)

		if materials := region.ResourceStocks[world.ResourceMaterials]; materials != nil {
			materials.Stock -= int64(materialsLoss)
		}
		if energy := region.ResourceStocks[world.ResourceEnergy]; energy != nil {
			energy.Stock -= int64(energyLoss)
		}

		if region.OwnerCivilizationID != nil {
			if civ := w.Civilizations[*region.OwnerCivilizationID]; civ != nil {
				civ.Cohesion -= cohesionLoss
			}
		}
	}

	return nil
}

func (s *BasicCatastropheService) hasActiveTypeInRegion(w *world.World, regionID uuid.UUID, t world.CatastropheType) bool {
	for _, c := range w.Catastrophes {
		if c.RegionID == regionID && c.CatastropheType == t && c.State == world.CatastropheStateActive {
			return true
		}
	}
	return false
}

func (s *BasicCatastropheService) shouldTrigger(
	w *world.World,
	regionID uuid.UUID,
	t world.CatastropheType,
	pressure int,
	threshold int,
) bool {
	if pressure < threshold {
		return false
	}

	chance := float64(pressure-threshold) / float64(100-threshold)
	if chance > 0.90 {
		chance = 0.90
	}
	return deterministicRoll(w.ID, regionID, w.CurrentTick, t) < chance
}

func catastropheSeverity(pressure int) int {
	switch {
	case pressure >= 92:
		return 3
	case pressure >= 82:
		return 2
	default:
		return 1
	}
}

func payloadInt(m map[string]any, key string) int {
	v, ok := m[key]
	if !ok {
		return 0
	}

	switch x := v.(type) {
	case int:
		return x
	case int64:
		return int(x)
	case float64:
		return int(x)
	default:
		return 0
	}
}

func deterministicRoll(worldID, regionID uuid.UUID, tick int64, t world.CatastropheType) float64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(worldID.String()))
	_, _ = h.Write([]byte(regionID.String()))
	_, _ = h.Write([]byte(t))

	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(tick))
	_, _ = h.Write(buf[:])

	return float64(h.Sum64()%10000) / 10000.0
}
