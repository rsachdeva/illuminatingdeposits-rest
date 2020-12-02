package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
)

func ListenAndServeWithShutdown(server *http.Server, log *log.Logger, shutdownCh chan os.Signal, cfg AppConfig) error {
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrorsCh := make(chan error, 1)

	// Start the route listening for requests.
	go func() {
		log.Printf("main : Register listening on %s", server.Addr)
		// send signal to serverErrorCh
		if cfg.Web.ServiceServerTLS {
			serverErrorsCh <- server.ListenAndServeTLS("", "")
			return
		}
		serverErrorsCh <- server.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdownCh.
	// send signal to shutdownCh
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
	err := quitGracefully(server, log, cfg, shutdownCh, serverErrorsCh)
	if err != nil {
		return err
	}
	return nil
}

func quitGracefully(server *http.Server, log *log.Logger, cfg AppConfig, shutdownCh chan os.Signal, serverErrorsCh chan error) error {
	// receive singal for serverErrorCh and shutdownCh
	select {
	case err := <-serverErrorsCh:
		return errors.Wrap(err, "starting server")

	case sig := <-shutdownCh:
		// syscall.Signal:suspended (signal) : Start shutdown -- integrity error
		// vs
		// syscall.Signal:interrupt : Start shutdown -- control c
		log.Printf("main : %T:%v : Start shutdown", sig, sig)

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
			return errors.New("Graceful shutdown as requested")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}
	return nil
}
