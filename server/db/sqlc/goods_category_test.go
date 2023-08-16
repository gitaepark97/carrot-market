package db

import (
	"context"
	"testing"

	"github.com/gitaepark/carrot-market/util"
	"github.com/stretchr/testify/require"
)

func TestCreateGoodsCategory(t *testing.T) {
	user, _ := createRandomUser(t)

	goods := createRandomGoods(t, user)

	createGoodsRandomCategory(t, goods)
}

func TestDeleteGoodsCategory(t *testing.T) {
	user, _ := createRandomUser(t)

	goods := createRandomGoods(t, user)

	categoryID := createGoodsRandomCategory(t, goods)

	arg := DeleteGoodsCategoryParams{
		GoodsID:    goods.GoodsID,
		CategoryID: categoryID,
	}

	err := testQueries.DeleteGoodsCategory(context.Background(), arg)
	require.NoError(t, err)
}

func createGoodsRandomCategory(t *testing.T, goods Good) int32 {
	categoryID := util.CreateRandomInt32(1, 12)

	arg := CreateGoodsCategoryParams{
		GoodsID:    goods.GoodsID,
		CategoryID: categoryID,
	}

	err := testQueries.CreateGoodsCategory(context.Background(), arg)
	require.NoError(t, err)

	return categoryID
}
