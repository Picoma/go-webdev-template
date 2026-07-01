-- name: GetCounter :one
SELECT value
FROM counter
WHERE id = 1;

-- name: IncrementCounter :exec
UPDATE counter
SET value = value + 1
WHERE id = 1;

-- name: IncrementAndGetCounter :one
UPDATE counter
SET value = value + 1
WHERE id = 1
RETURNING value;
