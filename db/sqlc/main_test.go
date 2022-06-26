package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/zohaibAsif/simple_bank_management_system/util"
)

var queries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config file :: error:", err)
	}

	testDb, err = sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to db :: error:", err)
	}

	queries = New(testDb)

	os.Exit(m.Run())
}
