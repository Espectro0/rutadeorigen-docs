# Encolar Tareas Pesadas

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-05

**Historia técnica:** Definir cómo se ejecutan las operaciones que no requieren una respuesta inmediata al usuario, como la generación del sello digital de verificación, sin bloquear la solicitud original.

## Contexto y planteamiento del problema

Multiples procesos necesitan trabajarse de manera asíncrona para evitar que las solicitudes tengan un tiempo de espera muy grande, permitiendo que el usuario siga utilizando las funcionalidades del software mientras por ejemplo, de manera concurrente se estan procesando los archivos de un lote que subió anteriormente.

¿Dónde se ejecutan las tareas que no requieren una respuesta inmediata, de forma coordinada entre las distintas instancias del backend?

## Impulsores de decisión

* QS-12 Generación de sello digital

## Opciones consideradas

* Cola de trabajos sobre Redis
* Cola de mensajería dedicada (RabbitMQ)
* Goroutines en segundo plano dentro del propio proceso del backend, sin cola persistente

## Resultado de la decisión

Opción elegida: **Cola de mensajería dedicada (RabbitMQ)**, porque es de código abierto (sin costo de licencia) y ofrece garantías de entrega más robustas: reintentos configurables, colas de mensajes muertos (dead-letter queues) y persistencia propia, separando por completo la responsabilidad de mensajería de la de caché. RabbitMQ recibe y distribuye las tareas entre las distintas instancias del backend habilitadas por el escalado horizontal, sin bloquear la solicitud original del usuario.

### Consecuencias positivas

* Ofrece garantías de entrega más robustas que una librería de colas apoyada sobre Redis.
* Separa por completo la responsabilidad de mensajería de la de caché, sin sobrecargar a Redis con un segundo rol.
* Las tareas quedan coordinadas entre las distintas instancias del escalado horizontal, evitando que se pierdan o se dupliquen al reiniciarse una instancia.
* La solicitud original del usuario no queda bloqueada esperando el resultado de la generación del sello.
* Es de código abierto, sin costo de licencia, acorde al presupuesto limitado del proyecto.

### Consecuencias negativas

* Introduce un componente de infraestructura completamente nuevo (RabbitMQ), adicional a Redis, que operar y presupuestar en cómputo y almacenamiento.
* Implica una curva de aprendizaje y una carga operativa adicional para un equipo de una sola persona, al sumar un sistema de mensajería distinto a los ya adoptados.
* Añade una dependencia interna más que mantener disponible, junto a Redis y PostgreSQL.

## Pros y contras de las opciones

### Cola de trabajos sobre Redis

Utilizar el propio Redis ya adoptado para encolar y procesar tareas en segundo plano.

* Bien, porque no añade un nuevo componente de infraestructura a operar.
* Bien, porque coordina las tareas entre las múltiples instancias del backend.
* Bien, porque cuenta con bibliotecas maduras en el ecosistema.
* Malo, porque acopla dos responsabilidades (caché y cola) sobre el mismo componente Redis.
* Malo, porque sus garantías de entrega dependen de la librería elegida, no de un sistema de mensajería dedicado.

### Cola de mensajería dedicada (RabbitMQ)

RabbitMQ, un sistema de mensajería de código abierto diseñado específicamente para colas, gestiona la generación y el consumo de tareas en segundo plano.

* Bien, porque es de código abierto, sin costo de licencia.
* Bien, porque ofrece garantías más robustas que una cola sobre Redis: reintentos configurables, dead-letter queues y persistencia propia.
* Bien, porque separa completamente la responsabilidad de mensajería de la de caché, sin sobrecargar a Redis con un segundo rol.
* Malo, porque introduce un componente de infraestructura completamente nuevo, adicional a Redis, que operar y presupuestar.
* Malo, porque implica una curva de aprendizaje y una carga operativa adicional para un equipo de una sola persona.

### Goroutines en segundo plano dentro del propio proceso del backend, sin cola persistente

El propio proceso del backend lanza goroutines para ejecutar tareas en segundo plano, sin persistirlas en un almacén externo.

* Bien, porque no requiere infraestructura adicional y su implementación es inmediata en Go.
* Bien, porque no añade latencia de red hacia un componente externo.
* Malo, porque las tareas se pierden si la instancia se reinicia o cae a mitad del proceso.
* Malo, porque no se coordina entre las múltiples instancias del escalado horizontal, pudiendo duplicar o perder tareas.

## Enlaces

* [Relacionado con] [ADR-01: Uso de Redis como memoria caché](./adr-01-gestionar-cache.md)
* [Relacionado con] [ADR-18: Notificar Invitaciones a una Organización](./adr-18-notificar-invitaciones.md)
