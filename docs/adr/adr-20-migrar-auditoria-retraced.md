# Migrar la Auditoría de Cambios a Retraced Autoalojado

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-12

**Historia técnica:** Reemplazar la decisión del ADR-11 de enviar la bitácora de auditoría a WorkOS Audit Logs por Retraced, el motor de auditoría de código abierto que WorkOS absorbió y dejó de ofrecer como servicio propio, ahora autoalojado dentro de la infraestructura propia del proyecto.

## Contexto y planteamiento del problema

El ADR-11 decidió enviar la bitácora de auditoría a WorkOS Audit Logs por ser un producto especializado en auditoría, distinto de una plataforma de observabilidad genérica adaptada para este fin. Sin embargo, esa decisión introduce una dependencia recurrente de un servicio externo de pago a partir de cierto volumen de eventos, en un proyecto que ya ha optado consistentemente por infraestructura propia y sin costo de licencia cuando existe una alternativa madura. Retraced era el motor original detrás de este tipo de producto de auditoría antes de que WorkOS lo absorbiera, y sigue disponible como proyecto de código abierto autoalojable, lo que permite conservar las mismas ventajas señaladas en el ADR-11 sin depender de un tercero ni de su capa de pago.


## Impulsores de decisión

* QS-08 Modificación o eliminación de un registro
* QS-18 Consulta del historial de cambios

## Opciones consideradas

* Mantener WorkOS Audit Logs (decisión original del ADR-11)
* Autoalojar Retraced (motor de auditoría de código abierto)
* Tabla de auditoría propia en PostgreSQL (append-only) — ya descartada en el ADR-11

## Resultado de la decisión

Opción elegida: **Autoalojar Retraced**, porque conserva el mismo beneficio central que motivó el ADR-11 sin el costo recurrente ni la dependencia de un proveedor externo, siendo coherente con el resto de la infraestructura propia ya adoptada en el proyecto. Este ADR reemplaza al ADR-11.

### Consecuencias positivas

* Elimina el costo recurrente de WorkOS Audit Logs a partir de cierto volumen de eventos, sin renunciar a un producto especializado en auditoría.
* Conserva el crecimiento de la bitácora separado de las tablas operativas de PostgreSQL, igual que con WorkOS Audit Logs.
* Se despliega junto al resto de la infraestructura propia, sin depender de una cuenta ni credenciales de un tercero.
* Al ser de código abierto, el equipo controla directamente su disponibilidad y su ciclo de actualizaciones.

### Consecuencias negativas

* Retraced ya no tiene el respaldo comercial ni el soporte activo de WorkOS: el proyecto quedó como código abierto mantenido por la comunidad (Replicated/BoxyHQ), con menor garantía de continuidad a largo plazo que un producto de pago activamente vendido.
* La disponibilidad, el respaldo y el escalado de la instancia de Retraced quedan bajo responsabilidad del propio proyecto, sin el SLA de un servicio administrado.
* Requiere desplegar y mantener un componente de infraestructura adicional, con su propio almacenamiento, a diferencia de simplemente consumir una API externa.

## Pros y contras de las opciones

### Mantener WorkOS Audit Logs

Continuar enviando la bitácora de auditoría al servicio externo decidido originalmente en el ADR-11.

* Bien, porque no requiere ningún cambio ni migración sobre lo ya implementado.
* Bien, porque es un servicio administrado, sin infraestructura propia que operar.
* Malo, porque introduce un costo recurrente a partir de cierto volumen de eventos, contrario al criterio de presupuesto limitado ya usado en el resto de decisiones del proyecto.
* Malo, porque depende de que WorkOS mantenga este producto específico disponible y con los mismos términos a futuro.

### Autoalojar Retraced

Desplegar la propia instancia de Retraced (código abierto) junto al resto de la infraestructura del proyecto, y emitir hacia ella los eventos de auditoría.

* Bien, porque es de código abierto, sin costo de licencia ni de volumen procesado.
* Bien, porque es el mismo motor que originalmente ofrecía este tipo de producto de auditoría especializado, antes de ser absorbido por WorkOS.
* Bien, porque se despliega junto al resto de la infraestructura propia ya adoptada (MinIO, RabbitMQ, Prometheus + Grafana).
* Malo, porque ya no cuenta con el respaldo comercial activo de un proveedor, solo con el mantenimiento de la comunidad.
* Malo, porque suma un componente más que operar, desplegar y respaldar, a diferencia de consumir una API externa administrada.

### Tabla de auditoría propia en PostgreSQL (append-only)

Ya evaluada y descartada en el ADR-11 por competir por recursos e índices con las tablas operativas a medida que crece; se mantiene descartada por el mismo motivo.

* Bien, porque no suma ningún componente de infraestructura nuevo.
* Malo, porque mezcla en la misma base de datos un flujo de solo lectura de crecimiento indefinido con las tablas operativas de escritura frecuente, el mismo problema identificado en el ADR-11.

## Enlaces

* [Reemplaza a] [ADR-11: Auditar Cambios](./adr-11-auditar-cambios.md)
* [Relacionado con] [ADR-02: Selección de PostgreSQL como Motor de Bases de Datos](./adr-02-almacenar-datos.md)
* [Relacionado con] [ADR-07: Almacenar Evidencias como Objetos](./adr-07-almacenar-evidencias.md)
* [Relacionado con] [ADR-12: Proteger Servicios Externos](./adr-12-proteger-servicios.md)
