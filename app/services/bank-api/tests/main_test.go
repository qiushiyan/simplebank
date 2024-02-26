package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type DataResponse[T any] struct {
	Data T `json:"data"`
}

// getResponseData is a helper for extracting the data field from a response body
func getResponseData[T any](t *testing.T, body *bytes.Buffer) T {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var got DataResponse[T]

	json.NewDecoder(bytes.NewReader(data)).Decode(&got)

	return got.Data
}
