package main

import (
	"fmt"
	"net/http"
	"os"
	"slices"

	"github.com/burberrymyshirt/shurl/db"
	"github.com/burberrymyshirt/shurl/utils"
	xxhash "github.com/cespare/xxhash/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
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
	// NOTE: only works with http/https
	var request struct {
		UrlToShorten string `json:"url_to_shorten" binding:"required"`
		TTL          string `json:"ttl" `
	}

	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "validation failed" + err.Error()})
		return
	}

	validUrl, err := utils.ValidateUrl(request.UrlToShorten)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	xxh := xxhash.NewWithSeed(69)
	xxh.WriteString(validUrl)
	hashedUrl := xxh.Sum64()

	hexString := fmt.Sprintf("%06x", hashedUrl)[:6]
	shortUrl := fmt.Sprintf("%s/%s", os.Getenv("SHORT_URL"), hexString)

	// TODO: Add database stuff, along with checking if it already exists

	c.JSON(http.StatusCreated, gin.H{"shortened_url": shortUrl})
}

// TODO: implement function
func RedirectUrl(c *gin.Context) {
	panic("not implemented")
}
