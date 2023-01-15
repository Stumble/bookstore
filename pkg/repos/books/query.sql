-- name: BulkInsert :copyfrom
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
