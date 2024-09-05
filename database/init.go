package database

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/surrealdb/surrealdb.go"
)

var DB *surrealdb.DB

type dbAuth struct {
	Namespace string `json:"NS,omitempty"`
	Database  string `json:"DB,omitempty"`
	Scope     string `json:"SC,omitempty"`
	Username  string `json:"user,omitempty"`
	Password  string `json:"pass,omitempty"`
}

func init() {
	dbAddr := os.Getenv("DB_ADDR")
	if dbAddr == "" {
		log.Fatal().Msg("DB_ADDR not set")
		return
	}

	ns := os.Getenv("DB_NS")
	if ns == "" {
		ns = "main"
	}

	// Initalize a connetion to the DB
	var err error
	DB, err = surrealdb.New(dbAddr)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
		return
	}

	_, err = DB.Signin(&dbAuth{
		Database:  "auth",
		Namespace: ns,
		Username:  os.Getenv("DB_USER"),
		Password:  os.Getenv("DB_PASS"),
	})

	if err != nil {
		log.Fatal().Err(err).Msg("DB authentication failed")
		return
	}
}
