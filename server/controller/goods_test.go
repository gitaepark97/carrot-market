package controller

import (
	"fmt"
	"testing"
	"time"

	db "github.com/gitaepark/carrot-market/db/sqlc"
	"github.com/gitaepark/carrot-market/util"
)

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
