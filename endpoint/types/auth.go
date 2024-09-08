package types

// AuthRequest is the request the Oomph client sends when it wants to authenticate with the API.
type AuthRequest struct {
	// Key is the authentication key of the client. This key should have
	// only a length of 64, and should be alpha-numeric.
	Key string `json:"key"`
}

// AuthResponse is the reponse to an authentication request when it is successful.
type AuthResponse struct {
	// Token is the authentication token the Oomph client must use for all other
	Token string `json:"token"`
}

// DBAuthData is the structure of the Oomph authentication data retrieved from the database.
type DBAuthData struct {
	// ID is the identifier for this particular authentication data in the database.
	// This is usually used to update this particular auth data from the database.
	ID string `json:"id"`
	// Admin is a boolean that is true if the specified user related to this
	// authentication key is able to run administrator actions on the API.
	Admin bool `json:"admin"`
	// Expiration is the Unix timestamp of when this authentication key expires.
	Expiration int64 `json:"expiration"`
	// IPList is a list of IP addresses allowed to use this authentication key.
	// If a user attempts to authenticate with this IP address, and the list is
	// empty, the address will automatically be added to the IP list.
	IPList []string `json:"ip_list"`
	// Key is the actual authentication key related to this authentication data.
	Key string `json:"key"`
	// Owner is the name or Discord tag of the owner of this authentication key.
	Owner string `json:"owner"`
}
