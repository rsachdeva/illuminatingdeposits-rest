package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/internal/platform/auth"
	"github.com/rsachdeva/illuminatingdeposits/internal/platform/web"
	"github.com/rsachdeva/illuminatingdeposits/internal/user"
	"go.opencensus.io/trace"
)

// Users holds handlers for dealing with user.
type Users struct {
	db *sqlx.DB
}

func (us *Users) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Users.Create")
	defer span.End()

	email, pass, ok := r.BasicAuth()
	fmt.Printf("r.Header.GetGet(\"Authorization\") is %s", r.Header.Get("Authorization"))
	if !ok {
		err := errors.New("must provide email and password in Basic auth")
		return web.NewRequestError(err, http.StatusUnauthorized)
	}

	nu := user.NewUser{
		Email:           email,
		Password:        pass,
		PasswordConfirm: pass,
		Roles:           []string{auth.RoleUser},
	}

	u, err := user.Create(ctx, us.db, nu, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, u, http.StatusCreated)
}
