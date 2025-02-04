package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"project/simplebank/gapi"
	"project/simplebank/pb"
	"project/simplebank/util"

	"github.com/rs/zerolog/log"

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
		log.Fatal().Err(err).Msg("cannot load config")
	}

	connPool, err := pgxpool.New(context.Background(), config.DB_SOURCE)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}
	//初始化数据库服务
	store := db.NewStore(connPool)
	//运行gin框架
	//RunGinServer(config, store)
	go func() {
		runGrpGatewayServer(config, store)
	}()
	runGrpcServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run gRPC server")
	}
	//创建gRPC服务器实例
	gprcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(gprcLogger)
	//注册服务
	pb.RegisterSimplebankServer(grpcServer, server)
	//注册反射服务
	reflection.Register(grpcServer)
	//创建监听器
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}
	//启动gRPC服务器
	log.Printf("start gRPC server at %s", listener.Addr().String())
	// err = grpcServer.Serve(listener)
	// if err != nil {
	// 	log.Fatal("cannot start gRPC server")
	// }
	grpcServer.Serve(listener)
}

func runGrpGatewayServer(
	//ctx context.Context,
	config util.Config,
	store db.Store,
) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}
	grpcmux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterSimplebankHandlerServer(ctx, grpcmux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcmux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	// 输出启动信息
	log.Printf("Starting gRPC Gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}

}

func RunGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
