package web

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dimfeld/httptreemux"
	"github.com/go-json-experiment/json"
	"github.com/google/uuid"
	"github.com/qiushiyan/simplebank/foundation/validate"
)

// A Handler is a type that handles an http request within App
// different than the http.Handler in that it accepts context as the first parameter and returns an error
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// RouteGroup defines the behavior that binds routes to App
type RouteGroup interface {
	Register(*App)
}

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
}

// AddGroup accepts a list of RouteGroup and registers them with the App
func (a *App) AddGroup(rg ...RouteGroup) {
	for i := range rg {
		rg[i].Register(a)
	}
}

func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	// execute handler-specific middleware first
	handler = wrapMiddleware(mw, handler)
	// then apply app-side middleware
	handler = wrapMiddleware(a.mw, handler)

	a.handle(method, path, handler)
}

func (a *App) handle(method string, path string, handler Handler) {
	h := func(w http.ResponseWriter, r *http.Request) {
		// inject the values into the context
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}

		ctx := context.WithValue(r.Context(), key, &v)
		if err := handler(ctx, w, r); err != nil {
			// all errors should be processed in the errors middleware
			// except for the shutdown error
			if validateShutdown(err) {
				a.SignalShutdown()
				return
			}
		}
	}

	a.ContextMux.Handle(method, path, h)
}

func (a *App) GET(path string, handler Handler, mw ...Middleware) {
	a.Handle(http.MethodGet, path, handler, mw...)
}

func (a *App) POST(path string, handler Handler, mw ...Middleware) {
	a.Handle(http.MethodPost, path, handler, mw...)
}

func (a *App) PATCH(path string, handler Handler, mw ...Middleware) {
	a.Handle(http.MethodPatch, path, handler, mw...)
}

// Param returns the web call parameters from the request.
func Param(r *http.Request, name string) string {
	params := httptreemux.ContextParams(r.Context())
	return params[name]
}

// Decode decodes the body of an HTTP request and stores the result in dst.
func Decode(r *http.Request, val any) error {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return NewError(
				errors.New("request body cannot be empty"),
				http.StatusBadRequest,
			)
		}
		return NewError(fmt.Errorf("unable to read request body: %w", err), http.StatusBadRequest)
	}

	if err := json.Unmarshal(bytes, val); err != nil {
		return NewError(
			fmt.Errorf("unable to parse request body: %w", err),
			http.StatusBadRequest,
		)
	}
	if err := validate.Check(val); err != nil {
		return err
	}

	return nil

}

// SignalShutdown is used to gracefully shut down the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// validateShutdown validates the error for special conditions that do not
// warrant an actual shutdown by the system.
func validateShutdown(err error) bool {

	// Ignore syscall.EPIPE and syscall.ECONNRESET errors which occurs
	// when a write operation happens on the http.ResponseWriter that
	// has simultaneously been disconnected by the client (TCP
	// connections is broken). For instance, when large amounts of
	// data is being written or streamed to the client.
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// https://gosamples.dev/broken-pipe/
	// https://gosamples.dev/connection-reset-by-peer/

	switch {
	case errors.Is(err, syscall.EPIPE):

		// Usually, you get the broken pipe error when you write to the connection after the
		// RST (TCP RST Flag) is sent.
		// The broken pipe is a TCP/IP error occurring when you write to a stream where the
		// other end (the peer) has closed the underlying connection. The first write to the
		// closed connection causes the peer to reply with an RST packet indicating that the
		// connection should be terminated immediately. The second write to the socket that
		// has already received the RST causes the broken pipe error.
		return false

	case errors.Is(err, syscall.ECONNRESET):

		// Usually, you get connection reset by peer error when you read from the
		// connection after the RST (TCP RST Flag) is sent.
		// The connection reset by peer is a TCP/IP error that occurs when the other end (peer)
		// has unexpectedly closed the connection. It happens when you send a packet from your
		// end, but the other end crashes and forcibly closes the connection with the RST
		// packet instead of the TCP FIN, which is used to close a connection under normal
		// circumstances.
		return false
	}

	return true
}
