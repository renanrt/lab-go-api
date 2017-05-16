package api

import (
	"encoding/json"
	"net/http"
)

// StatusCode defines an interface that will be used to determine if an error
// has meta data about the type of status code that a request should return.
// Errors that can be converted to this type will set the appropriate code
// on their api response.
type StatusCode interface {
	StatusCode() int
}

// ErrorCode defines an interface that will be used to display a user facing
// error code
type ErrorCode interface {
	ErrorCode() int
}

// ErrorFields defines an interface that returns a map of errors related to
// specific fields. This is generally used in validation errors.
type ErrorFields interface {
	ErrorFields() map[string]string
}

// ErrorLogger implements an interface for exposing information that could potentially
// be private and therefore not appropriate for the normal Error implementation.
type ErrorLogger interface {
	LoggableError() string
}

// ErrorContext implements an interface for returning structured context data
// about an error that can be used for logging purposes.
type ErrorContext interface {
	ErrorContext() map[string]interface{}
}

// ErrorResponse is the type that defines what an error payload will look like.
type ErrorResponse struct {
	Error  string            `json:"error"`
	Code   int               `json:"code,omitempty"`
	Fields map[string]string `json:"fields,omitempty"`
}

// RespondWithError takes an error, and determines the correct response code
// and body payload. This function is always the last to write to a response.
func RespondWithError(w http.ResponseWriter, r *http.Request, responseErr error) {
	w.Header().Add("Content-Type", "application/json")

	// Determine http status code
	code := http.StatusInternalServerError
	if statusCodeErr, ok := responseErr.(StatusCode); ok {
		// If this is an error which defines a status code, we should use that.
		code = statusCodeErr.StatusCode()
	}
	w.WriteHeader(code)

	// Create the response object
	response := ErrorResponse{Error: responseErr.Error()}

	// Add any detail fields to the response
	if fieldsErr, ok := responseErr.(ErrorFields); ok {
		response.Fields = fieldsErr.ErrorFields()
	}

	// Add a user error code to the response if possible
	if errorCodeErr, ok := responseErr.(ErrorCode); ok {
		response.Code = errorCodeErr.ErrorCode()
	} else {
		// Create a new private error as a fallback if the provided
		// error does not implement the ErrorCode interface
		response.Error = responseErr.Error()
	}

	// Generate a response body by creating an ErrorResponse instance with the
	// error message, and marshalling it to JSON.
	body, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		// This just won't happen, but I hate not assigning errors.
		panic(err)
	}

	w.Write(body)
}
