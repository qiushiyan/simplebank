package tests

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/lib/pq"
	. "github.com/qiushiyan/simplebank/business/db/generated"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:postgres@localhost:5433/bank?sslmode=disable"
)

var testQueries Querier
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDB)
	m.Run()
}