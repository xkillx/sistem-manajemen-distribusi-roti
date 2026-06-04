-- name: CreateProduct :one
INSERT INTO products (name, sku, price, created_by)
VALUES ($1, $2, $3, $4)
RETURNING id, name, sku, price, active, created_by, created_at, updated_at;

-- name: UpdateProduct :exec
UPDATE products
SET name = COALESCE($2, name),
    price = COALESCE($3, price),
    updated_at = NOW()
WHERE id = $1;

-- name: DeactivateProduct :exec
UPDATE products
SET active = FALSE, updated_at = NOW()
WHERE id = $1;

-- name: FindProductByID :one
SELECT id, name, sku, price, active, created_by, created_at, updated_at
FROM products
WHERE id = $1;

-- name: FindProductBySKU :one
SELECT id, name, sku, price, active, created_by, created_at, updated_at
FROM products
WHERE sku = $1;

-- name: ListProducts :many
SELECT id, name, sku, price, active, created_by, created_at, updated_at
FROM products
ORDER BY active DESC, name;

-- name: ListActiveProducts :many
SELECT id, name, sku, price, active, created_by, created_at, updated_at
FROM products
WHERE active = TRUE
ORDER BY name;

-- name: ActiveProductNameExists :one
SELECT EXISTS(
    SELECT 1 FROM products
    WHERE LOWER(name) = LOWER($1) AND active = TRUE AND ($2::bigint IS NULL OR id != $2)
);

-- name: GetProductPrice :one
SELECT price FROM products WHERE id = $1;
