package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	db "github.com/dreamcreeep/roflan_bank/db/sqlc"
	"github.com/dreamcreeep/roflan_bank/db/util"
	"github.com/dreamcreeep/roflan_bank/gapi"
	"github.com/dreamcreeep/roflan_bank/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	config := util.Config{
		DBDriver:          os.Getenv("DB_DRIVER"),
		DBSource:          os.Getenv("DB_SOURCE"),
		HTTPServerAddress: os.Getenv("HTTP_SERVER_ADDRESS"),
		GRPCServerAddress: os.Getenv("GRPC_SERVER_ADDRESS"),
	}

	// Эти переменные требуют парсинга
	accessTokenDuration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		logger.Error("invalid access token duration", slog.Any("error", err))
		os.Exit(1)
	}
	refreshTokenDuration, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))
	if err != nil {
		logger.Error("invalid refresh token duration", slog.Any("error", err))
		os.Exit(1)
	}

	config.TokenSymmetricKey = os.Getenv("TOKEN_SYMMETRIC_KEY")
	config.AccessTokenDuration = accessTokenDuration
	config.RefreshTokenDuration = refreshTokenDuration

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		logger.Error("cannot connect to db", slog.Any("error", err))
		os.Exit(1)
	}

	store := db.NewStore(conn)

	go runGatewayServer(config, store, logger)
	runGrpcServer(config, store, logger)
}

func runGatewayServer(config util.Config, store db.Store, logger *slog.Logger) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		logger.Error("cannot create server", slog.Any("error", err))
		os.Exit(1)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		logger.Error("cannot register handler server", slog.Any("error", err))
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		logger.Error("cannot create listener", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("start HTTP gateway server", slog.String("address", listener.Addr().String()))

	err = http.Serve(listener, mux)
	if err != nil {
		logger.Error("cannot start HTTP gateway server", slog.Any("error", err))
		os.Exit(1)
	}
}

func runGrpcServer(config util.Config, store db.Store, logger *slog.Logger) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		logger.Error("cannot create server", slog.Any("error", err))
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(server.GrpcLogger))

	pb.RegisterSimpleBankServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		logger.Error("cannot create listener", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("start gRPC server", slog.String("address", listener.Addr().String()))
	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Error("cannot start gRPC server", slog.Any("error", err))
		os.Exit(1)
	}
}
