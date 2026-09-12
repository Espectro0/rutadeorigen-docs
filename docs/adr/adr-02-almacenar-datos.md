# Selección de PostgreSQL como Motor de Bases de Datos

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-04

**Historia técnica:** Uso de un motor de bases de datos para almacenar información de trazabilidad.

## Contexto y planteamiento del problema

*Ruta de Origen* requiere persistir la información de la cadena de trazabilidad del café: fincas, lotes, procesos y evidencias multimedia, todas relacionadas entre sí. El sistema necesita una consistencia de datos fuerte sobre esa información y a la vez una respuesta inmediata en las consultas, ya que los usuarios registran procesos en campo y los consumidores consultan la trazabilidad completa de un lote mediante un código QR.

¿Qué motor de base de datos modela mejor esta estructura fuertemente relacionada, sostiene consultas con múltiples relaciones bajo concurrencia alta, y garantiza que una escritura falle por completo o se complete por completo?

## Impulsores de decisión

* QS-01 Consumidor consulta la información de trazabilidad de un café
* QS-09 Registro inmediato de procesos
* QS-13 Solicitudes concurrentes masivas
* QS-14 Falla durante una operación con alta concurrencia
* QS-15 Error de almacenamiento interno

## Opciones consideradas

* Base de Datos *No Relacional*
    * MongoDB

* Base de Datos *Relacional*
    * PostgreSQL
    * SQL Server

## Resultado de la decisión

Opción elegida: **PostgreSQL**, porque es un motor relacional maduro y de código abierto que modela de forma natural la estructura jerárquica y fuertemente relacionada de la trazabilidad sin necesidad de duplicar información, garantiza consistencia fuerte mediante transacciones ACID, y ofrece un motor de indexación avanzado capaz de sostener consultas con múltiples relaciones bajo alta concurrencia sin desnormalizar el modelo.

### Consecuencias positivas

* Modelo relacional natural para una jerarquía organizada, sin duplicar datos entre entidades.
* Transacciones ACID nativas, sin necesidad de lógica adicional para garantizar atomicidad.
* Motor de código abierto, sin costo de licencia, alineado con el presupuesto limitado del proyecto.
* Ecosistema maduro de indexación para escalar consultas de solo lectura.
* Amplio soporte de ORMs y de proveedores administrados en la nube.
* Soporta columnas `JSONB`, dando flexibilidad puntual sin abandonar el modelo relacional.

### Consecuencias negativas

* Escalar horizontalmente la escritura (sharding) no es nativo; si el volumen de escritura crece de forma extrema se necesitarían herramientas adicionales.
* Requiere definir el esquema y sus migraciones de forma explícita, lo que implica más diseño previo que un modelo sin esquema fijo.
* Una única instancia primaria de escritura es, en principio, un punto de contención si no se planea una estrategia de réplicas de lectura.

## Pros y contras de las opciones

### PostgreSQL

Sistema de gestión de bases de datos relacional de código abierto, con soporte completo de transacciones ACID e indexación extensa.

* Bien, porque modela de forma natural las relaciones jerárquicas de la trazabilidad sin duplicar información.
* Bien, porque ofrece transacciones ACID, indispensables para evitar registros parciales o corruptos.
* Bien, porque cuenta con un motor de indexación maduro que sostiene consultas con joins bajo alta concurrencia.
* Bien, porque es de código abierto y sin costo de licencia, alineado con el presupuesto limitado del proyecto.
* Bien, porque soporta columnas `JSONB`, permitiendo flexibilidad puntual sin abandonar el modelo relacional.
* Bien, porque cuenta con amplio soporte de ORMs y multiples proveedores en la nube.
* Malo, porque escalar horizontalmente la escritura (sharding) no es nativo y requiere herramientas adicionales si el volumen de escritura crece mucho.
* Malo, porque requiere definir el esquema y las migraciones de forma explícita, lo que implica más diseño previo que un modelo sin esquema.

### MongoDB

Base de datos documental NoSQL, orientada a documentos JSON/BSON sin esquema fijo, con escalamiento horizontal nativo mediante sharding.

* Bien, porque el modelo de documentos permite iterar rápido sin migraciones formales de esquema.
* Bien, porque el sharding horizontal es nativo, útil si el volumen de escritura llegara a crecer de forma extrema.
* Bien, porque los documentos anidados encajan bien con datos semi-estructurados, como lecturas variables de sensores IoT.
* Malo, porque para reconstruir una historia completa se requiere desnormalizar y duplicar datos entre colecciones, o resolver varias consultas que degradan el rendimiento bajo alta concurrencia.
* Malo, porque las transacciones multi-documento son más limitadas y costosas en rendimiento que las transacciones nativas de un motor relacional.
* Malo, porque al no tener esquema fijo, la responsabilidad de la integridad de los datos recae en la capa de aplicación, aumentando el riesgo de inconsistencias en un equipo de una sola persona.
* Malo, porque el dominio del negocio es naturalmente relacional; forzarlo a documentos añade complejidad innecesaria.

### SQL Server

Motor de base de datos relacional propietario de Microsoft, con soporte completo de transacciones ACID y herramientas de administración robustas.

* Bien, porque ofrece las mismas garantías relacionales y transaccionales que PostgreSQL: modelo natural para la jerarquía de trazabilidad y ACID multi-tabla.
* Bien, porque cuenta con herramientas de administración e integración maduras dentro del ecosistema Microsoft / .NET.
* Malo, porque requiere licenciamiento comercial.
* Malo, porque no aporta ninguna ventaja técnica adicional frente a PostgreSQL para este caso de uso que justifique su costo.

## Enlaces

* [Relacionado con] [ADR-01: Uso de Redis como memoria caché](./adr-01-gestionar-cache.md)
* [Relacionado con] [ADR-03: Selección de Go como Lenguaje y Runtime del Backend](./adr-03-definir-runtime.md)
* [Relacionado con] [ADR-04: Selección de Ent como Capa de Acceso a Datos (ORM)](./adr-04-seleccionar-orm.md)
