package route

import (
	"context"
	"net/http"
)

// Handler is the signature used by all in this route.
type Handler func(context.Context, http.ResponseWriter, *http.Request) error
