-- name: GetTopItems :many
select * from by_book_revenues
order by
  total
limit 3;

-- name: Refresh :exec
REFRESH MATERIALIZED VIEW CONCURRENTLY by_book_revenues;
