package rest

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rsachdeva/illuminatingdeposits/database"
	"github.com/rsachdeva/illuminatingdeposits/invest"
	"github.com/rsachdeva/illuminatingdeposits/rest/mux"
	"github.com/rsachdeva/illuminatingdeposits/user"
)

func registerCheckService(db *sqlx.DB, m *mux.ReqMux) {
	// Register health check handler. This route is not authenticated.
	c := database.Service{Db: db}
	m.Handle(http.MethodGet, "/v1/health", c.Health)
}

func registerUserService(db *sqlx.DB, m *mux.ReqMux) {
	// Register user interestsvc.
	u := user.Service{Db: db}

	// The route can't be authenticated because we need this route to
	// create the user in the first place.
	m.Handle(http.MethodPost, "/v1/users", u.Create)
}

func registerInvestService(log *log.Logger, m *mux.ReqMux) {
	i := invest.Service{Log: log}
	m.Handle(http.MethodPost, "/v1/interests", i.Create)
}
