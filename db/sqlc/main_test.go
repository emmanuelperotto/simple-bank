package db

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"os"
	"simplebank/util"
	"testing"
)

var (
	testQueries *Queries
	testDb      *sql.DB
)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}
	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}

	if err := testDb.Ping(); err != nil {
		log.Fatal("DB connection not alive: ", err)
	}

	driver, err := postgres.WithInstance(testDb, &postgres.Config{})
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting working directory: ", err)
	}

	migr, err := migrate.NewWithDatabaseInstance(
		"file://"+dir+"/migrations",
		"postgres", driver)

	if err != nil {
		log.Fatal("Error starting migrations process: ", err)
	}

	if err := migr.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("Error applying UP migrations: ", err)
	}

	testQueries = New(testDb)
	exitCode := m.Run()

	if err := testDb.Close(); err != nil {
		log.Fatal("Error closing db: ", err)
	}

	os.Exit(exitCode)
}
