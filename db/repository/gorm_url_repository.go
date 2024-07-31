package repository

import (
	"database/sql"
	"errors"

	"github.com/burberrymyshirt/shurl/db"
	"github.com/burberrymyshirt/shurl/model"
	"gorm.io/gorm"
)

type gormUrlRepository struct {
	db *gorm.DB
}

func NewGormUrlRepository() UrlRepository {
	return &gormUrlRepository{
		db: db.DBConnection(),
	}
}

func (repo *gormUrlRepository) CreateUrl(
	hash string,
	originalUrl string,
	shortenedUrl string,
	runsOutAt sql.NullTime,
) (model.Url, error) {
	urlModel := model.Url{
		OriginalUrl:  originalUrl,
		ShortenedUrl: shortenedUrl,
		RunsOutAt:    runsOutAt,
	}
	return urlModel, repo.db.Create(&urlModel).Error
}

func (repo *gormUrlRepository) GetUrlByHash(hash string) (model.Url, error) {
	var urlModel model.Url
	err := repo.db.Where("hash = ?", hash).First(&urlModel).Error
	return urlModel, err
}

func (repo *gormUrlRepository) UpdateRunsOutAtByHash(hash string, runsOutAt sql.NullTime) error {
	return errors.New("not implemented")
}
