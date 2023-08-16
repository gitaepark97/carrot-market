// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: goods_image.sql

package db

import (
	"context"
)

const createGoodsImage = `-- name: CreateGoodsImage :one
INSERT INTO goods_images(
  goods_id,
  image_url
) VALUES (
  $1, $2
) RETURNING goods_images.goods_image_id, goods_images.goods_id, goods_images.image_url, goods_images.created_at
`

type CreateGoodsImageParams struct {
	GoodsID  int32  `json:"goods_id"`
	ImageUrl string `json:"image_url"`
}

func (q *Queries) CreateGoodsImage(ctx context.Context, arg CreateGoodsImageParams) (GoodsImage, error) {
	row := q.db.QueryRowContext(ctx, createGoodsImage, arg.GoodsID, arg.ImageUrl)
	var i GoodsImage
	err := row.Scan(
		&i.GoodsImageID,
		&i.GoodsID,
		&i.ImageUrl,
		&i.CreatedAt,
	)
	return i, err
}

const deleteGoodsImage = `-- name: DeleteGoodsImage :exec
DELETE
FROM goods_images
WHERE goods_images.goods_image_id = $1
`

func (q *Queries) DeleteGoodsImage(ctx context.Context, goodsImageID int32) error {
	_, err := q.db.ExecContext(ctx, deleteGoodsImage, goodsImageID)
	return err
}

const getGoodsImageByUrl = `-- name: GetGoodsImageByUrl :one
SELECT
  goods_images.goods_image_id, goods_images.goods_id, goods_images.image_url, goods_images.created_at
FROM goods_images
WHERE goods_images.goods_id = $1
  AND goods_images.image_url = $2
FOR NO KEY UPDATE
`

type GetGoodsImageByUrlParams struct {
	GoodsID  int32  `json:"goods_id"`
	ImageUrl string `json:"image_url"`
}

func (q *Queries) GetGoodsImageByUrl(ctx context.Context, arg GetGoodsImageByUrlParams) (GoodsImage, error) {
	row := q.db.QueryRowContext(ctx, getGoodsImageByUrl, arg.GoodsID, arg.ImageUrl)
	var i GoodsImage
	err := row.Scan(
		&i.GoodsImageID,
		&i.GoodsID,
		&i.ImageUrl,
		&i.CreatedAt,
	)
	return i, err
}

const getGoodsImageList = `-- name: GetGoodsImageList :many
SELECT
  goods_images.goods_image_id, goods_images.goods_id, goods_images.image_url, goods_images.created_at
FROM goods_images
WHERE goods_images.goods_id = $1
ORDER BY goods_images.created_at DESC
`

func (q *Queries) GetGoodsImageList(ctx context.Context, goodsID int32) ([]GoodsImage, error) {
	rows, err := q.db.QueryContext(ctx, getGoodsImageList, goodsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GoodsImage{}
	for rows.Next() {
		var i GoodsImage
		if err := rows.Scan(
			&i.GoodsImageID,
			&i.GoodsID,
			&i.ImageUrl,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}