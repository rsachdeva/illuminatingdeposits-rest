package app

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rsachdeva/illuminatingdeposits-rest/dbconn"
	"github.com/rsachdeva/illuminatingdeposits-rest/invest"
	"github.com/rsachdeva/illuminatingdeposits-rest/responder"
	"github.com/rsachdeva/illuminatingdeposits-rest/user"
)

func registerDbHealthService(db *sqlx.DB, m *responder.ServeMux) {
	// Register health check handler. This responder is not authenticated.
	c := dbconn.Service{Db: db}
	m.Handle(http.MethodGet, "/v1/health", c.Health)
}

func registerUserService(db *sqlx.DB, m *responder.ServeMux) {
	// Register user interestsvc.
	u := user.Service{Db: db}

	// The responder can't be authenticated because we need this responder to
	// create the user in the first place.
	m.Handle(http.MethodPost, "/v1/users", u.Create)
}

func registerInvestService(log *log.Logger, m *responder.ServeMux) {
	i := invest.Service{Log: log}
	m.Handle(http.MethodPost, "/v1/interests", i.Create)
}
