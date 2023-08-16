package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitaepark/carrot-market/dto"
)

func (controller *Controller) setAuthRoutes() {
	authRoutes := controller.Router.Group("/api/auth")

	authRoutes.POST("/register", func(ctx *gin.Context) {
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
	})

	authRoutes.POST("/login", func(ctx *gin.Context) {
		var reqBody dto.LoginRequest
		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.BindingErrorResponse(err, &reqBody, "json"))
			return
		}

		rsp, cErr := controller.service.Login(ctx, reqBody, ctx.Request.UserAgent(), ctx.ClientIP())
		if cErr.Err != nil {
			ctx.AbortWithStatusJSON(cErr.StatusCode, dto.ErrorResponse(cErr.Err))
			return
		}

		ctx.JSON(http.StatusOK, rsp)
	})

	authRoutes.POST("/renew-access-token", func(ctx *gin.Context) {
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
	})
}
