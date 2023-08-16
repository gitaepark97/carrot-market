package dto

import (
	"strings"
	"time"

	db "github.com/gitaepark/carrot-market/db/sqlc"
)

type CreateGoodsBodyRequest struct {
	Title           string    `json:"title" binding:"required,max=50"`
	Price           int32     `json:"price" binding:"required,gte=0"`
	Description     string    `json:"description" binding:"required"`
	CategoryIDList  []int32   `json:"category_id_list" binding:"required"`
	DefaultImageUrl string    `json:"default_image_url" binding:"required"`
	ImageUrlList    *[]string `json:"image_url_list" binding:"omitempty"`
}

type GoodsResponse struct {
	GoodsID           int32           `json:"goods_id"`
	Title             string          `json:"title"`
	Price             int32           `json:"price"`
	Description       string          `json:"description"`
	DefaultImageUrl   string          `json:"default_image_url"`
	CategoryTitleList []string        `json:"category_title_list"`
	GoodsImageList    []db.GoodsImage `json:"goods_image_list"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

func NewGoodsResponse(goods db.GetGoodsRow, goodsImageList []db.GoodsImage) GoodsResponse {
	return GoodsResponse{
		GoodsID:           goods.GoodsID,
		Title:             goods.Title,
		Price:             goods.Price,
		Description:       goods.Description,
		DefaultImageUrl:   goods.DefaultImageUrl,
		CategoryTitleList: strings.Split(string(goods.CategoryTitles), ","),
		GoodsImageList:    goodsImageList,
		CreatedAt:         goods.CreatedAt,
		UpdatedAt:         goods.UpdatedAt,
	}
}

type GetGoodsListQueryRequest struct {
	PageID   int32 `form:"page_id" binding:"required,gte=1"`
	PageSize int32 `form:"page_size" binding:"required,gte=10"`
}

type GetGoodsListResponse struct {
	GoodsList []db.Good `json:"goods_list"`
}

func NewGetGoodsListResponse(goodsList []db.Good) GetGoodsListResponse {
	return GetGoodsListResponse{
		GoodsList: goodsList,
	}
}

type GoodsPathRequest struct {
	GoodsID int32 `uri:"goods_id" binding:"required,gte=1"`
}

type UpdateGoodsBodyRequest struct {
	Title                  string    `json:"title" binding:"omitempty,max=50"`
	Price                  int32     `json:"price" binding:"omitempty,gte=0"`
	Description            string    `json:"description" binding:"omitempty"`
	DefaultImageUrl        string    `json:"default_image_url" binding:"omitempty"`
	AddCategoryIDList      *[]int32  `json:"add_category_id_list" binding:"omitempty"`
	DeleteCategoryIDList   *[]int32  `json:"delete_category_id_list" binding:"omitempty"`
	AddImageUrlList        *[]string `json:"add_image_url_list" binding:"omitempty"`
	DeleteGoodsImageIDList *[]int32  `json:"delete_goods_image_id_list" binding:"omitempty"`
}
