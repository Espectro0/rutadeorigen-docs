#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

docker compose exec postgres pgbackrest --stanza=poc-demo --config=/etc/pgbackrest/pgbackrest.conf stanza-create
docker compose exec postgres pgbackrest --stanza=poc-demo --config=/etc/pgbackrest/pgbackrest.conf check

echo "Stanza creada y verificada contra MinIO."
