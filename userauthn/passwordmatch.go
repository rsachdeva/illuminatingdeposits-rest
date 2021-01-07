package userauthn

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"golang.org/x/crypto/bcrypt"
)

// Mock bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

type PasswordVerifier struct{}

func (cp PasswordVerifier) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

type PasswordVeriferInterface interface {
	CompareHashAndPassword(hashedPassword, password []byte) error
}

func PasswordMatch(hashedPassword []byte, password string) error {
	err := PasswordVerifier{}.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return muxhttp.NewRequestError(
			errors.New("NotFound Error: User email/password combination not found"),
			http.StatusNotFound)
	}
	return nil
}
