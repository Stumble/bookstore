-- name: Insert :exec
INSERT INTO activities (
   action, parameter, created_at
) VALUES (
  $1, $2, NOW()
);
