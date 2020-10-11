package rest

import (
	"github.com/jmoiron/sqlx"
	"github.com/rsachdeva/illuminatingdeposits/internal/platform/database"
)

func Db(cfg AppConfig) (*sqlx.DB, error) {
	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
