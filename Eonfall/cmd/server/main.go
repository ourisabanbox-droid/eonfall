package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"

	httpapi "project-eonfall/internal/api/http"
	"project-eonfall/internal/config"
	"project-eonfall/internal/db"
	"project-eonfall/internal/worldengine"
	"project-eonfall/internal/worldloader"
	"project-eonfall/internal/worldrepo"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config.Load: %v", err)
	}

	pool, err := db.NewPostgresPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db.NewPostgresPool: %v", err)
	}
	defer pool.Close()

	worldRepo := worldrepo.NewWorldRepository(pool)
	civRepo := worldrepo.NewCivilizationRepository(pool)
	regionRepo := worldrepo.NewRegionRepository(pool)
	actionRepo := worldrepo.NewActionRepository(pool)
	buildingRepo := worldrepo.NewBuildingRepository(pool)
	researchRepo := worldrepo.NewResearchRepository(pool)
	alertRepo := worldrepo.NewAlertRepository(pool)

	loader := worldloader.New(worldRepo, civRepo, regionRepo, researchRepo)

	worldID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	w, err := loader.LoadWorld(ctx, worldID)
	if err != nil {
		log.Fatalf("loader.LoadWorld: %v", err)
	}

	engine := worldengine.NewEngine(
		w,
		time.Duration(cfg.TickRateMs)*time.Millisecond,
		worldengine.NewQueuedActionService(actionRepo, buildingRepo, researchRepo, alertRepo),
		worldengine.NewBasicProductionService(),
		worldengine.NewBasicConsumptionService(),
		worldengine.NewBasicResearchService(researchRepo, alertRepo),
		worldengine.NewBasicRiskService(),
		worldengine.NewSimulationNormalizer(alertRepo),
		worldengine.NewBasicPersistenceService(worldRepo, civRepo, regionRepo, 10),
	)

	go func() {
		log.Printf("starting world engine world=%s", w.ID)
		if err := engine.Run(ctx); err != nil && ctx.Err() == nil {
			log.Fatalf("engine.Run: %v", err)
		}
	}()

	handler := httpapi.NewHandler(worldRepo, civRepo, regionRepo, actionRepo, researchRepo, alertRepo)
	router := httpapi.NewRouter(handler)

	server := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("HTTP server listening on :%s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http.ListenAndServe: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server.Shutdown: %v", err)
	}

	log.Println("server stopped")
}
