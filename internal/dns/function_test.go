package dns

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestIsValidHTTPURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "Valid HTTPS URL",
			url:      "https://www.example.com",
			expected: true,
		},
		{
			name:     "Valid HTTP URL",
			url:      "http://example.com",
			expected: true,
		},
		{
			name:     "Valid HTTP URL with IP",
			url:      "http://192.168.1.1",
			expected: true,
		},
		{
			name:     "Valid HTTP URL with port",
			url:      "http://localhost:8080",
			expected: true,
		},
		{
			name:     "Invalid URL - FTP scheme",
			url:      "ftp://example.com",
			expected: false,
		},
		{
			name:     "Invalid URL - Missing scheme",
			url:      "www.example.com",
			expected: false,
		},
		{
			name:     "Invalid URL - Missing host",
			url:      "https://",
			expected: false,
		},
		{
			name:     "Invalid URL - Empty string",
			url:      "",
			expected: false,
		},
		{
			name:     "Valid URL without Scheme",
			url:      "example.com",
			expected: true,
		},
		{
			name:     "Valid URL without Scheme",
			url:      "google.com",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result := URLValidator(tt.url)
			assert.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}
