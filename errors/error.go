package errors

import (
	"fmt"
	"net/http"
)

const (
	APIUserFault byte = iota
	APIUserFaultNeedsLog
	APIServerFault
	APITimedOut
	APINoCapacity
	APIUnexpectedValue
	APIDatabaseFailed
)

// APIError is the underlying error struct for when the API encounters an error.
type APIError struct {
	// Type is a uint8 that specifies the error type. If at any time, an error's type is
	// set to zero - it means we may not handle the error properly and it should be looked into.
	Type byte
	// Message is the message of the error.
	Message string
	// UnderlyingErr is the error that may have caused this error to be created. This is mainly used
	// when a function recovers, and we don't know what type the recover is.
	UnderlyingErr any
}

// New returns a new APIError
func New(t byte, msg string, uErr any) *APIError {
	return &APIError{
		Type:          t,
		Message:       msg,
		UnderlyingErr: uErr,
	}
}

func (err *APIError) Error() string {
	e := fmt.Sprintf("%s (%d)", err.Message, err.Type)
	if err.UnderlyingErr != nil {
		e += fmt.Sprintf(": %v", err.UnderlyingErr)
	}

	return e
}

func (err *APIError) StatusCode() int {
	switch err.Type {
	case APIUserFault, APIUserFaultNeedsLog:
		return http.StatusUnauthorized
	case APIServerFault, APIUnexpectedValue, APIDatabaseFailed:
		return http.StatusInternalServerError
	case APITimedOut:
		return http.StatusRequestTimeout
	case APINoCapacity:
		return http.StatusServiceUnavailable
	default:
		return http.StatusUnauthorized
	}
}
