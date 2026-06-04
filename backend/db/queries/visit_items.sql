-- name: CreateDepositItem :exec
INSERT INTO visit_deposit_items (visit_id, product_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (visit_id, product_id) DO UPDATE SET quantity = $3;

-- name: FindDepositItems :many
SELECT vdi.id, vdi.visit_id, vdi.product_id, vdi.quantity, p.name AS product_name
FROM visit_deposit_items vdi
JOIN products p ON vdi.product_id = p.id
WHERE vdi.visit_id = $1
ORDER BY p.name;

-- name: DeleteDepositItems :exec
DELETE FROM visit_deposit_items WHERE visit_id = $1;

-- name: CreateSaleItem :exec
INSERT INTO visit_sale_items (visit_id, product_id, quantity, price)
VALUES ($1, $2, $3, $4)
ON CONFLICT (visit_id, product_id) DO UPDATE SET quantity = $3, price = $4;

-- name: FindSaleItems :many
SELECT vsi.id, vsi.visit_id, vsi.product_id, vsi.quantity, vsi.price, p.name AS product_name
FROM visit_sale_items vsi
JOIN products p ON vsi.product_id = p.id
WHERE vsi.visit_id = $1
ORDER BY p.name;

-- name: DeleteSaleItems :exec
DELETE FROM visit_sale_items WHERE visit_id = $1;

-- name: CreateReturnItem :exec
INSERT INTO visit_return_items (visit_id, product_id, quantity, reason)
VALUES ($1, $2, $3, $4)
ON CONFLICT (visit_id, product_id) DO UPDATE SET quantity = $3, reason = $4;

-- name: FindReturnItems :many
SELECT vri.id, vri.visit_id, vri.product_id, vri.quantity, vri.reason, p.name AS product_name
FROM visit_return_items vri
JOIN products p ON vri.product_id = p.id
WHERE vri.visit_id = $1
ORDER BY p.name;

-- name: DeleteReturnItems :exec
DELETE FROM visit_return_items WHERE visit_id = $1;

-- name: UpsertPayment :exec
INSERT INTO visit_payments (visit_id, amount, method, note)
VALUES ($1, $2, $3, $4)
ON CONFLICT (visit_id) DO UPDATE SET amount = $2, method = $3, note = $4;

-- name: FindPayment :one
SELECT id, visit_id, amount, method, note
FROM visit_payments
WHERE visit_id = $1;

-- name: DeletePayment :exec
DELETE FROM visit_payments WHERE visit_id = $1;
