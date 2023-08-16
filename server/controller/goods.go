package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gitaepark/carrot-market/dto"
	"github.com/gitaepark/carrot-market/middleware"
	"github.com/gitaepark/carrot-market/token"
)

func (controller *Controller) setGoodsRoutes() {
	goodsRoutes := controller.Router.Group("/api/goods").Use(middleware.AuthMiddleware(controller.service.TokenMaker))

	goodsRoutes.POST("/", func(ctx *gin.Context) {
		var reqBody dto.CreateGoodsBodyRequest
		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.BindingErrorResponse(err, &reqBody, "json"))
			return
		}

		authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

		rsp, cErr := controller.service.CreateGoods(ctx, authPayload.UserID, reqBody)
		if cErr.Err != nil {
			ctx.AbortWithStatusJSON(cErr.StatusCode, dto.ErrorResponse(cErr.Err))
			return
		}

		ctx.JSON(http.StatusOK, rsp)
	})

	goodsRoutes.GET("/", func(ctx *gin.Context) {
		var reqQuery dto.GetGoodsListQueryRequest
		if err := ctx.ShouldBindQuery(&reqQuery); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.BindingErrorResponse(err, &reqQuery, "form"))
			return
		}

		rsp, cErr := controller.service.GetGoodsList(ctx, reqQuery)
		if cErr.Err != nil {
			ctx.AbortWithStatusJSON(cErr.StatusCode, dto.ErrorResponse(cErr.Err))
			return
		}

		ctx.JSON(http.StatusOK, rsp)
	})

	goodsRoutes.GET("/:goods_id", func(ctx *gin.Context) {
		var reqPath dto.GoodsPathRequest
		if err := ctx.ShouldBindUri(&reqPath); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.BindingErrorResponse(err, &reqPath, "uri"))
			return
		}

		rsp, cErr := controller.service.GetGoods(ctx, reqPath)
		if cErr.Err != nil {
			ctx.AbortWithStatusJSON(cErr.StatusCode, dto.ErrorResponse(cErr.Err))
			return
		}

		ctx.JSON(http.StatusOK, rsp)
	})

	goodsRoutes.PATCH("/:goods_id", func(ctx *gin.Context) {
		var reqPath dto.GoodsPathRequest
		if err := ctx.ShouldBindUri(&reqPath); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.BindingErrorResponse(err, &reqPath, "uri"))
			return
		}
		var reqBody dto.UpdateGoodsBodyRequest
		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.BindingErrorResponse(err, &reqBody, "json"))
			return
		}

		authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

		rsp, cErr := controller.service.UpdateGoods(ctx, authPayload.UserID, reqPath, reqBody)
		if cErr.Err != nil {
			ctx.AbortWithStatusJSON(cErr.StatusCode, dto.ErrorResponse(cErr.Err))
			return
		}

		ctx.JSON(http.StatusOK, rsp)
	})

	goodsRoutes.DELETE("/:goods_id", func(ctx *gin.Context) {
		var reqPath dto.GoodsPathRequest
		if err := ctx.ShouldBindUri(&reqPath); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.BindingErrorResponse(err, &reqPath, "uri"))
			return
		}

		authPayload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

		cErr := controller.service.DeleteGoods(ctx, authPayload.UserID, reqPath)
		if cErr.Err != nil {
			ctx.AbortWithStatusJSON(cErr.StatusCode, dto.ErrorResponse(cErr.Err))
			return
		}

		ctx.Status(http.StatusOK)
	})
}
