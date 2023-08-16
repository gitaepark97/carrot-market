package service

import (
	"context"
	"database/sql"

	db "github.com/gitaepark/carrot-market/db/sqlc"
	"github.com/gitaepark/carrot-market/dto"
	"github.com/gitaepark/carrot-market/util"
	"github.com/lib/pq"
)

func (service *Service) CreateGoods(ctx context.Context, userID int32, reqBody dto.CreateGoodsBodyRequest) (rsp dto.GoodsResponse, cErr util.CustomError) {
	arg := db.CreateGoodsTxParams{
		UserID:          userID,
		Title:           reqBody.Title,
		Price:           reqBody.Price,
		Description:     reqBody.Description,
		DefaultImageUrl: reqBody.DefaultImageUrl,
		CategoryIDList:  reqBody.CategoryIDList,
		ImageUrlList:    reqBody.ImageUrlList,
	}

	result, err := service.repository.CreateGoodsTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case util.DB_FK_ERROR.Name:
				cErr = util.ErrNotFoundCategory
				return
			}
		}

		cErr = util.NewInternalServerError(err)
		return
	}

	rsp = dto.NewGoodsResponse(result.Goods, result.GoodsImageList)
	return
}

func (service *Service) GetGoodsList(ctx context.Context, reqQuery dto.GetGoodsListQueryRequest) (rsp dto.GetGoodsListResponse, cErr util.CustomError) {
	arg := db.GetGoodsListParams{
		Limit:  reqQuery.PageSize,
		Offset: (reqQuery.PageID - 1) * reqQuery.PageSize,
	}

	goodsList, err := service.repository.GetGoodsList(ctx, arg)
	if err != nil {
		cErr = util.NewInternalServerError(err)
		return
	}

	rsp = dto.NewGetGoodsListResponse(goodsList)

	return
}

func (service *Service) GetGoods(ctx context.Context, reqPath dto.GoodsPathRequest) (rsp dto.GoodsResponse, cErr util.CustomError) {
	goods, err := service.repository.GetGoods(ctx, reqPath.GoodsID)
	if err != nil {
		if err == sql.ErrNoRows {
			cErr = util.ErrNotFoundGoods
			return
		}

		cErr = util.NewInternalServerError(err)
		return
	}

	goodsImageList, err := service.repository.GetGoodsImageList(ctx, goods.GoodsID)
	if err != nil {
		cErr = util.NewInternalServerError(err)
		return
	}

	rsp = dto.NewGoodsResponse(goods, goodsImageList)

	return
}

func (service *Service) UpdateGoods(ctx context.Context, userID int32, reqPath dto.GoodsPathRequest, reqBody dto.UpdateGoodsBodyRequest) (rsp dto.GoodsResponse, cErr util.CustomError) {
	goods, err := service.repository.GetGoods(ctx, reqPath.GoodsID)
	if err != nil {
		if err == sql.ErrNoRows {
			cErr = util.ErrNotFoundGoods
			return
		}

		cErr = util.NewInternalServerError(err)
		return
	}

	if goods.UserID != userID {
		cErr = util.ErrForbiddenUser
		return
	}

	arg := db.UpdateGoodsTxParams{
		GoodsID:                reqPath.GoodsID,
		Title:                  util.CreateNullableString(reqBody.Title),
		Price:                  util.CreateNullableInt32(&reqBody.Price),
		Description:            util.CreateNullableString(reqBody.Description),
		DefaultImageUrl:        util.CreateNullableString(reqBody.DefaultImageUrl),
		AddCategoryIDList:      reqBody.AddCategoryIDList,
		DeleteCategoryIDList:   reqBody.DeleteCategoryIDList,
		AddImageUrlList:        reqBody.AddImageUrlList,
		DeleteGoodsImageIDList: reqBody.DeleteGoodsImageIDList,
	}

	err = service.repository.UpdateGoodsTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case util.DB_UK_ERROR.Name:
				cErr = util.ErrDuplicateCategory
				return
			case util.DB_FK_ERROR.Name:
				cErr = util.ErrNotFoundCategory
				return
			}
		}

		cErr = util.NewInternalServerError(err)
		return
	}

	goods, err = service.repository.GetGoods(ctx, goods.GoodsID)
	if err != nil {
		cErr = util.NewInternalServerError(err)
		return
	}

	goodsImage, err := service.repository.GetGoodsImageList(ctx, goods.GoodsID)
	if err != nil {
		cErr = util.NewInternalServerError(err)
		return
	}

	rsp = dto.NewGoodsResponse(goods, goodsImage)
	return
}

func (service *Service) DeleteGoods(ctx context.Context, userID int32, reqPath dto.GoodsPathRequest) (cErr util.CustomError) {
	goods, err := service.repository.GetGoods(ctx, reqPath.GoodsID)
	if err != nil {
		if err == sql.ErrNoRows {
			cErr = util.ErrNotFoundGoods
			return
		}

		cErr = util.NewInternalServerError(err)
		return
	}

	if goods.UserID != userID {
		cErr = util.ErrForbiddenUser
		return
	}

	err = service.repository.DeleteGoods(ctx, goods.GoodsID)
	if err != nil {
		cErr = util.NewInternalServerError(err)
		return
	}

	return
}
