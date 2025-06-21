package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

	db "github.com/dreamcreeep/roflan_bank/db/sqlc"
	"github.com/dreamcreeep/roflan_bank/db/util"
	"github.com/dreamcreeep/roflan_bank/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	logger     *slog.Logger
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	logger := util.NewLogger(config.Environment, config.LogFormat, config.LogLevel)

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		logger:     logger,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.New() // Используем gin.New() вместо gin.Default()

	// Добавляем middleware для логирования и восстановления
	router.Use(LoggerMiddleware(server.logger))
	router.Use(RecoveryMiddleware(server.logger))

	// Health check endpoint
	router.GET("/health", server.healthCheck)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	srv := &http.Server{
		Addr:    address,
		Handler: server.router,
	}

	// Запускаем сервер в горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Ждем сигнала прерывания
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Даем 5 секунд на завершение текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
	return nil
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}

// healthCheck - endpoint для проверки состояния сервиса
func (server *Server) healthCheck(ctx *gin.Context) {
	server.logger.Info("Health check requested",
		slog.String("client_ip", ctx.ClientIP()),
		slog.String("user_agent", ctx.GetHeader("User-Agent")),
	)
}
