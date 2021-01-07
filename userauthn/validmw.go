package userauthn

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

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

			log.Println("Authentication Middleware now going to verify token...")
			err := valid(r.Header["Authorization"])
			if err != nil {
				return err
			}

			// Call the next Handler and set its return value in the err variable.
			return handler(ctx, w, r)
		}

		return h
	}

	return f
}

// valid validates the authorization.
func valid(authnHeader []string) error {
	fmt.Printf("authnHeader is %v and len(authnHeader) is %v\n", authnHeader, len(authnHeader))
	if len(authnHeader) < 1 {
		return muxhttp.NewRequestError(
			errors.New("no authorization header"),
			http.StatusUnauthorized)
	}
	token := strings.TrimPrefix(authnHeader[0], "Bearer ")
	fmt.Println("token extracted for verification is: ", token)
	claims, err := verify(token)
	if err != nil {
		return err
	}
	email := claims.Email
	if len(email) < 1 {
		return muxhttp.NewRequestError(
			errors.New("invalid token without email"),
			http.StatusUnauthorized)
	}
	return nil
}
