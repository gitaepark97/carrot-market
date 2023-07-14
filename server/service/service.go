package service

import (
	db "github.com/gitaepark/carrot-market/db/sqlc"
	"github.com/gitaepark/carrot-market/token"
	"github.com/gitaepark/carrot-market/util"
)

type Service struct {
	config     util.Config
	TokenMaker token.Maker
	repository db.Store
}

func NewService(config util.Config, repository db.Store) (*Service, error) {
	service := &Service{
		config:     config,
		repository: repository,
	}

	tokenMaker, err := token.NewJWTMaker(config.JWTSecret)
	if err != nil {
		return service, err
	}

	service.TokenMaker = tokenMaker

	return service, nil
}
