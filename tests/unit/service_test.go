package services_test

import (
	"errors"
	"testing"

	"github.com/IavilaGw/cat-api/internal/models"
	"github.com/IavilaGw/cat-api/internal/services"
)

// Mock del repositorio
type MockCatRepository struct {
	SaveFunc        func([]byte, string) (*models.CatImage, error)
	CountUniqueFunc func() (int64, error)
	GetStatsFunc    func() (*models.CatImageStats, error)
	FindByIDFunc    func(uint) (*models.CatImage, error)
}

func (m *MockCatRepository) Save(data []byte, contentType string) (*models.CatImage, error) {
	if m.SaveFunc != nil {
		return m.SaveFunc(data, contentType)
	}
	return nil, errors.New("not implemented")
}

func (m *MockCatRepository) CountUnique() (int64, error) {
	if m.CountUniqueFunc != nil {
		return m.CountUniqueFunc()
	}
	return 0, errors.New("not implemented")
}

func (m *MockCatRepository) GetStats() (*models.CatImageStats, error) {
	if m.GetStatsFunc != nil {
		return m.GetStatsFunc()
	}
	return nil, errors.New("not implemented")
}

func (m *MockCatRepository) FindByID(id uint) (*models.CatImage, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}
	return nil, errors.New("not implemented")
}

// Mock del cliente
type MockCataasClient struct {
	GetRandomCatFunc func() (*services.CatImageResponse, error)
	HealthCheckFunc  func() error
}

func (m *MockCataasClient) GetRandomCat() (*services.CatImageResponse, error) {
	if m.GetRandomCatFunc != nil {
		return m.GetRandomCatFunc()
	}
	return nil, errors.New("not implemented")
}

func (m *MockCataasClient) HealthCheck() error {
	if m.HealthCheckFunc != nil {
		return m.HealthCheckFunc()
	}
	return errors.New("not implemented")
}

// Tests

func TestFetchAndSaveRandomCat_Success(t *testing.T) {
	mockRepo := &MockCatRepository{
		SaveFunc: func(data []byte, contentType string) (*models.CatImage, error) {
			return &models.CatImage{
				ID:          1,
				ImageData:   data,
				ImageHash:   "test-hash",
				ContentType: contentType,
				Size:        int64(len(data)),
			}, nil
		},
	}

	mockClient := &MockCataasClient{
		GetRandomCatFunc: func() (*services.CatImageResponse, error) {
			return &services.CatImageResponse{
				Data:        []byte("fake-image-data"),
				ContentType: "image/jpeg",
				Size:        15,
			}, nil
		},
	}

	service := services.NewCatService(mockRepo, mockClient)

	catImage, imageData, err := service.FetchAndSaveRandomCat()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if catImage == nil {
		t.Fatal("Expected cat image, got nil")
	}

	if catImage.ID != 1 {
		t.Errorf("Expected ID 1, got %d", catImage.ID)
	}

	if string(imageData) != "fake-image-data" {
		t.Errorf("Expected 'fake-image-data', got %s", string(imageData))
	}
}

func TestGetUniqueImageCount(t *testing.T) {
	mockRepo := &MockCatRepository{
		CountUniqueFunc: func() (int64, error) {
			return 42, nil
		},
	}

	mockClient := &MockCataasClient{}
	service := services.NewCatService(mockRepo, mockClient)

	count, err := service.GetUniqueImageCount()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if count != 42 {
		t.Errorf("Expected count 42, got %d", count)
	}
}

func TestGetUniqueImageCountError(t *testing.T) {
	mockRepo := &MockCatRepository{
		CountUniqueFunc: func() (int64, error) {
			return 0, errors.New("database error")
		},
	}

	mockClient := &MockCataasClient{}
	service := services.NewCatService(mockRepo, mockClient)

	count, err := service.GetUniqueImageCount()

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

func TestGetStats(t *testing.T) {
	expectedStats := &models.CatImageStats{
		TotalImages:     10,
		TotalSize:       1024000,
		MostAccessedID:  5,
		MostAccessCount: 100,
	}

	mockRepo := &MockCatRepository{
		GetStatsFunc: func() (*models.CatImageStats, error) {
			return expectedStats, nil
		},
	}

	mockClient := &MockCataasClient{}
	service := services.NewCatService(mockRepo, mockClient)

	stats, err := service.GetStats()

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if stats == nil {
		t.Fatal("Expected stats, got nil")
	}

	if stats.TotalImages != 10 {
		t.Errorf("Expected 10 images, got %d", stats.TotalImages)
	}

	if stats.TotalSize != 1024000 {
		t.Errorf("Expected size 1024000, got %d", stats.TotalSize)
	}

	if stats.MostAccessedID != 5 {
		t.Errorf("Expected most accessed ID 5, got %d", stats.MostAccessedID)
	}
}

