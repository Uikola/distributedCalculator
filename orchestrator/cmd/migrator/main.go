package main

import (
	"flag"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var dbURL string

	flag.StringVar(&dbURL, "db-url", "", "url of the database")
	flag.Parse()

	if dbURL == "" {
		panic("db-url is required")
	}

	m, err := migrate.New(
		"file://migrations/",
		dbURL,
	)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		panic(err)
	}
}
