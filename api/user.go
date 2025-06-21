package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/dreamcreeep/roflan_bank/db/sqlc"
	"github.com/dreamcreeep/roflan_bank/db/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logger.Error("Failed to bind create user request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	server.logger.Info("Creating new user",
		zap.String("username", req.Username),
		zap.String("email", req.Email),
		zap.String("client_ip", ctx.ClientIP()),
	)

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		server.logger.Error("Failed to hash password", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				server.logger.Warn("User creation failed - duplicate username/email",
					zap.String("username", req.Username),
					zap.String("email", req.Email),
					zap.Error(err),
				)
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		server.logger.Error("Failed to create user in database", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.logger.Info("User created successfully",
		zap.String("username", user.Username),
		zap.String("email", user.Email),
	)

	resp := newUserResponse(user)
	ctx.JSON(http.StatusOK, resp)
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logger.Error("Failed to bind login request", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	server.logger.Info("User login attempt",
		zap.String("username", req.Username),
		zap.String("client_ip", ctx.ClientIP()),
		zap.String("user_agent", ctx.GetHeader("User-Agent")),
	)

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			server.logger.Warn("Login failed - user not found",
				zap.String("username", req.Username),
				zap.String("client_ip", ctx.ClientIP()),
			)
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		server.logger.Error("Failed to get user from database", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		server.logger.Warn("Login failed - invalid password",
			zap.String("username", req.Username),
			zap.String("client_ip", ctx.ClientIP()),
		)
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		server.logger.Error("Failed to create access token", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		server.logger.Error("Failed to create refresh token", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.GetHeader("User-Agent"),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		server.logger.Error("Failed to create session", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.logger.Info("User logged in successfully",
		zap.String("username", user.Username),
		zap.String("session_id", session.ID.String()),
		zap.String("client_ip", ctx.ClientIP()),
	)

	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}
