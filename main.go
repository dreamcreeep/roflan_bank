package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/hibiken/asynq"

	db "github.com/dreamcreeep/roflan_bank/db/sqlc"
	"github.com/dreamcreeep/roflan_bank/db/util"
	"github.com/dreamcreeep/roflan_bank/gapi"
	"github.com/dreamcreeep/roflan_bank/pb"
	"github.com/dreamcreeep/roflan_bank/worker"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func setConfig() util.Config {
	config := util.Config{
		DBDriver:          os.Getenv("DB_DRIVER"),
		DBSource:          os.Getenv("DBSOURCE"),
		HTTPServerAddress: os.Getenv("HTTP_SERVER_ADDRESS"),
		GRPCServerAddress: os.Getenv("GRPC_SERVER_ADDRESS"),
		RedisAddress:      os.Getenv("REDIS_ADDRESS"),
		TokenSymmetricKey: os.Getenv("TOKEN_SYMMETRIC_KEY"),
		Environment:       os.Getenv("ENVIRONMENT"),
		LogLevel:          os.Getenv("LOG_LEVEL"),
		LogFormat:         os.Getenv("LOG_FORMAT"),
	}

	return config
}

func main() {
	config := setConfig()

	logger := util.NewLogger(config.Environment, config.LogFormat, config.LogLevel)

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

	config.AccessTokenDuration = accessTokenDuration
	config.RefreshTokenDuration = refreshTokenDuration

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		logger.Error("cannot connect to db", slog.Any("error", err))
		os.Exit(1)
	}

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, logger, store)

	go runGatewayServer(config, store, taskDistributor, logger)
	runGrpcServer(config, store, taskDistributor, logger)
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, logger *slog.Logger, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, logger, store)
	slog.Info("start task processor")

	err := taskProcessor.Start()
	if err != nil {
		logger.Error("cannot start task processor", slog.Any("error", err))
		os.Exit(1)
	}

}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor, logger *slog.Logger) {
	server, err := gapi.NewServer(config, store, taskDistributor, logger)
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

func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor, logger *slog.Logger) {
	server, err := gapi.NewServer(config, store, taskDistributor, logger)
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
