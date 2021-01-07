package userauthn

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"golang.org/x/crypto/bcrypt"
)

type PasswordVerifier struct {
	hashedPassword []byte
	password       string
}

func (pv PasswordVerifier) CompareHashAndPassword() error {
	return bcrypt.CompareHashAndPassword(pv.hashedPassword, []byte(pv.password))
}

type PasswordVeriferInterface interface {
	CompareHashAndPassword() error
}

func PasswordMatch(pvi PasswordVeriferInterface) error {
	err := pvi.CompareHashAndPassword()
	if err != nil {
		return muxhttp.NewRequestError(
			errors.New("NotFound Error: User email/password combination not found"),
			http.StatusNotFound)
	}
	return nil
}
