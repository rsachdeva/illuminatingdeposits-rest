// Package relog provides request logging for all services and also request dump as needed
package reqlog

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"go.opencensus.io/trace"
)

// NewMiddleware writes some information about the request to the logs in the
// format: TraceID : (200) GET /foo -> IP ADDR (latency)
func NewMiddleware(log *log.Logger) muxhttp.Middleware {

	// This is the actual middlewarefunc function to be executed.
	f := func(handler muxhttp.Handler) muxhttp.Handler {

		// CreateInterest the handler that will be attached in the middlewarefunc chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.StartSpan(ctx, "reqlog.NewMiddleware")
			defer span.End()

			// If the context is missing this value, request the jsonfmt
			// to be shutdown gracefully.
			v, ok := ctx.Value(muxhttp.KeyValues).(*muxhttp.Values)
			if !ok {
				return muxhttp.NewShutdownError("in logger mid web value missing from context")
			}

			log.Printf("At the start reqLogmw: r.Header is %#v", r.Header)
			log.Printf("reqLogmw: r.Header[\"Authorization\"] is %#v", r.Header["Authorization"])

			log.Printf("handler func type being called %T", handler)
			err := handler(ctx, w, r)

			log.Printf("AT the end reqLogmw: %s : (%d) : %s %s -> %s (%s)",
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
