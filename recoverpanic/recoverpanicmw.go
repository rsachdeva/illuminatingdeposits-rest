// Package recoverpanic provides panic recovery for all services. It uses go built in recovery() function
package recoverpanic

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"go.opencensus.io/trace"
)

// NewMiddleware recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.
func NewMiddleware(log *log.Logger) muxhttp.Middleware {

	// This is the actual middlewarefunc function to be executed.
	f := func(after muxhttp.Handler) muxhttp.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			ctx, span := trace.StartSpan(ctx, "recoverpanic.NewMiddleware")
			defer span.End()

			// If the context is missing this value, request the jsonfmt
			// to be shutdown gracefully.
			v, ok := ctx.Value(muxhttp.KeyValues).(*muxhttp.Values)
			if !ok {
				return muxhttp.NewShutdownError("in panic mid web value missing from context")
			}

			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if r := recover(); r != nil {
					err = errors.Errorf("panic: %v", r)

					// Log the Go stack trace for this panic'd goroutine.
					log.Printf("%s :\n%s", v.TraceID, debug.Stack())
				}
			}()

			// Call the next Handler and set its return value in the err variable.
			return after(ctx, w, r)
		}

		return h
	}

	return f
}
