#!/bin/sh
# Daily PostgreSQL backup with 7-day retention for SMDR MVP

BACKUP_DIR="/var/backups/smdr"
RETENTION_DAYS=7
DB_NAME="${DB_NAME:-smdr}"
DB_USER="${DB_USER:-smdr}"
DB_PASSWORD="${DB_PASSWORD:-smdr}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/smdr_$TIMESTAMP.sql.gz"

mkdir -p "$BACKUP_DIR"

export PGPASSWORD="$DB_PASSWORD"
pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" | gzip > "$BACKUP_FILE"

find "$BACKUP_DIR" -name "smdr_*.sql.gz" -mtime +$RETENTION_DAYS -delete

echo "Backup completed: $BACKUP_FILE"
