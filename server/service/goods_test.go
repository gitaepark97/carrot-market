package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	mockdb "github.com/gitaepark/carrot-market/db/mock"
	db "github.com/gitaepark/carrot-market/db/sqlc"
	"github.com/gitaepark/carrot-market/dto"
	"github.com/gitaepark/carrot-market/util"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateGoods(t *testing.T) {
	user, _ := createRandomUser(t)
	goods := createRandomGoods(t, user)
	goodsCategoryTitles := createCategoryTitles(2)
	var goodsCategoryIDList []int32
	var goodsImageList []db.GoodsImage
	var goodsImageUrlList []string

	for i := 0; i < 2; i++ {
		goodsCategory := createRandomGoodsCategory(t, goods)
		goodsCategoryIDList = append(goodsCategoryIDList, goodsCategory.CategoryID)
	}

	for i := 0; i < 2; i++ {
		goodsImage := createRandomGoodsImage(t, goods)
		goodsImageList = append(goodsImageList, goodsImage)
		goodsImageUrlList = append(goodsImageUrlList, goodsImage.ImageUrl)
	}

	testCases := []struct {
		name          string
		userID        int32
		reqBody       dto.CreateGoodsBodyRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(rsp dto.GoodsResponse, err util.CustomError)
	}{
		{
			name:   "OK",
			userID: goods.UserID,
			reqBody: dto.CreateGoodsBodyRequest{
				Title:           goods.Title,
				Price:           goods.Price,
				Description:     goods.Description,
				DefaultImageUrl: goods.DefaultImageUrl,
				CategoryIDList:  goodsCategoryIDList,
				ImageUrlList:    &goodsImageUrlList,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateGoodsTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateGoodsTxResult{
						Goods:          createGetGoodsRow(goods, goodsCategoryTitles),
						GoodsImageList: goodsImageList,
					}, nil)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				requireMatchGoodsResponse(t, rsp, goods, goodsCategoryTitles, goodsImageList)
				require.Empty(t, err)
			},
		},
		{
			name:   "InternalServerError",
			userID: goods.UserID,
			reqBody: dto.CreateGoodsBodyRequest{
				Title:           goods.Title,
				Price:           goods.Price,
				Description:     goods.Description,
				DefaultImageUrl: goods.DefaultImageUrl,
				CategoryIDList:  goodsCategoryIDList,
				ImageUrlList:    &goodsImageUrlList,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateGoodsTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateGoodsTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireMatchError(t, err, util.NewInternalServerError(sql.ErrConnDone))
			},
		},
		{
			name:   "NotFoundCategory",
			userID: goods.UserID,
			reqBody: dto.CreateGoodsBodyRequest{
				Title:           goods.Title,
				Price:           goods.Price,
				Description:     goods.Description,
				DefaultImageUrl: goods.DefaultImageUrl,
				CategoryIDList:  []int32{0},
				ImageUrlList:    &goodsImageUrlList,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateGoodsTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateGoodsTxResult{}, &pq.Error{Code: pq.ErrorCode(util.DB_FK_ERROR.Code)})
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireMatchError(t, err, util.ErrNotFoundCategory)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			service := newTestService(t, store)

			tc.buildStubs(store)

			rsp, err := service.CreateGoods(context.Background(), tc.userID, tc.reqBody)
			tc.checkResponse(rsp, err)
		})
	}
}

func TestGetGoodsList(t *testing.T) {
	user, _ := createRandomUser(t)
	var goodsList []db.Good

	for i := 0; i < 10; i++ {
		goods := createRandomGoods(t, user)
		goodsList = append(goodsList, goods)
	}

	testCases := []struct {
		name          string
		reqQuery      dto.GetGoodsListQueryRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(rsp dto.GetGoodsListResponse, err util.CustomError)
	}{
		{
			name: "OK",
			reqQuery: dto.GetGoodsListQueryRequest{
				PageID:   1,
				PageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoodsList(gomock.Any(), gomock.Any()).
					Times(1).
					Return(goodsList, nil)
			},
			checkResponse: func(rsp dto.GetGoodsListResponse, err util.CustomError) {
				requireMatchGetGoodsListResponse(t, rsp, goodsList)
				require.Empty(t, err)
			},
		},
		{
			name: "InternalServerError",
			reqQuery: dto.GetGoodsListQueryRequest{
				PageID:   1,
				PageSize: 10,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoodsList(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Good{}, sql.ErrConnDone)
			},
			checkResponse: func(rsp dto.GetGoodsListResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireMatchError(t, err, util.NewInternalServerError(sql.ErrConnDone))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			service := newTestService(t, store)

			tc.buildStubs(store)

			rsp, err := service.GetGoodsList(context.Background(), tc.reqQuery)
			tc.checkResponse(rsp, err)
		})
	}
}

func TestGetGoods(t *testing.T) {
	user, _ := createRandomUser(t)
	goods := createRandomGoods(t, user)
	goodsCategoryTitles := createCategoryTitles(1)
	var goodsImageList []db.GoodsImage

	for i := 0; i < 2; i++ {
		goodsImage := createRandomGoodsImage(t, goods)
		goodsImageList = append(goodsImageList, goodsImage)
	}

	testCases := []struct {
		name          string
		reqPath       dto.GoodsPathRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(rsp dto.GoodsResponse, err util.CustomError)
	}{
		{
			name: "OK",
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, goodsCategoryTitles), nil)

				store.EXPECT().
					GetGoodsImageList(gomock.Any(), gomock.Any()).
					Times(1).
					Return(goodsImageList, nil)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				requireMatchGoodsResponse(t, rsp, goods, goodsCategoryTitles, goodsImageList)
				require.Empty(t, err)
			},
		},
		{
			name: "InternalServerError",
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetGoodsRow{}, sql.ErrConnDone)

				store.EXPECT().
					GetGoodsImageList(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireMatchError(t, err, util.NewInternalServerError(sql.ErrConnDone))
			},
		},
		{
			name: "NotFoundGoods",
			reqPath: dto.GoodsPathRequest{
				GoodsID: 0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetGoodsRow{}, sql.ErrNoRows)

				store.EXPECT().
					GetGoodsImageList(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireMatchError(t, err, util.ErrNotFoundGoods)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			service := newTestService(t, store)

			tc.buildStubs(store)

			rsp, err := service.GetGoods(context.Background(), tc.reqPath)
			tc.checkResponse(rsp, err)
		})
	}
}

func TestUpdateGoods(t *testing.T) {
	user, _ := createRandomUser(t)
	goods := createRandomGoods(t, user)
	goodsCategoryTitles := createCategoryTitles(3)
	var goodsCategoryIDList []int32
	var goodsImageList []db.GoodsImage

	for i := 0; i < 3; i++ {
		goodsCategory := createRandomGoodsCategory(t, goods)
		goodsCategoryIDList = append(goodsCategoryIDList, goodsCategory.CategoryID)
	}

	for i := 0; i < 3; i++ {
		goodsImage := createRandomGoodsImage(t, goods)
		goodsImageList = append(goodsImageList, goodsImage)
	}

	existGoodsCategoryTitles := strings.Join(strings.Split(goodsCategoryTitles, ",")[:1], ",")
	existGoodsImageList := goodsImageList[:1]

	newTitle := util.CreateRandomString(50)
	newPrice := util.CreateRandomInt32(1000, 100000)
	newDescription := util.CreateRandomString(100)
	newDefaulImageUrl := util.CreateRandomString(100)

	testCases := []struct {
		name          string
		userID        int32
		reqPath       dto.GoodsPathRequest
		reqBody       dto.UpdateGoodsBodyRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(rsp dto.GoodsResponse, err util.CustomError)
	}{
		{
			name:   "OnlyTitleOK",
			userID: goods.UserID,
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			reqBody: dto.UpdateGoodsBodyRequest{
				Title: newTitle,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					UpdateGoodsTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				goods.Title = newTitle

				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					GetGoodsImageList(gomock.Any(), gomock.Any()).
					Times(1).
					Return(existGoodsImageList, nil)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				requireMatchGoodsResponse(t, rsp, goods, existGoodsCategoryTitles, existGoodsImageList)
				require.Empty(t, err)
			},
		},
		{
			name:   "OnlyPriceOK",
			userID: goods.UserID,
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			reqBody: dto.UpdateGoodsBodyRequest{
				Price: newPrice,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					UpdateGoodsTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				goods.Price = newPrice

				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					GetGoodsImageList(gomock.Any(), gomock.Any()).
					Times(1).
					Return(existGoodsImageList, nil)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				requireMatchGoodsResponse(t, rsp, goods, existGoodsCategoryTitles, existGoodsImageList)
				require.Empty(t, err)
			},
		},
		{
			name:   "OnlyDescriptionOK",
			userID: goods.UserID,
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			reqBody: dto.UpdateGoodsBodyRequest{
				Description: newDescription,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					UpdateGoodsTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				goods.Description = newDescription

				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					GetGoodsImageList(gomock.Any(), gomock.Any()).
					Times(1).
					Return(existGoodsImageList, nil)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				requireMatchGoodsResponse(t, rsp, goods, existGoodsCategoryTitles, existGoodsImageList)
				require.Empty(t, err)
			},
		},
		{
			name:   "OnlyDefaultImageUrlOK",
			userID: goods.UserID,
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			reqBody: dto.UpdateGoodsBodyRequest{
				DefaultImageUrl: newDefaulImageUrl,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					UpdateGoodsTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				goods.DefaultImageUrl = newDefaulImageUrl

				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					GetGoodsImageList(gomock.Any(), gomock.Any()).
					Times(1).
					Return(existGoodsImageList, nil)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				requireMatchGoodsResponse(t, rsp, goods, existGoodsCategoryTitles, existGoodsImageList)
				require.Empty(t, err)
			},
		},
		{
			name:   "OnlyAddCategoryIDListOK",
			userID: goods.UserID,
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			reqBody: dto.UpdateGoodsBodyRequest{
				AddCategoryIDList: &[]int32{goodsCategoryIDList[2]},
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					UpdateGoodsTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, goodsCategoryTitles), nil)

				store.EXPECT().
					GetGoodsImageList(gomock.Any(), gomock.Any()).
					Times(1).
					Return(existGoodsImageList, nil)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				requireMatchGoodsResponse(t, rsp, goods, goodsCategoryTitles, existGoodsImageList)
				require.Empty(t, err)
			},
		},
		{
			name:   "OnlyDeleteCategoryIDListOK",
			userID: goods.UserID,
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			reqBody: dto.UpdateGoodsBodyRequest{
				DeleteCategoryIDList: &[]int32{goodsCategoryIDList[2]},
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					UpdateGoodsTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					GetGoodsImageList(gomock.Any(), gomock.Any()).
					Times(1).
					Return(existGoodsImageList, nil)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				requireMatchGoodsResponse(t, rsp, goods, existGoodsCategoryTitles, existGoodsImageList)
				require.Empty(t, err)
			},
		},
		{
			name:   "OnlyAddImageUrlListOK",
			userID: goods.UserID,
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			reqBody: dto.UpdateGoodsBodyRequest{
				AddImageUrlList: &[]string{goodsImageList[2].ImageUrl},
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					UpdateGoodsTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					GetGoodsImageList(gomock.Any(), gomock.Any()).
					Times(1).
					Return(goodsImageList, nil)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				requireMatchGoodsResponse(t, rsp, goods, existGoodsCategoryTitles, goodsImageList)
				require.Empty(t, err)
			},
		},
		{
			name:   "OnlyDeleteImageUrlListOK",
			userID: goods.UserID,
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			reqBody: dto.UpdateGoodsBodyRequest{
				DeleteGoodsImageIDList: &[]int32{goodsImageList[2].GoodsImageID},
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					UpdateGoodsTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)

				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					GetGoodsImageList(gomock.Any(), gomock.Any()).
					Times(1).
					Return(existGoodsImageList, nil)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				requireMatchGoodsResponse(t, rsp, goods, existGoodsCategoryTitles, existGoodsImageList)
				require.Empty(t, err)
			},
		},
		{
			name:   "InternalServerError",
			userID: goods.UserID,
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			reqBody: dto.UpdateGoodsBodyRequest{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetGoodsRow{}, sql.ErrConnDone)

				store.EXPECT().
					UpdateGoodsTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireMatchError(t, err, util.NewInternalServerError(sql.ErrConnDone))
			},
		},
		{
			name:   "NotFoundGoods",
			userID: goods.UserID,
			reqPath: dto.GoodsPathRequest{
				GoodsID: 0,
			},
			reqBody: dto.UpdateGoodsBodyRequest{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetGoodsRow{}, sql.ErrNoRows)

				store.EXPECT().
					UpdateGoodsTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireMatchError(t, err, util.ErrNotFoundGoods)
			},
		},
		{
			name:   "ForbiddenUser",
			userID: 0,
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			reqBody: dto.UpdateGoodsBodyRequest{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					UpdateGoodsTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireMatchError(t, err, util.ErrForbiddenUser)
			},
		},
		{
			name:   "NotFoundCategory",
			userID: goods.UserID,
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			reqBody: dto.UpdateGoodsBodyRequest{
				AddCategoryIDList: &[]int32{0},
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					UpdateGoodsTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&pq.Error{Code: pq.ErrorCode(util.DB_FK_ERROR.Code)})
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireMatchError(t, err, util.ErrNotFoundCategory)
			},
		},
		{
			name:   "DuplicateCategory",
			userID: goods.UserID,
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			reqBody: dto.UpdateGoodsBodyRequest{
				AddCategoryIDList: &[]int32{goodsCategoryIDList[1]},
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, existGoodsCategoryTitles), nil)

				store.EXPECT().
					UpdateGoodsTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(&pq.Error{Code: pq.ErrorCode(util.DB_UK_ERROR.Code)})
			},
			checkResponse: func(rsp dto.GoodsResponse, err util.CustomError) {
				require.Empty(t, rsp)
				requireMatchError(t, err, util.ErrDuplicateCategory)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			service := newTestService(t, store)

			tc.buildStubs(store)

			rsp, err := service.UpdateGoods(context.Background(), tc.userID, tc.reqPath, tc.reqBody)
			tc.checkResponse(rsp, err)
		})
	}
}

func TestDeleteGoods(t *testing.T) {
	user, _ := createRandomUser(t)
	goods := createRandomGoods(t, user)
	goodsCategoryTitles := createCategoryTitles(1)

	testCases := []struct {
		name          string
		userID        int32
		reqPath       dto.GoodsPathRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(err util.CustomError)
	}{
		{
			name:   "OK",
			userID: goods.UserID,
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, goodsCategoryTitles), nil)

				store.EXPECT().
					DeleteGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(err util.CustomError) {
				require.Empty(t, err)
			},
		},
		{
			name:   "InternalServerError",
			userID: goods.UserID,
			reqPath: dto.GoodsPathRequest{
				GoodsID: goods.GoodsID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetGoodsRow{}, sql.ErrConnDone)

				store.EXPECT().
					DeleteGoods(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(err util.CustomError) {
				requireMatchError(t, err, util.NewInternalServerError(sql.ErrConnDone))
			},
		},
		{
			name:   "NotFoundGoods",
			userID: goods.UserID,
			reqPath: dto.GoodsPathRequest{
				GoodsID: 0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetGoodsRow{}, sql.ErrNoRows)

				store.EXPECT().
					DeleteGoods(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(err util.CustomError) {
				requireMatchError(t, err, util.ErrNotFoundGoods)
			},
		},
		{
			name:   "ForbiddenUser",
			userID: 0,
			reqPath: dto.GoodsPathRequest{
				GoodsID: 0,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGoods(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createGetGoodsRow(goods, goodsCategoryTitles), nil)

				store.EXPECT().
					DeleteGoods(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(err util.CustomError) {
				requireMatchError(t, err, util.ErrForbiddenUser)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			service := newTestService(t, store)

			tc.buildStubs(store)

			err := service.DeleteGoods(context.Background(), tc.userID, tc.reqPath)
			tc.checkResponse(err)
		})
	}
}

func createRandomGoods(t *testing.T, user db.User) db.Good {
	return db.Good{
		GoodsID:         util.CreateRandomInt32(1, 30),
		UserID:          user.UserID,
		Title:           util.CreateRandomString(50),
		Price:           util.CreateRandomInt32(1000, 100000),
		Description:     util.CreateRandomString(100),
		DefaultImageUrl: util.CreateRandomString(100),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

func createCategoryTitles(n int) string {
	var categoryTitles string

	for i := 0; i < n; i++ {
		if i == 0 {
			categoryTitles += util.CreateRandomString(10)
		} else {
			categoryTitles += fmt.Sprintf(",%s", util.CreateRandomString(10))
		}
	}

	return categoryTitles
}

func createGetGoodsRow(goods db.Good, goodsCategoryTitles string) db.GetGoodsRow {
	return db.GetGoodsRow{
		GoodsID:         goods.GoodsID,
		UserID:          goods.UserID,
		Title:           goods.Title,
		Price:           goods.Price,
		Description:     goods.Description,
		DefaultImageUrl: goods.DefaultImageUrl,
		CreatedAt:       goods.CreatedAt,
		UpdatedAt:       goods.UpdatedAt,
		CategoryTitles:  []byte(goodsCategoryTitles),
	}
}

func createRandomGoodsCategory(t *testing.T, goods db.Good) db.GoodsCategory {
	return db.GoodsCategory{
		GoodsID:    goods.GoodsID,
		CategoryID: util.CreateRandomInt32(1, 12),
		CreatedAt:  time.Now(),
	}
}

func createRandomGoodsImage(t *testing.T, goods db.Good) db.GoodsImage {
	return db.GoodsImage{
		GoodsImageID: util.CreateRandomInt32(1, 30),
		GoodsID:      goods.GoodsID,
		ImageUrl:     util.CreateRandomString(100),
		CreatedAt:    time.Now(),
	}
}

func requireMatchGoodsResponse(t *testing.T, rsp dto.GoodsResponse, goods db.Good, goodsCategoryTitles string, goodsImageList []db.GoodsImage) {
	require.NotEmpty(t, rsp)
	require.Equal(t, goods.GoodsID, rsp.GoodsID)
	require.Equal(t, goods.Title, rsp.Title)
	require.Equal(t, goods.Price, rsp.Price)
	require.Equal(t, goods.Description, rsp.Description)
	require.Equal(t, goods.DefaultImageUrl, rsp.DefaultImageUrl)
	require.Equal(t, strings.Split(goodsCategoryTitles, ","), rsp.CategoryTitleList)
	require.WithinDuration(t, goods.CreatedAt, rsp.CreatedAt, time.Second)
	require.WithinDuration(t, goods.UpdatedAt, rsp.UpdatedAt, time.Second)

	for idx, goodsImage := range rsp.GoodsImageList {
		require.Equal(t, goodsImageList[idx].GoodsImageID, goodsImage.GoodsImageID)
		require.Equal(t, goodsImageList[idx].GoodsID, goodsImage.GoodsID)
		require.Equal(t, goodsImageList[idx].ImageUrl, goodsImage.ImageUrl)
		require.WithinDuration(t, goodsImageList[idx].CreatedAt, goodsImage.CreatedAt, time.Second)
	}
}

func requireMatchGetGoodsListResponse(t *testing.T, rsp dto.GetGoodsListResponse, goodsList []db.Good) {
	require.NotEmpty(t, rsp)

	for idx, goods := range rsp.GoodsList {
		require.Equal(t, goodsList[idx].GoodsID, goods.GoodsID)
		require.Equal(t, goodsList[idx].Title, goods.Title)
		require.Equal(t, goodsList[idx].Price, goods.Price)
		require.Equal(t, goodsList[idx].Description, goods.Description)
		require.Equal(t, goodsList[idx].DefaultImageUrl, goods.DefaultImageUrl)
		require.WithinDuration(t, goodsList[idx].CreatedAt, goods.CreatedAt, time.Second)
		require.WithinDuration(t, goodsList[idx].UpdatedAt, goods.UpdatedAt, time.Second)
	}
}
