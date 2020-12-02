package responder

import (
	"context"
	"net/http"
)

// Handler is the signature used by all in this responder.
type Handler func(context.Context, http.ResponseWriter, *http.Request) error
