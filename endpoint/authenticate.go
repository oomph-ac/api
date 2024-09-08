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
	PathAuthentication = "/authenticate"
)

// Authenticate is the HTTP endpoint for giving authentication tokens to Oomph clients. These
// authentication tokens are used to access other endpoints, mainly resources such as detections,
// configuration, etc.
func Authenticate(w http.ResponseWriter, r *http.Request) {
	// Though this should never happen - better to be safe than sorry.
	defer utils.CaptureAndRecover(r, PathAuthentication)

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
			errors.APIInternalServer,
			"failed to read content of HTTP request body",
			ioErr,
		), PathAuthentication)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Attempt to unmarshal the body's content into an authentication request. If we are
	// unable to decode the response, it is invalid and we discard it.
	var request types.AuthRequest
	if err := json.Unmarshal(content, &request); err != nil {
		utils.EndpointWarning(r, PathAuthentication, "sent non-json request")
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
		utils.EndpointError(r, err, PathAuthentication)
		return
	}

	// Create a JWT token that contains claims for this authentication request.
	token, tkErr := jwt.NewAuthToken(res.Key, r.Header.Get("CF-Connecting-IP"))
	if tkErr != nil {
		err = errors.New(
			errors.APIInternalServer,
			"cannot create JWT token",
			tkErr,
		)
		w.WriteHeader(err.StatusCode())
		enc.Encode(types.NewErrorResponse(err.Message))
		utils.EndpointError(r, err, PathAuthentication)
		return
	}

	w.WriteHeader(http.StatusOK)
	enc.Encode(types.AuthResponse{Token: token})
}
