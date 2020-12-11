package app

import (
	"log"
	"net/http"

	"github.com/rsachdeva/illuminatingdeposits-rest/interestcal"
	"github.com/rsachdeva/illuminatingdeposits-rest/responder"
)

func RegisterInvestService(log *log.Logger, m *responder.ServeMux) {
	i := interestcal.Service{Log: log}
	m.Handle(http.MethodPost, "/v1/interests", i.Create)
}
