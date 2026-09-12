# PoC 01: Sistema de tareas mediante una cola (RabbitMQ)

Prueba de concepto (PoC) que implementa un sistema de procesamiento de tareas utilizando **RabbitMQ** y **Go**. Simula la producción de tareas, su consumo mediante workers (1+), una estrategia de reintentos con un **retraso** y cun mecanismo de **Dead Letter Queue (DLQ)** que permite aislar los fallos persistentes.

## Arquitectura y Justificación

El uso de un gestor de colas externo como **RabbitMQ** desacopla el flujo de trabajo y previene la saturación del sistema, evitando los límites y riesgos asociados al uso exclusivo de concurrencia en memoria con _Goroutines_.

### Estructura del proyecto.

```
poc-01-rabbitqm-colas/
├── cmd/
│   └── main.go          # Inicializa conexión, publica y ejecuta los workers
├── internal/
│   ├── queue/
│   │   └── connection.go # Conexión/Canal AMQP y helpers de topología
│   └── task/
│       └── task.go       # Definición de tareas, lógica del worker y reintentos
├── go.mod
└── go.sum 
```

### Parámetros de Configuración (`internal/task/task.go`)

| Parámetro | Valor | Descripción |
| :--- | :---: | :--- |
| `NormalCount` | `20` | Tareas estándar simuladas para procesamiento exitoso. |
| `PoisonCount` | `3` | Tareas inyectadas con el propósito de fallar. |
| `MaxRetries` | `3` | Número máximo de reintentos permitidos antes de mover la tarea a la `DLQ`. |

## Guía de Levantamiento

### 1. Levantar RabbitMQ (Docker)

```bash
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:4-management
```

* Panel de administración: http://localhost:15672 (Credenciales por defecto: `guest`/`guest`).

### 2. Instalar dependencias

```bash
go get github.com/rabbitmq/amqp091-go
go mod tidy
```

### 3. Ejecución

```bash
go run ./cmd
```

## Comportamiento Observado

- **Procesamiento Exitoso:** Las 20 tareas normales son consumidas, procesadas y confirmadas correctamente.
- **Estrategia de Reintentos:** Las 3 tareas inyectadas fallan de forma intencional, esperan un intervalo de 3 segundos y regresan a la cola.
- **Manejo de Errores Críticos (DLQ):** Al agotar los 3 intentos permitidos, las tareas fallidas son derivadas automáticamente a una cola de procesos muertos (Dead Letter Queue).
- **Trazabilidad:** Cada etapa del ciclo de vida de la tarea emite un resultado estructurado que permite identificar qué instantcia del worker proceso la carga.


## Stack Tecnológico

- **Lenguage:** Go (Golang)
- **Broker de Mensajes:** RabbitMQ (`amqp091-go`)
- **Contenedores:** Docker