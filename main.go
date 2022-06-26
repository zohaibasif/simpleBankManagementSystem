package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/zohaibAsif/simple_bank_management_system/api"
	db "github.com/zohaibAsif/simple_bank_management_system/db/sqlc"
	"github.com/zohaibAsif/simple_bank_management_system/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config file :: error:", err)
	}

	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("cannot connect to db :: error:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("Cannot start server :: error:", err)
	}

}
