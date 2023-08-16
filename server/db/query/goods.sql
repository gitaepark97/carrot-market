-- name: CreateGoods :one
INSERT INTO goods(
  user_id,
  title,
  price,
  description,
  default_image_url
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING goods.*;

-- name: GetGoods :one
SELECT
  goods.*,
  string_agg(categories.title, ',' ORDER BY categories.created_at DESC) as category_titles
FROM goods
JOIN goods_categories USING(goods_id)
JOIN categories USING(category_id)
WHERE goods.goods_id = $1
GROUP BY goods.goods_id;

-- name: GetGoodsList :many
SELECT
  goods.*
FROM goods
ORDER BY goods.created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateGoods :exec
UPDATE goods
SET
  title = coalesce(sqlc.narg(title), title),
  price = coalesce(sqlc.narg(price), price),
  description = coalesce(sqlc.narg(description), description),
  default_image_url = coalesce(sqlc.narg(default_image_url), default_image_url),
  updated_at = now()
WHERE goods_id = sqlc.arg(goods_id);

-- name: UpdateGoodsOnlyUpdatedAt :exec
UPDATE goods
SET
  updated_at = now()
WHERE goods_id = $1;

-- name: DeleteGoods :exec
DELETE
FROM goods
WHERE goods.goods_id = $1;