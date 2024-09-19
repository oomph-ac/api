package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/getsentry/sentry-go"
	_ "github.com/oomph-ac/api/endpoint"
	"github.com/rs/zerolog/log"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal().Msg("Usage: ./binary <host_addr>")
		return
	}

	// Initalize sentry if enabled.
	if dsn := os.Getenv("SENTRY_DSN"); dsn != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn: dsn,
		}); err != nil {
			panic(err)
		}
	}

	go func() {
		hostAddr := os.Args[1]
		log.Info().Msg("Oomph API serving on " + hostAddr)

		if err := http.ListenAndServe(hostAddr, http.DefaultServeMux); err != nil {
			panic(err)
		}
	}()

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)
	<-shutdownChan

	log.Info().Msg("Oomph REST API shuts down with grace :)")
}
