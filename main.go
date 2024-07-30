package main

import (
	"fmt"
	"net/http"
	"os"
	"slices"

	"github.com/CoolRunner-dk/shurl/utils"
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

	hexString := fmt.Sprintf("%06x", hashedUrl) // "499602d2"
	c.JSON(http.StatusCreated, gin.H{"shortened_url": hexString})
}
