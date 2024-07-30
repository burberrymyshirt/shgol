package main

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"unicode"

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
	// TODO: Make work with other types of urls, that are not using http/https
	var request struct {
		UrlToShorten string `json:"url_to_shorten" binding:"required"`
		TTL          string `json:"ttl" `
	}

	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "validation failed" + err.Error()})
		return
	}

	validUrl, err := validateUrl(request.UrlToShorten)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	xxh := xxhash.NewWithSeed(69)
	xxh.WriteString(validUrl)
	hashedUrl := xxh.Sum64()

	c.JSON(http.StatusCreated, gin.H{"shortened_url": hashedUrl})
}

func validateUrl(u string) (string, error) {
	// this fixes urls being passed without a scheme
	if !strings.HasPrefix(u, "http://") &&
		!strings.HasPrefix(u, "https://") {
		// prepending http, as it will work in more cases. It's the requestees job to redirect from http to https
		u = "http://" + u
	}

	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", errors.New("invalid URL provided")
	}

	// Check if the URL has a scheme and host
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", errors.New("incomplete URL provided")
	}

	// Check if the host contains at least one period
	if !strings.Contains(parsedURL.Host, ".") {
		return "", errors.New("invalid URL host")
	}

	// Optionally, check if the host contains only valid characters
	if !isValidHost(parsedURL.Host) {
		return "", errors.New("invalid URL host")
	}

	return u, nil
}

// isValidHost checks if the host contains only valid characters (alphanumeric and hyphens)
func isValidHost(host string) bool {
	for _, r := range host {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '.' && r != '-' {
			return false
		}
	}
	return true
}
