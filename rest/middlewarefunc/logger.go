package middlewarefunc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rsachdeva/illuminatingdeposits/router"
	"go.opencensus.io/trace"
)

// Logger writes some information about the request to the logs in the
// format: TraceID : (200) GET /foo -> IP ADDR (latency)
func Logger(log *log.Logger) router.Middleware {

	// This is the actual middlewarefunc function to be executed.
	f := func(before router.Handler) router.Handler {

		// Create the handler that will be attached in the middlewarefunc chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			fmt.Printf("Entering Logger before handler is %T\n", before)
			defer fmt.Printf("Exiting Logger before handler is %T\n", before)

			ctx, span := trace.StartSpan(ctx, "internal.mid.RequestLogger")
			defer span.End()

			// If the context is missing this value, request the router
			// to be shutdown gracefully.
			v, ok := ctx.Value(router.KeyValues).(*router.Values)
			if !ok {
				return router.NewShutdownError("in logger mid web value missing from context")
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
