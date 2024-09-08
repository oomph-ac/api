package utils

import (
	"net/http"

	"github.com/oomph-ac/api/errors"
	"github.com/rs/zerolog/log"
)

// EndpointWarning is used when a certain API endpoint wants to log a warning. These warnings should be
// used when there is an invalid request, or a rejected request.
func EndpointWarning(r *http.Request, endpoint, msg string) {
	log.Warn().Str("endpoint", endpoint).Str("ip", r.Header.Get("CF-Connecting-IP")).Msg(msg)
}

// EndpointError is used when a certain API endpoint encounters an error. This function will log the error,
// and report the error to Sentry as well. However, the error will not be logged or reported if the
// error is because of a "user fault" (e.g - the DB cannot find a key in the database, but is able to query it)
func EndpointError(r *http.Request, err *errors.APIError, endpoint string) {
	if err.Type == errors.APIUserFault {
		return
	}

	log.Error().Err(err).Str("endpoint", endpoint).Str("ip", r.Header.Get("CF-Connecting-IP"))
}
