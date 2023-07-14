package service

import (
	"testing"

	db "github.com/gitaepark/carrot-market/db/sqlc"
	"github.com/gitaepark/carrot-market/util"
	"github.com/stretchr/testify/require"
)

func newTestService(t *testing.T, store db.Store) *Service {
	config := util.Config{
		JWTSecret: util.CreateRandomString(32),
	}

	service, err := NewService(config, store)
	require.NoError(t, err)

	return service
}

func requireErrorMatch(t *testing.T, err, expectErr util.CustomError) {
	require.Equal(t, err, expectErr)
}
