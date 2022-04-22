package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Error loading initial config: ", err)
	}
	sqlDB, err := sql.Open(config.DBDriver, config.DBSource)
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

	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("Can't start server: ", err)
	}
}
