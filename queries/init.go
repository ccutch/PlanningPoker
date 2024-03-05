package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/pkg/errors"
)

var (
	dbURL = os.Getenv("DATABASE_URL")
	db    *sql.DB
)

func init() {
	var err error
	if db, err = sql.Open("postgres", dbURL); err != nil {
		err = errors.Wrap(err, "failed to connect to db: "+dbURL)
		log.Fatal(err)
	}

	m, err := migrate.New("github://ccutch/PlanningPoker/migrations", dbURL)
	if err == nil {
		m.Up()
	}
}

func genID(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Fatal("failed to generate random secret")
	}
	return base64.URLEncoding.EncodeToString(b)[:n]
}
