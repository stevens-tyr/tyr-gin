package tyrgin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidEmail(t *testing.T) {
	// Test for invalid
	invalidEmails := []string{
		"",
		"test",
		"test@",
		"testing.domain",
		"noDotAfterAtSymbol@domain",
		"endswithdot@domain.",
		"test..twodotsinvalid@gmail.com",
		".firstchardot@domain.com",
		"specialCharsNotAllowedInDomain@domain!#.com",
	}
	for _, s := range invalidEmails {
		isValid := IsValidEmail(s) == nil
		assert.Falsef(t, isValid, "Invalid email detected as valid: '%s'", s)
	}
	// Valid emails, check for all valid chars
	validEmails := []string{
		"test@test.com",
		"test@test.subdomain.com",
		"test-testerson@testdoma.in",
		"!#$%&'*+-/=?^_`{|}~@domain.com", // All valid special chars
	}
	for _, s := range validEmails {
		isValid := IsValidEmail(s) == nil
		assert.Truef(t, isValid, "Valid email detected as invalid: '%s'")
	}
}
