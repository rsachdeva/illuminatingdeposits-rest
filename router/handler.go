package router

import (
	"context"
	"net/http"
)

// Handler is the signature used by all in this router.
type Handler func(context.Context, http.ResponseWriter, *http.Request) error
