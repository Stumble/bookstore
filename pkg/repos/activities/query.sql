-- name: Insert :exec
-- -- timeout : 1s
INSERT INTO activities (
   action, parameter, created_at
) VALUES (
  $1, $2, NOW()
);
