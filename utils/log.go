package utils

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// EndpointWarning is used when a certain API endpoint wants to log a warning. These warnings should be
// used when there is an invalid request, or a rejected request.
func EndpointWarning(r *http.Request, endpoint, msg string) {
	log.Warn().Str("endpoint", endpoint).Str("ip", r.RemoteAddr).Msg(msg)
}

// EndpointError is used when a certain API endpoint encounters an error. This function will log the error,
// and report the error to Sentry as well.
func EndpointError(r *http.Request, err error, endpoint string) {
	log.Error().Err(err).Str("endpoint", endpoint).Str("ip", r.RemoteAddr)
}
