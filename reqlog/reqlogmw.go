// Package relog provides request logging for all services and also request dump as needed
package reqlog

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-rest/appmux"
	"go.opencensus.io/trace"
)

// NewMiddleware writes some information about the request to the logs in the
// format: TraceID : (200) GET /foo -> IP ADDR (latency)
func NewMiddleware(log *log.Logger) appmux.Middleware {

	// This is the actual middlewarefunc function to be executed.
	f := func(before appmux.Handler) appmux.Handler {

		// ListCalculations the handler that will be attached in the middlewarefunc chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			fmt.Printf("Entering reqlog NewMiddleware handler is %T\n", before)
			defer fmt.Printf("Exiting reqlog NewMiddleware handler is %T\n", before)

			ctx, span := trace.StartSpan(ctx, "reqlog.NewMiddleware")
			defer span.End()

			// If the context is missing this value, request the appjson
			// to be shutdown gracefully.
			v, ok := ctx.Value(appmux.KeyValues).(*appmux.Values)
			if !ok {
				return appmux.NewShutdownError("in logger mid web value missing from context")
			}

			err := before(ctx, w, r)

			log.Printf("%s : (%d) : %s %s -> %s (%s)",
				v.TraceID, v.StatusCode,
				r.Method, r.URL.Path,
				r.RemoteAddr, time.Since(v.Start),
			)

			// Return the error so it can be handled further up the chain.
			return err
		}

		return h
	}

	return f
}
