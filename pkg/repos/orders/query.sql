-- name: GetOrderByID :one
-- -- cache : 10m
SELECT
  orders.ID,
  orders.user_id,
  orders.book_id,
  orders.created_at,
  users.name AS user_name,
  users.image AS user_thumbnail,
  books.name AS book_name,
  books.price As book_price,
  books.metadata As book_metadata
FROM
  orders
  INNER JOIN books ON orders.book_id = books.id
  INNER JOIN users ON orders.user_id = users.id
WHERE
  orders.is_deleted = FALSE;

-- name: ListOrdersByUser :many
SELECT * FROM orders
WHERE
  user_id = @user_id AND created_at < @after
ORDER BY created_at DESC
LIMIT @first;

-- name: ListOrdersByUserAndBook :many
SELECT * FROM orders
WHERE
  (user_id, book_id) IN (
  SELECT
    UNNEST(@user_id::int[]),
    UNNEST(@book_id::int[])
);

-- name: BulkUpdate :exec
UPDATE orders
SET
  price=temp.price,
  book_id=temp.book_id
FROM
  (
    SELECT
      UNNEST(@id::int[]) as id,
      UNNEST(@price::bigint[]) as price,
      UNNEST(@book_id::int[]) as book_id
  ) AS temp
WHERE
  orders.id=temp.id;

-- name: CreateAuthor :one
INSERT INTO orders (
  user_id, book_id, price, is_deleted
) VALUES (
  $1, $2, $3, FALSE
)
RETURNING *;

-- name: DeleteOrder :exec
UPDATE orders
SET
  is_deleted = TRUE
WHERE
  id = $1;

-- name: ListOrdersByGender :many
-- -- cache : 1m
-- This is just an example for using type annotation for JSON field and 'with clause'.
WITH users_by_gender AS (
  SELECT * FROM users WHERE users.metadata->>'gender' = @gender::text
)
SELECT * FROM orders
WHERE
  user_id IN (SELECT id FROM users_by_gender) AND orders.id > @after
LIMIT @first;
