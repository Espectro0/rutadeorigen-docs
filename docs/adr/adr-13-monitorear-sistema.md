# Monitorear el Sistema

* **Estado:** Propuesto
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-05

**Historia técnica:** Definir el mecanismo de métricas y alertas que permite detectar de forma proactiva fallas de conexión, fallas de servicios externos, picos de tráfico y accesos sospechosos.

## Contexto y planteamiento del problema


Dentro del sistema, es importante observar su comportamiento en tiempo real ya que esto nos permite detectar errores o posibles fallos, y detectar manualmente de manera más sencilla que ocurrió exactamente.

¿Qué mecanismo de métricas y alertas permite observar en tiempo real el comportamiento del sistema frente a estos escenarios?

## Impulsores de decisión

* QS-02 Pérdida de conexión durante un registro
* QS-10 Fallo de un servicio externo
* QS-13 Solicitudes concurrentes masivas
* QS-20 Inicio de sesión desde contexto sospechoso

## Opciones consideradas

* Prometheus + Grafana
* Servicio externo
* Solo logs estructurados en archivos/stdout, sin métricas ni alertas dedicadas

## Resultado de la decisión

Opción elegida: **Prometheus + Grafana**, porque se despliega junto al resto de la infraestructura propia ya definida en los contenedores orquestados de la plataforma, sin costo recurrente de licencia ni una nueva dependencia externa crítica, y da control total sobre las métricas y alertas configuradas.

### Consecuencias positivas

* Sin costo recurrente de licencia, a diferencia de un servicio externo administrado.
* Se despliega junto al resto de la infraestructura propia, bajo el mismo modelo operativo.
* Control total sobre qué métricas se recolectan y cómo se configuran las alertas para cada escenario de calidad.
* No añade una nueva dependencia externa crítica que, a su vez, necesitaría protegerse con el Circuit Breaker.

### Consecuencias negativas

* Añade carga operativa adicional: mantener, respaldar y dimensionar el propio stack de monitoreo, para un equipo de una sola persona.
* El stack de observabilidad se vuelve, a su vez, un componente que puede fallar y que también requiere ser monitoreado.
* No ofrece dashboards ni alertas preconfiguradas de fábrica; su configuración inicial recae completamente en el equipo del proyecto.

## Pros y contras de las opciones

### Prometheus + Grafana

Prometheus recolecta métricas del backend y de la infraestructura, y Grafana las visualiza y dispara alertas configuradas manualmente.

* Bien, porque no tiene costo recurrente de licencia.
* Bien, porque se despliega junto al resto de la infraestructura propia, bajo el mismo modelo de contenedores de la plataforma.
* Bien, porque da control total sobre las métricas y alertas configuradas.
* Malo, porque añade carga operativa adicional para mantenerlo, respaldarlo y dimensionarlo.
* Malo, porque su configuración inicial (dashboards, alertas) recae completamente en el equipo del proyecto.

### Servicio externo

Un proveedor externo recolecta métricas y gestiona dashboards y alertas.

* Bien, porque no requiere infraestructura propia que mantener.
* Bien, porque ofrece dashboards y alertas listas de fábrica, con dimensionado y escalado a cargo del proveedor.
* Malo, porque genera un costo recurrente adicional.
* Malo, porque introduce una nueva dependencia externa que, si se usa para decisiones críticas, también debería protegerse con el Circuit Breaker.

### Solo logs estructurados en archivos/stdout, sin métricas ni alertas dedicadas

El sistema únicamente registra eventos como logs estructurados, sin un mecanismo de métricas ni alertas en tiempo real.

* Bien, porque no requiere infraestructura adicional y su complejidad es mínima.
* Malo, porque no permite detectar de forma proactiva los errores, ni cuantificar en tiempo real las métricas.
* Malo, porque obliga a revisar logs manualmente después de que el problema ya ocurrió, en contra del propósito de este ADR.

## Enlaces

* [Relacionado con] [ADR-05: Reintentar y Persistir Temporalmente ante Pérdida de Conexión](./adr-05-reintentar-conexion.md)
* [Relacionado con] [ADR-09: Autenticar Usuarios](./adr-09-autenticar-usuarios.md)
* [Relacionado con] [ADR-12: Proteger Servicios Externos](./adr-12-proteger-servicios.md)
* [Relacionado con] [ADR-15: Orquestar Contenedores para Escalado Horizontal](./adr-15-orquestar-contenedores.md)
