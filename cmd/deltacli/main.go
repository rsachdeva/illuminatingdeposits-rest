package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/cmd/deltacli/internal/handlers"
	"github.com/rsachdeva/illuminatingdeposits/internal/invest"
	"github.com/rsachdeva/illuminatingdeposits/internal/platform/auth"
	"github.com/rsachdeva/illuminatingdeposits/internal/platform/conf"
	"github.com/rsachdeva/illuminatingdeposits/internal/platform/database"
	"github.com/rsachdeva/illuminatingdeposits/internal/platform/inout"
	"github.com/rsachdeva/illuminatingdeposits/internal/platform/schema"
	"github.com/rsachdeva/illuminatingdeposits/internal/user"
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
	case "createInterest":
		err = createInterest()
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

func createInterest() error {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not createInterest CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	log := log.New(os.Stdout, "DEPOSITS : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	hi := handlers.Interest{Log: log}
	var ni invest.NewInterestBanks

	fmt.Println("flag.Arg(1) is", flag.Arg(1))
	if err := inout.InputJSON(flag.Arg(1), &ni); err != nil {
		return errors.Wrap(err, "parsing json file for interest")
	}
	executionTimes := 1
	if *memprofile != "" {
		executionTimes = 100000
	}
	fmt.Println("executionTimes is", executionTimes)
	if err := hi.Create(os.Stdout, ni, executionTimes); err != nil {
		return errors.Wrap(err, "printing all banks calculations")
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not createInterest memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
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
