package endpoint

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/oomph-ac/api/endpoint/types"
	"github.com/oomph-ac/api/errors"
	"github.com/oomph-ac/api/utils"
)

const (
	PathAuthentication = "/authenticate"
)

// Authenticate is the HTTP endpoint for giving authentication tokens to Oomph clients. These
// authentication tokens are used to access other endpoints, mainly resources such as detections,
// and other
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

	// Validate that the fields in the authentication request are valid. If they are not,
	// then we discard it.
	if !request.Validate() {
		utils.EndpointWarning(r, PathAuthentication, "fields in auth request are invalid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: Make a request to the database for the authentication data.

	// Send a response back to the client.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"You made it here! What now?\"}"))
}
