package loader

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitaepark/carrot-market/controller"
)

type Server struct {
	controller *controller.Controller
}

func NewServer(controller *controller.Controller) (*Server, error) {
	server := &Server{
		controller: controller,
	}

	server.setHealthCheck()

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.controller.Router.Run(address)
}

func (server *Server) setHealthCheck() {
	server.controller.Router.GET("/api/health-check", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
	})
}
