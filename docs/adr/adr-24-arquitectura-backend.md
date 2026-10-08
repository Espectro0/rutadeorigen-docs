# Definir arquitectura concreta para el Backend

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-10-08

**Historia técnica:** Definir el estilo arquitectónico del backend: cómo se divide en unidades desplegables, cómo se organiza cada una por dentro, cómo se comunican entre sí y cómo se separan sus datos. Las decisiones de infraestructura sobre las que corre (PostgreSQL en el ADR-02, RabbitMQ en el ADR-08, k3s en el ADR-15 y Traefik en el ADR-21) ya están tomadas; este ADR define la forma del software que se despliega sobre ellas. La estrategia de repositorio y de CI/CD se define en el ADR-26.

## Contexto y planteamiento del problema

El backend concentra capacidades con requisitos muy distintos entre sí. La consulta pública de un lote a partir de un QR impreso debe estar disponible siempre y atender picos de escaneos anónimos (QS-01, QS-13, TC-02). El registro de trazabilidad es transaccional, autenticado y aislado por organización (TC-07). El procesamiento de evidencias consume CPU y E/S de forma intensiva (ADR-19). El sello digital y el anclaje son asíncronos por naturaleza (QS-12). Las notificaciones y el asistente de IA dependen de proveedores externos que pueden fallar (QS-10). Mantener todo en una sola unidad desplegable obliga a escalar, actualizar y hacer fallar estas capacidades en conjunto, aunque solo una de ellas lo necesite.

¿Cómo se divide y organiza el backend para que cada capacidad pueda desplegarse, escalarse y fallar de forma independiente, sin perder la consistencia de la información de trazabilidad?

## Impulsores de decisión

* QS-01 Consumidor consulta la información de trazabilidad de un café
* QS-10 Fallo de un servicio externo
* QS-12 Generación de sello digital
* QS-13 Solicitudes concurrentes masivas
* QS-14 Falla durante una operación con alta concurrencia
* TC-02 Permanencia de QR
* TC-07 Seguridad (aislamiento de datos entre organizaciones)
* TC-08 Patrones de diseño
* TC-09 Prácticas de código limpio
* TC-10 Prácticas DevOps

## Opciones consideradas

* Monolito por capas
* Monolito modular
* Híbrido (monolito modular con servicios aparte solo donde los requisitos divergen)
* Microservicios

## Resultado de la decisión

Opción elegida: **Microservicios**, porque permite que cada capacidad del backend se despliegue, escale y falle de forma independiente según sus propios requisitos: la consulta pública por QR sigue respondiendo aunque el servicio de trazabilidad esté caído o en actualización (TC-02, TC-10), el procesamiento de evidencias escala por separado sin competir por recursos con las solicitudes de los usuarios (QS-13), y una caída de un proveedor externo queda contenida en el servicio que lo usa (QS-10). El autoescalado horizontal por servicio que ofrece k3s (ADR-15) y el enrutamiento por dominio y ruta de Traefik (ADR-21) se aprovechan directamente con esta división.

El backend se divide en los siguientes servicios:

| Servicio | Responsabilidad |
| :--- | :--- |
| **Trazabilidad** | Organizaciones, fincas, lotes, eventos de etapa, transferencias y cadena de integridad por lote. Es el lado transaccional y la fuente de verdad de la trazabilidad. |
| **QR** | Códigos, asignaciones, tabla de redirecciones (ADR-16), generación de las imágenes (ADR-25), resolución de escaneos, registro de escaneos y página pública del lote, renderizada desde una copia local de lectura. |
| **Procesamiento de archivos** | Workers que validan, procesan y almacenan las evidencias en MinIO (ADR-07, ADR-19). |
| **Sello y anclaje** | Generación del sello digital de verificación, firma y anclaje periódico de la cadena de integridad (QS-12). |
| **Notificaciones** | Consumo de eventos y envío de correos transaccionales mediante Resend (ADR-18). |
| **Asistente IA** | Integración con el proveedor de LLM y RAG sobre los datos de trazabilidad (ADR-14). |

Además, se definen las siguientes reglas para todos los servicios:

* **Organización interna con Clean Architecture:** cada servicio separa `domain` (entidades y reglas), `app` (casos de uso y las interfaces que necesitan) y `adapters` (HTTP, PostgreSQL, RabbitMQ y servicios externos), con la regla de dependencias apuntando siempre hacia el dominio. El código propio de cada servicio vive bajo `internal/` para que ningún otro servicio pueda importarlo.
* **Comunicación mixta:** por defecto, los servicios se comunican mediante eventos de dominio publicados en RabbitMQ (ADR-08), usando el patrón *outbox* para que el cambio y su evento se guarden en la misma transacción. Las llamadas síncronas (HTTP interno, no expuesto a internet) se reservan para consultas puntuales que necesitan respuesta inmediata, como reconstruir una copia local de datos.
* **Consumidores idempotentes:** todo evento lleva un identificador único y una versión del agregado, de modo que un mensaje repetido o fuera de orden no altera el estado.
* **Datos separados por esquema:** una sola instancia de PostgreSQL (ADR-02) con un esquema y un usuario de base de datos por servicio, sin permisos sobre los esquemas de los demás. Ningún servicio lee directamente los datos de otro; si los necesita, mantiene una copia propia alimentada por eventos.
* **Consistencia transaccional dentro de un servicio:** los eventos de etapa de un lote y su eslabón en la cadena de integridad se escriben en la misma transacción dentro del servicio de Trazabilidad. El servicio de Sello y anclaje consume esa cadena ya construida y nunca participa en su escritura.

### Consecuencias positivas

* La consulta pública por QR queda aislada del resto del sistema: los escaneos siguen resolviendo aunque Trazabilidad esté caído o en actualización, cumpliendo directamente TC-02 y QS-01.
* Cada servicio escala de forma independiente con su propio autoescalado en k3s (ADR-15): los picos de escaneos o de procesamiento de evidencias no consumen los recursos del registro de trazabilidad.
* Una falla en un proveedor externo (Resend, proveedor de LLM) queda contenida en su servicio y no degrada el registro de trazabilidad ni la consulta pública (QS-10).
* Cada servicio se puede actualizar y desplegar por separado, reduciendo el alcance y el riesgo de cada cambio (TC-10).
* Clean Architecture mantiene la lógica de negocio independiente de la infraestructura, de modo que un cambio de proveedor (por ejemplo, del servicio de correo o del proveedor de LLM) solo reemplaza un adaptador (TC-08, TC-09).
* Separar por esquema en una sola instancia de PostgreSQL conserva la propiedad de los datos por servicio sin multiplicar la operación de bases de datos ni los respaldos ya definidos en el ADR-06 (BC-06).

### Consecuencias negativas

* Implica una carga operativa y de diseño considerable para un equipo de una sola persona: seis servicios que desplegar, monitorear y versionar, en lugar de una sola unidad, en tensión con el principio KISS del TC-08.
* La comunicación por eventos introduce consistencia eventual: un cambio en Trazabilidad tarda un tiempo en reflejarse en la copia local del servicio de QR, y el sistema debe tolerar ese desfase.
* Exige implementar y mantener patrones adicionales que un monolito no necesita: *outbox*, consumidores idempotentes, colas de mensajes muertos, versionado de eventos y autenticación entre servicios.
* Depurar una falla requiere seguir una solicitud a través de varios servicios y de RabbitMQ, lo que hace indispensable propagar un identificador de correlación y centralizar logs y trazas sobre la plataforma de monitoreo del ADR-13.
* Los contratos de los eventos se vuelven una interfaz pública entre servicios: un cambio incompatible exige versionar el evento y migrar a sus consumidores en lugar de modificarlo directamente.
* Compartir una sola instancia de PostgreSQL significa que una caída de esa instancia afecta a todos los servicios a la vez; el aislamiento es lógico, no de disponibilidad.

## Pros y contras de las opciones

### Monolito por capas

Una sola unidad desplegable organizada por capas técnicas (controladores, servicios y repositorios), sin módulos de dominio aislados entre sí.

* Bien, porque es la opción más simple de construir, desplegar y depurar.
* Bien, porque todas las operaciones comparten una misma transacción, sin consistencia eventual.
* Malo, porque, al no aislar el dominio, la lógica de trazabilidad, QR y evidencias termina acoplada entre capas, dificultando cualquier separación futura (TC-08, TC-09).
* Malo, porque la consulta pública por QR comparte disponibilidad con todo el sistema: una falla o actualización del backend interrumpe los escaneos (TC-02).

### Monolito modular

Una sola unidad desplegable dividida en módulos de dominio con interfaces públicas explícitas, esquemas de base de datos separados y dependencias controladas.

* Bien, porque conserva la simplicidad operativa de un solo despliegue y la consistencia transaccional entre módulos.
* Bien, porque los límites claros entre módulos permiten extraer un servicio más adelante sin reescribirlo.
* Malo, porque todos los módulos escalan, se actualizan y fallan juntos: los picos de escaneos o de procesamiento de archivos compiten por los mismos recursos que el registro de trazabilidad.
* Malo, porque la consulta pública por QR sigue sin poder aislarse del resto del sistema.

### Híbrido

Un monolito modular para la parte transaccional y servicios aparte solo para las capacidades con requisitos claramente distintos, como la consulta pública por QR.

* Bien, porque aísla la disponibilidad de los escaneos sin asumir la carga de operar todos los servicios por separado.
* Bien, porque mantiene en un mismo despliegue las capacidades que comparten transacciones.
* Malo, porque combina dos modelos de despliegue y comunicación distintos dentro del mismo sistema.
* Malo, porque las capacidades que permanecen en el monolito, como las notificaciones y el asistente de IA, siguen compartiendo recursos y fallas con el registro de trazabilidad.

### Microservicios

Cada capacidad del backend se despliega como un servicio independiente, con su propio esquema de datos, comunicado con los demás principalmente mediante eventos.

* Bien, porque cada servicio se despliega, escala y falla de forma independiente, según sus propios requisitos.
* Bien, porque la consulta pública por QR queda aislada del resto del sistema, cumpliendo TC-02.
* Bien, porque aprovecha directamente el autoescalado por servicio de k3s (ADR-15) y el enrutamiento de Traefik (ADR-21).
* Malo, porque introduce consistencia eventual y patrones adicionales (*outbox*, idempotencia, versionado de eventos).
* Malo, porque implica la mayor carga operativa y de observabilidad de las cuatro opciones para un equipo pequeño.

## Enlaces

* [Relacionado con] [ADR-02: Selección de PostgreSQL como Motor de Bases de Datos](./adr-02-almacenar-datos.md)
* [Relacionado con] [ADR-07: Almacenar Evidencias como Objetos](./adr-07-almacenar-evidencias.md)
* [Relacionado con] [ADR-08: Encolar Tareas Pesadas](./adr-08-encolar-tareas.md)
* [Relacionado con] [ADR-13: Monitorear el Sistema](./adr-13-monitorear-sistema.md)
* [Relacionado con] [ADR-14: Integrar Asistente de IA de Trazabilidad](./adr-14-integrar-asistente-ia.md)
* [Relacionado con] [ADR-15: Orquestar Contenedores para Escalado Horizontal](./adr-15-orquestar-contenedores.md)
* [Relacionado con] [ADR-16: Versionar las URLs de los Códigos QR](./adr-16-versionar-urls.md)
* [Relacionado con] [ADR-18: Notificar Invitaciones a una Organización](./adr-18-notificar-invitaciones.md)
* [Relacionado con] [ADR-19: Procesar Archivos antes de Almacenarlos en MinIO](./adr-19-procesar-archivos.md)
* [Relacionado con] [ADR-21: Exponer los Servicios mediante un Reverse Proxy](./adr-21-exponer-servicios.md)
* [Relacionado con] [ADR-25: Migrar la Generación de Códigos QR al Servicio de QR](./adr-25-migrar-generacion-qr.md)
* [Relacionado con] [ADR-26: Automatizar la Integración y el Despliegue Continuos](./adr-26-automatizar-despliegues.md)
