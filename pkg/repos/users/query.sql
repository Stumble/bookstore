-- name: UpsertUsers :exec
-- -- timeout : 1s
insert into users
  (name, metadata, image)
select
        unnest(@name::VARCHAR(255)[]),
        unnest(@metadata::JSON[]),
        unnest(@image::TEXT[])
on conflict ON CONSTRAINT users_lower_name_key do
update set
    metadata = excluded.metadata,
    image = excluded.image;

-- name: GetUserByID :one
-- -- cache : 30s
-- -- timeout : 1s
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByName :one
-- -- timeout : 1s
-- -- cache : 5m
SELECT * FROM users
WHERE Name = $1 LIMIT 1;

-- name: IncorrectQuery :one
-- -- cache : 5m
-- -- timeout : 1s
SELECT * FROM users
WHERE Name = sqlc.narg('name_pointer') LIMIT 1;

-- name: ListUsers :many
-- -- timeout : 1s
SELECT * FROM users
WHERE id > @after
ORDER BY id
LIMIT @first;

-- name: UpdateNameByID :one
-- -- timeout : 1s
UPDATE users
SET
  name = $1
WHERE
  id = $2
RETURNING ID;

-- name: UpdateMetaByID :execrows
-- -- timeout : 1s
UPDATE users
SET
  metadata = $1
WHERE
  id = $2;

-- name: ListUserNames :many
-- -- timeout : 1s
SELECT id, name FROM users
WHERE id > @after
ORDER BY id
LIMIT @first;

-- name: CreateUser :one
-- -- invalidate : [GetUserByID, GetUserByName]
-- -- timeout : 1s
INSERT INTO Users (
  name, metadata, image
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteUser :exec
-- -- invalidate : [GetUserByID, GetUserByName]
-- -- timeout : 1s
DELETE FROM Users
WHERE id = $1;

-- name: UpdateUserGrade :execrows
-- -- invalidate : [GetUserByID]
-- -- timeout : 1s
UPDATE users
  SET metadata = jsonb_set(metadata, '{grade}', @grade::text, true)
WHERE
  Name LIKE @name;

-- name: DeleteBadUsers :execresult
-- -- invalidate : [GetUserByID]
-- -- timeout : 1s
DELETE FROM Users
WHERE NAME LIKE $1;

-- name: Complicated :one
-- -- cache : 1m
-- example of sqlc cannot handle recursive query.
-- -- timeout : 1s
WITH RECURSIVE fibonacci(n,x,y) AS (
	SELECT
    	1 AS n ,
  		0 :: int AS x,
    	1 :: int AS y
  	UNION ALL
  	SELECT
    	n + 1 AS n,
  		y AS x,
    	x + y AS y
  	FROM fibonacci
  	WHERE n < @n::int
	)
SELECT
	x
FROM fibonacci;
