package model

import (
	"time"

	"gorm.io/gorm"
)

type Url struct {
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Hash         string         `gorm:"primaryKey"`
	UrlToShorten string
	ShortenedUrl string
	RunsOutAt    time.Time
}
