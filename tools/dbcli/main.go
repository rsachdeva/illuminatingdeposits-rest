package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/platform/auth"
	"github.com/rsachdeva/illuminatingdeposits/platform/conf"
	"github.com/rsachdeva/illuminatingdeposits/platform/database"
	"github.com/rsachdeva/illuminatingdeposits/platform/schema"
	"github.com/rsachdeva/illuminatingdeposits/user"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: quitting app: %+v", err)
		os.Exit(1)
	}

}

func run() error {
	var err error

	//DEPOSITS_DB_DISABLE_TLS=true for local testing with db
	//DEPOSITS_DB_HOST=192.168.254.33
	var cfg struct {
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:db"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:false"`
		}
		Args conf.Args
	}

	if err := conf.Parse(os.Args[1:], "DEPOSITS", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("DEPOSITS", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "error: parsing config")
	}

	// This is used for multiple commands below.
	dbConfig := database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	}

	switch cfg.Args.Num(0) {
	case "createAdmin":
		err = createAdmin(dbConfig, cfg.Args.Num(1), cfg.Args.Num(2))
	case "migrate":
		err = migrate(dbConfig)
	case "seed":
		err = seed(dbConfig)
	default:
		err = errors.New("Must specify a command with correct option createAdmin, createInterest, migrate and seed")
	}

	if err != nil {
		return err
	}

	return nil
}

func createAdmin(cfg database.Config, email, password string) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if email == "" || password == "" {
		return errors.New("createAdmin command must be called with two additional arguments for email and password")
	}

	fmt.Printf("Admin user will be created with email %q and password %q\n", email, password)
	fmt.Print("Continue? (1/0) ")

	var confirm bool
	if _, err := fmt.Scanf("%t\n", &confirm); err != nil {
		return errors.Wrap(err, "processing response")
	}

	if !confirm {
		fmt.Println("Canceling")
		return nil
	}

	ctx := context.Background()

	nu := user.NewUser{
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Roles:           []string{auth.RoleAdmin, auth.RoleUser},
	}

	u, err := user.Create(ctx, db, nu, time.Now())
	if err != nil {
		return err
	}

	fmt.Println("User created with id:", u.ID)
	return nil
}

func migrate(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := schema.Migrate(db); err != nil {
		return err
	}

	fmt.Println("Migrations complete")
	return nil
}

func seed(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := schema.Seed(db); err != nil {
		return err
	}

	fmt.Println("Seed data complete")
	return nil
}
