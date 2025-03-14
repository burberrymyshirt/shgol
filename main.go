package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/burberrymyshirt/shgol/db"
	"github.com/burberrymyshirt/shgol/db/repository"
	"github.com/burberrymyshirt/shgol/utils"
	xxhash "github.com/cespare/xxhash/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// TODO: fix this temporary cli argument handling
	if slices.Contains(os.Args, "--prod") {
		godotenv.Load(".env.prod")
		gin.SetMode(gin.ReleaseMode)
	} else {
		godotenv.Load(".env.local")
	}

	// TODO: fix ugly database initialization, that only really works well with gorm orm
	dbf := db.DatabaseConnectionFactory()
	dbf.DatabaseInit()
	dbf.DatabaseMigrator()

	router := gin.Default()

	// test route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	api := router.Group("/api")
	api.POST("/shorten", ShortenURL)

	router.GET("/:path", RedirectUrl)

	router.Run("0.0.0.0:9090")
}

func ShortenURL(c *gin.Context) {
	var request struct {
		UrlToShorten string         `json:"url_to_shorten" binding:"required"`
		RunsOutAt    utils.NullTime `json:"runs_out_at" `
	}

	// NOTE: RunsOutAt is overwritten, as i have not decided on how to implement the ttl yet.
	request.RunsOutAt = utils.NullTime{
		NullTime: sql.NullTime{Time: time.Now().Add(time.Hour * 24 * 14), Valid: true},
	}

	if err := utils.BindJSON(c, &request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "validation failed: " + err.Error()})
		return
	}

	validInputUrl, err := utils.ValidateUrl(request.UrlToShorten)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	repo := repository.NewUrlRepository()
	const maxHashAttempts = 5
	var hexHashString, fullShortenedUrl string
	uniqueHashFound := false

	for attempt := 0; attempt < maxHashAttempts; attempt++ {
		xxh := xxhash.NewWithSeed(uint64(attempt))
		xxh.WriteString(validInputUrl)
		// Hash and convert to 6 character string
		hexHashString = fmt.Sprintf("%06x", xxh.Sum64())[:6]

		// Build url
		shortUrl := os.Getenv("SHORT_URL")
		fullShortenedUrl = fmt.Sprintf("%s/%s", shortUrl, hexHashString)

		existingUrl, err := repo.GetUrlByHash(hexHashString)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				uniqueHashFound = true
				break
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
			return
		}

		if existingUrl.OriginalUrl == validInputUrl {
			if err := repo.UpdateRunsOutAtByHash(hexHashString, request.RunsOutAt.NullTime); err != nil {
				c.JSON(
					http.StatusInternalServerError,
					gin.H{"error": "could not update the URL TTL: " + err.Error()},
				)
				return
			}
			c.JSON(http.StatusOK, gin.H{"shortened_url": fullShortenedUrl})
			return
		}
	}

	if !uniqueHashFound {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "could not generate a unique hash for the URL"},
		)
		return
	}

	// No existing URL found or URL is unique, create a new record
	if _, err := repo.CreateUrl(hexHashString, validInputUrl, fullShortenedUrl, request.RunsOutAt.NullTime); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "could not save the URL: " + err.Error()},
		)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"shortened_url": fullShortenedUrl})
}

func RedirectUrl(c *gin.Context) {
	hash, ok := c.Params.Get("path")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	repo := repository.NewUrlRepository()
	url, err := repo.GetUrlByHash(hash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusMovedPermanently, url.OriginalUrl)
}
