package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	views "planner/templates"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

var port = os.Getenv("PORT")

func init() {
	if port == "" {
		port = "8000"
	}
}

func main() {
	host := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Serving Baleen @ http://%s", host)
	if err := http.ListenAndServe(host, views.Routes); err != nil {
		log.Fatal(err)
	}
}
