#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

mkdir -p certs
openssl req -x509 -nodes -days 365 \
  -newkey rsa:2048 \
  -keyout certs/private.key \
  -out certs/public.crt \
  -subj "/CN=minio" \
  -addext "subjectAltName=DNS:minio,DNS:localhost"

echo "Certificado autofirmado generado en certs/ (solo para este PoC, MinIO lo detecta automaticamente al arrancar)."
