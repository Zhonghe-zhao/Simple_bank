package main

import (
	"context"
	"log"
	"project/simplebank/util"

	api "project/simplebank/api"
	db "project/simplebank/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DB_SOURCE)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	//初始化数据库服务
	store := db.NewStore(connPool)
	//运行gin框架
	RunGinServer(config, store)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}

func RunGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
