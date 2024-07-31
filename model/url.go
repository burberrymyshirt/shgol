package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Url struct {
	Hash         string `gorm:"primaryKey"`
	OriginalUrl  string
	ShortenedUrl string
	RunsOutAt    sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
