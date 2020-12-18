// Package appserver provides appserver configuration for db, tracing, tls and env variables.
// It also provides regsitration for services including starting the server
package main

import (
	"crypto/tls"
	_ "expvar" // Register the expvar interestsvc
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // Register the pprof interestsvc
	"os"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/conf"
	"github.com/rsachdeva/illuminatingdeposits-rest/errconv"
	"github.com/rsachdeva/illuminatingdeposits-rest/interestcal"
	"github.com/rsachdeva/illuminatingdeposits-rest/metriccnt"
	"github.com/rsachdeva/illuminatingdeposits-rest/muxhttp"
	"github.com/rsachdeva/illuminatingdeposits-rest/postgreshealth"
	"github.com/rsachdeva/illuminatingdeposits-rest/readenv"
	"github.com/rsachdeva/illuminatingdeposits-rest/recoverpanic"
	"github.com/rsachdeva/illuminatingdeposits-rest/reqlog"
	"github.com/rsachdeva/illuminatingdeposits-rest/usermgmt"
)

func NewServer(cfg AppConfig, tlc *tls.Config) *http.Server {
	server := http.Server{
		Addr:         cfg.Web.Address,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}
	if tlc != nil {
		server.TLSConfig = tlc
	}
	fmt.Println("server.TLSConfig is ", server.TLSConfig)
	fmt.Println("DEPOSITS_DB_HOST is ", cfg.DB.Host)
	fmt.Println("DEPOSITS_TRACE_URL is", cfg.Trace.URL)
	return &server
}

func ConfigureAndServe() error {
	cfg, err := ParsedConfig(AppConfig{})
	if err != nil {
		return err
	}

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	log := log.New(os.Stdout, "DEPOSITS : ", log.LstdFlags|log.Lmicroseconds|log.Llongfile)

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main : Config :\n%v\n", out)

	db, err := Db(cfg)
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer db.Close()

	// =========================================================================
	// Start Tracing Support

	closer, err := RegisterTracer(
		cfg.Trace.Service,
		cfg.Web.Address,
		cfg.Trace.URL,
		cfg.Trace.Probability,
	)
	if err != nil {
		return err
	}
	defer func() {
		err := closer()
		if err != nil {
			log.Println("could not close reporter", err)
		}
	}()

	// =========================================================================
	// Start Debug service
	//
	// /debug/pprof - Added to the default jsonfmt by importing the net/http/pprof package.
	// /debug/vars - Added to the default jsonfmt by importing the expvar package.
	//
	// Not concerned with shutting this down when the application is shutdownCh.
	go func() {
		Debug(log, cfg)
	}()

	tlsEnabled := readenv.TlsEnabled()
	fmt.Println("tls is", tlsEnabled)
	var tlc *tls.Config
	if tlsEnabled {
		tlc, err = tlsConfig()
		if err != nil {
			return err
		}
	}
	s := NewServer(cfg, tlc)
	// =========================================================================
	// Register API Services and Start Server

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	// https://golang.org/pkg/os/signal/#Notify
	shutdownCh := make(chan os.Signal, 1)

	m := appmux.NewRouter(shutdownCh, log,
		reqlog.NewMiddleware(log),
		errconv.NewMiddleware(log),
		metriccnt.NewMiddleware(),
		recoverpanic.NewMiddleware(log))
	s.Handler = m
	postgreshealth.RegisterSvc(db, m)
	usermgmt.RegisterSvc(db, m)
	interestcal.RegisterSvc(log, m)

	err = ListenAndServeWithShutdown(s, log, shutdownCh, cfg)
	if err != nil {
		return err
	}

	return nil
}
