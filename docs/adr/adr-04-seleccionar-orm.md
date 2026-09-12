# Selección de Ent como Capa de Acceso a Datos (ORM)

* **Estado:** Propuesto
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-05

**Historia técnica:** Selección de la capa de acceso a datos para interactuar con PostgreSQL desde el backend.

## Contexto y planteamiento del problema

Ya se decidió Go como lenguaje/runtime y PostgreSQL como motor de base de datos. Falta decidir cómo el backend interactúa con esa base de datos: cuánto código repetitivo se escribe a mano, qué tan disciplinada queda la gestión de migraciones a medida que el modelo de trazabilidad crece, y qué tan segura es esa interacción frente a errores de tipos o de consultas mal construidas.

¿Qué capa de acceso a datos modela de forma segura las relaciones de trazabilidad, mantiene las migraciones de esquema de forma disciplinada, y sostiene las consultas con múltiples relaciones bajo la concurrencia exigida sin sacrificar seguridad de tipos?

## Impulsores de decisión

* QS-01 Consumidor consulta la información de trazabilidad de un café
* QS-06 Validación de datos en formularios
* QS-09 Registro inmediato de procesos
* QS-13 Solicitudes concurrentes masivas
* QS-14 Falla durante una operación con alta concurrencia
* QS-15 Error de almacenamiento interno

## Opciones consideradas

* Ent (`entgo.io`)
* GORM
* `database/sql` + SQL plano (`sqlc`)

## Resultado de la decisión

Opción elegida: **Ent**, porque genera un esquema fuertemente tipado a partir de la definición de entidades en Go, generando automáticamente las migraciones y detectando en tiempo de compilación errores de relaciones o de tipos que en GORM o en SQL plano solo aparecerían en tiempo de ejecución. Su motor de consultas, basado en un grafo de relaciones, modela de forma natural la jerarquía de la trazabilidad sin necesidad de escribir joins a mano, y soporta el manejo transaccional explícito que la siguiente decisión necesita.

### Consecuencias positivas

* Migraciones de esquema generadas y versionadas automáticamente a partir del modelo de entidades.
* Seguridad de tipos en tiempo de compilación: un cambio de campo o de relación rompe el build, no la producción.
* API de consultas basada en el grafo de relaciones, evitando escribir joins SQL a mano para reconstruir la trazabilidad completa de un lote.
* Soporte nativo de transacciones explícitas.
* La generación de código repetitivo reduce el mantenimiento manual para un equipo de una sola persona.

### Consecuencias negativas

* Curva de aprendizaje inicial mayor que GORM, al requerir generar código antes de poder usar las entidades.
* Menor flexibilidad que SQL plano para consultas muy puntuales de optimización fina.
* Comunidad y cantidad de ejemplos disponibles más pequeña que la de GORM.

## Pros y contras de las opciones

### Ent (entgo.io)

ORM de Go, basado en la definición de entidades como grafo, con generación de código para el esquema, las migraciones y las consultas.

* Bien, porque genera migraciones de esquema automáticamente a partir de las entidades definidas en Go.
* Bien, porque su seguridad de tipos detecta en tiempo de compilación errores de relaciones o consultas que en otras opciones solo se ven en tiempo de ejecución.
* Bien, porque modela de forma natural relaciones jerárquicas mediante su API de grafo, sin necesidad de joins escritos a mano.
* Bien, porque soporta transacciones explícitas de forma nativa.
* Malo, porque requiere un paso de generación de código antes de poder compilar, añadiendo un paso extra al flujo de desarrollo.
* Malo, porque su comunidad y documentación de ejemplos son más pequeñas que las de GORM.

### GORM

ORM de Go de uso extendido en la comunidad, basado en structs con etiquetas y reflexión en tiempo de ejecución.

* Bien, porque cuenta con una comunidad amplia y abundante documentación y ejemplos.
* Bien, porque su curva de aprendizaje inicial es más baja, sin pasos de generación de código.
* Malo, porque su detección de errores ocurre en tiempo de ejecución mediante reflexión, no en tiempo de compilación, un riesgo mayor para un equipo de una sola persona.
* Malo, porque es más fácil incurrir en el problema N+1 al no forzar explícitamente la carga de relaciones, arriesgando el rendimiento bajo escenarios de alta concurrencia.
* Malo, porque sus migraciones automáticas son más limitadas y suelen complementarse con herramientas externas para esquemas complejos.

### `database/sql` + SQL plano (`sqlc`)

Acceso directo a la base de datos mediante la librería estándar de Go, escribiendo las consultas SQL a mano o generándolas desde SQL con herramientas como `sqlc`.

* Bien, porque ofrece control total sobre cada consulta, permitiendo optimización fina para casos puntuales.
* Bien, porque no depende de ninguna librería de terceros como ORM, reduciendo la superficie de dependencias.
* Malo, porque cada relación debe reconstruirse a mano con joins SQL explícitos, aumentando el código repetitivo y el riesgo de error humano en un equipo de una sola persona.
* Malo, porque las migraciones de esquema deben gestionarse con una herramienta aparte, sin garantía de sincronía automática con el modelo en Go.
* Malo, porque no ofrece ninguna verificación de tipos entre el modelo en Go y el esquema de la base de datos; un cambio de columna solo se detecta en tiempo de ejecución.

## Enlaces

* [Relacionado con] [ADR-02: Selección de PostgreSQL como Motor de Bases de Datos](./adr-02-almacenar-datos.md)
* [Relacionado con] [ADR-03: Selección de Go como Lenguaje y Runtime del Backend](./adr-03-definir-runtime.md)
