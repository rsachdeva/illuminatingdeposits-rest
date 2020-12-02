package mux

// Middleware is a function designed to run some code before and/or after
// another Handler. It is designed to remove boilerplate or other concerns not
// direct to any given Handler.
type Middleware func(Handler) Handler

// wrapMiddleware creates a new handler by wrapping middleware around a final
// handler. The middlewares' Handlers will be executed by requests in the order
// they are provided.
func wrapMiddleware(mws []Middleware, handler Handler) Handler {

	// Loop backwards through the middleware invoking each one. Replace the
	// handler with the new wrapped handler. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mws) - 1; i >= 0; i-- {
		mw := mws[i]
		if mw != nil {
			// last middle ware wraps first
			// so in this list
			// 	app := web.NewReqMux(log, mid.Errors(log), mid.Metrics())
			// the last middleware mid.Metrics h is wrapped initially
			// and then mid.Errors h is wrapped
			handler = mw(handler)
		}
	}

	return handler
}
