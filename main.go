package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	"github.com/dreamcreeep/roflan_bank/api"
	db "github.com/dreamcreeep/roflan_bank/db/sqlc"
	"github.com/dreamcreeep/roflan_bank/db/util"
	"github.com/dreamcreeep/roflan_bank/gapi"
	"github.com/dreamcreeep/roflan_bank/pb"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables from system")
	}

	config := util.Config{
		DBDriver:          os.Getenv("DB_DRIVER"),
		DBSource:          os.Getenv("DB_SOURCE"),
		HTTPServerAddress: os.Getenv("HTTP_SERVER_ADDRESS"),
		GRPCServerAddress: os.Getenv("GRPC_SERVER_ADDRESS"),
	}

	// Эти переменные требуют парсинга
	accessTokenDuration, err := time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	if err != nil {
		log.Fatalf("invalid access token duration: %v", err)
	}
	refreshTokenDuration, err := time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))
	if err != nil {
		log.Fatalf("invalid refresh token duration: %v", err)
	}

	config.TokenSymmetricKey = os.Getenv("TOKEN_SYMMETRIC_KEY")
	config.AccessTokenDuration = accessTokenDuration
	config.RefreshTokenDuration = refreshTokenDuration

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	runGrpcServer(config, store)

}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterSimpleBankServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server")
	}

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
