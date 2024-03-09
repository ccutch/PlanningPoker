package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"planner"
)

func main() {
	// Check for required environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Missing Environment Variable: DATABASE_URL")
	}

	// Check that port is given or default to 8000
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Serve http and log server start time
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Serving Baleen.app @ http://%s\n", addr)
	if err := http.ListenAndServe(addr, planner.Routes); err != nil {
		log.Fatal("Failed to server http:", err)
	}
}
