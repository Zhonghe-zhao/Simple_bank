package main

import (
	"context"
	"log"
	"net"
	"project/simplebank/gapi"
	"project/simplebank/pb"
	"project/simplebank/util"

	api "project/simplebank/api"
	db "project/simplebank/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	//RunGinServer(config, store)
	runGrpcServer(config, store)

}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create gRPCserver %v", err)
	}
	//创建gRPC服务器实例
	grpcServer := grpc.NewServer()
	//注册服务
	pb.RegisterSimplebankServer(grpcServer, server)
	//注册反射服务
	reflection.Register(grpcServer)
	//创建监听器
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create Listener")
	}
	//启动gRPC服务器
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}
}

func RunGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create ginserver: %v", err)
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}
