package repository

import "github.com/burberrymyshirt/shurl/model"

type UrlRepository interface {
	CreateUrl(url string, ttl uint64) error
	GetUrlByHash(url string, ttl uint64) (model.Url, error)
}
