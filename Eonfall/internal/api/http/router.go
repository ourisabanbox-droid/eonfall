package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", h.Health)
	r.Get("/worlds/{worldID}", h.GetWorldSummary)
	r.Get("/worlds/{worldID}/map", h.GetWorldMap)
	r.Get("/worlds/{worldID}/alerts", h.GetWorldAlerts)
	r.Get("/worlds/{worldID}/actions", h.GetWorldActions)
	r.Get("/worlds/{worldID}/dashboard/{civilizationID}", h.GetCivilizationDashboard)
	r.Get("/worlds/{worldID}/regions/{regionID}", h.GetRegionDetail)
	r.Get("/worlds/{worldID}/regions/{regionID}/actions", h.GetRegionActions)
	r.Get("/worlds/{worldID}/regions/{regionID}/available-actions", h.GetRegionAvailableActions)
	r.Get("/worlds/{worldID}/priority-regions", h.GetPriorityRegions)
	r.Post("/worlds/{worldID}/actions/build", h.PostBuildAction)
	r.Post("/worlds/{worldID}/actions/research", h.PostResearchAction)
	r.Get("/civilizations/{civilizationID}/researches", h.GetCivilizationResearches)
	r.Post("/worlds/{worldID}/actions/stabilize", h.PostStabilizeRegionAction)
	r.Post("/worlds/{worldID}/actions/drought-relief", h.PostDroughtReliefAction)

	return r
}
