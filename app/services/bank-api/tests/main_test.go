package tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/qiushiyan/simplebank/app/services/bank-api/routes"
	"github.com/qiushiyan/simplebank/business/auth/token"
	mockdb "github.com/qiushiyan/simplebank/business/db/mock"
	simplemanager "github.com/qiushiyan/simplebank/business/task/simple"
	"github.com/qiushiyan/simplebank/foundation/logger"
	"github.com/qiushiyan/simplebank/foundation/web"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var (
	adminToken string
	userToken  string
	log        *logger.Logger
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
	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "SALES", traceIDFn, events)
	fmt.Printf("log at %s\n", logPath)

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

	cfg := routes.Config{
		Shutdown: make(chan os.Signal, 1),
		Log:      log,
		Store:    store,
		Task: simplemanager.New(simplemanager.Config{
			Log: log,
		}),
		Build: "develop",
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
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)
	return got.Data
}
