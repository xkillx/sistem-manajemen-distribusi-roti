-- name: GetDashboardMetrics :one
SELECT
    (SELECT COUNT(*)::bigint FROM warungs WHERE active = TRUE) AS total_warung,
    (SELECT COUNT(*)::bigint FROM users WHERE role = 'sales' AND active = TRUE) AS total_sales,
    COALESCE((SELECT SUM(vsi.quantity)::bigint
     FROM visit_sale_items vsi
     JOIN visits v ON vsi.visit_id = v.id
     WHERE v.status = 'selesai' AND v.business_date = (CURRENT_DATE AT TIME ZONE 'Asia/Jakarta')), 0) AS pcs_sold_today,
    COALESCE((SELECT SUM(vri.quantity)::bigint
     FROM visit_return_items vri
     JOIN visits v ON vri.visit_id = v.id
     WHERE v.status = 'selesai' AND v.business_date = (CURRENT_DATE AT TIME ZONE 'Asia/Jakarta')), 0) AS pcs_returned_today,
    COALESCE((SELECT SUM(vp.amount)::bigint
     FROM visit_payments vp
     JOIN visits v ON vp.visit_id = v.id
     WHERE v.status = 'selesai' AND v.business_date = (CURRENT_DATE AT TIME ZONE 'Asia/Jakarta')), 0) AS total_pemasukan_today;

-- name: GetDashboardRecentActivity :many
SELECT v.id, v.check_in_time, v.distance_meters, v.status, v.notes,
       u.name AS sales_name,
       w.name AS warung_name,
       COALESCE(payment.total, 0)::bigint AS pembayaran,
       COALESCE(sales.nilai, 0)::bigint AS nilai_penjualan
FROM visits v
JOIN users u ON v.sales_id = u.id
JOIN warungs w ON v.warung_id = w.id
LEFT JOIN (
    SELECT vp.visit_id, vp.amount AS total
    FROM visit_payments vp
) payment ON v.id = payment.visit_id
LEFT JOIN (
    SELECT vsi.visit_id, SUM(vsi.quantity * vsi.price) AS nilai
    FROM visit_sale_items vsi
    GROUP BY vsi.visit_id
) sales ON v.id = sales.visit_id
WHERE v.business_date = (CURRENT_DATE AT TIME ZONE 'Asia/Jakarta')
ORDER BY v.check_in_time DESC
LIMIT 20;

-- name: GetSalesReport :many
SELECT
    v.business_date,
    u.id AS sales_id,
    u.name AS sales_name,
    w.id AS warung_id,
    w.name AS warung_name,
    p.id AS product_id,
    p.name AS product_name,
    p.sku AS product_sku,
    vsi.quantity,
    vsi.price,
    (vsi.quantity * vsi.price)::bigint AS nilai
FROM visit_sale_items vsi
JOIN visits v ON vsi.visit_id = v.id
JOIN users u ON v.sales_id = u.id
JOIN warungs w ON v.warung_id = w.id
JOIN products p ON vsi.product_id = p.id
WHERE v.status = 'selesai'
  AND v.business_date >= $1 AND v.business_date <= $2
ORDER BY v.business_date, u.name, w.name, p.name;

-- name: GetSalesReportTotal :one
SELECT
    COALESCE(SUM(vsi.quantity), 0)::bigint AS total_pcs,
    COALESCE(SUM(vsi.quantity * vsi.price), 0)::bigint AS total_nilai
FROM visit_sale_items vsi
JOIN visits v ON vsi.visit_id = v.id
WHERE v.status = 'selesai'
  AND v.business_date >= $1 AND v.business_date <= $2;

-- name: GetReturnsReport :many
SELECT
    v.business_date,
    u.id AS sales_id,
    u.name AS sales_name,
    w.id AS warung_id,
    w.name AS warung_name,
    p.id AS product_id,
    p.name AS product_name,
    vri.quantity,
    vri.reason
FROM visit_return_items vri
JOIN visits v ON vri.visit_id = v.id
JOIN users u ON v.sales_id = u.id
JOIN warungs w ON v.warung_id = w.id
JOIN products p ON vri.product_id = p.id
WHERE v.status = 'selesai'
  AND v.business_date >= $1 AND v.business_date <= $2
  AND ($3::bigint = 0 OR p.id = $3)
  AND ($4::bigint = 0 OR w.id = $4)
  AND ($5::bigint = 0 OR u.id = $5)
ORDER BY v.business_date, u.name, w.name, p.name;

-- name: GetWarungReport :many
SELECT
    w.id AS warung_id,
    w.name AS warung_name,
    COALESCE(SUM(deposits.quantity), 0)::bigint AS total_titipan,
    COALESCE(SUM(sales.quantity), 0)::bigint AS total_terjual,
    COALESCE(SUM(retur.quantity), 0)::bigint AS total_retur,
    COALESCE(SUM(sales.nilai), 0)::bigint AS nilai_penjualan,
    COALESCE(SUM(payments.amount), 0)::bigint AS pembayaran
FROM warungs w
LEFT JOIN visits v ON w.id = v.warung_id AND v.status = 'selesai' AND v.business_date >= $1 AND v.business_date <= $2
LEFT JOIN (
    SELECT visit_id, SUM(quantity)::bigint AS quantity
    FROM visit_deposit_items
    GROUP BY visit_id
) deposits ON v.id = deposits.visit_id
LEFT JOIN (
    SELECT vsi.visit_id, SUM(vsi.quantity)::bigint AS quantity, SUM(vsi.quantity * vsi.price)::bigint AS nilai
    FROM visit_sale_items vsi
    GROUP BY vsi.visit_id
) sales ON v.id = sales.visit_id
LEFT JOIN (
    SELECT visit_id, SUM(quantity)::bigint AS quantity
    FROM visit_return_items
    GROUP BY visit_id
) retur ON v.id = retur.visit_id
LEFT JOIN visit_payments payments ON v.id = payments.visit_id
WHERE v.id IS NOT NULL
GROUP BY w.id, w.name
ORDER BY w.name;

-- name: GetSalespersonReport :many
SELECT
    u.id AS sales_id,
    u.name AS sales_name,
    COUNT(v.id)::bigint AS total_kunjungan,
    COUNT(v.id) FILTER (WHERE EXISTS (
        SELECT 1 FROM visit_deposit_items WHERE visit_id = v.id
        UNION ALL
        SELECT 1 FROM visit_sale_items WHERE visit_id = v.id
        UNION ALL
        SELECT 1 FROM visit_return_items WHERE visit_id = v.id
        UNION ALL
        SELECT 1 FROM visit_payments WHERE visit_id = v.id
    ))::bigint AS kunjungan_dengan_transaksi,
    COALESCE(SUM(sales.quantity), 0)::bigint AS produk_terjual,
    COALESCE(SUM(sales.nilai), 0)::bigint AS nilai_penjualan,
    COALESCE(SUM(payments.amount), 0)::bigint AS pembayaran,
    COALESCE(SUM(retur.quantity), 0)::bigint AS retur
FROM users u
LEFT JOIN visits v ON u.id = v.sales_id AND v.status = 'selesai' AND v.business_date >= $1 AND v.business_date <= $2
LEFT JOIN (
    SELECT vsi.visit_id, SUM(vsi.quantity)::bigint AS quantity, SUM(vsi.quantity * vsi.price)::bigint AS nilai
    FROM visit_sale_items vsi
    GROUP BY vsi.visit_id
) sales ON v.id = sales.visit_id
LEFT JOIN visit_payments payments ON v.id = payments.visit_id
LEFT JOIN (
    SELECT visit_id, SUM(quantity)::bigint AS quantity
    FROM visit_return_items
    GROUP BY visit_id
) retur ON v.id = retur.visit_id
GROUP BY u.id, u.name
HAVING u.role = 'sales' AND COUNT(v.id) > 0
ORDER BY u.name;
