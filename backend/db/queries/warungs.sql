-- name: CreateWarung :one
INSERT INTO warungs (name, owner_name, phone, address, latitude, longitude, created_by, acquired_from_checkin)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, name, owner_name, phone, address, latitude, longitude, active, created_by, acquired_from_checkin, created_at, updated_at;

-- name: UpdateWarung :exec
UPDATE warungs
SET name = COALESCE($2, name),
    owner_name = COALESCE($3, owner_name),
    phone = COALESCE($4, phone),
    address = COALESCE($5, address),
    latitude = COALESCE($6, latitude),
    longitude = COALESCE($7, longitude),
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateWarungGPS :exec
UPDATE warungs
SET latitude = $2, longitude = $3, updated_at = NOW()
WHERE id = $1;

-- name: DeactivateWarung :exec
UPDATE warungs
SET active = FALSE, updated_at = NOW()
WHERE id = $1;

-- name: FindWarungByID :one
SELECT id, name, owner_name, phone, address, latitude, longitude, active, created_by, acquired_from_checkin, created_at, updated_at
FROM warungs WHERE id = $1;

-- name: ListWarungs :many
SELECT id, name, owner_name, phone, address, latitude, longitude, active, created_by, acquired_from_checkin, created_at, updated_at
FROM warungs
ORDER BY active DESC, name;

-- name: ListActiveWarungs :many
SELECT id, name, owner_name, phone, address, latitude, longitude, active, created_by, acquired_from_checkin, created_at, updated_at
FROM warungs
WHERE active = TRUE
ORDER BY name;

-- name: SearchWarungs :many
SELECT id, name, owner_name, phone, address, latitude, longitude, active, created_by, acquired_from_checkin, created_at, updated_at
FROM warungs
WHERE active = TRUE AND (LOWER(name) LIKE LOWER('%' || $1 || '%') OR LOWER(address) LIKE LOWER('%' || $1 || '%'))
ORDER BY name
LIMIT 50;

-- name: CountWarungs :one
SELECT COUNT(*)::bigint FROM warungs WHERE active = TRUE;

-- name: FindDuplicateWarungs :many
SELECT id, name, owner_name, phone, address, latitude, longitude, active, created_by, created_at, updated_at
FROM warungs
WHERE LOWER(name) = LOWER($1) AND active = TRUE AND id != $2
LIMIT 10;
