-- name: Insert :exec
INSERT INTO books (
   name, description, metadata, category, price
) VALUES (
  $1, $2, $3, $4, $5
);

-- name: ListByCategory :many
SELECT *
FROM
  books
WHERE
  category = @category AND id > @after
ORDER BY
  id
LIMIT @first;

-- name: GetAllBooks :many
-- -- cache : 10m
SELECT * FROM books;

-- name: SearchBooks :many
-- -- cache : 10m
SELECT * FROM books WHERE name like $1;

-- name: GetBookByID :one
-- -- cache : 10m
SELECT * FROM books WHERE id = @id;

-- name: UpdateBookByID :exec
-- -- invalidate : [GetBookByID]
UPDATE books
SET
  description = @description, metadata = @meta, price = @price, updated_at = NOW()
WHERE
  id = @id;

-- name: PartialUpdateByID :exec
-- -- invalidate : [GetBookByID]
UPDATE books
SET
  description = coalesce(sqlc.narg('description'), description),
  metadata = coalesce(sqlc.narg('meta'), metadata),
  price = coalesce(sqlc.narg('price'), price),
  dummy_field = coalesce(sqlc.narg('dummy_field'), dummy_field),
  updated_at = NOW()
WHERE
  id = @id;