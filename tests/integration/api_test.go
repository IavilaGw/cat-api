
package integration_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/IavilaGw/cat-api/internal/config"
	"github.com/IavilaGw/cat-api/internal/database"
	"github.com/IavilaGw/cat-api/internal/handlers"
	"github.com/IavilaGw/cat-api/internal/repositories"
	"github.com/IavilaGw/cat-api/internal/services"
	"github.com/IavilaGw/cat-api/pkg/client"
)

var (
	router       *gin.Engine
	db           *database.Database
	testDBConfig = &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "catdb_test",
		SSLMode:  "disable",
	}
)

func setupTestRouter() (*gin.Engine, error) {
	gin.SetMode(gin.TestMode)

	// Conectar a base de datos de prueba
	testDB, err := database.NewDatabase(testDBConfig)
	if err != nil {
		return nil, err
	}

	// Ejecutar migraciones
	if err := testDB.AutoMigrate(); err != nil {
		return nil, err
	}

	db = testDB

	// Setup servicios
	cataasClient := client.NewCataasClient("https://cataas.com", 30)
	catRepo := repositories.NewCatRepository(testDB.DB)
	catService := services.NewCatServiceWithConcrete(catRepo, cataasClient)

	catHandler := handlers.NewCatHandler(catService)
	healthHandler := handlers.NewHealthHandler(testDB, cataasClient)

	// Setup router
	r := gin.New()
	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)

	api := r.Group("/api")
	{
		api.GET("/cat", catHandler.GetRandomCat)
		api.GET("/count", catHandler.GetCount)
		api.GET("/stats", catHandler.GetStats)
	}

	return r, nil
}

func teardownTest() {
	if db != nil {
		db.DB.Exec("DELETE FROM cat_images")
		db.Close()
	}
}

func TestHealthEndpoint(t *testing.T) {
	router, err := setupTestRouter()
	if err != nil {
		t.Skip("Database not available:", err)
		return
	}
	defer teardownTest()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}



func TestGetCatEndpoint(t *testing.T) {
	router, err := setupTestRouter()
	if err != nil {
		t.Skip("Database not available:", err)
		return
	}
	defer teardownTest()

	req, _ := http.NewRequest("GET", "/api/cat", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Header().Get("Content-Type") == "" {
		t.Error("Expected Content-Type header")
	}

	if len(w.Body.Bytes()) == 0 {
		t.Error("Expected image data, got empty response")
	}
}

func TestGetCountEndpoint_Empty(t *testing.T) {
	router, err := setupTestRouter()
	if err != nil {
		t.Skip("Database not available:", err)
		return
	}
	defer teardownTest()

	req, _ := http.NewRequest("GET", "/api/count", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}


func TestGetStatsEndpoint(t *testing.T) {
	router, err := setupTestRouter()
	if err != nil {
		t.Skip("Database not available:", err)
		return
	}
	defer teardownTest()

	req, _ := http.NewRequest("GET", "/api/stats", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}