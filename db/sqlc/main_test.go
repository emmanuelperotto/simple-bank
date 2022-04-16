package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret123@localhost:5432/simple_bank?sslmode=disable"
)

var (
	testQueries *Queries
	testDb      *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	testDb, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
