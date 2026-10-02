-- name: CreateUser :one
INSERT INTO users (
    id,
    created_at,
    updated_at,
    email,
    hash)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1; 

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1; 

-- name: UpdateUser :one
UPDATE users
SET email = $2, hash = $3
WHERE id = $1
RETURNING *;

-- name: UpgradeUser :one
UPDATE users
SET is_chirpy_red = true
WHERE id = $1
RETURNING *;