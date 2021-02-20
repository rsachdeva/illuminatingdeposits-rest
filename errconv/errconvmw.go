// Package errconv provides conversion for error response for all services
package errconv

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/rsachdeva/illuminatingdeposits-rest/jsonfmt"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"go.opencensus.io/trace"
)

// NewMiddleware handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the cli in a uniform way.
// Unexpected errors (status >= 500) are logged.
func NewMiddleware(log *log.Logger) muxhttp.Middleware {

	// This is the actual middlewarefunc function to be executed.
	f := func(before muxhttp.Handler) muxhttp.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.StartSpan(ctx, "errconv.NewMiddleware")
			log.Println("span.SpanContext().IsSampled is", span.SpanContext().IsSampled())
			defer span.End()

			// If the context is missing this value, request the jsonfmt
			// to be shutdown gracefully.
			v, ok := ctx.Value(muxhttp.KeyValues).(*muxhttp.Values)
			if !ok {
				return muxhttp.NewShutdownError("in error mid web value missing from context")
			}

			// Run the handler chain and catch any propagated error.
			if err := before(ctx, w, r); err != nil {
				fmt.Println("NewMiddleware err is", err)
				// Log the error.
				log.Printf("TraceID %s : \n ERROR :\n %+v  web.IsShutdown(err) is %v", v.TraceID, err, muxhttp.IsShutdown(err))

				// Respond to the error.
				if err := jsonfmt.RespondError(ctx, w, err); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shutdown the jsonfmt.
				if ok := muxhttp.IsShutdown(err); ok {
					return err
				}
			}

			// Return nil to indicate the error has been handled.
			return nil
		}

		return h
	}

	return f
}
