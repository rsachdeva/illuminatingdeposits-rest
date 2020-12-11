// Package postgreshealth provides postgress health status check service
package postgreshealth

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rsachdeva/illuminatingdeposits-rest/postgreshealth/healthvalue"
	"github.com/rsachdeva/illuminatingdeposits-rest/responder"
	"go.opencensus.io/trace"
)

// Service provides support for orchestration health checks.
type Service struct {
	Db *sqlx.DB

	// ADD OTHER STATE LIKE THE LOGGER IF NEEDED.
}

// Health validates the responder is healthy and ready to accept requests.
func (c *Service) Health(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "postgresconn.Service.Health")
	defer span.End()

	var health struct {
		Status string `json:"status"`
	}

	// Service if the postgresconn is ready.
	if err := healthvalue.StatusCheck(ctx, c.Db); err != nil {

		// If the postgresconn is not ready we will tell the cli and use a 500
		// status. Do not respond by just returning an error because further up in
		// the call stack will interpret that as an unhandled error.
		health.Status = "Db not ready"
		return responder.Respond(ctx, w, health, http.StatusInternalServerError)
	}

	health.Status = "ok"
	return responder.Respond(ctx, w, health, http.StatusOK)
}

func RegisterPostgresHealthService(db *sqlx.DB, m *responder.ServeMux) {
	// Register health check handler. This responder is not authenticated.
	c := Service{Db: db}
	m.Handle(http.MethodGet, "/v1/health", c.Health)
}