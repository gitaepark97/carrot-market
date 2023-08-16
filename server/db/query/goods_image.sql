-- name: CreateGoodsImage :one
INSERT INTO goods_images(
  goods_id,
  image_url
) VALUES (
  $1, $2
) RETURNING goods_images.*;

-- name: GetGoodsImageList :many
SELECT
  goods_images.*
FROM goods_images
WHERE goods_images.goods_id = $1
ORDER BY goods_images.created_at DESC;

-- name: GetGoodsImageByUrl :one
SELECT
  goods_images.*
FROM goods_images
WHERE goods_images.goods_id = $1
  AND goods_images.image_url = $2
FOR NO KEY UPDATE;

-- name: DeleteGoodsImage :exec
DELETE
FROM goods_images
WHERE goods_images.goods_image_id = $1;