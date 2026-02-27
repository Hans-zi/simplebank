package main

import (
	"database/sql"
	"log"

	"github.com/Hans-zi/simple_bank/api"
	db "github.com/Hans-zi/simple_bank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:xiaohan1234@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = ":8080"
)

var store *db.Store

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Can not connect to db:", err)
	}

	store = db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		return
	}
}
