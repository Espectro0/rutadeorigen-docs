#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

TARGET_TIME="${1:?Uso: ./scripts/04-restaurar-pitr.sh 'YYYY-MM-DD HH:MM:SS'}"

docker compose run --rm --entrypoint sh postgres -c \
  "pgbackrest --stanza=poc-demo --config=/etc/pgbackrest/pgbackrest.conf --type=time --target='${TARGET_TIME}' --target-action=promote restore"

echo "Restauracion completada hasta ${TARGET_TIME}. Levantando Postgres..."
docker compose up -d postgres
