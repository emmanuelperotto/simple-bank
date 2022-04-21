package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret123@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	sqlDB, err := sql.Open(dbDriver, dbSource)
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Fatal("Error closing DB: ", err)
		}
	}()
	if err != nil {
		log.Fatal("Cannot open connection to DB: ", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("DB connection not alive: ", err)
	}

	store := db.NewStore(sqlDB)
	server := api.NewServer(store)

	if err := server.Start(serverAddress); err != nil {
		log.Fatal("Can't start server: ", err)
	}
}
