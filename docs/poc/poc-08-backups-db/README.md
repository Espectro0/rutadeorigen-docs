# PoC 08: Respaldo y Recuperación de Postgres con pgBackRest hacia MinIO

Prueba de concepto (PoC) que implementa la estrategia de **respaldo físico continuo con archivado de WAL (Point-in-Time Recovery)** definida en el ADR-06, usando **pgBackRest** contra **MinIO** — el mismo almacenamiento de objetos elegido en el ADR-07 — como repositorio S3. Levanta Postgres y MinIO por Docker, toma un backup base, sigue insertando datos con un programa en **Go**, simula la pérdida total del disco, y restaura la base de datos a un punto exacto en el tiempo posterior al backup base, para confirmar que se recupera más allá del último snapshot, tal como exige el ADR-06.

## Arquitectura y Justificación

pgBackRest corre instalado dentro del mismo contenedor de Postgres, porque necesita acceso directo al `PGDATA` para leer y archivar cada segmento de WAL a medida que se completa (`archive_mode=on` + `archive_command` apuntando a `pgbackrest archive-push`), sin esperar al siguiente backup base. MinIO actúa como repositorio S3 remoto, reutilizando el mismo mecanismo de almacenamiento de objetos que ya se decidió para las evidencias en el ADR-07, en vez de introducir un segundo backend de almacenamiento solo para los backups.

El driver S3 de pgBackRest solo habla TLS, así que MinIO se levanta con un certificado autofirmado generado localmente (`scripts/00-generar-certificados.sh`); para este PoC se desactiva la verificación del certificado (`repo1-storage-verify-tls=n`) en vez de gestionar una CA de prueba, algo que en producción sí debería resolverse con un certificado válido o una CA propia confiable.

### Estructura del proyecto

```
poc-08-backups-db/
├── cmd/
│   └── main.go                    # Programa Go: siembra filas de prueba y las lista, para generar y verificar actividad de WAL
├── scripts/
│   ├── 00-generar-certificados.sh # Genera el certificado autofirmado que usa MinIO
│   ├── 01-crear-stanza.sh         # Crea y verifica la stanza de pgBackRest contra MinIO
│   ├── 02-backup-completo.sh      # Toma el backup base completo
│   ├── 03-simular-desastre.sh     # Detiene Postgres y borra el data directory (simula el desastre)
│   └── 04-restaurar-pitr.sh       # Restaura a un punto en el tiempo específico
├── docker-compose.yml             # Levanta MinIO (+ init del bucket) y Postgres con pgBackRest instalado
├── Dockerfile.postgres            # Imagen de Postgres 16 con pgBackRest instalado
├── pgbackrest.conf                # Configuración de la stanza y del repositorio S3 (MinIO)
├── go.mod
└── (sin go.sum hasta correr `go mod tidy`)
```

## Guía de Levantamiento

### 1. Generar el certificado de MinIO

```bash
./scripts/00-generar-certificados.sh
```

### 2. Levantar los contenedores

```bash
docker compose up -d --build
```

Esto levanta MinIO, crea el bucket `pgbackrest` (contenedor `minio-init`), y construye y levanta Postgres con `archive_mode` activado apuntando a pgBackRest.

* Consola de MinIO: https://localhost:9001 (`minioadmin` / `minioadmin123`, certificado autofirmado).

### 3. Crear la stanza de pgBackRest

```bash
./scripts/01-crear-stanza.sh
```

### 4. Sembrar datos base y tomar el backup completo

```bash
go get github.com/jackc/pgx/v5
go mod tidy

go run ./cmd seed 3
./scripts/02-backup-completo.sh
```

### 5. Insertar datos después del backup y anotar el punto de recuperación

```bash
go run ./cmd seed 3   # filas 4-6: estas SÍ deben sobrevivir al restore
date -u +"%Y-%m-%d %H:%M:%S"   # anota este timestamp (UTC) como punto objetivo

go run ./cmd seed 2   # filas 7-8: estas NO deben aparecer tras el restore
```

### 6. Simular el desastre

```bash
./scripts/03-simular-desastre.sh
```

### 7. Restaurar al punto en el tiempo anotado

```bash
./scripts/04-restaurar-pitr.sh "2026-09-20 15:04:05"
```

### 8. Verificar

```bash
go run ./cmd list
```

## Comportamiento Observado

- **Archivado continuo:** cada segmento de WAL se sube a MinIO apenas Postgres lo completa (`archive-push`), sin esperar al siguiente backup base — visible en los logs de pgBackRest (`/tmp/pgbackrest-log` dentro del contenedor) y en los objetos del bucket `pgbackrest`.
- **Backup base:** `pgbackrest backup --type=full` produce un backup completo en MinIO, listado luego con `pgbackrest info`.
- **Desastre simulado:** al borrar el data directory y reiniciar Postgres, el contenedor arranca sin datos: se pierde absolutamente todo lo que no esté en MinIO.
- **Restauración a un punto en el tiempo:** `pgbackrest restore --type=time --target='<timestamp>'` recupera el backup base y reproduce el WAL archivado hasta el instante indicado — las filas 1 a 6 reaparecen (`go run ./cmd list`), pero las filas 7 y 8 (insertadas después del timestamp objetivo) no, confirmando que la recuperación va más allá del último backup completo, que era justamente el punto débil del respaldo lógico periódico descartado en el ADR-06.

## Stack Tecnológico

- **Lenguaje:** Go (Golang), driver `jackc/pgx/v5`
- **Base de datos:** PostgreSQL 16
- **Respaldo y PITR:** pgBackRest
- **Repositorio S3:** MinIO
- **Contenedores:** Docker / Docker Compose
