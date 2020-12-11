// Package errconv provides conversion for error response for all services
package errconv

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/rsachdeva/illuminatingdeposits-rest/appmux"
	"github.com/rsachdeva/illuminatingdeposits-rest/appjson"
	"go.opencensus.io/trace"
)

// Middleware handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the cli in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Middleware(log *log.Logger) appmux.Middleware {

	// This is the actual middlewarefunc function to be executed.
	f := func(before appmux.Handler) appmux.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			fmt.Printf("\tEntering errconv middleware handler is %T\n", before)
			defer fmt.Printf("\tExiting errconv middleware handler is %T\n", before)

			ctx, span := trace.StartSpan(ctx, "errconv.Middleware")
			defer span.End()

			// If the context is missing this value, request the appjson
			// to be shutdown gracefully.
			v, ok := ctx.Value(appmux.KeyValues).(*appmux.Values)
			if !ok {
				return appmux.NewShutdownError("in error mid web value missing from context")
			}

			// Run the handler chain and catch any propagated error.
			if err := before(ctx, w, r); err != nil {
                fmt.Println("Middleware err is", err)
				// Log the error.
				log.Printf("TraceID %s : \n ERROR :\n %+v  web.IsShutdown(err) is %v", v.TraceID, err, appmux.IsShutdown(err))

				// Respond to the error.
				if err := appjson.RespondError(ctx, w, err); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shutdown the appjson.
				if ok := appmux.IsShutdown(err); ok {
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
