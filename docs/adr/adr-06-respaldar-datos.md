# Respaldar y Recuperar Datos ante Desastres

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-05

**Historia técnica:** Definir la estrategia de respaldo y recuperación de la base de datos ante un desastre, independiente de las garantías transaccionales a nivel de consulta.

## Contexto y planteamiento del problema

Las transacciones atómicas protegen cada operación individual, pero no protegen contra un desastre que destruya la base de datos completa: un borrado accidental, una corrupción del disco, o una falla catastrófica de la instancia. Ante ese escenario, el sistema necesita poder restaurar la información de trazabilidad sin registros corruptos y en un tiempo acotado, ya que esos registros representan el historial oficial del producto.

¿Qué estrategia de respaldo permite restaurar la base de datos con la menor pérdida de información posible, dentro del tiempo de recuperación exigido?

## Impulsores de decisión

* QS-14 Falla durante una operación con alta concurrencia
* QS-15 Error de almacenamiento interno

## Opciones consideradas

* Respaldo lógico periódico
* Respaldo físico continuo con archivado de WAL (Point-in-Time Recovery)
* Solo replicación en caliente, sin respaldos históricos

## Resultado de la decisión

Opción elegida: **Respaldo físico continuo con archivado de WAL (Point-in-Time Recovery)**, porque permite restaurar la base de datos a un punto exacto en el tiempo, minimizando la pérdida de información frente a un respaldo lógico periódico que solo cubre hasta el último snapshot tomado. Además, reproducir el WAL desde un respaldo base es, en general, más rápido que restaurar un volcado completo a medida que la base de datos crece, ayudando a cumplir la meta de recuperación en menos de 60 minutos.

### Consecuencias positivas

* Pérdida de información mínima ante un desastre, al poder recuperar el estado casi hasta el momento exacto de la falla.
* Tiempo de restauración que escala mejor que un volcado lógico completo a medida que crece el volumen de datos.
* Complementa (no reemplaza) las garantías transaccionales, cubriendo el caso de una falla catastrófica de infraestructura.

### Consecuencias negativas

* Requiere más configuración inicial: archivado continuo de WAL y almacenamiento adicional para conservarlo.
* El proceso de restauración es más técnico de lo que se puede pensar, al requerir ubicar el punto exacto de recuperación.
* Necesita respaldos base periódicos sobre los cuales aplicar el WAL; sin ellos, el tiempo de recuperación crece indefinidamente.

## Pros y contras de las opciones

### Respaldo lógico periódico

Un volcado completo de la base de datos se genera de forma periódica y se restaura por completo ante un desastre.

* Bien, porque es simple de configurar y de restaurar, con una sola herramienta.
* Bien, porque el archivo resultante es portátil entre versiones o entornos de PostgreSQL.
* Malo, porque solo protege hasta el último volcado tomado; cualquier cambio posterior a ese momento se pierde ante un desastre.
* Malo, porque restaurar el volcado completo puede tardar más que la meta de 60 minutos a medida que la base de datos crece.

### Respaldo físico continuo con archivado de WAL (Point-in-Time Recovery)

Se archivan de forma continua los registros de escritura (WAL) de PostgreSQL sobre un respaldo base periódico, permitiendo reconstruir el estado de la base de datos hasta un instante exacto.

* Bien, porque minimiza la pérdida de información, pudiendo recuperar el estado casi hasta el momento exacto del desastre.
* Bien, porque reproducir el WAL desde un respaldo base suele ser más rápido que restaurar un volcado completo en bases de datos grandes.
* Malo, porque requiere más configuración inicial y almacenamiento adicional para el archivado continuo del WAL.
* Malo, porque la restauración es un proceso más técnico, al requerir ubicar el punto de recuperación exacto.

### Solo replicación en caliente, sin respaldos históricos

Una réplica de la base de datos se mantiene sincronizada en tiempo real con la instancia primaria, permitiendo un failover casi inmediato ante su caída.

* Bien, porque el failover ante la caída física del nodo primario es casi inmediato.
* Malo, porque replica también los errores lógicos.
* Malo, porque no cumple con la necesidad de un punto de recuperación histórico independiente que es necesario ante corrupción o errores internos.

## Enlaces

* [Relacionado con] [ADR-02: Selección de PostgreSQL como Motor de Bases de Datos](./adr-02-almacenar-datos.md)
