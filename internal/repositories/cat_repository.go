package repositories

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/IavilaGw/cat-api/internal/models"
	"gorm.io/gorm"
)

type CatRepository struct {
	db *gorm.DB
}

func NewCatRepository(db *gorm.DB) *CatRepository {
	return &CatRepository{db: db}
}

func (r *CatRepository) Save(imageData []byte, contentType string) (*models.CatImage, error) {
	hash := calculateHash(imageData)

	var existing models.CatImage
	if err := r.db.Where("image_hash = ?", hash).First(&existing).Error; err == nil {
		existing.UpdateLastAccessed()
		if err := r.db.Save(&existing).Error; err != nil {
			return nil, fmt.Errorf("failed to update image: %w", err)
		}
		return &existing, nil
	}

	catImage := &models.CatImage{
		ImageData:   imageData,
		ImageHash:   hash,
		ContentType: contentType,
		Size:        int64(len(imageData)),
		AccessCount: 1,
	}

	if err := r.db.Create(catImage).Error; err != nil {
		return nil, fmt.Errorf("failed to save image: %w", err)
	}

	return catImage, nil
}

func (r *CatRepository) FindByID(id uint) (*models.CatImage, error) {
	var catImage models.CatImage
	if err := r.db.First(&catImage, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("image not found")
		}
		return nil, fmt.Errorf("failed to find image: %w", err)
	}

	catImage.UpdateLastAccessed()
	if err := r.db.Save(&catImage).Error; err != nil {
		return nil, fmt.Errorf("failed to update access: %w", err)
	}

	return &catImage, nil
}

func (r *CatRepository) CountUnique() (int64, error) {
	var count int64
	if err := r.db.Model(&models.CatImage{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count images: %w", err)
	}
	return count, nil
}

func (r *CatRepository) GetStats() (*models.CatImageStats, error) {
	stats := &models.CatImageStats{}

	if err := r.db.Model(&models.CatImage{}).Count(&stats.TotalImages).Error; err != nil {
		return nil, fmt.Errorf("failed to count: %w", err)
	}

	var totalSize sql.NullInt64
	if err := r.db.Model(&models.CatImage{}).Select("SUM(size)").Scan(&totalSize).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate size: %w", err)
	}
	if totalSize.Valid {
		stats.TotalSize = totalSize.Int64
	}

	var mostAccessed models.CatImage
	if err := r.db.Order("access_count DESC").First(&mostAccessed).Error; err == nil {
		stats.MostAccessedID = mostAccessed.ID
		stats.MostAccessCount = mostAccessed.AccessCount
	}

	return stats, nil
}

func calculateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}