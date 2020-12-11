// Package relog provides request logging for all services
package reqlog

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-rest/responder"
	"go.opencensus.io/trace"
)

// Middleware writes some information about the request to the logs in the
// format: TraceID : (200) GET /foo -> IP ADDR (latency)
func Middleware(log *log.Logger) responder.Middleware {

	// This is the actual middlewarefunc function to be executed.
	f := func(before responder.Handler) responder.Handler {

		// ListCalculations the handler that will be attached in the middlewarefunc chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			fmt.Printf("Entering reqlog Middleware handler is %T\n", before)
			defer fmt.Printf("Exiting reqlog Middleware handler is %T\n", before)

			ctx, span := trace.StartSpan(ctx, "reqlog.Middleware")
			defer span.End()

			// If the context is missing this value, request the responder
			// to be shutdown gracefully.
			v, ok := ctx.Value(responder.KeyValues).(*responder.Values)
			if !ok {
				return responder.NewShutdownError("in logger mid web value missing from context")
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
