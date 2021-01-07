package userauthn_test

import (
	"testing"

	"github.com/rsachdeva/illuminatingdeposits-rest/userauthn"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// https://github.com/stretchr/testify#mock-package
// // Mock bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
type MockedPasswordVerifier struct {
	mock.Mock
	hashedPassword []byte
	password       string
}

func (mpv *MockedPasswordVerifier) CompareHashAndPassword() error {
	args := mpv.Called()
	return args.Error(0)
}

func TestPasswordNotMatching(t *testing.T) {
	// create an instance of our test object
	mpv := MockedPasswordVerifier{hashedPassword: []byte("hello"), password: "bye"}

	// setup expectations
	mpv.On("CompareHashAndPassword").Return(bcrypt.ErrMismatchedHashAndPassword)

	// call the code we are testing
	err := userauthn.PasswordMatch(&mpv)
	require.NotNil(t, err, "should be an error for password mismatch")

}

func TestPasswordMatches(t *testing.T) {
	// create an instance of our test object
	mpv := MockedPasswordVerifier{hashedPassword: []byte("hello"), password: "hello"}

	// setup expectations
	mpv.On("CompareHashAndPassword").Return(nil)

	// call the code we are testing
	err := userauthn.PasswordMatch(&mpv)
	require.Nil(t, err)
}
