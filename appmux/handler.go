package appmux

import (
	"context"
	"net/http"
)

// Handler is the signature used by all in this appjson.
type Handler func(context.Context, http.ResponseWriter, *http.Request) error
