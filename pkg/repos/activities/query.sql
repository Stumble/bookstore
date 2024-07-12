-- name: Insert :exec
-- -- timeout : 1s
-- -- invalidate : [GetActivitiesByAction]
INSERT INTO activities (
   action, parameter, created_at
) VALUES (
  $1, $2, NOW()
);

-- name: GetActivitiesByAction :many
-- -- timeout : 250ms
-- -- cache : 1m
select * from activities where action = @action;
