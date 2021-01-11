package uservalue_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/testserver"
	"github.com/rsachdeva/illuminatingdeposits-rest/usermgmt/uservalue"
	"github.com/stretchr/testify/require"
)

func TestAddUser(t *testing.T) {
	t.Parallel()

	db := testserver.PostgresConnect(t, true)
	nu := uservalue.NewUser{
		Name:            "Rohit Sachdeva",
		Email:           "growth@drinnovations.us",
		Password:        "kubernetes",
		PasswordConfirm: "kubernetes",
		Roles:           []string{"Admin", "User"},
	}

	u, err := uservalue.AddUser(context.Background(), db, nu, time.Now(), uservalue.PasswordGenerator{})
	require.Nil(t, err)
	require.Equal(t, len(u.Uuid), 36)
}

func TestAddUserDBClientConnectionFailure(t *testing.T) {
	t.Parallel()

	db := testserver.PostgresConnect(t, true)
	nu := uservalue.NewUser{
		Name:            "Rohit Sachdeva",
		Email:           "growth@drinnovations.us",
		Password:        "kubernetes",
		PasswordConfirm: "kubernetes",
		Roles:           []string{"Admin", "User"},
	}
	db.Close()
	_, err := uservalue.AddUser(context.Background(), db, nu, time.Now(), uservalue.PasswordGenerator{})
	require.NotNil(t, err)
	require.Regexp(t, regexp.MustCompile("inserting user: sql: database is closed"), err)
}

type MockedPasswordGenerator struct{}

func (pf MockedPasswordGenerator) Hash(password []byte) ([]byte, error) {
	return nil, errors.New("some weird error when hashing")
}

func TestAdduserHashingPasswordfails(t *testing.T) {
	t.Parallel()

	db := testserver.PostgresConnect(t, true)
	nu := uservalue.NewUser{
		Name:            "Rohit Sachdeva",
		Email:           "growth@drinnovations.us",
		Password:        "",
		PasswordConfirm: "",
		Roles:           []string{"Admin", "User"},
	}
	_, err := uservalue.AddUser(context.Background(), db, nu, time.Now(), MockedPasswordGenerator{})
	require.NotNil(t, err)
	require.Regexp(t, regexp.MustCompile("generating password hash: some weird error when hashing"), err)
}

func TestFindByEmail(t *testing.T) {
	t.Parallel()

	db := testserver.PostgresConnect(t, true)
	nu := uservalue.NewUser{
		Name:            "Rohit Sachdeva",
		Email:           "growth@drinnovations.us",
		Password:        "kubernetes",
		PasswordConfirm: "kubernetes",
		Roles:           []string{"Admin", "User"},
	}
	_, err := uservalue.AddUser(context.Background(), db, nu, time.Now(), uservalue.PasswordGenerator{})
	require.Nil(t, err)

	u, err := uservalue.FindByEmail(context.Background(), db, "growth@drinnovations.us")
	require.Nil(t, err)
	require.Equal(t, nu.Email, u.Email)

}

func TestFindByEmailNotFound(t *testing.T) {
	t.Parallel()

	db := testserver.PostgresConnect(t, true)
	nu := uservalue.NewUser{
		Name:            "Rohit Sachdeva",
		Email:           "growth@drinnovations.us",
		Password:        "kubernetes",
		PasswordConfirm: "kubernetes",
		Roles:           []string{"Admin", "User"},
	}
	_, err := uservalue.AddUser(context.Background(), db, nu, time.Now(), uservalue.PasswordGenerator{})
	require.Nil(t, err)

	_, err = uservalue.FindByEmail(context.Background(), db, "growth@drinnova.us")
	require.NotNil(t, err)
	require.Regexp(t, regexp.MustCompile("no rows in result set"), err)
}
