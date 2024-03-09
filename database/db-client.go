package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	dbURL = os.Getenv("DATABASE_URL")
	db    *sql.DB
)

func init() {
	var err error
	if db, err = sql.Open("postgres", dbURL); err != nil {
		log.Fatal(errors.Wrap(err, "failed to connect to db: "+dbURL))
	}
}

func genID(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		log.Fatal("failed to generate random secret")
	}
	nonce := base64.URLEncoding.EncodeToString(b)[:n-1]
	return "x" + strings.ReplaceAll(nonce, "-", "_")
}
