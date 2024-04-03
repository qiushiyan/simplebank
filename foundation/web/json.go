package web

import (
	"context"
	"encoding/json"
	"net/http"
)

type responseData struct {
	Data interface{} `json:"data,omitempty"`
}

// RespondJson wraps data according to the status code and sends json to the response writer.
func RespondJson(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	// Set the status code for the request logger middleware.
	SetStatusCode(ctx, statusCode)

	// Set the status code for the response writer.
	w.WriteHeader(statusCode)

	if statusCode == http.StatusNoContent {
		return nil
	}

	// If the data is not nil, encode it to the response writer.
	if data != nil {
		switch statusCode {
		case http.StatusOK, http.StatusCreated:
			return json.NewEncoder(w).Encode(responseData{Data: data})
		default:
			return json.NewEncoder(w).Encode(data)
		}

	}

	return nil
}

// RespondJsonPlain sends json to the response writer without wrapping the data.
func RespondJsonPlain(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}
