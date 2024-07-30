package utils

import (
	"errors"
	"net"
	"net/url"
	"strings"
	"unicode"
)

// validateUrl checks if the given string is a valid complete URL, and either returns a absolute url, or an error.
// If the input url is not an absolute url, but just a fully qualified domain, http:// will be prepended.
func ValidateUrl(u string) (string, error) {
	// Ensure the URL has a scheme; prepend "http://" if missing
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		u = "http://" + u
	}

	// Parse the URL
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", errors.New("invalid URL provided")
	}

	// Check if the URL has both a scheme and host
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", errors.New("incomplete URL provided")
	}

	// Check if the host contains at least one period
	if !strings.Contains(parsedURL.Host, ".") {
		return "", errors.New("invalid URL host")
	}

	// Check if the host is a valid IP address (IPv4 or IPv6)
	if net.ParseIP(parsedURL.Host) != nil {
		return u, nil
	}

	// Handle IPv6 addresses enclosed in square brackets
	if strings.HasPrefix(parsedURL.Host, "[") && strings.HasSuffix(parsedURL.Host, "]") {
		if net.ParseIP(parsedURL.Host[1:len(parsedURL.Host)-1]) != nil {
			return u, nil
		}
	}

	// Validate domain names
	host := parsedURL.Host
	lastChar := ' '
	for i, r := range host {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '.' && r != '-' {
			return "", errors.New("invalid characters in URL host")
		}
		// Check for consecutive periods or trailing hyphens
		if r == '.' && (lastChar == '.' || i == len(host)-1) {
			return "", errors.New("invalid URL host")
		}
		if r == '-' && (lastChar == '-' || i == len(host)-1) {
			return "", errors.New("invalid URL host")
		}
		lastChar = r
	}

	return u, nil
}
