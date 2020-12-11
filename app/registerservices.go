package app

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rsachdeva/illuminatingdeposits-rest/interestcal"
	"github.com/rsachdeva/illuminatingdeposits-rest/responder"
	"github.com/rsachdeva/illuminatingdeposits-rest/usermgmt"
)

func RegisterUserService(db *sqlx.DB, m *responder.ServeMux) {
	// Register usermgmt interestsvc.
	u := usermgmt.Service{Db: db}

	// The responder can't be authenticated because we need this responder to
	// create the usermgmt in the first place.
	m.Handle(http.MethodPost, "/v1/users", u.Create)
}

func RegisterInvestService(log *log.Logger, m *responder.ServeMux) {
	i := interestcal.Service{Log: log}
	m.Handle(http.MethodPost, "/v1/interests", i.Create)
}
