package db

import (
	"context"
	"testing"

	"github.com/gitaepark/carrot-market/util"
	"github.com/stretchr/testify/require"
)

func TestCreateGoodsImage(t *testing.T) {
	user, _ := createRandomUser(t)

	goods := createRandomGoods(t, user)

	createGoodsRandomImage(t, goods)
}

func TestGetGoodsImagList(t *testing.T) {
	user, _ := createRandomUser(t)

	goods := createRandomGoods(t, user)

	for i := 0; i < 10; i++ {
		createGoodsRandomImage(t, goods)
	}

	goodsImageList, err := testQueries.GetGoodsImageList(context.Background(), goods.GoodsID)
	require.NoError(t, err)
	require.NotEmpty(t, goodsImageList)

	for _, goodsImage := range goodsImageList {
		require.NotZero(t, goodsImage.GoodsImageID)
		require.Equal(t, goods.GoodsID, goodsImage.GoodsID)
		require.NotEmpty(t, goodsImage.ImageUrl)
	}
}

func TestDeleteGoodsImage(t *testing.T) {
	user, _ := createRandomUser(t)

	goods := createRandomGoods(t, user)

	goodsImage := createGoodsRandomImage(t, goods)

	err := testQueries.DeleteGoodsImage(context.Background(), goodsImage.GoodsImageID)
	require.NoError(t, err)
}

func createGoodsRandomImage(t *testing.T, goods Good) GoodsImage {
	imageUrl := util.CreateRandomString(100)

	arg := CreateGoodsImageParams{
		GoodsID:  goods.GoodsID,
		ImageUrl: imageUrl,
	}

	goodsImage, err := testQueries.CreateGoodsImage(context.Background(), arg)
	require.NoError(t, err)

	return goodsImage
}
