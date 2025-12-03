package models

import (
	"time"

	"gorm.io/gorm"
)

type CatImage struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ImageData      []byte         `gorm:"type:bytea;not null" json:"-"`
	ImageHash      string         `gorm:"type:varchar(64);uniqueIndex;not null" json:"image_hash"`
	ContentType    string         `gorm:"type:varchar(50);not null" json:"content_type"`
	Size           int64          `gorm:"not null" json:"size"`
	CreatedAt      time.Time      `gorm:"not null" json:"created_at"`
	LastAccessedAt time.Time      `gorm:"not null" json:"last_accessed_at"`
	AccessCount    int            `gorm:"default:0" json:"access_count"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CatImage) TableName() string {
	return "cat_images"
}

func (c *CatImage) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	c.CreatedAt = now
	c.LastAccessedAt = now
	return nil
}

func (c *CatImage) UpdateLastAccessed() {
	c.LastAccessedAt = time.Now()
	c.AccessCount++
}

type CatImageStats struct {
	TotalImages     int64 `json:"total_images"`
	TotalSize       int64 `json:"total_size_bytes"`
	MostAccessedID  uint  `json:"most_accessed_id,omitempty"`
	MostAccessCount int   `json:"most_access_count,omitempty"`
}