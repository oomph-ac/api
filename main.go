package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/oomph-ac/api/endpoint"
	"github.com/rs/zerolog/log"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal().Msg("Usage: ./binary <host_addr>")
		return
	}

	go func() {
		hostAddr := os.Args[1]
		log.Info().Msg("Oomph REST API serving on " + hostAddr)

		http.HandleFunc(endpoint.PathAuthentication, endpoint.Authenticate)
		if err := http.ListenAndServe(hostAddr, http.DefaultServeMux); err != nil {
			panic(err)
		}
	}()

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)
	<-shutdownChan

	log.Info().Msg("Oomph REST API shuts down with grace :)")
}
