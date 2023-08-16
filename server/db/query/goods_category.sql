-- name: CreateGoodsCategory :exec
INSERT INTO goods_categories(
  goods_id,
  category_id
) VALUES (
  $1, $2
);

-- name: DeleteGoodsCategory :exec
DELETE
FROM goods_categories
WHERE goods_categories.goods_id = $1
  AND goods_categories.category_id = $2;