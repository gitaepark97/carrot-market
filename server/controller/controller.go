package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gitaepark/carrot-market/service"
)

type Controller struct {
	service *service.Service
	Router  *gin.Engine
}

func NewController(service *service.Service) (*Controller, error) {
	controller := &Controller{
		service: service,
	}

	controller.setupRouter()

	return controller, nil
}

func (controller *Controller) setupRouter() {
	router := gin.Default()

	controller.Router = router

	controller.setAuthRoutes()
	controller.setGoodsRoutes()
}
