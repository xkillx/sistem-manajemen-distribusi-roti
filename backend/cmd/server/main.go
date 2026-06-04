package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smdr/backend/internal/config"
	"github.com/smdr/backend/internal/handler"
	"github.com/smdr/backend/internal/middleware"
	"github.com/smdr/backend/internal/repository"
	"github.com/smdr/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	pool, err := repository.Connect(cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	store := repository.NewStore(pool)

	authService := service.NewAuthService(store.Queries, cfg)
	seedService := service.NewSeedService(store.Queries, cfg)
	if err := seedService.SeedOwner(context.Background()); err != nil {
		log.Printf("Warning: Owner seed: %v", err)
	}

	authMiddleware := middleware.NewAuthMiddleware(authService, store.Queries)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.SecurityHeaders())

	api := r.Group("/api")
	handler.RegisterAuthRoutes(api, authService, authMiddleware)
	handler.RegisterProductRoutes(api, store, authMiddleware)
	handler.RegisterWarungRoutes(api, store, authMiddleware)
	handler.RegisterSalesRoutes(api, store, authMiddleware, authService)
	handler.RegisterVisitRoutes(api, store, authMiddleware)
	handler.RegisterDashboardRoutes(api, store, authMiddleware)
	handler.RegisterReportRoutes(api, store, authMiddleware)
	handler.RegisterSummaryRoutes(api, store, authMiddleware)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port()),
		Handler: r,
	}

	go func() {
		log.Printf("SMDR backend starting on :%s", cfg.Port())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
