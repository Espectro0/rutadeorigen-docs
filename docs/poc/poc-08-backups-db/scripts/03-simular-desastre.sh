#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

echo "Deteniendo Postgres para simular el desastre..."
docker compose stop postgres

echo "Borrando el data directory para simular la perdida total del disco..."
docker compose run --rm --entrypoint sh postgres -c "rm -rf /var/lib/postgresql/data/*"

echo "Data directory borrado. Ahora corre 04-restaurar-pitr.sh con el timestamp objetivo."
