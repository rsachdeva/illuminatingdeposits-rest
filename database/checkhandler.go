package database

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rsachdeva/illuminatingdeposits/json"
	"github.com/rsachdeva/illuminatingdeposits/service"
	"go.opencensus.io/trace"
)

// CheckHandler provides support for orchestration health checks.
type CheckHandler struct {
	db *sqlx.DB

	// ADD OTHER STATE LIKE THE LOGGER IF NEEDED.
}

func RegisterCheckHandler(db *sqlx.DB, app *service.App) {
	{
		// Register health check handler. This route is not authenticated.
		c := CheckHandler{db: db}
		app.Handle(http.MethodGet, "/v1/health", c.Health)
	}
}

// Health validates the service is healthy and ready to accept requests.
func (c *CheckHandler) Health(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "database.CheckHandler.Health")
	defer span.End()

	var health struct {
		Status string `json:"status"`
	}

	// CheckHandler if the database is ready.
	if err := StatusCheck(ctx, c.db); err != nil {

		// If the database is not ready we will tell the cli and use a 500
		// status. Do not respond by just returning an error because further up in
		// the call stack will interpret that as an unhandled error.
		health.Status = "db not ready"
		return json.Respond(ctx, w, health, http.StatusInternalServerError)
	}

	health.Status = "ok"
	return json.Respond(ctx, w, health, http.StatusOK)
}
