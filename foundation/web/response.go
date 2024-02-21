package web

import (
	"context"
	"encoding/json"
	"net/http"
)

type ResponseData struct {
	Data interface{} `json:"data,omitempty"`
}

func RespondJson(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	// Set the status code for the request logger middleware.
	SetStatusCode(ctx, statusCode)

	// Set the status code for the response writer.
	w.WriteHeader(statusCode)

	// If the data is not nil, encode it to the response writer.
	if data != nil {
		switch statusCode {
		case http.StatusOK, http.StatusCreated:
			return json.NewEncoder(w).Encode(ResponseData{Data: data})
		default:
			return json.NewEncoder(w).Encode(data)
		}

	}

	return nil
}
