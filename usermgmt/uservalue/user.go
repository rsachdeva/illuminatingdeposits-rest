// Package interestvalue provides struct values and associated operations for user handling including persistence in postgres
package uservalue

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User represents someone with access to our system.
type User struct {
	Uuid         string         `db:"uuid" json:"uuid"`
	Name         string         `db:"name" json:"name"`
	Email        string         `db:"email" json:"email"`
	Roles        pq.StringArray `db:"roles" json:"roles"`
	PasswordHash []byte         `db:"password_hash" json:"-"`
	DateCreated  time.Time      `db:"date_created" json:"date_created"`
	DateUpdated  time.Time      `db:"date_updated" json:"date_updated"`
}

// NewUser contains information needed to create a new User.
type NewUser struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required"`
	Roles           []string `json:"roles" validate:"required"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"password_confirm" validate:"eqfield=Password"`
}

func AddUser(ctx context.Context, db *sqlx.DB, n NewUser, now time.Time) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(n.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generating password hash")
	}

	u := User{
		Uuid:         uuid.New().String(),
		Name:         n.Name,
		Email:        n.Email,
		PasswordHash: hash,
		Roles:        n.Roles,
		DateCreated:  now.UTC(),
		DateUpdated:  now.UTC(),
	}

	const q = `INSERT INTO users
		(uuid, name, email, password_hash, roles, date_created, date_updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.ExecContext(ctx, q, u.Uuid, u.Name, u.Email, u.PasswordHash, u.Roles, u.DateCreated, u.DateUpdated)
	if err != nil {
		fmt.Printf("\ndb err is %T %v", err, err)
		return nil, errors.Wrap(err, "inserting usermgmt")
	}

	return &u, nil
}

func FindByEmail(ctx context.Context, db *sqlx.DB, email string) (User, error) {
	usr := User{}

	err := db.GetContext(ctx, &usr, db.Rebind("SELECT u.uuid, u.password_hash, u.roles, u.Email FROM users u WHERE u.email=?"), email)
	fmt.Printf("\nusr just after Postgres query %#v", usr)
	if err != nil {
		return User{}, err
	}
	if usr.Email != email {
		return User{}, errors.Wrap(err, "error finding user by email")
	}

	return usr, nil
}
