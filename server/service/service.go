package service

import (
	db "github.com/gitaepark/carrot-market/db/sqlc"
	"github.com/gitaepark/carrot-market/token"
	"github.com/gitaepark/carrot-market/util"
)

type Service struct {
	config     util.Config
	tokenMaker token.Maker
	repository db.Store
}

func NewService(config util.Config, tokenMaker token.Maker, repository db.Store) *Service {
	service := &Service{
		config:     config,
		tokenMaker: tokenMaker,
		repository: repository,
	}

	return service
}
