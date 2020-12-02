package rest

import (
	"crypto/tls"
	_ "expvar" // Register the expvar interestsvc
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // Register the pprof interestsvc
	"os"

	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/conf"
	"github.com/rsachdeva/illuminatingdeposits/database"
	"github.com/rsachdeva/illuminatingdeposits/invest"
	"github.com/rsachdeva/illuminatingdeposits/rest/middleware"
	"github.com/rsachdeva/illuminatingdeposits/service"
	"github.com/rsachdeva/illuminatingdeposits/user"
)

func tlsConfig() (*tls.Config, error) {
	certFile := "conf/tls/servercrtto.pem"
	keyFile := "conf/tls/serverkeyto.pem"
	// _, err := ioutil.ReadFile(certFile)
	// if err != nil {
	// 	log.Fatalf("Error in reading cert file %v", certFile)
	// }
	// _, err = ioutil.ReadFile(keyFile)
	// if err != nil {
	// 	log.Fatalf("Error in reading key file %v", keyFile)
	// }
	// fmt.Println("Ok to load cert and key files")
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		// https://github.com/stellar/go/issues/64 shutting down, error: LoadX509KeyPair error: tls: failed to parse private key
		// From https://github.com/pulumi/pulumi-kafka/issues/15
		// // that redirected to https://golang.org/pkg/crypto/tls/#LoadX509KeyPair
		// LoadX509KeyPair reads and parses a public/private key pair from a pair of files. The files must contain PEM encoded data.
		// so we need PEM files jmd
		// Next error private key does not match public key
		// Based on https://stackoverflow.com/questions/991758/how-to-get-pem-file-from-key-and-crt-files JMD
		// 2020/11/03 16:16:51 shutting down, error: LoadX509KeyPair error: tls: failed to find certificate PEM data in certificate input, but did find a private key; PEM inputs may have been switched
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

	// =========================================================================
	// ReqHandler Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main : Config :\n%v\n", out)

	// =========================================================================
	// Start Database

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
	// /debug/pprof - Added to the default mux by importing the net/http/pprof package.
	// /debug/vars - Added to the default mux by importing the expvar package.
	//
	// Not concerned with shutting this down when the application is shutdownCh.
	go func() {
		Debug(log, cfg)
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
	server := NewServer(cfg, tl)
	h := service.NewReqHandler(shutdownCh, log, middleware.Logger(log), middleware.Errors(log), middleware.Metrics(), middleware.Panics(log))
	server.Handler = h
	database.RegisterCheckService(db, h)
	user.RegisterUserService(db, h)
	invest.RegisterInvestService(log, h)

	err = ListenAndServeWithShutdown(server, log, shutdownCh, cfg)
	if err != nil {
		return err
	}

	return nil
}
