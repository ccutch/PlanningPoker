package main

import (
	"flag"
	"log"
	"net/http"
	"planner"
)

var (
	addr  = flag.String("http", "0.0.0.0:8000", "http address for listening")
	dbURL = flag.String("db-url", "", "database url for data storage")
)

func main() {
	flag.Parse()
	log.Printf("Serving Baleen @ http://%s\n", *addr)

	// First lets establish a connection pool to db
	if err := planner.GoOnline(*dbURL); err != nil {
		// And fail early if we encounter an error.
		log.Fatal("failed to open db:", err)
	}

	if err := http.ListenAndServe(*addr, planner.Routes); err != nil {
		log.Fatal(err)
	}
}
