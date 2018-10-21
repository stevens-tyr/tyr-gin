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
		"testing.gmail",
		"noDotAfterAtSymbol@gmail",
		"endswithdot@domain.",
		"test..twodotsinvalid@gmail.com",
		"hostinvalid@gmail.co.uk",
		".firstchardot@gmail.com",
		"specialCharsNotAllowedInDomain@gmail!#.com",
	}
	for _, s := range invalidEmails {
		err := IsValidEmail(s)
		assert.Truef(t, err != nil, "Invalid email detected as valid: '%s'", s)
	}
	// Valid emails, check for all valid chars
	validEmails := []string{
		"test@lists.stevens.edu", // Subdomains are valid
		"test@gmail.com",
		"someguy@stevens.edu",
		"!#$%&'*+-/=?^_`{|}~@yahoo.com", // All valid special chars
	}
	for _, s := range validEmails {
		err := IsValidEmail(s)
		assert.Truef(t, err == nil, "Valid email detected as invalid: '%s', err: '%s'", s, err)
	}
}
