# Centralizar los Logs del Sistema

* **Estado:** Propuesto
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-10-11

**Historia técnica:** Definir dónde se recolectan, almacenan y consultan los logs de los servicios del backend (ADR-24) desplegados en k3s (ADR-15), de modo que una falla pueda diagnosticarse siguiendo una solicitud a través de varios servicios, aun después de que los pods involucrados se hayan reiniciado o reemplazado.

## Contexto y planteamiento del problema

El ADR-13 definió Prometheus + Grafana para métricas y alertas, pero dejó fuera los logs: las métricas indican *que* algo falla, no *por qué*. El ADR-24 dividió el backend en microservicios que se comunican por HTTP y por eventos en RabbitMQ, y señaló como consecuencia que depurar una falla exige propagar un identificador de correlación y centralizar los logs sobre la plataforma de monitoreo.

Hoy cada servicio escribe logs estructurados en JSON por la salida estándar, con el nombre del servicio y el identificador de correlación de cada solicitud. Sin embargo, esos logs solo viven dentro de cada contenedor: se consultan pod por pod con `kubectl logs`, se pierden cuando el pod se reinicia o se reemplaza en un despliegue (que con el ADR-26 ocurre en cada merge a `main`), y no permiten seguir una misma solicitud a través de varios servicios.

¿Dónde se centralizan los logs del sistema para poder consultarlos de forma unificada, conservarlos más allá de la vida de cada pod y correlacionarlos con las métricas ya definidas?

## Impulsores de decisión

* QS-10 Fallo de un servicio externo
* QS-14 Falla durante una operación con alta concurrencia
* QS-02 Pérdida de conexión durante un registro
* TC-10 Prácticas DevOps
* BC-06 Costo de operación sostenible

## Opciones consideradas

* Grafana Loki
* Elasticsearch / OpenSearch (stack EFK)
* Servicio externo administrado
* Solo los logs de los contenedores (`kubectl logs`), sin centralizar

## Resultado de la decisión

Opción elegida: **Grafana Loki**, porque indexa únicamente un conjunto pequeño de etiquetas y guarda el contenido de los logs comprimido, lo que le permite correr con pocos recursos en el mismo servidor que atiende producción. Además, se consulta desde el mismo Grafana del ADR-13, de modo que métricas y logs quedan en una sola interfaz, sin costo de licencia ni una nueva dependencia externa.

La implementación queda así:

* **Despliegue:** Loki en modo monolítico (un solo proceso), suficiente para un clúster de un nodo. Se instala como una Application de Argo CD separada de la de los servicios, igual que el resto de la infraestructura autoalojada.
* **Recolección:** Grafana Alloy corre como DaemonSet, lee los logs de los contenedores del nodo y los envía a Loki. Los servicios no cambian: siguen escribiendo JSON por la salida estándar.
* **Etiquetas:** solo etiquetas de baja cardinalidad: `namespace`, `pod`, `container` y `service` (desde `app.kubernetes.io/name`). El identificador de correlación **no** es una etiqueta; se filtra en la consulta extrayéndolo del JSON (`| json | correlation_id="…"`), para no disparar el tamaño del índice.
* **Almacenamiento:** un volumen persistente local del nodo.
* **Retención:** 14 días; los logs más antiguos se eliminan automáticamente.
* **Consulta:** Grafana agrega Loki como fuente de datos junto a Prometheus, y se expone según el ADR-21.

### Consecuencias positivas

* Los logs sobreviven a reinicios y despliegues de los pods, y se conservan durante el periodo de retención.
* Una sola consulta por identificador de correlación muestra el recorrido de una solicitud a través de todos los servicios.
* Métricas y logs se ven en el mismo Grafana: desde un pico en un panel se puede saltar a los logs de ese intervalo.
* Bajo consumo de recursos frente a alternativas que indexan el texto completo.
* No requiere cambios en el código: los servicios ya escriben logs estructurados por la salida estándar.

### Consecuencias negativas

* Suma dos componentes más que operar (Loki y Alloy) y que también pueden fallar.
* Las búsquedas de texto libre son más lentas que en un motor que indexa todo el contenido, porque Loki recorre los logs del intervalo consultado.
* Con almacenamiento en el disco del nodo, los logs se pierden si se pierde el servidor, y su volumen compite por disco con PostgreSQL.
* Un mal uso de las etiquetas (por ejemplo, poner identificadores únicos como etiqueta) degrada el rendimiento; requiere disciplina al configurar Alloy.

## Pros y contras de las opciones

### Grafana Loki

Sistema de agregación de logs que indexa solo etiquetas y guarda el contenido comprimido; se consulta con LogQL desde Grafana.

* Bien, porque consume pocos recursos y puede correr en un servidor de un solo nodo.
* Bien, porque se integra de forma nativa con el Grafana del ADR-13.
* Bien, porque no tiene costo de licencia y se despliega con el mismo modelo que el resto de la infraestructura.
* Malo, porque la búsqueda de texto libre es más lenta al no indexar el contenido.
* Malo, porque añade carga operativa: mantener Loki y su recolector.

### Elasticsearch / OpenSearch (stack EFK)

Motor de búsqueda que indexa todo el contenido de los logs, con un recolector (Fluent Bit) y Kibana u OpenSearch Dashboards como interfaz.

* Bien, porque permite búsquedas de texto completo muy rápidas y análisis complejos.
* Bien, porque es un estándar ampliamente conocido para centralizar logs.
* Malo, porque consume mucha memoria (del orden de gigabytes solo para el motor), difícil de sostener en el mismo servidor de producción.
* Malo, porque añade una interfaz distinta de Grafana, separando métricas y logs.
* Malo, porque su operación (índices, réplicas, retención) es más compleja para un equipo de una persona.

### Servicio externo administrado

Un proveedor recibe, almacena y permite consultar los logs.

* Bien, porque no requiere infraestructura propia que mantener ni dimensionar.
* Bien, porque ofrece búsquedas y paneles listos.
* Malo, porque genera un costo recurrente, o queda limitado por los topes de volumen y retención del plan gratuito.
* Malo, porque saca los logs, que pueden contener datos de los usuarios, a un tercero, lo que debe evaluarse frente a la Ley 1581/2012 (BC-03).
* Malo, porque introduce una nueva dependencia externa.

### Solo los logs de los contenedores (`kubectl logs`), sin centralizar

Los logs permanecen en cada contenedor y se consultan pod por pod.

* Bien, porque no requiere infraestructura adicional.
* Malo, porque los logs se pierden al reiniciarse o reemplazarse un pod, lo que ocurre en cada despliegue.
* Malo, porque no permite seguir una solicitud a través de varios servicios.
* Malo, porque exige acceso directo al clúster para depurar, en contra de las prácticas DevOps del TC-10.

## Enlaces

* [Complementa] [ADR-13: Monitorear el Sistema](./adr-13-monitorear-sistema.md)
* [Relacionado con] [ADR-24: Definir arquitectura concreta para el Backend](./adr-24-arquitectura-backend.md)
* [Relacionado con] [ADR-15: Orquestar Contenedores para Escalado Horizontal](./adr-15-orquestar-contenedores.md)
* [Relacionado con] [ADR-21: Exponer los Servicios mediante un Reverse Proxy](./adr-21-exponer-servicios.md)
* [Relacionado con] [ADR-26: Automatizar la Integración y el Despliegue Continuos](./adr-26-automatizar-despliegues.md)
