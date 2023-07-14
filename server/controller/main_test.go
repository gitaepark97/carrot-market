package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	db "github.com/gitaepark/carrot-market/db/sqlc"
	"github.com/gitaepark/carrot-market/service"
	"github.com/gitaepark/carrot-market/util"
	"github.com/stretchr/testify/require"
)

func newTestController(t *testing.T, store db.Store) *Controller {
	config := util.Config{
		JWTSecret: util.CreateRandomString(32),
	}

	service, _ := service.NewService(config, store)
	constroller, err := NewController(service)
	require.NoError(t, err)

	return constroller
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

type errrorResponse struct {
	Error string `json:"error"`
}

func requireErrorMatch(t *testing.T, body *bytes.Buffer, expectErr error) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotError errrorResponse
	err = json.Unmarshal(data, &gotError)

	require.NoError(t, err)
	require.Equal(t, expectErr.Error(), gotError.Error)
}
