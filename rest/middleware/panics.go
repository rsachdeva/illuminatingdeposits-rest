package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/transport"
	"go.opencensus.io/trace"
)

// Panics recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.
func Panics(log *log.Logger) transport.Middleware {

	// This is the actual middleware function to be executed.
	f := func(after transport.Handler) transport.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			fmt.Printf("\t\t\t\t\tEntering Panics after handler is %T\n", after)
			defer fmt.Printf("\t\t\t\t\tExiting Panics before handler is %T\n", after)

			ctx, span := trace.StartSpan(ctx, "internal.mid.Panics")
			defer span.End()

			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, ok := ctx.Value(transport.KeyValues).(*transport.Values)
			if !ok {
				return transport.NewShutdownError("in panic mid web value missing from context")
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
