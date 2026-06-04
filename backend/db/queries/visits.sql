-- name: CreateVisit :one
INSERT INTO visits (sales_id, warung_id, check_in_lat, check_in_lng, check_in_accuracy, check_in_time, business_date)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, sales_id, warung_id, status, check_in_lat, check_in_lng, check_in_accuracy, check_in_time, distance_meters, business_date, notes, cancelled_at, cancelled_by, cancel_reason, created_at, updated_at;

-- name: FindVisitByID :one
SELECT id, sales_id, warung_id, status, check_in_lat, check_in_lng, check_in_accuracy, check_in_time, distance_meters, business_date, notes, cancelled_at, cancelled_by, cancel_reason, created_at, updated_at
FROM visits
WHERE id = $1;

-- name: UpdateVisitStatus :exec
UPDATE visits
SET status = $2, updated_at = NOW()
WHERE id = $1;

-- name: FinalizeVisit :exec
UPDATE visits
SET status = 'selesai', notes = $2, updated_at = NOW()
WHERE id = $1;

-- name: CancelVisit :exec
UPDATE visits
SET status = 'dibatalkan', cancelled_at = NOW(), cancelled_by = $2, cancel_reason = $3, updated_at = NOW()
WHERE id = $1;

-- name: UpdateVisitDistance :exec
UPDATE visits
SET distance_meters = $2
WHERE id = $1;

-- name: FindActiveDraft :one
SELECT id, sales_id, warung_id, status, check_in_lat, check_in_lng, check_in_accuracy, check_in_time, distance_meters, business_date, notes, cancelled_at, cancelled_by, cancel_reason, created_at, updated_at
FROM visits
WHERE sales_id = $1 AND status = 'draft';

-- name: AutoCancelStaleDrafts :exec
UPDATE visits
SET status = 'dibatalkan', cancelled_at = NOW(), cancel_reason = 'Auto-cancelled: tidak diselesaikan sebelum akhir Hari Bisnis', updated_at = NOW()
WHERE status = 'draft' AND business_date < (CURRENT_DATE AT TIME ZONE 'Asia/Jakarta');

-- name: ListVisitsBySales :many
SELECT id, sales_id, warung_id, status, check_in_lat, check_in_lng, check_in_accuracy, check_in_time, distance_meters, business_date, notes, cancelled_at, cancelled_by, cancel_reason, created_at, updated_at
FROM visits
WHERE sales_id = $1
ORDER BY created_at DESC
LIMIT 50;

-- name: ListVisitsBySalesAndDate :many
SELECT id, sales_id, warung_id, status, check_in_lat, check_in_lng, check_in_accuracy, check_in_time, distance_meters, business_date, notes, cancelled_at, cancelled_by, cancel_reason, created_at, updated_at
FROM visits
WHERE sales_id = $1 AND business_date >= $2 AND business_date <= $3
ORDER BY created_at DESC;

-- name: ListVisitsByWarung :many
SELECT id, sales_id, warung_id, status, check_in_lat, check_in_lng, check_in_accuracy, check_in_time, distance_meters, business_date, notes, cancelled_at, cancelled_by, cancel_reason, created_at, updated_at
FROM visits
WHERE warung_id = $1 AND status = 'selesai'
ORDER BY business_date, created_at;

-- name: ListRecentVisits :many
SELECT v.id, v.sales_id, v.warung_id, v.status, v.check_in_lat, v.check_in_lng, v.check_in_accuracy, v.check_in_time, v.distance_meters, v.business_date, v.notes, v.cancelled_at, v.cancelled_by, v.cancel_reason, v.created_at, v.updated_at,
       u.name AS sales_name, w.name AS warung_name
FROM visits v
JOIN users u ON v.sales_id = u.id
JOIN warungs w ON v.warung_id = w.id
WHERE v.business_date = (CURRENT_DATE AT TIME ZONE 'Asia/Jakarta')
ORDER BY v.check_in_time DESC
LIMIT 20;

-- name: CorrectVisit :exec
UPDATE visits
SET notes = COALESCE($2, notes), updated_at = NOW()
WHERE id = $1;
