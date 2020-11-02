package interestsvc

import (
	"context"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/internal/invest"
	"github.com/rsachdeva/illuminatingdeposits/internal/platform/debug"
	"github.com/rsachdeva/illuminatingdeposits/internal/platform/web"
	"go.opencensus.io/trace"
)

//Interest handler
type Interest struct {
	log *log.Logger
}

// Create decodes the body of a request to create interest calculations. The full
// banks and deposit details with generated 30 days Interest fields are sent back in the response.
func (*Interest) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "interestsvc.Invest.Create")
	defer span.End()

	debug.Dump(r)
	var nin invest.NewInterest
	if err := web.Decode(r, &nin); err != nil {
		return errors.Wrap(err, "decoding new interest calculation request with banks and deposits")
	}

	in, err := nin.ComputeDelta()
	if err != nil {
		return errors.Wrap(err, "creating new interest calculations")
	}

	return web.Respond(ctx, w, &in, http.StatusCreated)
}
