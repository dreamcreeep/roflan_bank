package gapi

import (
	"fmt"
	"log/slog"

	"github.com/dreamcreeep/roflan_bank/pb"
	"github.com/dreamcreeep/roflan_bank/worker"

	db "github.com/dreamcreeep/roflan_bank/db/sqlc"
	"github.com/dreamcreeep/roflan_bank/db/util"
	"github.com/dreamcreeep/roflan_bank/token"
)

// Server обслуживает gRPC запросы нашего банковского сервиса.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	logger          *slog.Logger
	taskDistributor worker.TaskDistributor
}

// NewServer создаёт новый HTTP сервер и настраивает маршрутизацию.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	logger := util.NewLogger(config.Environment, config.LogFormat, config.LogLevel)

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		logger:          logger,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
