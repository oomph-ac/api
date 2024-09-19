package endpoint

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/oomph-ac/api/database"
	"github.com/oomph-ac/api/endpoint/types"
	"github.com/oomph-ac/api/errors"
	"github.com/oomph-ac/api/jwt"
	"github.com/oomph-ac/api/utils"
)

const (
	PathAuthenticate  = "/auth/login"
	PathVerifySession = "/auth/verify"
)

func init() {
	http.HandleFunc(PathAuthenticate, Authenticate)
	http.HandleFunc(PathVerifySession, VerifySession)
}

// Authenticate is the HTTP endpoint for giving authentication tokens to Oomph clients. These
// authentication tokens are used to access other endpoints, mainly resources such as detections,
// configuration, etc.
func Authenticate(w http.ResponseWriter, r *http.Request) {
	// Though this should never happen - better to be safe than sorry.
	defer utils.CaptureAndRecover(w, r, PathAuthenticate)

	// Check if this is a POST request, and reject all other requests.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Check if the body's content type is of JSON, and reject it if not.
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Read the content sent by the client.
	content, ioErr := io.ReadAll(r.Body)
	if ioErr != nil {
		// This should never happen - we'd want to report this error to Sentry.
		// TODO: Report the error to Sentry.
		utils.EndpointError(r, errors.New(
			errors.APIServerFault,
			"failed to read content of HTTP request body",
			ioErr,
		), PathAuthenticate)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Attempt to unmarshal the body's content into an authentication request. If we are
	// unable to decode the response, it is invalid and we discard it.
	var request types.AuthRequest
	if err := json.Unmarshal(content, &request); err != nil {
		utils.EndpointWarning(r, PathAuthenticate, "sent non-json request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create a new JSON encoder for the HTTP response.
	enc := json.NewEncoder(w)

	// Obtain the authentication data from the database.
	res, err := database.ObtainAuth(request.Key)
	if err != nil {
		w.WriteHeader(err.StatusCode())
		enc.Encode(types.NewErrorResponse(err.Message))
		utils.EndpointError(r, err, PathAuthenticate)
		return
	}

	if !slices.Contains(res.IPList, utils.ClientIP(r)) { // Validate that this IP address is allowed to use the authentication key.
		err = errors.New(
			errors.APIUserFaultNeedsLog,
			"IP address not allowed to use this authentication key - this incident has been reported",
			nil,
		)
		utils.EndpointError(r, err, PathAuthenticate)
		w.WriteHeader(err.StatusCode())
		enc.Encode(types.NewErrorResponse(err.Message))
		return
	}

	// Create a JWT token that contains claims for this authentication request.
	token, tkErr := jwt.NewAuthToken(res, utils.ClientIP(r))
	if tkErr != nil {
		err = errors.New(
			errors.APIServerFault,
			"cannot create JWT token",
			tkErr,
		)
		w.WriteHeader(err.StatusCode())
		enc.Encode(types.NewErrorResponse(err.Message))
		utils.EndpointError(r, err, PathAuthenticate)
		return
	}

	w.WriteHeader(http.StatusOK)
	enc.Encode(types.AuthResponse{
		Token:     token,
		RefreshAt: time.Now().Add(time.Minute * 30).Unix(),
	})
}

// VerifySession ensures that a token the Oomph client has is valid. It does not send any content, only
// writing the status code to represent wether the JWT is valid or not.
func VerifySession(w http.ResponseWriter, r *http.Request) {
	defer utils.CaptureAndRecover(w, r, PathVerifySession)

	// Make sure the request is a POST request.
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Make sure we have an auth token to verify
	token := r.Header.Get(types.HeaderAuthToken)
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Verify that the authentication token is valid.
	if _, err := jwt.ValidateAuthToken(token, utils.ClientIP(r)); err != nil {
		utils.EndpointError(r, err, PathVerifySession)
		w.WriteHeader(err.StatusCode())
		return
	}
	w.WriteHeader(http.StatusOK)
}
