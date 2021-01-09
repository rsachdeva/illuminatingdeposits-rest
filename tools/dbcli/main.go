package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/postgresconn"
	"github.com/rsachdeva/illuminatingdeposits-rest/tools/dbcli/schema"
	"github.com/rsachdeva/illuminatingdeposits-rest/usermgmt/uservalue"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := run(); err != nil {
		log.Printf("error: quitting appserver: %+v", err)
		os.Exit(1)
	}

}

type AppConfig struct {
	DB struct {
		User       string `conf:"default:postgres"`
		Password   string `conf:"default:postgres,noprint"`
		Host       string `conf:"default:db"`
		Name       string `conf:"default:postgres"`
		DisableTLS bool   `conf:"default:true"`
	}
	Args conf.Args
}

func ParsedConfig(cfg AppConfig) (AppConfig, error) {
	if err := conf.Parse(os.Args[1:], "DEPOSITS", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("DEPOSITS", &cfg)
			if err != nil {
				return AppConfig{}, errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return AppConfig{}, nil
		}
		return AppConfig{}, errors.Wrap(err, "parsing config")
	}
	log.Printf("AppConfig is %v", cfg)
	return cfg, nil
}

func run() error {
	var err error

	cfg, err := ParsedConfig(AppConfig{})
	if err != nil {
		return err
	}
	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main : Config :\n%v\n", out)

	log.Printf("cfg is %v", cfg)

	fmt.Println("cli cfg.DB.Host is", cfg.DB.Host)
	fmt.Println("cfg.Args is", cfg.Args)
	// This is used for multiple commands below.
	dbConfig := postgresconn.Config{
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

func createAdmin(cfg postgresconn.Config, email, password string) error {
	db, err := postgresconn.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if email == "" || password == "" {
		return errors.New("createAdmin command must be called with two additional arguments for email and password")
	}

	fmt.Printf("Admin usermgmt will be created with email %q and password %q\n", email, password)
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

	nu := uservalue.NewUser{
		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Roles:           []string{"Admin", "User"},
	}

	u, err := uservalue.AddUser(ctx, db, nu, time.Now(), bcrypt.GenerateFromPassword)
	if err != nil {
		return err
	}

	fmt.Println("User created with id:", u.Uuid)
	return nil
}

func migrate(cfg postgresconn.Config) error {
	db, err := postgresconn.Open(cfg)
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

func seed(cfg postgresconn.Config) error {
	db, err := postgresconn.Open(cfg)
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
