-- name: GetUser :one
SELECT * FROM "user"
WHERE id = $1 LIMIT 1;

-- name: CountUser :one
SELECT COUNT(*) FROM "user";

-- name: ListUser :many
SELECT * FROM "user"
ORDER BY username
LIMIT($1)
OFFSET($2);

-- name: CreateUser :one
INSERT INTO "user" (
  username, password
) VALUES (
  $1, $2
)
RETURNING id, username, created_at, updated_at;

-- name: UpdateUser :one
UPDATE "user"
SET username = CASE WHEN sqlc.arg(username)::VARCHAR <> ''
               THEN sqlc.arg(username)::VARCHAR
               ELSE username END,
    password = CASE WHEN sqlc.arg(password)::VARCHAR <> ''
               THEN sqlc.arg(password)::VARCHAR
               ELSE password END
WHERE id = $1
RETURNING id, username, created_at, updated_at;

-- name: DeleteUser :one
DELETE FROM "user"
WHERE id = $1
RETURNING *;

-- name: GetByUsernameAndPassword :one
SELECT * FROM "user"
WHERE username = $1 AND password = $2 LIMIT 1;