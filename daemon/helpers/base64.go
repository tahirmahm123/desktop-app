
package helpers

import "regexp"

// ValidateBase64 is validating if a string fits base64 format (contains only base64 characters set and have correct format)
func ValidateBase64(base64str string) bool {
	// In base64 encoding, the character set is [A-Z, a-z, 0-9, and + /].
	// If the rest length is less than 4, the string is padded with '=' characters.
	base64Regexp := regexp.MustCompile("^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$")
	return base64Regexp.MatchString(base64str)
}
