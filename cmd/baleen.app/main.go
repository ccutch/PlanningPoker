package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"planner"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	host := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Serving Baleen @ http://%s\n", host)
	if err := http.ListenAndServe(host, planner.Routes); err != nil {
		log.Fatal(err)
	}
}
