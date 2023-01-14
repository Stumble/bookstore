-- name: ListByCategory :many
-- -- cache : 30s
SELECT *
FROM
  books
WHERE
  category = @category AND id > @after
ORDER BY
  id
LIMIT @first;
