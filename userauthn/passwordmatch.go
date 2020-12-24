package userauthn

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"golang.org/x/crypto/bcrypt"
)

func passwordMatch(hashedPassword []byte, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return muxhttp.NewRequestError(
			errors.New("NotFound Error: User email/password combination not found"),
			http.StatusNotFound)
	}
	return nil
}
