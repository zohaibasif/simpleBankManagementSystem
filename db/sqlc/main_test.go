package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbdriver = "postgres"
	dbSource = "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable"
)

var queries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDb, err = sql.Open(dbdriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db :: error:", err)
	}

	queries = New(testDb)

	os.Exit(m.Run())
}
