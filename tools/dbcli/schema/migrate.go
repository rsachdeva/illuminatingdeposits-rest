package schema

import (
	"log"

	"github.com/GuiaBolso/darwin"
	"github.com/jmoiron/sqlx"
)

// migrations contains the queries needed to construct the postgresconn schema.
// Entries should never be removed from this slice once they have been ran in
// production.
//
// Including the queries directly in this file has the same pros/cons mentioned
// in seeds.go

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add users",
		Script: `
CREATE TABLE users (
	uuid          UUID,
	name          TEXT,
	email         TEXT UNIQUE,
	roles         TEXT[],
	password_hash TEXT,

	date_created TIMESTAMP,
	date_updated TIMESTAMP,

	PRIMARY KEY (uuid)
);`,
	},
}

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {
	log.Println("Going to migrate database..")
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)
	err := d.Migrate()
	log.Println("migrate err", err)
	if err != nil {
		return err
	}
	log.Println("Migration completed!")
	return nil
}
