// Package muxhttp provides mux, middleware, handler with regards to routing capabilities
package appmux

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Values carries information about each request.
type Values struct {
	TraceID    string
	StatusCode int
	Start      time.Time
}

// Router is the entrypoint into our application and what controls the context of
// each request. Feel free to add any configuration data/logic on this type.
type Router struct {
	log      *log.Logger
	mux      *chi.Mux
	mws      []Middleware
	och      *ochttp.Handler
	shutdown chan os.Signal
}

// NewRouter constructs an Router to handle a set of routes. Any Middleware provided
// will be ran for every request.
func NewRouter(shutdownCh chan os.Signal, log *log.Logger, mw ...Middleware) *Router {
	m := Router{
		log:      log,
		mux:      chi.NewRouter(),
		mws:      mw,
		shutdown: shutdownCh,
	}

	// ListCalculations an OpenCensus HTTP Handler which wraps the jsonfmt. This will start
	// the initial span and annotate it with information about the request/response.
	//
	// This is configured to use the W3C TraceContext standard to set the remote
	// parent if an cli request includes the appropriate headers.
	// https://w3c.github.io/trace-context/
	m.och = &ochttp.Handler{
		Handler:     m.mux,
		Propagation: &tracecontext.HTTPFormat{},
	}

	return &m
}

// Handle associates a handler function with an HTTP Method and URL pattern.
//
// It converts our custom handler type to the std lib Handler type. It captures
// errors from the handler and serves them to the cli in a uniform way.
func (a *Router) Handle(method, url string, h Handler, mw ...Middleware) {

	// First wrap handler specific middlewarefunc around this handler.
	slicemws := append(a.mws, mw...)
	fmt.Println("slicemws is", slicemws)
	h = wrapMiddleware(slicemws, h)

	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx, span := trace.StartSpan(r.Context(), "muxhttp.ServerMux.Handle")
		defer span.End()

		// ListCalculations a Values struct to record state for the request. Store the
		// address in the request's context so it is sent down the call chain.
		v := Values{
			TraceID: span.SpanContext().TraceID.String(),
			Start:   time.Now(),
		}
		ctx = context.WithValue(ctx, KeyValues, &v)

		// Run the handler chain and catch any propagated error.
		if err := h(ctx, w, r); err != nil {
			a.log.Printf("%s : unhandled error: %+v", v.TraceID, err)
			if IsShutdown(err) {
				a.SignalShutdown()
			}
		}
	}

	a.mux.MethodFunc(method, url, fn)
}

// ServeHTTP implements the http.Handler interface.
func (a *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.och.ServeHTTP(w, r)
}

// SignalShutdown is used to gracefully shutdown the appserver when an integrity
// issue is identified.
func (a *Router) SignalShutdown() {
	a.log.Println("error returned from handler indicated integrity issue, shutting down jsonfmt")
	a.shutdown <- syscall.SIGSTOP
}
