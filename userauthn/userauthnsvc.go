package userauthn

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/jsonfmt"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"github.com/rsachdeva/illuminatingdeposits-rest/reqlog"
	"github.com/rsachdeva/illuminatingdeposits-rest/userauthn/userauthnvalue"
	"go.opencensus.io/trace"
)

// Users holds interestsvc for dealing with usermgmt.
type service struct {
	db *sqlx.DB
}

func (svc service) CreateToken(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "usermgmt.Create")
	defer span.End()

	reqlog.Dump(r, "userauthn.CreateToken")

	var ctreq userauthnvalue.CreateTokenRequest
	if err := jsonfmt.Decode(r, &ctreq); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	ctresp, err := generateAccessToken(ctx, svc.db, &ctreq)
	if err != nil {
		return err
	}

	return jsonfmt.Respond(ctx, w, ctresp, http.StatusCreated)
}

func RegisterSvc(db *sqlx.DB, rt *muxhttp.Router) {
	u := service{db: db}
	rt.Handle(http.MethodPost, "/v1/users/token", u.CreateToken)
}
