package middleware

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gitaepark/carrot-market/token"
	"github.com/gitaepark/carrot-market/util"
)

type Server struct {
	tokenMaker token.Maker
	router     *gin.Engine
}

func newServer() Server {
	tokenMaker, _ := token.NewJWTMaker(util.CreateRandomString(32))

	return Server{
		tokenMaker: tokenMaker,
		router:     gin.Default(),
	}
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
