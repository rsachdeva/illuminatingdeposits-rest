package userauthn

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func passwordMatch(hashedPassword []byte, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.Wrap(err, "NotFound Error: User email/password combination not found")
	}
	return nil
}
