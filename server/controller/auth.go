package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitaepark/carrot-market/dto"
)

func (controller *Controller) setAuthRouter() {
	controller.Router.POST("/api/auth/register", controller.Register)
}

func (controller *Controller) Register(ctx *gin.Context) {
	var reqBody dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.BindingErrorResponse(err, &reqBody, "json"))
		return
	}

	rsp, err := controller.service.Register(ctx, reqBody)
	if err.Err != nil {
		ctx.AbortWithStatusJSON(err.StatusCode, dto.ErrorResponse(err.Err))
		return
	}

	ctx.JSON(http.StatusOK, rsp)
}
