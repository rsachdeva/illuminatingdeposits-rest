package userauthn

import (
	"context"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"go.opencensus.io/trace"
)

func NewMiddleware(log *log.Logger) muxhttp.Middleware {

	// This is the actual middlewarefunc function to be executed.
	f := func(handler muxhttp.Handler) muxhttp.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.StartSpan(ctx, "userauthn.NewMiddleware")
			defer span.End()

			log.Println("Authentication Middleware")
			return muxhttp.NewRequestError(
				errors.New("no authorization header"),
				http.StatusUnauthorized)

			// Call the next Handler and set its return value in the err variable.
			// return handler(ctx, w, r)
		}

		return h
	}

	return f
}
