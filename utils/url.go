package utils

import (
	"net/url"
	"strings"
)

// NormalizeURL extracts the base URL without query parameters
// Example: https://www.airbnb.com/rooms/123?search_mode=... -> https://www.airbnb.com/rooms/123
func NormalizeURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}

	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		// If parsing fails, try to extract room ID manually
		return extractRoomURL(rawURL)
	}

	// Return scheme + host + path
	baseURL := parsedURL.Scheme + "://" + parsedURL.Host + parsedURL.Path

	// Remove trailing slash if present
	return strings.TrimSuffix(baseURL, "/")
}

// extractRoomURL extracts base room URL as fallback
func extractRoomURL(rawURL string) string {
	// Find "/rooms/" and extract up to the next "?" or end
	roomsIndex := strings.Index(rawURL, "/rooms/")
	if roomsIndex == -1 {
		return rawURL // Return as-is if not a room URL
	}

	// Find where query params start
	queryIndex := strings.Index(rawURL[roomsIndex:], "?")
	if queryIndex == -1 {
		return rawURL // return as-is
	}

	// Extract base URL, everything before "?"
	return rawURL[:roomsIndex+queryIndex]
}
