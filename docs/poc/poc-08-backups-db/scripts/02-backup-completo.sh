#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

docker compose exec postgres pgbackrest --stanza=poc-demo --config=/etc/pgbackrest/pgbackrest.conf --type=full backup

echo "Backup completo (base) tomado y subido a MinIO."
