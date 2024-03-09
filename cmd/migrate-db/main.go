package main

import (
	"log"
	"os"
	"planner/database"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Missing Env Variable: DATABASE_URL")
	}

	// Run migration script embedded in binary
	if err := database.UpgradeDatabase(); err != nil {
		log.Println("Failed to upgrade database:", err)
	}
}
