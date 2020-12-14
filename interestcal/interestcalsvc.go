// Package interestcal provides interest calculation service
package interestcal

import (
	"context"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/jsonfmt"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"github.com/rsachdeva/illuminatingdeposits-rest/interestcal/interestvalue"
	"github.com/rsachdeva/illuminatingdeposits-rest/reqlog"
	"go.opencensus.io/trace"
)

// service handler
type service struct {
	log *log.Logger
}

// ListCalculations decodes the body of a request to create interest calculations. The full
// banks and deposit details with generated 30 days service fields are sent back in the response.
func (s *service) ListCalculations(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "interestcal.service.ListCalculations")
	defer span.End()

	reqlog.Dump(r)
	var nin interestvalue.InterestRequest
	if err := jsonfmt.Decode(r, &nin); err != nil {
		return errors.Wrap(err, "decoding new interest calculation request with banks and deposits")
	}

	s.log.Println("Starting interest and 30day interest calculations(also called as delta)")
	in, err := nin.CalculateDelta()
	if err != nil {
		return errors.Wrap(err, "creating new interest calculations")
	}

	return jsonfmt.Respond(ctx, w, &in, http.StatusCreated)
}

func RegisterSvc(log *log.Logger, m *appmux.Router) {
	i := service{log: log}
	m.Handle(http.MethodPost, "/v1/interests", i.ListCalculations)
}
