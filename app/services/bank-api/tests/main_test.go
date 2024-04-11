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

	"github.com/qiushiyan/simplebank/app/services/bank-api/routes"
	"github.com/qiushiyan/simplebank/business/auth/token"
	mockdb "github.com/qiushiyan/simplebank/business/db/mock"
	simplemanager "github.com/qiushiyan/simplebank/business/task/simple"
	loggerlib "github.com/qiushiyan/simplebank/foundation/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

var (
	adminToken string
	userToken  string
	logger     *zap.SugaredLogger
)

type DataResponse[T any] struct {
	Data T `json:"data"`
}

func TestMain(m *testing.M) {

	rolesAdmin := []token.Role{token.RoleAdmin}
	rolesUser := []token.Role{token.RoleUser}

	t, _ := token.NewToken("admin", rolesAdmin, 0)
	adminToken = "Bearer " + t.Value

	t, _ = token.NewToken("user", rolesUser, 0)
	userToken = "Bearer " + t.Value

	logPath := fmt.Sprintf("%s/simplebank-log.txt", os.TempDir())
	logger, _ = loggerlib.New("bank-api", logPath)
	fmt.Printf("log at %s\n", logPath)
	defer logger.Sync()

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

	cfg := routes.MuxConfig{
		Shutdown: make(chan os.Signal, 1),
		Log:      logger,
		Store:    store,
		Task:     simplemanager.New(logger),
		Build:    "develop",
	}

	app := routes.NewMux(cfg)
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
