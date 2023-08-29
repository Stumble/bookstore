-- name: Insert :exec
INSERT INTO books (
   name, description, metadata, category, price
) VALUES (
  $1, $2, $3, $4, $5
);

-- name: BulkInsertByCopyfrom :copyfrom
INSERT INTO books (
   name, description, metadata, category, price
) VALUES (
  $1, $2, $3, $4, $5
);

-- name: InsertAndReturnID :one
INSERT INTO books (
   name, description, metadata, category, price
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING id;


-- name: RefreshIDSerial :exec
SELECT setval(seq_name, (SELECT MAX(id) FROM books)+1, false)
FROM PG_GET_SERIAL_SEQUENCE('books', 'id') as seq_name;

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

-- name: GetBookBySpec :many
-- -- cache : 10m
SELECT * FROM books WHERE
  name LIKE coalesce(sqlc.narg('name'), name) AND
  price = coalesce(sqlc.narg('price'), price) AND
  (sqlc.narg('dummy')::int is NULL or dummy_field = sqlc.narg('dummy'));

-- name: GetBookByNameMaybe :many
-- -- cache : 10m
SELECT * FROM books WHERE
  name LIKE coalesce(sqlc.narg('name'), name);

-- name: UpdateBookByID :exec
-- -- invalidate : [GetBookByID]
UPDATE books
SET
  description = @description, metadata = @meta, price = @price, updated_at = NOW()
WHERE
  id = @id;

-- name: PartialUpdateByID :exec
UPDATE books
SET
  description = coalesce(sqlc.narg('description'), description),
  metadata = coalesce(sqlc.narg('meta'), metadata),
  price = coalesce(sqlc.narg('price'), price),
  dummy_field = coalesce(sqlc.narg('dummy_field'), dummy_field),
  updated_at = NOW()
WHERE
  id = @id;

-- name: InsertWithInvalidate :exec
-- -- invalidate : [GetBookByNameMaybe, GetBookBySpec]
INSERT INTO books (
   id, name, description, metadata, category, dummy_field, price
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
);
