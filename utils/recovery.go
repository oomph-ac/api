package utils

import (
	"net/http"

	"github.com/oomph-ac/api/errors"
)

// CaptureAndRecover is deffered in a function to prevent panics from crashing the API server.
// Every API endpoint should have this function deffered to avoid crashes.
func CaptureAndRecover(w http.ResponseWriter, r *http.Request, endpoint string) {
	// Capture whatever caused the function calling CaptureAndRecover to function.
	// If there is no error, the function ran successfully and we don't have to do anything.
	// However, if there is - we will log it and send the error to sentry.
	if v := recover(); v != nil {
		EndpointError(r, errors.New(errors.APIServerFault, endpoint+" crashed", v), endpoint)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"message\": \"crash occurred when processing request, no further info\"}"))
	}
}
