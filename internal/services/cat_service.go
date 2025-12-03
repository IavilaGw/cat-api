package services

import (
	"fmt"
	"log"

	"github.com/IavilaGw/cat-api/internal/models"
	"github.com/IavilaGw/cat-api/internal/repositories"
	"github.com/IavilaGw/cat-api/pkg/client"
)

type CatService struct {
	repo         CatRepositoryInterface
	cataasClient CataasClientInterface
}

func NewCatService(repo CatRepositoryInterface, cataasClient CataasClientInterface) *CatService {
	return &CatService{
		repo:         repo,
		cataasClient: cataasClient,
	}
}

func NewCatServiceWithConcrete(repo *repositories.CatRepository, cataasClient *client.CataasClient) *CatService {
	return &CatService{
		repo:         repo,
		cataasClient: &cataasClientAdapter{cataasClient},
	}
}

type cataasClientAdapter struct {
	client *client.CataasClient
}

func (a *cataasClientAdapter) GetRandomCat() (*CatImageResponse, error) {
	resp, err := a.client.GetRandomCat()
	if err != nil {
		return nil, err
	}
	return &CatImageResponse{
		Data:        resp.Data,
		ContentType: resp.ContentType,
		Size:        resp.Size,
	}, nil
}

func (a *cataasClientAdapter) HealthCheck() error {
	return a.client.HealthCheck()
}

func (s *CatService) FetchAndSaveRandomCat() (*models.CatImage, []byte, error) {
	log.Println("Fetching random cat image...")

	response, err := s.cataasClient.GetRandomCat()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch image: %w", err)
	}

	log.Printf("Fetched image: %d bytes, type: %s", response.Size, response.ContentType)

	catImage, err := s.repo.Save(response.Data, response.ContentType)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to save image: %w", err)
	}

	log.Printf("Image saved with ID: %d (hash: %s)", catImage.ID, catImage.ImageHash)

	return catImage, response.Data, nil
}

func (s *CatService) GetCatImageByID(id uint) (*models.CatImage, error) {
	return s.repo.FindByID(id)
}

func (s *CatService) GetUniqueImageCount() (int64, error) {
	count, err := s.repo.CountUnique()
	if err != nil {
		return 0, fmt.Errorf("failed to get count: %w", err)
	}
	return count, nil
}

func (s *CatService) GetStats() (*models.CatImageStats, error) {
	stats, err := s.repo.GetStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}
	return stats, nil
}
