package utils

import (
	"regexp"
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	const expectedLength = 6
	// Regex to match only alphanumeric characters (A-Z, a-z, 0-9), 6 characters long
	regex := regexp.MustCompile(`^[a-zA-Z0-9]{6}$`)

	code := GenerateShortCode()

	// Check length
	if len(code) != expectedLength {
		t.Errorf("Expected code length %d, got %d (%s)", expectedLength, len(code), code)
	}

	// Check regex match
	if !regex.MatchString(code) {
		t.Errorf("Generated code %s does not match expected format", code)
	}
}
