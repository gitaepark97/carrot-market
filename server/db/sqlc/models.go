// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	CategoryID int32     `json:"category_id"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Good struct {
	GoodsID         int32     `json:"goods_id"`
	UserID          int32     `json:"user_id"`
	Title           string    `json:"title"`
	Price           int32     `json:"price"`
	Description     string    `json:"description"`
	DefaultImageUrl string    `json:"default_image_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type GoodsCategory struct {
	GoodsID    int32     `json:"goods_id"`
	CategoryID int32     `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type GoodsImage struct {
	GoodsImageID int32     `json:"goods_image_id"`
	GoodsID      int32     `json:"goods_id"`
	ImageUrl     string    `json:"image_url"`
	CreatedAt    time.Time `json:"created_at"`
}

type Session struct {
	SessionID    uuid.UUID `json:"session_id"`
	UserID       int32     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiredAt    time.Time `json:"expired_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type User struct {
	UserID         int32     `json:"user_id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	Nickname       string    `json:"nickname"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
