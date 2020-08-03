package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/rsachdeva/illuminatingdeposits/internal/mid"
	"github.com/rsachdeva/illuminatingdeposits/internal/platform/web"
)

// API constructs an http.Handler with all application routes defined.
func API(shutdown chan os.Signal, db *sqlx.DB, log *log.Logger) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown, log, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	{
		// Register health check handler. This route is not authenticated.
		c := Check{db: db}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}

	{
		i := Interest{log: log}
		app.Handle(http.MethodPost, "/v1/interests", i.Create)
	}

	return app
}
