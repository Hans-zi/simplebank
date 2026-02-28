package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Hans-zi/simple_bank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("unable to load config", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Can not connect to db:", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
