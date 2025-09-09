package main

import (
	"database/sql"
	"log"

	"github.com/JeongWoo-Seo/simpleBank/api"
	db "github.com/JeongWoo-Seo/simpleBank/db/sqlc"
	"github.com/JeongWoo-Seo/simpleBank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}

	con, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("❌ cannot connect to db: %v", err)
	}

	store := db.NewStore(con)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("can not create server", err)
	}

	err = server.StartServer(config.ServerAddress)
	if err != nil {
		log.Fatal("fail start server", err)
	}
}
