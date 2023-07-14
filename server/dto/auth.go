package dto

import (
	"time"

	db "github.com/gitaepark/carrot-market/db/sqlc"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname" binding:"required,max=50"`
}

type UserResponse struct {
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserResponse(user db.User) UserResponse {
	userResponse := UserResponse{
		Email:     user.Email,
		Nickname:  user.Nickname,
		CreatedAt: user.CreatedAt.Local(),
		UpdatedAt: user.UpdatedAt.Local(),
	}

	return userResponse
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken  string       `json:"access_token"`
}