package middlewarefunc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/router"
	"go.opencensus.io/trace"
)

// Panics recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.
func Panics(log *log.Logger) router.Middleware {

	// This is the actual middlewarefunc function to be executed.
	f := func(after router.Handler) router.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			fmt.Printf("\t\t\t\t\tEntering Panics after handler is %T\n", after)
			defer fmt.Printf("\t\t\t\t\tExiting Panics before handler is %T\n", after)

			ctx, span := trace.StartSpan(ctx, "internal.mid.Panics")
			defer span.End()

			// If the context is missing this value, request the router
			// to be shutdown gracefully.
			v, ok := ctx.Value(router.KeyValues).(*router.Values)
			if !ok {
				return router.NewShutdownError("in panic mid web value missing from context")
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
