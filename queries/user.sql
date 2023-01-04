-- name: GetUser :one
SELECT id, username FROM "user" WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: CountUser :one
SELECT COUNT(id) FROM "user";

-- name: ListUser :many
SELECT id, username FROM "user" ORDER BY username LIMIT($1) OFFSET($2);

-- name: CreateUser :one
INSERT INTO "user" (username, password) VALUES ($1, $2) RETURNING id, username;

-- name: UpdateUser :one
UPDATE "user"
SET username = CASE WHEN sqlc.arg(username)::VARCHAR <> ''
               THEN sqlc.arg(username)::VARCHAR
               ELSE username 
               END,
    password = CASE WHEN sqlc.arg(password)::VARCHAR <> ''
               THEN sqlc.arg(password)::VARCHAR
               ELSE password 
               END
WHERE id = $1
RETURNING id, username;

-- name: DeleteUser :exec
DELETE FROM "user" WHERE id = $1;

-- name: SoftDeleteUser :one
UPDATE "user" SET deleted_at = (current_timestamp AT TIME ZONE 'UTC') WHERE id = $1 AND deleted_at IS NULL RETURNING id, username;

-- name: GetByUsername :one
SELECT id, username, password FROM "user" WHERE username = $1 AND deleted_at IS NULL LIMIT 1;
