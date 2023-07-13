package service

import (
	"context"
	"database/sql"

	db "github.com/gitaepark/carrot-market/db/sqlc"
	"github.com/gitaepark/carrot-market/dto"
	"github.com/gitaepark/carrot-market/util"
	"github.com/lib/pq"
)

func (service *Service) Register(ctx context.Context, reqBody dto.RegisterRequest) (dto.UserResponse, util.CustomError) {
	hashedPassword, err := util.HashPassword(reqBody.Password)
	if err != nil {
		return dto.UserResponse{}, util.NewInternalServerError(err)
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
					return dto.UserResponse{}, util.ErrDuplicateEmail
				case util.DB_UK_USER_NICKNAME:
					return dto.UserResponse{}, util.ErrDuplicateNickname
				}
			}
		}

		return dto.UserResponse{}, util.NewInternalServerError(err)
	}

	rsp := dto.NewUserResponse(user)

	return rsp, util.CustomError{}
}

func (service *Service) Login(ctx context.Context, reqBody dto.LoginRequest) (dto.LoginResponse, util.CustomError) {
	user, err := service.repository.GetUserByEmail(ctx, reqBody.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.LoginResponse{}, util.ErrNotFoundUser
		}

		return dto.LoginResponse{}, util.NewInternalServerError(err)
	}

	err = util.CheckPassword(reqBody.Password, user.HashedPassword)
	if err != nil {
		return dto.LoginResponse{}, util.ErrInvalidPassword
	}

	accessToken, _, err := service.tokenMaker.CreateToken(user.UserID, service.config.AccessTokenDuration)
	if err != nil {
		return dto.LoginResponse{}, util.NewInternalServerError(err)
	}

	refreshToken, refreshPayload, err := service.tokenMaker.CreateToken(user.UserID, service.config.RefreshTokenDuration)
	if err != nil {
		return dto.LoginResponse{}, util.NewInternalServerError(err)
	}

	arg := db.CreateSessionParams{
		SessionID:    refreshPayload.ID,
		UserID:       user.UserID,
		RefreshToken: refreshToken,
	}

	_, err = service.repository.CreateSession(ctx, arg)
	if err != nil {
		return dto.LoginResponse{}, util.NewInternalServerError(err)
	}

	rsp := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         dto.NewUserResponse(user),
	}

	return rsp, util.CustomError{}
}
