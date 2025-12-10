#!/bin/bash
set -e

BACKUP_NAME="backup-$(date +%Y-%m-%d_%H-%M-%S).sql"

echo "Activating Service Account..."
gcloud auth activate-service-account --key-file="$GOOGLE_APPLICATION_CREDENTIALS"

echo "Dumping database..."
pg_dump "$POSTGRES_URL" > "$BACKUP_NAME"

echo "Uploading to GCS Bucket: $GCS_BUCKET..."
gsutil cp "$BACKUP_NAME" "gs://$GCS_BUCKET/$BACKUP_NAME"

echo "Backup done."
