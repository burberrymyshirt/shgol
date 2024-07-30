package main

import (
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"

	xxhash "github.com/cespare/xxhash/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if slices.Contains(os.Args, "--prod") {
		godotenv.Load(".env.prod")
	} else {
		godotenv.Load(".env.local")
	}

	router := gin.Default()

	router.GET("/ping", Ping)
	router.POST("/", ShortenURL)

	router.Run("0.0.0.0:9090")
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func ShortenURL(c *gin.Context) {
	var request struct {
		UrlToShorten string `json:"url_to_shorten" binding:"required"`
		TTL          string `json:"ttl" `
	}

	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "validation failed" + err.Error()})
		return
	}

	if !strings.HasPrefix(request.UrlToShorten, "http://") &&
		!strings.HasPrefix(request.UrlToShorten, "https://") {
		request.UrlToShorten = "http://" + request.UrlToShorten
	}

	if _, err := url.Parse(request.UrlToShorten); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "please provide a valid url"})
	}

	xxh := xxhash.NewWithSeed(69)
	xxh.WriteString(request.UrlToShorten)
	hash := xxh.Sum64()
}
