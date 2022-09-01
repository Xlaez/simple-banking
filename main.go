package main

import (
	"database/sql"
	"log"
	"simple-bank/api"
	db "simple-bank/db/sqlc"
	"simple-bank/util"

	_ "github.com/lib/pq"
)


func main(){
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load env", err)
	}
	conn, err := sql.Open(config.DbDriver, config.DbSource)

	if err != nil {
		log.Fatal("go cannot connect to db ..... exiting:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store) 

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("go cannot start server ...... exiting", err)
	}
}