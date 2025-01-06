package check

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainValidator(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		expected bool
	}{
		{
			name:     "Valid domain",
			domain:   "example.com",
			expected: true,
		},
		{
			name:     "Valid subdomain",
			domain:   "sub.example.com",
			expected: true,
		},
		{
			name:     "Invalid domain - hyphen at start",
			domain:   "-invalid.com",
			expected: false,
		},
		{
			name:     "Invalid domain - hyphen at end",
			domain:   "invalid-.com",
			expected: false,
		},
		{
			name:     "Invalid domain - missing top-level domain",
			domain:   "example",
			expected: false,
		},
		{
			name:     "Invalid domain - double dots",
			domain:   "invalid..com",
			expected: false,
		},
		{
			name:     "Valid domain with hyphens",
			domain:   "valid-domain.org",
			expected: true,
		},
		{
			name:     "Invalid domain - too long",
			domain:   "toolongdomainnamethatiswaylongerthanthemaximumallowedlengthof253charactersandshouldfailvalidationbecauseitistoolongandexceedsthelimit.toolongdomainnamethatiswaylongerthanthemaximumallowedlengthof253charactersandshouldfailvalidationbecauseitistoolongandexceedsthelimit.toolongdomainnamethatiswaylongerthanthemaximumallowedlengthof253charactersandshouldfailvalidationbecauseitistoolongandexceedsthelimit.toolongdomainnamethatiswaylongerthanthemaximumallowedlengthof253charactersandshouldfailvalidationbecauseitistoolongandexceedsthelimit",
			expected: false,
		},
		{
			name:     "Invalid domain - starts with dot",
			domain:   ".invalid",
			expected: false,
		},
		{
			name:     "Invalid domain - ends with dot",
			domain:   "invalid.",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DomainValidator(tt.domain)
			assert.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}