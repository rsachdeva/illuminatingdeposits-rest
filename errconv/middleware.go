// Package errconv provides conversion for errors
package errconv

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/rsachdeva/illuminatingdeposits-rest/responder"
	"go.opencensus.io/trace"
)

// Middleware handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the cli in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Middleware(log *log.Logger) responder.Middleware {

	// This is the actual middlewarefunc function to be executed.
	f := func(before responder.Handler) responder.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			fmt.Printf("\tEntering Middleware before handler is %T\n", before)
			defer fmt.Printf("\tExiting Middleware before handler is %T\n", before)

			ctx, span := trace.StartSpan(ctx, "internal.mid.Middleware")
			defer span.End()

			// If the context is missing this value, request the responder
			// to be shutdown gracefully.
			v, ok := ctx.Value(responder.KeyValues).(*responder.Values)
			if !ok {
				return responder.NewShutdownError("in error mid web value missing from context")
			}

			// Run the handler chain and catch any propagated error.
			if err := before(ctx, w, r); err != nil {
                fmt.Println("Middleware err is", err)
				// Log the error.
				log.Printf("TraceID %s : \n ERROR :\n %+v  web.IsShutdown(err) is %v", v.TraceID, err, responder.IsShutdown(err))

				// Respond to the error.
				if err := responder.RespondError(ctx, w, err); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shutdown the responder.
				if ok := responder.IsShutdown(err); ok {
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
