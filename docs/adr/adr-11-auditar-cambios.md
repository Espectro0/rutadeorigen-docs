# Auditar Cambios

* **Estado:** Reemplazado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-05

**Historia técnica:** Definir dónde y cómo se registra la bitácora de auditoría de las modificaciones y eliminaciones realizadas sobre la información de trazabilidad.

## Contexto y planteamiento del problema

El 100% de las modificaciones o eliminaciones de un registro queden auditadas con usuario, fecha, hora, lugar y acción, y ese historial pueda consultarse mostrando esos mismos datos para el 100% de los cambios almacenados. Este registro de auditoría es un flujo de datos append-only, de crecimiento constante, con un propósito y un patrón de consulta distintos a los de las tablas operativas de trazabilidad.

¿Dónde se almacena y cómo se consulta la bitácora de auditoría de cambios, sin comprometer el rendimiento ni el crecimiento de la base de datos operativa?

## Impulsores de decisión

* QS-08 Modificación o eliminación de un registro
* QS-18 Consulta del historial de cambios

## Opciones consideradas

* Tabla de auditoría propia en PostgreSQL (append-only)
* Event sourcing completo (el estado se reconstruye a partir de una secuencia de eventos)
* Bitácora enviada a un servicio externo especializado en auditoría (WorkOS Audit Logs)

## Resultado de la decisión

Opción elegida: **Bitácora enviada a WorkOS Audit Logs**, porque es un producto construido específicamente para eventos de auditoría, en lugar de una plataforma de observabilidad genérica adaptada para este fin. Esto evita además que la tabla de auditoría compita por recursos e índices con las tablas operativas a medida que crece. Cada modificación o eliminación registrada por la capa transaccional del backend se emite como un evento hacia WorkOS Audit Logs, que sirve como fuente para las consultas de historial.

### Consecuencias positivas

* El crecimiento de la bitácora de auditoría no compite por recursos ni índices con las tablas operativas.
* Se apoya en un producto construido específicamente para auditoría, no una plataforma de observabilidad genérica adaptada.
* Las consultas de historial quedan desacopladas de la carga de la base de datos operativa.
* Su capa gratuita es razonable para el volumen esperado de un proyecto de este tamaño.

### Consecuencias negativas

* Introduce una dependencia de un servicio externo para un requisito con implicaciones de confianza y trazabilidad, en lugar de mantenerlo en la infraestructura propia.
* Es un servicio más nuevo y especializado, con menos herramientas de análisis y visualización que una plataforma de observabilidad completa.
* La indisponibilidad del servicio externo puede dejar cambios sin auditar temporalmente si no se define una estrategia de reintento o cola intermedia.

## Pros y contras de las opciones

### Tabla de auditoría propia en PostgreSQL (append-only)

Cada modificación o eliminación se inserta como un registro en una tabla dedicada de auditoría dentro de la misma base de datos operativa.

* Bien, porque reutiliza la base de datos y el contexto transaccional ya definido.
* Bien, porque mantiene una sola fuente de verdad, consultable directamente con SQL para el historial.
* Malo, porque el crecimiento constante de la tabla de auditoría compite por recursos con las tablas operativas si no se gestiona con cuidado.
* Malo, porque mezcla en la misma base de datos un flujo de solo lectura de crecimiento indefinido con las tablas operativas de escritura frecuente.

### Event sourcing completo (el estado se reconstruye a partir de una secuencia de eventos)

El estado de cada proceso de trazabilidad se reconstruye a partir de la secuencia completa de eventos ocurridos sobre él, en lugar de guardarse como un registro mutable.

* Bien, porque la auditoría queda garantizada por diseño: cada cambio es, en sí mismo, un evento inmutable.
* Bien, porque ofrece trazabilidad total del historial de cualquier entidad.
* Malo, porque implica un cambio arquitectónico mayor sobre lo ya decidido.
* Malo, porque es un costo de implementación demasiado alto para esta altura del proyecto, solo para resolver el requisito de auditoría.

### Bitácora enviada a un servicio externo especializado en auditoría (WorkOS Audit Logs)

Cada modificación o eliminación se emite como un evento de auditoría hacia WorkOS Audit Logs, un producto construido específicamente para este propósito.

* Bien, porque es un producto hecho específicamente para auditoría, no un logger genérico adaptado.
* Bien, porque no compite por recursos con las tablas operativas de la base de datos propia.
* Bien, porque su capa gratuita es razonable para el volumen esperado del proyecto.
* Malo, porque introduce una dependencia externa para un requisito con implicaciones de confianza y trazabilidad.
* Malo, porque al ser más especializado y reciente, ofrece menos herramientas de análisis/visualización que una plataforma de observabilidad completa.

## Enlaces

* [Reemplazado por] [ADR-20: Migrar la Auditoría de Cambios a Retraced Autoalojado](./adr-20-migrar-auditoria-retraced.md)
* [Relacionado con] [ADR-02: Selección de PostgreSQL como Motor de Bases de Datos](./adr-02-almacenar-datos.md)
* [Relacionado con] [ADR-12: Proteger Servicios Externos](./adr-12-proteger-servicios.md)
