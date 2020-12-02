package invest

import (
	"context"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/debug"
	"github.com/rsachdeva/illuminatingdeposits/route"
	"go.opencensus.io/trace"
)

// Service handler
type Service struct {
	Log *log.Logger
}

// Create decodes the body of a request to create interest calculations. The full
// banks and deposit details with generated 30 days Service fields are sent back in the response.
func (*Service) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "interestsvc.Invest.Create")
	defer span.End()

	debug.Dump(r)
	var nin NewInterest
	if err := route.Decode(r, &nin); err != nil {
		return errors.Wrap(err, "decoding new interest calculation request with banks and deposits")
	}

	in, err := nin.ComputeDelta()
	if err != nil {
		return errors.Wrap(err, "creating new interest calculations")
	}

	return route.Respond(ctx, w, &in, http.StatusCreated)
}
