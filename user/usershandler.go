package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/auth"
	"github.com/rsachdeva/illuminatingdeposits/json"
	"github.com/rsachdeva/illuminatingdeposits/web"
	"go.opencensus.io/trace"
)

// Users holds interestsvc for dealing with user.
type UsersHandler struct {
	db *sqlx.DB
}

func RegisterUserHandler(db *sqlx.DB, app *web.App) {
	{
		// Register user interestsvc.
		u := UsersHandler{db: db}

		// The route can't be authenticated because we need this route to
		// create the user in the first place.
		app.Handle(http.MethodPost, "/v1/users", u.Create)
	}
}

func (us *UsersHandler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "interestsvc.Users.Create")
	defer span.End()

	email, pass, ok := r.BasicAuth()
	fmt.Printf("r.Header.GetGet(\"Authorization\") is %s", r.Header.Get("Authorization"))
	if !ok {
		err := errors.New("must provide email and password in Basic auth")
		return web.NewRequestError(err, http.StatusUnauthorized)
	}

	nu := NewUser{
		Email:           email,
		Password:        pass,
		PasswordConfirm: pass,
		Roles:           []string{auth.RoleUser},
	}

	u, err := Create(ctx, us.db, nu, time.Now())
	if err != nil {
		return err
	}

	return json.Respond(ctx, w, u, http.StatusCreated)
}
