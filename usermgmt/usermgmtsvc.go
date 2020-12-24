// Package usermgmt provides user management service to add user to the system
package usermgmt

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/jsonfmt"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"github.com/rsachdeva/illuminatingdeposits-rest/reqlog"
	"github.com/rsachdeva/illuminatingdeposits-rest/usermgmt/uservalue"
	"go.opencensus.io/trace"
)

// Users holds interestsvc for dealing with usermgmt.
type service struct {
	db *sqlx.DB
}

func (us service) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "usermgmt.Create")
	defer span.End()

	reqlog.Dump(r, "usermgmt.Create")

	var nu uservalue.NewUser
	if err := jsonfmt.Decode(r, &nu); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	u, err := uservalue.AddUser(ctx, us.db, nu, time.Now())
	if err != nil {
		return muxhttp.NewRequestError(err, http.StatusConflict)
	}

	return jsonfmt.Respond(ctx, w, u, http.StatusCreated)
}

func RegisterSvc(db *sqlx.DB, rt *muxhttp.Router) {
	// Register usermgmt interestsvc.
	u := service{db: db}

	// The jsonfmt can't be authenticated because we need this jsonfmt to
	// create the usermgmt in the first place.
	rt.Handle(http.MethodPost, "/v1/users", u.Create)
}
