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

	"github.com/gin-gonic/gin"
	"github.com/IavilaGw/cat-api/internal/config"
	"github.com/IavilaGw/cat-api/internal/database"
	"github.com/IavilaGw/cat-api/internal/handlers"
	"github.com/IavilaGw/cat-api/internal/repositories"
	"github.com/IavilaGw/cat-api/internal/services"
	"github.com/IavilaGw/cat-api/pkg/client"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)

	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Error database: %v", err)
	}
	defer db.Close()

	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	cataasClient := client.NewCataasClient(cfg.App.CataasAPIURL, cfg.App.TimeoutSeconds)
	catRepo := repositories.NewCatRepository(db.DB)
	catService := services.NewCatServiceWithConcrete(catRepo, cataasClient)

	catHandler := handlers.NewCatHandler(catService)
	healthHandler := handlers.NewHealthHandler(db, cataasClient)

	router := setupRouter(catHandler, healthHandler)

	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit


	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

}

func setupRouter(catHandler *handlers.CatHandler, healthHandler *handlers.HealthHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)

	api := router.Group("/api")
	{
		api.GET("/cat", catHandler.GetRandomCat)
		api.GET("/count", catHandler.GetCount)
		api.GET("/stats", catHandler.GetStats)
		api.GET("/image/:id", catHandler.GetImageByID)
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "cat-api",
			"version": "1.0.0",
			"endpoints": gin.H{
				"cat":   "/api/cat",
				"count": "/api/count",
				"stats": "/api/stats",
			},
		})
	})

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}