// Package interestcal provides interest calculation service for the server api
package interestcal

import (
	"context"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/interestcal/interestvalue"
	"github.com/rsachdeva/illuminatingdeposits-rest/jsonfmt"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"github.com/rsachdeva/illuminatingdeposits-rest/reqlog"
	"github.com/rsachdeva/illuminatingdeposits-rest/userauthn"
	"go.opencensus.io/trace"
)

// service handler
type service struct {
	log *log.Logger
}

// CreateInterest decodes the body of a request to create interest calculations. The full
// banks and deposit details with generated 30 days service fields are sent back in the response.
// https://cloud.google.com/apis/design/standard_methods
func (s service) CreateInterest(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "interestcal.service.CreateInterest")
	defer span.End()
	reqlog.Dump(r, "interestcal.CreateInterest")
	var cireq interestvalue.CreateInterestRequest
	if err := jsonfmt.Decode(r, &cireq); err != nil {
		return errors.Wrap(err, "decoding new interest calculation request with banks and deposits")
	}
	s.log.Printf("\nDecoded json is %+v\n", cireq)
	s.log.Println("Starting interest and 30day interest calculations(also called as delta)")
	ciresp, err := cireq.CalculateDelta()
	if err != nil {
		return errors.Wrap(err, "creating new interest calculations")
	}

	return jsonfmt.Respond(ctx, w, &ciresp, http.StatusCreated)
}

func RegisterSvc(log *log.Logger, rt *muxhttp.Router) {
	i := service{log: log}
	rt.Handle(http.MethodPost, "/v1/interests", i.CreateInterest, userauthn.NewMiddleware(log))
}
