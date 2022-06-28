
package types

import "fmt"

const (
	// CodeSuccess - success
	CodeSuccess int = 200

	// Unauthorized - Invalid Credentials	(Username or Password is not valid)
	Unauthorized int = 401

	// WGPublicKeyNotFound - WireGuard Public Key not found
	WGPublicKeyNotFound int = 424

	// SessionNotFound - Session not found Session not found
	SessionNotFound int = 601

	// CodeSessionsLimitReached - You've reached the session limit, log out from other device
	CodeSessionsLimitReached int = 602

	// AccountNotActive - account should be purchased
	AccountNotActive int = 702

	CaptchaRequired int = 70001
	CaptchaInvalid  int = 70002

	// Account has two-factor authentication enabled. Please enter TOTP token to login
	The2FARequired int = 70011
	// Specified two-factor authentication token is not valid
	The2FAInvalidToken int = 70012
)

// APIError - error, user not logged in into account
type APIError struct {
	ErrorCode int
	Message   string
}

// CreateAPIError creates new API error object
func CreateAPIError(errorCode int, message string) APIError {
	return APIError{
		ErrorCode: errorCode,
		Message:   message}
}

func (e APIError) Error() string {
	return fmt.Sprintf("API error: [%d] %s", e.ErrorCode, e.Message)
}
