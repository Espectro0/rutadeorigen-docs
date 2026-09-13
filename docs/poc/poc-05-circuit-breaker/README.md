# PoC 05: Circuit Breaker para Proteger Servicios Externos

Prueba de concepto (PoC) que implementa un **circuit breaker** con timeouts para proteger las llamadas del backend a servicios externos, usando **Go**. Simula un servicio externo que falla temporalmente y luego se recupera, para observar en vivo cómo el circuito pasa de cerrado a abierto, a medio abierto, y de vuelta a cerrado.

## Arquitectura y Justificación

Proteger cada llamada a un servicio externo directamente en el código, sin añadir un componente de infraestructura adicional como un proxy o service mesh dedicado, evita que una falla o una demora de un tercero se propague y bloquee al backend. La lógica se separa en tres responsabilidades: el servicio externo simulado, la configuración del circuit breaker y el cliente que envuelve la llamada real con el breaker. Así, cualquier llamada a un servicio externo pasa siempre por el mismo mecanismo de protección, sin que el resto de la aplicación dependa de los detalles internos del breaker.

### Estructura del proyecto

```
poc-05-circuit-breaker/
├── cmd/
│   └── main.go               # Levanta el servicio externo simulado y hace varias llamadas para forzar los estados del breaker
├── internal/
│   ├── external/
│   │   └── server.go          # Servidor HTTP de prueba: responde bien, luego falla, y luego se recupera
│   └── resiliencia/
│       ├── breaker.go          # Configura el circuit breaker
│       └── client.go           # Cliente HTTP con timeout que envuelve la llamada al servicio externo con el breaker
├── go.mod
└── go.sum
```

## Guía de Levantamiento

### 1. Instalar dependencias

```bash
go get github.com/sony/gobreaker
go mod tidy
```

### 2. Ejecución

```bash
go run ./cmd
```

## Comportamiento Observado

- **Circuito cerrado:** las primeras solicitudes llegan bien al servicio simulado y se resuelven con normalidad.
- **Apertura del circuito:** tras 2 fallos consecutivos del servicio simulado, el breaker se abre y pasa a rechazar las siguientes solicitudes de inmediato (`circuit breaker is open`), sin siquiera intentar la llamada real.
- **Medio abierto:** pasado el tiempo de espera configurado, el breaker deja pasar una única solicitud de prueba; si el servicio simulado sigue fallando, se reabre y vuelve a esperar.
- **Recuperación:** una vez que el servicio simulado ya está respondiendo bien, la solicitud de prueba en medio abierto tiene éxito y el breaker cierra el circuito, volviendo al comportamiento normal.
- **Trazabilidad:** cada cambio de estado del breaker se imprime en consola (`[breaker] servicio-externo-prueba pasó de closed a open`, etc.), permitiendo seguir la transición completa en vivo.

## Stack Tecnológico

- **Lenguaje:** Go (Golang)
- **Circuit Breaker:** `sony/gobreaker`
