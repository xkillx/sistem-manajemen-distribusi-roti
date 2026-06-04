-- name: FindUserByUsername :one
SELECT id, username, hashed_password, name, phone, role, active, must_change_password, created_at, updated_at
FROM users
WHERE LOWER(username) = LOWER(sqlc.arg('username'));

-- name: FindUserByID :one
SELECT id, username, hashed_password, name, phone, role, active, must_change_password, created_at, updated_at
FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (username, hashed_password, name, phone, role, must_change_password)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, username, hashed_password, name, phone, role, active, must_change_password, created_at, updated_at;

-- name: UpdateUserPassword :exec
UPDATE users
SET hashed_password = $2, must_change_password = $3, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users
SET name = COALESCE($2, name),
    phone = COALESCE($3, phone),
    updated_at = NOW()
WHERE id = $1;

-- name: DeactivateUser :exec
UPDATE users
SET active = FALSE, updated_at = NOW()
WHERE id = $1;

-- name: ListSales :many
SELECT id, username, name, phone, role, active, must_change_password, created_at, updated_at
FROM users
WHERE role = 'sales'
ORDER BY created_at DESC;

-- name: ListActiveSales :many
SELECT id, username, name, phone, role, active, must_change_password, created_at, updated_at
FROM users
WHERE role = 'sales' AND active = TRUE
ORDER BY name;

-- name: CountUsers :one
SELECT COUNT(*)::bigint FROM users;

-- name: CountActiveSales :one
SELECT COUNT(*)::bigint FROM users WHERE role = 'sales' AND active = TRUE;

-- name: UserExistsByUsername :one
SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) = LOWER($1));
