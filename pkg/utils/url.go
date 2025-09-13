package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	// ShortCodeLength is the default length for generated short codes
	ShortCodeLength = 6
	// MaxShortCodeLength is the maximum length for custom short codes
	MaxShortCodeLength = 8
	// MinShortCodeLength is the minimum length for short codes
	MinShortCodeLength = 3
)

// GenerateShortCode generates a random short code of specified length
func GenerateShortCode(length int) (string, error) {
	if length < MinShortCodeLength || length > MaxShortCodeLength {
		length = ShortCodeLength
	}

	// Generate random bytes
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Convert to base64 and clean up
	shortCode := base64.URLEncoding.EncodeToString(bytes)
	shortCode = strings.TrimRight(shortCode, "=")
	shortCode = shortCode[:length]

	// Ensure it's alphanumeric
	shortCode = regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(shortCode, "")

	// If too short after cleaning, pad with random characters
	if len(shortCode) < length {
		additional := make([]byte, length-len(shortCode))
		if _, err := rand.Read(additional); err != nil {
			return "", fmt.Errorf("failed to generate additional random bytes: %w", err)
		}
		additionalStr := base64.URLEncoding.EncodeToString(additional)
		additionalStr = regexp.MustCompile(`[^a-zA-Z0-9]`).ReplaceAllString(additionalStr, "")
		shortCode += additionalStr[:length-len(shortCode)]
	}

	return shortCode, nil
}

// ValidateURL validates if a string is a valid URL
func ValidateURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// Add protocol if missing
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must have a scheme (http/https)")
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a host")
	}

	// Check for malicious patterns
	if isMaliciousURL(parsedURL) {
		return fmt.Errorf("URL appears to be malicious")
	}

	return nil
}

// isMaliciousURL checks for common malicious URL patterns
func isMaliciousURL(parsedURL *url.URL) bool {
	host := strings.ToLower(parsedURL.Host)

	// Check for suspicious domains
	suspiciousPatterns := []string{
		"localhost",
		"127.0.0.1",
		"0.0.0.0",
		"::1",
		"file://",
		"javascript:",
		"data:",
		"vbscript:",
	}

	for _, pattern := range suspiciousPatterns {
		if strings.Contains(host, pattern) || strings.HasPrefix(parsedURL.Scheme, pattern) {
			return true
		}
	}

	// Check for IP addresses (might be suspicious)
	ipRegex := regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+$`)
	if ipRegex.MatchString(host) {
		return true
	}

	return false
}

// ValidateShortCode validates if a short code meets the requirements
func ValidateShortCode(shortCode string) error {
	if shortCode == "" {
		return fmt.Errorf("short code cannot be empty")
	}

	if len(shortCode) < MinShortCodeLength || len(shortCode) > MaxShortCodeLength {
		return fmt.Errorf("short code must be between %d and %d characters", MinShortCodeLength, MaxShortCodeLength)
	}

	// Check if it contains only alphanumeric characters
	matched, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, shortCode)
	if err != nil {
		return fmt.Errorf("failed to validate short code format: %w", err)
	}

	if !matched {
		return fmt.Errorf("short code must contain only alphanumeric characters")
	}

	// Check for reserved short codes
	reservedCodes := []string{
		"admin", "api", "www", "mail", "ftp", "blog", "shop", "help",
		"about", "contact", "terms", "privacy", "login", "register",
		"dashboard", "profile", "settings", "logout", "search",
	}

	lowerCode := strings.ToLower(shortCode)
	for _, reserved := range reservedCodes {
		if lowerCode == reserved {
			return fmt.Errorf("short code '%s' is reserved", shortCode)
		}
	}

	return nil
}

// IsURLExpired checks if a URL has expired
func IsURLExpired(expiresAt *time.Time) bool {
	if expiresAt == nil {
		return false
	}
	return time.Now().After(*expiresAt)
}

// NormalizeURL normalizes a URL by adding protocol if missing
func NormalizeURL(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return "https://" + rawURL
	}
	return rawURL
}

// GetDomainFromURL extracts the domain from a URL
func GetDomainFromURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}
	return parsedURL.Host, nil
}
