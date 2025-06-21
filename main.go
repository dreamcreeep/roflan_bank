package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/dreamcreeep/roflan_bank/api"
	db "github.com/dreamcreeep/roflan_bank/db/sqlc"
	"github.com/dreamcreeep/roflan_bank/db/util"
	_ "github.com/lib/pq"
)

func main() {
	config := util.Config{
		DBDriver:      os.Getenv("DB_DRIVER"),
		DBSource:      os.Getenv("DB_SOURCE"),
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
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
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
