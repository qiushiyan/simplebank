package web

import (
	"context"
	"net/http"

	"github.com/go-json-experiment/json"
)

type dataWrapper struct {
	Data interface{} `json:"data"`
}

// RespondJson wraps the return the value as {data: value} according to the status code and sends json to the response writer
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
			bytes, err := json.Marshal(dataWrapper{Data: data})
			if err != nil {
				return err
			}
			w.Write(bytes)
			return nil
		default:
			bytes, err := json.Marshal(data)
			if err != nil {
				return err
			}
			w.Write(bytes)
			return nil
		}

	}

	return nil
}

// RespondJsonPlain sends json to the response writer without wrapping
func RespondJsonPlain(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	w.WriteHeader(statusCode)
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Write(bytes)
	return nil
}
