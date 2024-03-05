package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"planner"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

func main() {

	port, dbURL := os.Getenv("PORT"), os.Getenv("DATABASE_URL")
	if port == "" {
		port = "8000"
	}
	host := fmt.Sprintf("0.0.0.0:%s", port)

	m, err := migrate.New("github://ccutch/PlanningPoker/migrations", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	if err = m.Up(); err != nil {
		log.Fatal(err)
	}

	// First let's establish a connection pool to db
	if err := planner.GoOnline(dbURL); err != nil {
		// And fail early if we encounter an error.
		log.Fatal("failed to open db:", err)
	}

	// Second let's listen for incoming http requests.
	log.Printf("Serving Baleen @ http://%s\n", host)
	if err := http.ListenAndServe(host, planner.Routes); err != nil {
		log.Fatal(err)
	}
}
