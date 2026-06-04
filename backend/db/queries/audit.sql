-- name: CreateAuditLog :exec
INSERT INTO audit_logs (actor_id, action, entity, entity_id, summary)
VALUES ($1, $2, $3, $4, $5);

-- name: ListAuditLogs :many
SELECT al.id, al.actor_id, al.action, al.entity, al.entity_id, al.summary, al.created_at,
       u.name AS actor_name
FROM audit_logs al
JOIN users u ON al.actor_id = u.id
ORDER BY al.created_at DESC
LIMIT 100;

-- name: ListAuditLogsByEntity :many
SELECT al.id, al.actor_id, al.action, al.entity, al.entity_id, al.summary, al.created_at,
       u.name AS actor_name
FROM audit_logs al
JOIN users u ON al.actor_id = u.id
WHERE al.entity = $1 AND al.entity_id = $2
ORDER BY al.created_at DESC;
