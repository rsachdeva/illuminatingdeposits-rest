package main

import (
	_ "expvar" // Register the expvar interestsvc
	"log"
	"net/http"
	_ "net/http/pprof" // Register the pprof interestsvc
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/cmd/server/internal/rest"
	"github.com/rsachdeva/illuminatingdeposits/internal/platform/conf"
)

func RunServerWithRegisteredService() error {

	// =========================================================================
	// Logging

	log := log.New(os.Stdout, "DEPOSITS : ", log.LstdFlags|log.Lmicroseconds|log.Llongfile)

	// =========================================================================
	// Configuration

	cfg, err := rest.ParsedConfig(rest.AppConfig{})
	if err != nil {
		return err
	}

	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main : Config :\n%v\n", out)

	// =========================================================================
	// Start Database

	db, err := rest.Db(cfg)
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer db.Close()

	// =========================================================================
	// Start Tracing Support

	closer, err := rest.RegisterTracer(
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
	// /debug/pprof - Added to the default mux by importing the net/http/pprof package.
	// /debug/vars - Added to the default mux by importing the expvar package.
	//
	// Not concerned with shutting this down when the application is shutdownCh.
	go func() {
		rest.Debug(log, cfg)
	}()

	// fmt.Println("hi there")
	// lis, err := net.Listen("tcp", "0.0.0.0:50051")
	// if err != nil {
	// 	log.Fatalf("could not listen %v", err)
	// }
	//
	// // since execution happens from root of project per the go.mod file
	// tls := true
	// var opts []grpc.ServerOption
	// if tls {
	// 	opts = tlsOpts(opts)
	// }
	// // https://golang.org/ref/spec#Passing_arguments_to_..._parameters
	// s := grpc.NewServer(opts...)
	// // s := grpc.NewServer()
	// greetpb.RegisterGreetServiceServer(s, server{})
	//
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("error is %#v", err)
	// }

	// =========================================================================
	// Start API Service

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	server := NewServer(cfg)
	rest.RegisterInterestService(server, log, db, shutdownCh)

	err = rest.ListenAndServeWithShutdown(server, log, shutdownCh, cfg)
	if err != nil {
		return err
	}

	return nil
}

func NewServer(cfg rest.AppConfig) *http.Server {
	server := http.Server{
		Addr:         cfg.Web.Address,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}
	return &server
}
