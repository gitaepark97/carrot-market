package db

import (
	"context"
	"database/sql"
)

type CreateGoodsTxParams struct {
	UserID          int32     `json:"user_id"`
	Title           string    `json:"title"`
	Price           int32     `json:"price"`
	Description     string    `json:"description"`
	DefaultImageUrl string    `json:"default_image_url"`
	CategoryIDList  []int32   `json:"category_id_list"`
	ImageUrlList    *[]string `json:"image_url_list"`
}

type CreateGoodsTxResult struct {
	Goods          GetGoodsRow  `json:"goods"`
	GoodsImageList []GoodsImage `json:"goods_image_list"`
}

func (store *SQLStore) CreateGoodsTx(ctx context.Context, arg CreateGoodsTxParams) (result CreateGoodsTxResult, err error) {
	err = store.execTx(ctx, func(q *Queries) error {
		goodsArg := CreateGoodsParams{
			UserID:          arg.UserID,
			Title:           arg.Title,
			Price:           arg.Price,
			Description:     arg.Description,
			DefaultImageUrl: arg.DefaultImageUrl,
		}

		goods, err := q.CreateGoods(ctx, goodsArg)
		if err != nil {
			return err
		}

		for _, categoryID := range arg.CategoryIDList {
			arg := CreateGoodsCategoryParams{
				GoodsID:    goods.GoodsID,
				CategoryID: categoryID,
			}

			err = q.CreateGoodsCategory(ctx, arg)
			if err != nil {
				return err
			}
		}

		for _, imageUrl := range *arg.ImageUrlList {
			arg := CreateGoodsImageParams{
				GoodsID:  goods.GoodsID,
				ImageUrl: imageUrl,
			}

			_, err = q.CreateGoodsImage(ctx, arg)
			if err != nil {
				return err
			}
		}

		err = q.UpdateGoodsOnlyUpdatedAt(ctx, goods.GoodsID)
		if err != nil {
			return err
		}

		result.Goods, err = q.GetGoods(ctx, goods.GoodsID)
		if err != nil {
			return err
		}

		result.GoodsImageList, err = q.GetGoodsImageList(ctx, goods.GoodsID)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

type UpdateGoodsTxParams struct {
	GoodsID                int32          `json:"goods_id"`
	Title                  sql.NullString `json:"title"`
	Price                  sql.NullInt32  `json:"price"`
	Description            sql.NullString `json:"description"`
	DefaultImageUrl        sql.NullString `json:"default_image_url"`
	AddCategoryIDList      *[]int32       `json:"add_category_id_list"`
	DeleteCategoryIDList   *[]int32       `json:"delete_category_id_list"`
	AddImageUrlList        *[]string      `json:"add_image_url_list"`
	DeleteGoodsImageIDList *[]int32       `json:"delete_goods_image_id_list"`
}

func (store *SQLStore) UpdateGoodsTx(ctx context.Context, arg UpdateGoodsTxParams) (err error) {
	err = store.execTx(ctx, func(q *Queries) error {
		goodsArg := UpdateGoodsParams{
			Title:           arg.Title,
			Price:           arg.Price,
			Description:     arg.Description,
			DefaultImageUrl: arg.DefaultImageUrl,
		}

		err := q.UpdateGoods(ctx, goodsArg)
		if err != nil {
			return err
		}

		if arg.DefaultImageUrl.Valid {
			arg := GetGoodsImageByUrlParams{
				GoodsID:  arg.GoodsID,
				ImageUrl: arg.DefaultImageUrl.String,
			}

			goodsImage, err := q.GetGoodsImageByUrl(ctx, arg)
			if err != nil && err != sql.ErrNoRows {
				return err
			}

			err = q.DeleteGoodsImage(ctx, goodsImage.GoodsImageID)
			if err != nil {
				return err
			}
		}

		if arg.AddCategoryIDList != nil && len(*arg.AddCategoryIDList) > 0 {
			for _, categoryID := range *arg.AddCategoryIDList {
				arg := CreateGoodsCategoryParams{
					GoodsID:    arg.GoodsID,
					CategoryID: categoryID,
				}

				err = q.CreateGoodsCategory(ctx, arg)
				if err != nil {
					return err
				}
			}
		}

		if arg.DeleteCategoryIDList != nil && len(*arg.DeleteCategoryIDList) > 0 {
			for _, categoryID := range *arg.DeleteCategoryIDList {
				arg := DeleteGoodsCategoryParams{
					GoodsID:    arg.GoodsID,
					CategoryID: categoryID,
				}

				err := q.DeleteGoodsCategory(ctx, arg)
				if err != nil {
					return err
				}
			}
		}

		if arg.AddImageUrlList != nil && len(*arg.AddImageUrlList) > 0 {
			for _, imageUrl := range *arg.AddImageUrlList {
				arg := CreateGoodsImageParams{
					GoodsID:  arg.GoodsID,
					ImageUrl: imageUrl,
				}

				_, err = q.CreateGoodsImage(ctx, arg)
				if err != nil {
					return err
				}
			}
		}

		if arg.DeleteGoodsImageIDList != nil && len(*arg.DeleteGoodsImageIDList) > 0 {
			for _, goodsImageID := range *arg.DeleteGoodsImageIDList {
				err = q.DeleteGoodsImage(ctx, goodsImageID)
				if err != nil {
					return err
				}
			}
		}

		err = q.UpdateGoodsOnlyUpdatedAt(ctx, arg.GoodsID)
		if err != nil {
			return err
		}

		return nil
	})

	return
}
