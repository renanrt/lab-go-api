package api

import (
	"encoding/json"
	"io"
	"net/http"
)

// CollectionResponse represents a collection payload from the API.
type CollectionResponse struct {
	Data interface{} `json:"data"`
}

// RespondWithData will write the supplied data in the standard API response
// format, and set the appropriate status code.
func RespondWithData(w http.ResponseWriter, r *http.Request, data interface{}, code int) {
	body, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		// This just won't happen, but I hate not assigning errors.
		panic(err)
	}

	// Write the content type first, before headers are locked in..
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
}

// RespondWithCollection send the supplied data as a collection to the
// RespondWithData func.
func RespondWithCollection(w http.ResponseWriter, r *http.Request, data interface{}, code int) {
	RespondWithData(w, r, CollectionResponse{data}, code)
}

// RespondWithStatusCode writes the response code and any application headers.
func RespondWithStatusCode(w http.ResponseWriter, r *http.Request, code int) {
	w.WriteHeader(code)
}

// RespondWithIOReader writes the response code and anything coming from the io.Reader.
func RespondWithIOReader(w http.ResponseWriter, r *http.Request, reader io.Reader, code int) {
	w.WriteHeader(code)
	io.Copy(w, reader)
}
