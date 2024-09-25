package endpoint

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/oomph-ac/api/database"
	"github.com/oomph-ac/api/endpoint/types"
	"github.com/oomph-ac/api/errors"
	"github.com/oomph-ac/api/internal"
	"github.com/oomph-ac/api/jwt"
	"github.com/oomph-ac/api/utils"
	"golang.org/x/exp/maps"
)

const (
	PathDownloadProxy = "/binary/download"
	PathUploadProxy   = "/binary/upload"
)

func init() {
	http.HandleFunc(PathDownloadProxy, DownloadProxy)
	http.HandleFunc(PathUploadProxy, UploadProxy)
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
	binData, dbErr := database.SearchForBinary(req.OS, req.Arch, req.Branch)
	if dbErr != nil {
		utils.EndpointError(r, dbErr, PathDownloadProxy)
		return
	}

	w.WriteHeader(http.StatusOK)
	enc.Encode(types.ProxyDownloadResponse{
		Data: binData,
	})
}

// UploadProxy is the endpoint meant to be reached by an internal tool to update binaries stored
// on the database. This endpoint requires an authentication JWT where the admin field is set to true
// which can be obtained via. the authentication endpoint. If an attempt to make an upload without
// having the sufficient permissions is made, the authentication key in question will be revoked.
func UploadProxy(w http.ResponseWriter, r *http.Request) {
	defer utils.CaptureAndRecover(w, r, PathUploadProxy)
	enc := json.NewEncoder(w)

	// Check that this is a POST request.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get(types.HeaderAuthToken)
	if token == "" {
		w.WriteHeader(http.StatusForbidden)
		enc.Encode(types.NewErrorResponse("missing " + types.HeaderAuthToken + " from header"))
		return
	}

	// Validate the authentication token given to us in the header.
	claims, err := jwt.ValidateAuthToken(token, utils.ClientIP(r))
	if err != nil {
		w.WriteHeader(err.StatusCode())
		enc.Encode(types.NewErrorResponse(err.Message))
		return
	}

	// Ensure that in the claims, the admin field is set to true. If a non-admin is trying to
	// upload a binary for whatever reason, say bye-bye to your authentication key - no refunds :>
	if !claims.Admin {
		data := internal.InfoPool.Get().(map[string]any)
		defer internal.InfoPool.Put(data)

		maps.Clear(data)
		data["key"] = claims.OomphKey

		database.DB.Query("DELETE oomphAuth WHERE key=$key;", data)
		w.WriteHeader(http.StatusTeapot)
		enc.Encode(types.NewErrorResponse(errors.MessageKeyInvalidatedDueToMalice))
		return
	}

	dat, ioErr := io.ReadAll(r.Body)
	if ioErr != nil {
		utils.EndpointError(r, errors.New(
			errors.APIServerFault,
			"failed to read body of HTTP request",
			ioErr,
		), PathUploadProxy)
		w.WriteHeader(http.StatusInternalServerError)
		enc.Encode(types.NewErrorResponse("failed to read request data"))
		return
	}

	var req types.ProxyUploadRequest
	if err := json.Unmarshal(dat, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(types.NewErrorResponse("invalid request"))
		return
	}

	if err := database.UpdateBinary(req.OS, req.Arch, req.Branch, req.Data); err != nil {
		utils.EndpointError(r, err, PathUploadProxy)
		w.WriteHeader(err.StatusCode())
		enc.Encode(types.NewErrorResponse(err.Message))
		return
	}

	w.WriteHeader(http.StatusOK)
}
