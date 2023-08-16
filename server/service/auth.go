package service

import (
	"context"
	"database/sql"
	"time"

	db "github.com/gitaepark/carrot-market/db/sqlc"
	"github.com/gitaepark/carrot-market/dto"
	"github.com/gitaepark/carrot-market/util"
	"github.com/lib/pq"
)

func (service *Service) Register(ctx context.Context, reqBody dto.RegisterRequest) (rsp dto.UserResponse, cErr util.CustomError) {
	hashedPassword, err := util.HashPassword(reqBody.Password)
	if err != nil {
		cErr = util.NewInternalServerError(err)
		return
	}

	arg := db.CreateUserParams{
		Email:          reqBody.Email,
		HashedPassword: hashedPassword,
		Nickname:       reqBody.Nickname,
	}

	user, err := service.repository.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == util.DB_UK_ERROR.Name {
				switch pqErr.Constraint {
				case util.DB_UK_USER_EMAIL:
					cErr = util.ErrDuplicateEmail
					return
				case util.DB_UK_USER_NICKNAME:
					cErr = util.ErrDuplicateNickname
					return
				}
			}
		}

		cErr = util.NewInternalServerError(err)
		return
	}

	rsp = dto.NewUserResponse(user)
	return
}

func (service *Service) Login(ctx context.Context, reqBody dto.LoginRequest, userAgent, clientIp string) (rsp dto.LoginResponse, cErr util.CustomError) {
	user, err := service.repository.GetUserByEmail(ctx, reqBody.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			cErr = util.ErrNotFoundUser
			return
		}

		cErr = util.NewInternalServerError(err)
		return
	}

	err = util.CheckPassword(reqBody.Password, user.HashedPassword)
	if err != nil {
		cErr = util.ErrInvalidPassword
		return
	}

	accessToken, _, err := service.TokenMaker.CreateToken(user.UserID, service.config.AccessTokenDuration)
	if err != nil {
		cErr = util.NewInternalServerError(err)
		return
	}

	refreshToken, refreshPayload, err := service.TokenMaker.CreateToken(user.UserID, service.config.RefreshTokenDuration)
	if err != nil {
		cErr = util.NewInternalServerError(err)
		return
	}

	arg := db.CreateSessionParams{
		SessionID:    refreshPayload.ID,
		UserID:       user.UserID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIp,
		IsBlocked:    false,
		ExpiredAt:    refreshPayload.ExpiredAt,
	}

	_, err = service.repository.CreateSession(ctx, arg)
	if err != nil {
		cErr = util.NewInternalServerError(err)
		return
	}

	rsp = dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         dto.NewUserResponse(user),
	}
	return
}

func (service *Service) RenewAccessToken(ctx context.Context, reqBody dto.RenewAccessTokenRequest) (rsp dto.RenewAccessTokenResponse, cErr util.CustomError) {
	refreshPayload, err := service.TokenMaker.VerifyToken(reqBody.RefreshToken)
	if err != nil {
		cErr = util.NewInternalServerError(err)
		return
	}

	session, err := service.repository.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			cErr = util.ErrNotFoundSession
			return
		}

		cErr = util.NewInternalServerError(err)
		return
	}

	if session.IsBlocked {
		cErr = util.ErrBlockedSession
		return
	}
	if session.UserID != refreshPayload.UserID {
		cErr = util.ErrIncorrectSessionUser
		return
	}
	if session.RefreshToken != reqBody.RefreshToken {
		cErr = util.ErrMismatchedSessionToken
		return
	}
	if time.Now().After(session.ExpiredAt) {
		cErr = util.ErrExpiredSession
		return
	}

	accessToken, _, err := service.TokenMaker.CreateToken(refreshPayload.UserID, service.config.AccessTokenDuration)
	if err != nil {
		cErr = util.NewInternalServerError(err)
		return
	}

	rsp = dto.RenewAccessTokenResponse{
		AccessToken: accessToken,
	}
	return
}
