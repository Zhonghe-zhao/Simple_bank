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

	api "project/simplebank/api"
	db "project/simplebank/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		if err := runGrpGatewayServer(config, store); err != nil {
			log.Fatal(err)
		}

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

func runGrpGatewayServer(config util.Config, store db.Store) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	grpcmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterSimplebankHandlerFromEndpoint(ctx, grpcmux, *grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalf("cannot register handler from endpoint: %v", err)
	}
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	// 创建 HTTP 监听器
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Failed to create listener on port 8081: %v", err)
	}

	// 输出启动信息
	log.Printf("Starting gRPC Gateway server at %s", listener.Addr().String())

	return http.Serve(listener, grpcmux)
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
