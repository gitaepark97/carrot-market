package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitaepark/carrot-market/dto"
)

func (controller *Controller) setAuthRouter() {
	controller.Router.POST("/api/auth/register", controller.register)
	controller.Router.POST("/api/auth/login", controller.login)
	controller.Router.POST("/api/auth/renew-access-token", controller.renewAccessToken)
}

func (controller *Controller) register(ctx *gin.Context) {
	var reqBody dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.BindingErrorResponse(err, &reqBody, "json"))
		return
	}

	rsp, cErr := controller.service.Register(ctx, reqBody)
	if cErr.Err != nil {
		ctx.AbortWithStatusJSON(cErr.StatusCode, dto.ErrorResponse(cErr.Err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (controller *Controller) login(ctx *gin.Context) {
	var reqBody dto.LoginRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.BindingErrorResponse(err, &reqBody, "json"))
		return
	}

	rsp, cErr := controller.service.Login(ctx, reqBody)
	if cErr.Err != nil {
		ctx.AbortWithStatusJSON(cErr.StatusCode, dto.ErrorResponse(cErr.Err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (controller *Controller) renewAccessToken(ctx *gin.Context) {
	var reqBody dto.RenewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.BindingErrorResponse(err, &reqBody, "json"))
		return
	}

	rsp, cErr := controller.service.RenewAccessToken(ctx, reqBody)
	if cErr.Err != nil {
		ctx.AbortWithStatusJSON(cErr.StatusCode, dto.ErrorResponse(cErr.Err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
