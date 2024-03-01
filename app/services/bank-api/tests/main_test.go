package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/qiushiyan/simplebank/app/services/bank-api/handlers"
	"github.com/qiushiyan/simplebank/business/auth/token"
	mockdb "github.com/qiushiyan/simplebank/business/db/mock"
	"github.com/qiushiyan/simplebank/foundation/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

var (
	adminToken string
	userToken  string
	Log        *zap.SugaredLogger
)

type DataResponse[T any] struct {
	Data T `json:"data"`
}

func TestMain(m *testing.M) {
	t, _ := token.NewToken("admin", []string{"ADMIN"}, 0)
	adminToken = "Bearer " + t.GetToken()

	t, _ = token.NewToken("user", []string{"USER"}, 0)
	userToken = "Bearer " + t.GetToken()

	logPath := fmt.Sprintf("%s/simplebank-log.txt", os.TempDir())
	Log, _ = logger.New("bank-api", logPath)
	fmt.Printf("logs at: %s\n", logPath)

	m.Run()
}

func serveRequest(
	t *testing.T,
	request *http.Request,
	buildStubs func(*mockdb.MockStore),
) *httptest.ResponseRecorder {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)

	buildStubs(store)

	cfg := handlers.APIMuxConfig{
		Shutdown: make(chan os.Signal, 1),
		Log:      Log,
		Store:    store,
	}

	app := handlers.NewMux(cfg)
	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, request)

	return recorder
}

// getResponseData is a helper for extracting the data field from a response body
func getResponseData[T any](t *testing.T, body *bytes.Buffer) T {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var got DataResponse[T]

	json.NewDecoder(bytes.NewReader(data)).Decode(&got)

	return got.Data
}
