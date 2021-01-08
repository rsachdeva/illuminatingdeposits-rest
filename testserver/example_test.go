package testserver_test

import (
	"fmt"

	"github.com/rsachdeva/illuminatingdeposits-rest/testserver"
	"github.com/rsachdeva/illuminatingdeposits-rest/tools/dbcli/schema"
)

func ExampleStartPostgresAndMigrate() {
	db, pool, resource := testserver.StartPostgresAndMigrate()
	// this is just for this test to make sure can populate db; we don't require this for tests in general
	err := schema.Seed(db)
	if err != nil {
		fmt.Println("Cound not seed test database")
	}
	err = pool.Purge(resource)
	fmt.Println("ExamplePostgresConnect err is", err)
	// Output:
	// ExamplePostgresConnect err is <nil>
}
