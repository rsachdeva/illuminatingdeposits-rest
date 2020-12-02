package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rsachdeva/illuminatingdeposits/rest/service"
	"go.opencensus.io/trace"
)

// Logger writes some information about the request to the logs in the
// format: TraceID : (200) GET /foo -> IP ADDR (latency)
func Logger(log *log.Logger) service.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before service.Handler) service.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			fmt.Printf("Entering Logger before handler is %T\n", before)
			defer fmt.Printf("Exiting Logger before handler is %T\n", before)

			ctx, span := trace.StartSpan(ctx, "internal.mid.RequestLogger")
			defer span.End()

			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, ok := ctx.Value(service.KeyValues).(*service.Values)
			if !ok {
				return service.NewShutdownError("in logger mid web value missing from context")
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
