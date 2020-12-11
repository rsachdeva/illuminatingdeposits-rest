// Package interestcal provides interest calculation service
package interestcal

import (
	"context"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/appjson"
	"github.com/rsachdeva/illuminatingdeposits-rest/appmux"
	"github.com/rsachdeva/illuminatingdeposits-rest/interestcal/interestvalue"
	"github.com/rsachdeva/illuminatingdeposits-rest/reqlog"
	"go.opencensus.io/trace"
)

// Service handler
type Service struct {
	Log *log.Logger
}

// ListCalculations decodes the body of a request to create interest calculations. The full
// banks and deposit details with generated 30 days Service fields are sent back in the response.
func (*Service) ListCalculations(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "interestcal.Service.ListCalculations")
	defer span.End()

	reqlog.Dump(r)
	var nin interestvalue.NewInterest
	if err := appjson.Decode(r, &nin); err != nil {
		return errors.Wrap(err, "decoding new interest calculation request with banks and deposits")
	}

	in, err := nin.CalculateDelta()
	if err != nil {
		return errors.Wrap(err, "creating new interest calculations")
	}

	return appjson.Respond(ctx, w, &in, http.StatusCreated)
}

func RegisterSvc(log *log.Logger, m *appmux.Router) {
	i := Service{Log: log}
	m.Handle(http.MethodPost, "/v1/interests", i.ListCalculations)
}
