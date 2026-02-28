package main

import (
	"database/sql"
	"log"

	"github.com/Hans-zi/simple_bank/api"
	db "github.com/Hans-zi/simple_bank/db/sqlc"
	"github.com/Hans-zi/simple_bank/util"
	_ "github.com/lib/pq"
)

var store db.Store

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("unable to load config, %v", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Can not connect to db:", err)
	}

	store = db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		return
	}
}
