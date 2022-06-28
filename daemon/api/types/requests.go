
package types

// SessionNewRequest request to create new session
type SessionNewRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SessionDeleteRequest request to delete session
type SessionDeleteRequest struct {
	Session string `json:"session_token"`
}

// SessionStatusRequest request to get session status
type SessionStatusRequest struct {
	Session string `json:"session_token"`
}

// SessionStatusRequest request to get session status
type ServersListRequest struct {
	Session string `json:"session_token"`
}

// SessionWireGuardKeySetRequest request to set new WK key for a session
type SessionWireGuardKeySetRequest struct {
	Session            string `json:"session_token"`
	PublicKey          string `json:"public_key"`
	ConnectedPublicKey string `json:"connected_public_key"`
}
