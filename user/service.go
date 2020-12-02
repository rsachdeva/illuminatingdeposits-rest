package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/auth"
	"github.com/rsachdeva/illuminatingdeposits/route"
	"go.opencensus.io/trace"
)

// Users holds interestsvc for dealing with user.
type Service struct {
	Db *sqlx.DB
}

func (us *Service) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "interestsvc.Users.Create")
	defer span.End()

	email, pass, ok := r.BasicAuth()
	fmt.Printf("r.Header.GetGet(\"Authorization\") is %s", r.Header.Get("Authorization"))
	if !ok {
		err := errors.New("must provide email and password in Basic auth")
		return route.NewRequestError(err, http.StatusUnauthorized)
	}

	nu := NewUser{
		Email:           email,
		Password:        pass,
		PasswordConfirm: pass,
		Roles:           []string{auth.RoleUser},
	}

	u, err := Create(ctx, us.Db, nu, time.Now())
	if err != nil {
		return err
	}

	return route.Respond(ctx, w, u, http.StatusCreated)
}
