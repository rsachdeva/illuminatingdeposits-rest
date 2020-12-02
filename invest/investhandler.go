package invest

import (
	"context"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/debug"
	"github.com/rsachdeva/illuminatingdeposits/service"
	"go.opencensus.io/trace"
)

// InvestHandler handler
type InvestHandler struct {
	log *log.Logger
}

func RegisterInvestService(log *log.Logger, h *service.ReqHandler) {
	{
		i := InvestHandler{log: log}
		h.Handle(http.MethodPost, "/v1/interests", i.Create)
	}
}

// Create decodes the body of a request to create interest calculations. The full
// banks and deposit details with generated 30 days InvestHandler fields are sent back in the response.
func (*InvestHandler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "interestsvc.Invest.Create")
	defer span.End()

	debug.Dump(r)
	var nin NewInterest
	if err := service.Decode(r, &nin); err != nil {
		return errors.Wrap(err, "decoding new interest calculation request with banks and deposits")
	}

	in, err := nin.ComputeDelta()
	if err != nil {
		return errors.Wrap(err, "creating new interest calculations")
	}

	return service.Respond(ctx, w, &in, http.StatusCreated)
}
