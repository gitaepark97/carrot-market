-- name: GetCategoryList :many
SELECT
  categories.*
FROM categories
ORDER BY created_at;