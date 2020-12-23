package userauthn

import (
	"context"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/jsonfmt"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"github.com/rsachdeva/illuminatingdeposits-rest/userauthn/userauthnvalue"
)

// Users holds interestsvc for dealing with usermgmt.
type service struct {
	db *sqlx.DB
}

func (svc *service) CreateToken(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var ctreq userauthnvalue.CreateTokenRequest
	ctresp, err := generateAccessToken(ctx, svc.db, &ctreq)
	if err != nil {
		return errors.Wrap(err, "generating access token")
	}

	return jsonfmt.Respond(ctx, w, &ctresp, http.StatusCreated)
}

func RegisterSvc(db *sqlx.DB, m *appmux.Router) {
	u := service{db: db}
	m.Handle(http.MethodPost, "/v1/token", u.CreateToken)
}
