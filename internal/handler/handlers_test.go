package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateLongURL(t *testing.T) {
	testCases := []struct {
		name      string
		rawURL    string
		wantError bool
	}{
		{
			name:      "valid https URL",
			rawURL:    "https://example.com/path",
			wantError: false,
		},
		{
			name:      "valid http URL",
			rawURL:    "http://example.com/path",
			wantError: false,
		},
		{
			name:      "invalid URL format",
			rawURL:    "not a url",
			wantError: true,
		},
		{
			name:      "invalid scheme",
			rawURL:    "ftp://example.com/path",
			wantError: true,
		},
		{
			name:      "missing host",
			rawURL:    "https:///path-only",
			wantError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateLongURL(testCase.rawURL)
			if testCase.wantError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
