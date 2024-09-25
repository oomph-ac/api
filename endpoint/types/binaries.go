package types

// ProxyDownloadRequest is the request the Oomph authenticator makes when it wants to download the binary,
// or check if there are updates for it's current binary.
type ProxyDownloadRequest struct {
	// OS is the operating system for the client requesting the Oomph binary.
	OS string `json:"os"`
	// Arch is the architecture type for the client requesting the Oomph binary.
	Arch string `json:"arch"`
	// Branch is the branch that Oomph's binary should be downloaded from (e.g - stable, beta, dev, etc.)
	Branch string `json:"branch"`
}

// ProxyDownloadResponse is the response for a request to download the Oomph proxy.
type ProxyDownloadResponse struct {
	// Data contains the data of the Oomph proxy binary.
	Data string `json:"data"`
}

// ProxyUploadRequest is the request made by an internal Oomph tool makes when it wants
// to update a binary.
type ProxyUploadRequest struct {
	// OS is the operating system of this binary.
	OS string `json:"os"`
	// Arch is the architecture type for the client requesting the Oomph binary
	Arch string `json:"arch"`
	// Branch is the branch of the binary that should be updated (e.g - stable, beta, dev, etc.)
	Branch string `json:"branch"`
	// Data is the raw data of the Oomph binary that should be uploaded.
	Data string `json:"data"`
}

// DBProxyResponse is the database response for a binary request.
type DBProxyBinaryResponse struct {
	// Data is the raw data of the Oomph proxy to be sent back to the Oomph client.
	Data string `json:"data"`
}
