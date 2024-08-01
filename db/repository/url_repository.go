package repository

import (
	"database/sql"

	"github.com/burberrymyshirt/shgol/model"
)

type UrlRepository interface {
	CreateUrl(
		hash string,
		originalUrl string,
		shortenedUrl string,
		RunsOutAt sql.NullTime,
	) (model.Url, error)
	GetUrlByHash(hash string) (model.Url, error)
	UpdateRunsOutAtByHash(hash string, runsOutAt sql.NullTime) error
}
