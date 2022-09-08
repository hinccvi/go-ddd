-- name: GetUser :one
SELECT * FROM "user"
WHERE id = $1 LIMIT 1;

-- name: CountUser :one
SELECT COUNT(*) FROM "user";

-- name: ListUser :many
SELECT * FROM "user"
ORDER BY name;

-- name: CreateUser :one
INSERT INTO "user" (
  username, password
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteUser :one
DELETE FROM "user"
WHERE id = $1
RETURNING *;

-- name: GetByUsernameAndPassword :one
SELECT * FROM "user"
WHERE username = $1 AND password = $2 LIMIT 1;