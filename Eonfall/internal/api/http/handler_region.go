package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) GetRegionDetail(w http.ResponseWriter, r *http.Request) {
	worldID, err := uuid.Parse(chi.URLParam(r, "worldID"))
	if err != nil {
		http.Error(w, "invalid worldID", http.StatusBadRequest)
		return
	}

	regionID, err := uuid.Parse(chi.URLParam(r, "regionID"))
	if err != nil {
		http.Error(w, "invalid regionID", http.StatusBadRequest)
		return
	}

	region, err := h.regionRepo.GetByID(r.Context(), worldID, regionID)
	if err != nil {
		http.Error(w, "region not found", http.StatusNotFound)
		return
	}

	buildings, err := h.regionRepo.ListBuildingsByRegionID(r.Context(), regionID)
	if err != nil {
		http.Error(w, "failed to load region buildings", http.StatusInternalServerError)
		return
	}

	resourceStocks, err := h.regionRepo.ListResourceStocksByRegionID(r.Context(), regionID)
	if err != nil {
		http.Error(w, "failed to load region resources", http.StatusInternalServerError)
		return
	}

	catastrophes, err := h.catastropheRepo.ListActiveByRegionID(r.Context(), regionID)
	if err != nil {
		http.Error(w, "failed to load region catastrophes", http.StatusInternalServerError)
		return
	}

	type buildingView struct {
		ID            uuid.UUID `json:"id"`
		BuildingType  string    `json:"building_type"`
		Level         int       `json:"level"`
		State         string    `json:"state"`
		StartedTick   int64     `json:"started_tick"`
		CompletedTick *int64    `json:"completed_tick"`
		Durability    int       `json:"durability"`
	}

	type resourceView struct {
		ResourceType    string `json:"resource_type"`
		Stock           int64  `json:"stock"`
		ProductionRate  int    `json:"production_rate"`
		ConsumptionRate int    `json:"consumption_rate"`
		Capacity        int    `json:"capacity"`
	}

	type catastropheView struct {
		ID              uuid.UUID `json:"id"`
		CatastropheType string    `json:"catastrophe_type"`
		State           string    `json:"state"`
		Severity        int       `json:"severity"`
		StartsAtTick    int64     `json:"starts_at_tick"`
		EndsAtTick      *int64    `json:"ends_at_tick"`
		Payload         any       `json:"payload"`
	}

	bv := make([]buildingView, 0, len(buildings))
	for _, b := range buildings {
		bv = append(bv, buildingView{
			ID:            b.ID,
			BuildingType:  string(b.BuildingType),
			Level:         b.Level,
			State:         b.State,
			StartedTick:   b.StartedTick,
			CompletedTick: b.CompletedTick,
			Durability:    b.Durability,
		})
	}

	rv := make([]resourceView, 0, len(resourceStocks))
	for _, rs := range resourceStocks {
		rv = append(rv, resourceView{
			ResourceType:    string(rs.ResourceType),
			Stock:           rs.Stock,
			ProductionRate:  rs.ProductionRate,
			ConsumptionRate: rs.ConsumptionRate,
			Capacity:        rs.Capacity,
		})
	}

	cv := make([]catastropheView, 0, len(catastrophes))
	for _, c := range catastrophes {
		cv = append(cv, catastropheView{
			ID:              c.ID,
			CatastropheType: string(c.CatastropheType),
			State:           string(c.State),
			Severity:        c.Severity,
			StartsAtTick:    c.StartsAtTick,
			EndsAtTick:      c.EndsAtTick,
			Payload:         c.Payload,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":                    region.ID,
		"world_id":              region.WorldID,
		"q":                     region.Q,
		"r":                     region.R,
		"biome":                 region.Biome,
		"terrain_type":          region.TerrainType,
		"climate_type":          region.ClimateType,
		"owner_civilization_id": region.OwnerCivilizationID,
		"population":            region.Population,
		"stability":             region.Stability,
		"development_level":     region.DevelopmentLevel,
		"pollution":             region.Pollution,
		"drought_risk":          region.DroughtRisk,
		"revolt_risk":           region.RevoltRisk,
		"fire_risk":             region.FireRisk,
		"seismic_risk":          region.SeismicRisk,
		"energy_fragility":      region.EnergyFragility,
		"logistic_fragility":    region.LogisticFragility,
		"has_capital":           region.HasCapital,
		"buildings":             bv,
		"resource_stocks":       rv,
		"active_catastrophes":   cv,
	})
}
