package services

import (
	"github.com/IavilaGw/cat-api/internal/models"
)

type CatRepositoryInterface interface {
	Save(imageData []byte, contentType string) (*models.CatImage, error)
	FindByID(id uint) (*models.CatImage, error)
	CountUnique() (int64, error)
	GetStats() (*models.CatImageStats, error)
}

type CataasClientInterface interface {
	GetRandomCat() (*CatImageResponse, error)
	HealthCheck() error
}

type CatImageResponse struct {
	Data        []byte
	ContentType string
	Size        int64
}
