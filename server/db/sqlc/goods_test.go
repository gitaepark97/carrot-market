package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/gitaepark/carrot-market/util"
	"github.com/stretchr/testify/require"
)

func TestCreateGoods(t *testing.T) {
	user, _ := createRandomUser(t)

	createRandomGoods(t, user)
}

func TestGetGoods(t *testing.T) {
	user, _ := createRandomUser(t)

	goods1 := createRandomGoods(t, user)

	goods2, err := testQueries.GetGoods(context.Background(), goods1.GoodsID)
	require.ErrorIs(t, sql.ErrNoRows, err)
	require.Empty(t, goods2)

	createGoodsRandomCategory(t, goods1)

	goods2, err = testQueries.GetGoods(context.Background(), goods1.GoodsID)
	require.NoError(t, err)
	require.NotEmpty(t, goods2)

	require.Equal(t, goods1.GoodsID, goods2.GoodsID)
	require.Equal(t, goods1.UserID, goods2.UserID)
	require.Equal(t, goods1.Title, goods2.Title)
	require.Equal(t, goods1.Price, goods2.Price)
	require.Equal(t, goods1.Description, goods2.Description)
	require.Equal(t, goods1.DefaultImageUrl, goods2.DefaultImageUrl)
	require.NotEmpty(t, string(goods2.CategoryTitles))
	require.WithinDuration(t, goods1.CreatedAt, goods2.CreatedAt, time.Second)
	require.WithinDuration(t, goods1.UpdatedAt, goods2.UpdatedAt, time.Second)
}

func TestGetGoodsList(t *testing.T) {
	user, _ := createRandomUser(t)

	for i := 0; i < 10; i++ {
		createRandomGoods(t, user)
	}

	arg := GetGoodsListParams{
		Limit:  10,
		Offset: 0,
	}

	goodsList, err := testQueries.GetGoodsList(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, goodsList)

	require.Equal(t, 10, len(goodsList))

	for _, goods := range goodsList {
		require.NotZero(t, goods.GoodsID)
		require.NotZero(t, goods.UserID)
		require.NotEmpty(t, goods.Title)
		require.NotZero(t, goods.Price)
		require.NotEmpty(t, goods.Description)
		require.NotEmpty(t, goods.DefaultImageUrl)
		require.NotZero(t, goods.CreatedAt)
		require.NotZero(t, goods.UpdatedAt)
	}
}

func TestUpdateGoodsOnlyTitle(t *testing.T) {
	user, _ := createRandomUser(t)

	goods := createRandomGoods(t, user)

	arg := UpdateGoodsParams{
		Title: sql.NullString{
			String: util.CreateRandomString(30),
			Valid:  true,
		},
		GoodsID: goods.GoodsID,
	}

	err := testQueries.UpdateGoods(context.Background(), arg)
	require.NoError(t, err)
}

func TestUpdateGoodsOnlyPrice(t *testing.T) {
	user, _ := createRandomUser(t)

	goods := createRandomGoods(t, user)

	arg := UpdateGoodsParams{
		Price: sql.NullInt32{
			Int32: util.CreateRandomInt32(1000, 100000),
			Valid: true,
		},
		GoodsID: goods.GoodsID,
	}

	err := testQueries.UpdateGoods(context.Background(), arg)
	require.NoError(t, err)
}

func TestUpdateGoodsOnlyDescription(t *testing.T) {
	user, _ := createRandomUser(t)

	goods := createRandomGoods(t, user)

	arg := UpdateGoodsParams{
		Description: sql.NullString{
			String: util.CreateRandomString(100),
			Valid:  true,
		},
		GoodsID: goods.GoodsID,
	}

	err := testQueries.UpdateGoods(context.Background(), arg)
	require.NoError(t, err)
}

func TestUpdateGoodsOnlyDefaultImageUrl(t *testing.T) {
	user, _ := createRandomUser(t)

	goods := createRandomGoods(t, user)

	arg := UpdateGoodsParams{
		DefaultImageUrl: sql.NullString{
			String: util.CreateRandomString(100),
			Valid:  true,
		},
		GoodsID: goods.GoodsID,
	}

	err := testQueries.UpdateGoods(context.Background(), arg)
	require.NoError(t, err)
}

func TestUpdateGoodsAll(t *testing.T) {
	user, _ := createRandomUser(t)

	goods := createRandomGoods(t, user)

	arg := UpdateGoodsParams{
		Title: sql.NullString{
			String: util.CreateRandomString(30),
			Valid:  true,
		},
		Price: sql.NullInt32{
			Int32: util.CreateRandomInt32(1000, 100000),
			Valid: true,
		},
		Description: sql.NullString{
			String: util.CreateRandomString(100),
			Valid:  true,
		},
		DefaultImageUrl: sql.NullString{
			String: util.CreateRandomString(100),
			Valid:  true,
		},
		GoodsID: goods.GoodsID,
	}

	err := testQueries.UpdateGoods(context.Background(), arg)
	require.NoError(t, err)
}

func TestUpdateGoodsOnlyUpdatedAt(t *testing.T) {
	user, _ := createRandomUser(t)

	goods := createRandomGoods(t, user)

	err := testQueries.UpdateGoodsOnlyUpdatedAt(context.Background(), goods.GoodsID)
	require.NoError(t, err)
}

func createRandomGoods(t *testing.T, user User) Good {
	arg := CreateGoodsParams{
		UserID:          user.UserID,
		Title:           util.CreateRandomString(30),
		Price:           util.CreateRandomInt32(1000, 100000),
		Description:     util.CreateRandomString(100),
		DefaultImageUrl: util.CreateRandomString(100),
	}

	goods, err := testQueries.CreateGoods(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, goods)

	require.NotZero(t, goods.GoodsID)
	require.Equal(t, arg.UserID, goods.UserID)
	require.Equal(t, arg.Title, goods.Title)
	require.Equal(t, arg.Price, goods.Price)
	require.Equal(t, arg.Description, goods.Description)
	require.Equal(t, arg.DefaultImageUrl, goods.DefaultImageUrl)
	require.NotZero(t, goods.CreatedAt)
	require.NotZero(t, goods.UpdatedAt)

	return goods
}
