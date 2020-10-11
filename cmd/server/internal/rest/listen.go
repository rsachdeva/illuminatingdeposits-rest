package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rsachdeva/illuminatingdeposits/cmd/server/internal/interestsvc"
)

func ListenAndServeWithShutdown(server *http.Server, log *log.Logger, shutdownCh chan os.Signal, cfg AppConfig) error {
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrorsCh := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main : Register listening on %s", server.Addr)
		serverErrorsCh <- server.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdownCh.
	err := quitGracefully(server, log, cfg, shutdownCh, serverErrorsCh)
	if err != nil {
		return err
	}
	return nil
}

func RegisterInterestService(server *http.Server, log *log.Logger, db *sqlx.DB, shutdownCh chan os.Signal) {
	server.Handler = interestsvc.Register(shutdownCh, db, log)
}

func quitGracefully(server *http.Server, log *log.Logger, cfg AppConfig, shutdownCh chan os.Signal, serverErrorsCh chan error) error {
	select {
	case err := <-serverErrorsCh:
		return errors.Wrap(err, "starting server")

	case sig := <-shutdownCh:
		// syscall.Signal:suspended (signal) : Start shutdown -- integrity error
		// vs
		// syscall.Signal:interrupt : Start shutdown -- control c
		log.Printf("main : %T:%+v : Start shutdown", sig, sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		fmt.Println("main : As part of quitGracefully, Now attempting Graceful shutdown")
		log.Println("main : As part of quitGracefully, attempting Graceful shutdown")
		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			err = server.Close()
		}

		// Log the status of this shutdown.
		fmt.Println("sig == syscall.SIGSTOP is", sig == syscall.SIGSTOP)
		fmt.Println("sig == syscall.SIGINT is", sig == syscall.SIGINT)

		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("Integrity issue caused self shutdown")
		case sig == syscall.SIGINT:
			return errors.New("Graceful shutsdown as requested")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}
	return nil
}
