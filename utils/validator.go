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

	// Remove port if present for further validation
	host := parsedURL.Host
	if strings.Contains(host, ":") {
		host = host[:strings.Index(host, ":")]
	}

	// Check for trailing dot
	if strings.HasSuffix(host, ".") {
		return "", errors.New("invalid URL host: trailing dot not allowed")
	}

	// Check if the host is a valid IP address (IPv4 or IPv6)
	if net.ParseIP(host) != nil {
		return u, nil
	}

	// Handle IPv6 addresses enclosed in square brackets
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		if net.ParseIP(host[1:len(host)-1]) != nil {
			return u, nil
		}
	}

	// Validate domain names
	if !isValidDomainName(host) {
		return "", errors.New("invalid URL host")
	}

	return u, nil
}

// isValidDomainName checks if the domain name is properly formatted
func isValidDomainName(host string) bool {
	// Ensure the domain name contains at least one period
	if !strings.Contains(host, ".") {
		return false
	}

	// Split domain into labels and validate each label
	labels := strings.Split(host, ".")
	if len(labels) < 2 {
		return false
	}

	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return false
		}
		for i, r := range label {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' {
				return false
			}
			if r == '-' && (i == 0 || i == len(label)-1) {
				return false
			}
		}
	}

	// Check for consecutive periods or trailing hyphens in the full domain name
	lastChar := ' '
	for i, r := range host {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '.' && r != '-' {
			return false
		}
		if r == '.' && (lastChar == '.' || i == len(host)-1) {
			return false
		}
		if r == '-' && (lastChar == '-' || i == len(host)-1) {
			return false
		}
		lastChar = r
	}

	return true
}
