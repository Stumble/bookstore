-- name: GetTopItems :many
-- -- timeout : 250ms
select * from by_book_revenues
order by
  total
limit 3;

-- name: Refresh :exec
-- -- timeout : 250ms
REFRESH MATERIALIZED VIEW CONCURRENTLY by_book_revenues;
