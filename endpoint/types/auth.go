package types

// AuthRequest is the request the Oomph client sends when it wants to authenticate with the API.
type AuthRequest struct {
	// Key is the authentication key of the client. This key should have
	// only a length of 64, and should be alpha-numeric.
	Key string `json:"key"`
}

func (r AuthRequest) Validate() bool {
	return len(r.Key) == 64
}
