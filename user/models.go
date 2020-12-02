package user

import (
	"time"

	"github.com/lib/pq"
)

// User represents someone with access to our system.
type User struct {
	ID           string         `Db:"user_id" json:"id"`
	Name         string         `Db:"name" json:"name"`
	Email        string         `Db:"email" json:"email"`
	Roles        pq.StringArray `Db:"roles" json:"roles"`
	PasswordHash []byte         `Db:"password_hash" json:"-"`
	DateCreated  time.Time      `Db:"date_created" json:"date_created"`
	DateUpdated  time.Time      `Db:"date_updated" json:"date_updated"`
}

// NewUser contains information needed to create a new User.
type NewUser struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required"`
	Roles           []string `json:"roles" validate:"required"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"password_confirm" validate:"eqfield=Password"`
}
