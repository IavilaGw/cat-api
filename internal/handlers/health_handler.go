package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/IavilaGw/cat-api/internal/database"
	"github.com/IavilaGw/cat-api/pkg/client"
)

type HealthHandler struct {
	db           *database.Database
	cataasClient *client.CataasClient
}

func NewHealthHandler(db *database.Database, cataasClient *client.CataasClient) *HealthHandler {
	return &HealthHandler{
		db:           db,
		cataasClient: cataasClient,
	}
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "cat-api",
		"version": "1.0.0",
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	checks := make(map[string]string)
	allHealthy := true

	//health
	if err := h.db.HealthCheck(); err != nil {
		checks["database"] = "unhealthy: " + err.Error()
		allHealthy = false
	} else {
		checks["database"] = "healthy"
	}

	if err := h.cataasClient.HealthCheck(); err != nil {
		checks["cataas_api"] = "unhealthy: " + err.Error()
		allHealthy = false
	} else {
		checks["cataas_api"] = "healthy"
	}

	status := "ready"
	statusCode := http.StatusOK
	
	if !allHealthy {
		status = "not_ready"
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status":  status,
		"service": "cat-api",
		"version": "1.0.0",
		"checks":  checks,
	})
}