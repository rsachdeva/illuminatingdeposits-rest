package rest

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rsachdeva/illuminatingdeposits/database"
	"github.com/rsachdeva/illuminatingdeposits/invest"
	"github.com/rsachdeva/illuminatingdeposits/route"
	"github.com/rsachdeva/illuminatingdeposits/user"
)

func registerCheckService(db *sqlx.DB, m *route.ServeMux) {
	// Register health check handler. This route is not authenticated.
	c := database.Service{Db: db}
	m.Handle(http.MethodGet, "/v1/health", c.Health)
}

func registerUserService(db *sqlx.DB, m *route.ServeMux) {
	// Register user interestsvc.
	u := user.Service{Db: db}

	// The route can't be authenticated because we need this route to
	// create the user in the first place.
	m.Handle(http.MethodPost, "/v1/users", u.Create)
}

func registerInvestService(log *log.Logger, m *route.ServeMux) {
	i := invest.Service{Log: log}
	m.Handle(http.MethodPost, "/v1/interests", i.Create)
}
