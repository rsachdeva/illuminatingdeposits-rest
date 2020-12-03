package app

import (
	"crypto/tls"
	_ "expvar" // Register the expvar interestsvc
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // Register the pprof interestsvc
	"os"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits-rest/app/middlewarefunc"
	"github.com/rsachdeva/illuminatingdeposits-rest/conf"
	"github.com/rsachdeva/illuminatingdeposits-rest/responder"
)

func tlsConfig() (*tls.Config, error) {
	certFile := "conf/tls/servercrtto.pem"
	keyFile := "conf/tls/serverkeyto.pem"
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, errors.Wrap(err, "LoadX509KeyPair error")
	}
	fmt.Println("No errors with LoadX509KeyPair")
	tl := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return &tl, nil
}

func NewServer(cfg AppConfig, tl *tls.Config) *http.Server {
	server := http.Server{
		Addr:         cfg.Web.Address,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}
	if tl != nil {
		server.TLSConfig = tl
	}
	fmt.Println("DEPOSITS_WEB_SERVICE_SERVER_TLS is ", cfg.Web.ServiceServerTLS)
	fmt.Println("server.TLSConfig is ", server.TLSConfig)
	return &server
}

func ConfigureAndServe() error {

	// =========================================================================
	// Logging

	log := log.New(os.Stdout, "DEPOSITS : ", log.LstdFlags|log.Lmicroseconds|log.Llongfile)

	// =========================================================================
	// Configuration

	cfg, err := ParsedConfig(AppConfig{})
	if err != nil {
		return err
	}

	log.Printf("main : Started")
	defer log.Println("main : Completed")

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
	// Start Debug Service
	//
	// /debug/pprof - Added to the default responder by importing the net/http/pprof package.
	// /debug/vars - Added to the default responder by importing the expvar package.
	//
	// Not concerned with shutting this down when the application is shutdownCh.
	go func() {
		Debug(log, cfg)
	}()

	// =========================================================================
	// Register API Services and Start Server

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	// https://golang.org/pkg/os/signal/#Notify
	shutdownCh := make(chan os.Signal, 1)

	var tl *tls.Config
	if cfg.Web.ServiceServerTLS {
		tl, err = tlsConfig()
		if err != nil {
			return err
		}
	}
	s := NewServer(cfg, tl)
	m := responder.NewServeMux(shutdownCh, log,
		middlewarefunc.Logger(log),
		middlewarefunc.Errors(log),
		middlewarefunc.Metrics(),
		middlewarefunc.Panics(log))
	s.Handler = m
	registerDbHealthService(db, m)
	registerUserService(db, m)
	registerInvestService(log, m)

	err = ListenAndServeWithShutdown(s, log, shutdownCh, cfg)
	if err != nil {
		return err
	}

	return nil
}
