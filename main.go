package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"project/simplebank/gapi"
	"project/simplebank/pb"
	"project/simplebank/util"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"

	api "project/simplebank/api"
	db "project/simplebank/db/sqlc"

	_ "project/simplebank/doc/statik"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
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
	go func() {
		runGrpGatewayServer(config, store)
	}()

	if err := runGrpcServer(config, store); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

}

func runGrpcServer(config util.Config, store db.Store) error {
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
	// err = grpcServer.Serve(listener)
	// if err != nil {
	// 	log.Fatal("cannot start gRPC server")
	// }
	return grpcServer.Serve(listener)
}

func runGrpGatewayServer(
	//ctx context.Context,
	config util.Config,
	store db.Store,
) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create server: %v", err)
	}
	grpcmux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimplebankHandlerServer(ctx, grpcmux, server)
	if err != nil {
		log.Fatalf("cannot register handler server: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcmux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatalf("cannot create statik fs: %v", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatalf("Failed to create listener on port 8080: %v", err)
	}

	// 输出启动信息
	log.Printf("Starting gRPC Gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
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
