package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"planner"
	"planner/database"
)

var (
	port  = os.Getenv("PORT")
	dbURL = os.Getenv("DATABASE_URL")

	migrateDB = flag.Bool("migrate-db", false, "flag will use golang/migrate to migrate db schema")
)

func main() {
	flag.Parse()
	if port == "" {
		port = "8000"
	}
	if *migrateDB && dbURL != "" {
		database.UpgradeDatabase()
	}
	http.Handle("/", planner.Routes)
	log.Printf("Serving Baleen @ http://0.0.0.0:%s\n", port)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil); err != nil {
		log.Fatal(err)
	}
}
