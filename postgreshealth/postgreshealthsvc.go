// Package postgreshealth provides postgress health status check service
package postgreshealth

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"github.com/rsachdeva/illuminatingdeposits-rest/postgreshealth/healthvalue"
	"github.com/rsachdeva/illuminatingdeposits-rest/jsonfmt"
	"go.opencensus.io/trace"
)

// service provides support for orchestration health checks.
type service struct {
	db *sqlx.DB

	// ADD OTHER STATE LIKE THE LOGGER IF NEEDED.
}

// Health validates the jsonfmt is healthy and ready to accept requests.
func (c *service) Health(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "postgresconn.service.Health")
	defer span.End()

	var health struct {
		Status string `json:"status"`
	}

	// service if the postgresconn is ready.
	if err := healthvalue.StatusCheck(ctx, c.db); err != nil {

		// If the postgresconn is not ready we will tell the cli and use a 500
		// status. Do not respond by just returning an error because further up in
		// the call stack will interpret that as an unhandled error.
		health.Status = "Db not ready"
		return jsonfmt.Respond(ctx, w, health, http.StatusInternalServerError)
	}

	health.Status = "ok"
	return jsonfmt.Respond(ctx, w, health, http.StatusOK)
}

func RegisterSvc(db *sqlx.DB, m *appmux.Router) {
	// Register health check handler. This jsonfmt is not authenticated.
	c := service{db: db}
	m.Handle(http.MethodGet, "/v1/health", c.Health)
}