# PoC 07: Auditoría de Cambios con Retraced Autoalojado

Prueba de concepto (PoC) que levanta **Retraced** (el motor de auditoría de código abierto detrás del ADR-20) por Docker, y usa un programa mínimo en **Go** para publicar un evento de auditoría y consultarlo de vuelta, sin depender de ningún SDK.

## Arquitectura y Justificación

Retraced ya trae su propio bootstrap de desarrollo: al levantar el `docker-compose` del proyecto, crea automáticamente un proyecto llamado `dev` con un token de API `dev`, listo para usarse sin configuración manual. En vez de depender del SDK oficial de Go (`retracedhq/retraced-go`, que según lo investigado está bastante menos mantenido que el de JavaScript), esta PoC llama directamente a la API REST de Retraced con la librería estándar de Go (`net/http`), sin ninguna dependencia externa — esto también sirve para confirmar que el contrato HTTP de Retraced es sencillo de consumir por sí solo, sin atarse a un cliente de terceros.

### Estructura del proyecto

```
poc-07-retraced-auditoria/
├── cmd/
│   └── main.go            # Publica un evento de auditoría de prueba y luego lo consulta de vuelta
├── go.mod
└── (sin go.sum: no usa ninguna dependencia externa, solo la librería estándar)
```

## Guía de Levantamiento

### 1. Levantar Retraced con Docker

```bash
git clone https://github.com/retracedhq/retraced.git
cd retraced
docker compose up -d
```

Esto levanta varios contenedores (Elasticsearch, NSQ, PostgreSQL, la API de Retraced, su procesador, y el portal admin de BoxyHQ/Jackson), así que la primera vez tarda un poco en descargar las imágenes. Cuando el contenedor `retraced-dev-bootstrap` termina, ya existe un proyecto `dev` con API token `dev` listo para usar. La API queda escuchando en `http://localhost:3000` y el portal admin en `http://localhost:5225`.

### 2. Ejecución de la PoC

Desde la carpeta de esta PoC (no la de Retraced):

```bash
go run ./cmd
```

## Comportamiento Observado

- **Publicación del evento:** el programa arma un evento (`lote.actualizado`, con actor, grupo y target) y lo envía por `POST` a `/auditlog/publisher/v1/project/dev/event` con el header `Authorization: token=dev`. Retraced responde con el `id` del evento creado.
- **Consulta del historial:** un segundo `POST` a `/auditlog/admin/v1/project/dev/events/search?environment_id=dev` trae de vuelta el evento recién publicado, igual que lo haría una pantalla de historial de cambios.
- **Sin SDK de por medio:** toda la comunicación es HTTP + JSON plano, confirmando que integrar Retraced no exige depender del cliente de Go menos mantenido.
- **Carga operativa real:** levantar el `docker-compose` de Retraced trae consigo Elasticsearch, NSQ, PostgreSQL y un portal admin aparte — bastante más pesado que un simple contenedor, tal como se advirtió al decidir el ADR-20.

## Stack Tecnológico

- **Lenguaje:** Go (Golang), solo librería estándar (`net/http`, `encoding/json`)
- **Auditoría:** Retraced autoalojado (`docker compose`)
