package endpoint

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/oomph-ac/api/database"
	"github.com/oomph-ac/api/endpoint/types"
	"github.com/oomph-ac/api/errors"
	"github.com/oomph-ac/api/jwt"
	"github.com/oomph-ac/api/utils"
)

const (
	PathDownloadProxy = "/binaries/download"
	PathUploadProxy   = "/binaries/upload"
)

func init() {
	http.HandleFunc(PathDownloadProxy, DownloadProxy)
}

// DownloadProxy is the endpoint meant to be reached by the Oomph authenticator to download an
// Oomph proxy binary. This endpoint requires an authentication JWT to be passed in the header,
// which can be obtained via. the authentication endpoint.
func DownloadProxy(w http.ResponseWriter, r *http.Request) {
	// Capture any crashes that may happen.
	defer utils.CaptureAndRecover(w, r, PathDownloadProxy)

	// Make sure that this is a POST request.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Make sure that the client has a JWT token to access this endpoint.
	token := r.Header.Get(types.HeaderAuthToken)
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(types.NewErrorResponse(types.HeaderAuthToken + " missing from header"))
		return
	}

	// Validate the JWT token and ensure it's valid & has not expired. This will also check if the IP
	// address associated in the JWT's claims matches the client's IP address.
	enc := json.NewEncoder(w)
	if _, err := jwt.ValidateAuthToken(token, utils.ClientIP(r)); err != nil {
		utils.EndpointError(r, err, PathDownloadProxy)
		w.WriteHeader(err.StatusCode())
		enc.Encode(types.NewErrorResponse(err.Message))
		return
	}

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		utils.EndpointError(r, errors.New(
			errors.APIServerFault,
			"failed to decode JSON response",
			err,
		), PathDownloadProxy)

		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(types.ErrorResponse{
			Message: "unable to read request body",
		})
		return
	}

	// Decode the download request sent by the client.
	var req types.ProxyDownloadRequest
	if err := json.Unmarshal(dat, &req); err != nil {
		utils.EndpointError(r, errors.New(
			errors.APIUserFaultNeedsLog,
			"sent invalid JSON request",
			err,
		), PathDownloadProxy)
		return
	}

	// Search for the binary data from the database.
	binData, dbErr := database.SearchForBinary(req.OS, req.Arch)
	if dbErr != nil {
		utils.EndpointError(r, dbErr, PathDownloadProxy)
		return
	}

	w.WriteHeader(http.StatusOK)
	enc.Encode(types.ProxyDownloadResponse{
		Data: binData,
	})
}
