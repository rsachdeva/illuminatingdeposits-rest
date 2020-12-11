// Package usermgmt provides user management service to add user to the system
package usermgmt

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/auth/authvalue"
	"github.com/rsachdeva/illuminatingdeposits-rest/responder"
	"github.com/rsachdeva/illuminatingdeposits-rest/usermgmt/uservalue"
	"go.opencensus.io/trace"
)

// Users holds interestsvc for dealing with usermgmt.
type Service struct {
	Db *sqlx.DB
}

func (us *Service) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "usermgmt.Service.ListCalculations")
	defer span.End()

	email, pass, ok := r.BasicAuth()
	fmt.Printf("r.Header.GetGet(\"Authorization\") is %s", r.Header.Get("Authorization"))
	if !ok {
		err := errors.New("must provide email and password in Basic auth")
		return responder.NewRequestError(err, http.StatusUnauthorized)
	}

	nu := uservalue.NewUser{
		Email:           email,
		Password:        pass,
		PasswordConfirm: pass,
		Roles:           []string{authvalue.RoleUser},
	}

	u, err := uservalue.AddUser(ctx, us.Db, nu, time.Now())
	if err != nil {
		return responder.NewRequestError(err, http.StatusConflict)
	}

	return responder.Respond(ctx, w, u, http.StatusCreated)
}

func RegisterSvc(db *sqlx.DB, m *responder.ServeMux) {
	// Register usermgmt interestsvc.
	u := Service{Db: db}

	// The responder can't be authenticated because we need this responder to
	// create the usermgmt in the first place.
	m.Handle(http.MethodPost, "/v1/users", u.Create)
}
