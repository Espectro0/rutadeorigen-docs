# Uso de Redis como memoria caché

* **Estado:** Propuesto
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-02

**Historia técnica:** Optimización del acceso a la información pública de trazabilidad consultada mediante códigos QR.

## Contexto y planteamiento del problema

Uno o varios consumidores ingresan a la información pública de trazabilidad de café mediante el código QR impreso en el empaque. Esta información es solo de lectura, y no tiene cambios significativos o frecuentes una vez publicado el lote, pero puede ser consultada de forma masiva y simultánea: un solo lote puede tener cientos o miles de bolsas en circulación, y una campa de la marca en redes sociales puede generar picos de trafico concentrados en pocos minutos.

Realizar las consultas directamente a la base de datos implica ejecutar múltiples consultas relacionadas para reconstruir la misma respuesta una y otra vez. Esto incrementa la latencia percibida por el consumidor, consume recursos de la base de datos de manera innecesaria y limita la capacidad de atender solicitudes concurrentes.ta la capacidad de atender solicitudes concurrentes.

¿Es posible incorporar una capa intermedia entre la base de datos y el servidor que almacene temporalmente las respuestas ya construidas, de forma que las consultas repetidas se resuelvan sin llegar hasta la base de datos?

## Impulsores de decisión

* QS-01 Consumidor consulta la información de trazabilidad de un café
* QS-13 Solicitudes concurrentes masivas

## Opciones consideradas

* Redis
* Memcached
* Memoria local ( In-Memory )

## Resultado de la decisión

Opción elegida: **Redis**, porque es la única alternativa que combina persistencia opcional, estructuras de datos avanzadas y caché compartida entre múltiples instancias del servidor, lo cual es indispensable dado que la plataforma se despliega en contenedores orquestados con escalado horizontal automático donde no existe garantía de que dos solicitudes consecutivas sean atendidas por la misma instancia.

Adicionalmente, Redis permite implementar invalidación selectiva de la caché cuando un lote es editado, control de expiración por clave (TTL), y a futuro habilita casos de uso complementarios como limitación de tasa de peticiones (rate limiting) para proteger la API pública, y almacenamiento temporal de sesiones o trabajos en segundo plano.

### Consecuencias positivas

* Reducción de la respuesta en la consulta pública de trazabilidad, al retornar respuesta previamente construidas sin recalcularlas de la base de datos.
* Disminución de la carga de peticiones sobre la base de datos.
* Cpacidad de atender picos altos de tráfico.
* Caché compartida entre diferentes instancias del servidor.

### Consecuencias negativas

* Introducción de un componente adicional en la arquitectura, aumentando complejidad operativa y superficie de fallo.
* Costo adicional de infraestructura.
* Requiere actualización de la información en cache cada que la información sufra algún cambio.
* Dependencia externa cuya indisponibilidad debe ser manejada de manera controlada.

## Pros y contras de las opciones

### Redis

Almacén con datos en memoria de código abierto, es utilizado de mayor manera dentro de la industria, es disponible como servicio mediante diferentes proveedores o asi mismo es posible desplegarlo de manera automática en un servidor propio. Soporta estructuras de datos complejas, persistencia y replicación.

* Bien, porque ofrece cache compartida entre múltiples instancias, lo que es indispensable para un escalamiento horizontal.
* Bien, porque soporta expiración automática por clave (TTL), lo que permite definir políticas de frecura distintas para cada contenido.
* Bien, porque ofrece persistencia, evitando que en un reinicio del servicio provoque multiples consultas inmediatas a la base de datos.
* Bien, porque nos permite realizar rate limiting, colas de trabajo, contadores de analíticas sin necesidad de otra estructura.
* Bien, porque cuenta con multiples proveedores y una opción de autodespliegue adecuados para el crecimiento con el software.
* Bien, porque su adopción implica documentación extensa, y una comunidad activa.
* Malo, porque introduce un componente adicional que se necesita monitorear, mantener y presupuestar.
* Malo, porque implica una curva de aprendizaje mayor que otras alternativas más simples. 

### Memcached

Sistema de caché distribuida en memoria, diseñado especificamente para almacenamiento de pares clave-valor simples con alto rendimiento.

* Bien, porque es rápido y eficiente en el uso de memoria para casos de caché simples.
* Bien, porque su operación tiene una curva de aprendizaje baja, lo que lo hace sencillo de comprender y operar. 
* Bien, porque también ofrece caché compartida entre múltiples instancias del servidor.
* Malo, porque no ofrece persistencia ante reinicios, cualquier fallo puede hacer que se pierda la caché y desestabilice la base de datos con multiples consultas.
* Malo, porque solo soporta pares simples (clave-valor) sin estructuras que permitan mecanismos como rate limit, colas o contadores.
* Malo, porque no permite realizar validaciones de un elemento ya en caché al ser actualizado.
* Malo, porque no se cuenta con un ecosistema amplio de proveedores.

### Memoria local (In-Memory)

Almacenamiento de cache dentro de la propia memoria del servidor, utilizando estructuras de datos en memoria y lenguajes de programación.

* Bien, porque no requiere infraestructura adicional ni genera costos externos.
* Bien, porque ofrece la menor latencia posible al no involucrar servidores externos.
* Bien, porque tiene una implementación inmediata y no añade dependencias complejas al software.
* Malo, porque al escalar el sistema horizontalmente cada servicio tendra su memoria aislada lo que daña completamente la persistencia y genera inconsistencias.
* Malo, porque si se generan varias instancias no se permite actualizar la caché inmediatamente una información sea modificada.
* Malo, porque se genera un mayor consumo de recursos (memoria RAM) lo que puede generar bajo rendimiento del servidor o agotamiento de memoria.
* Malo, porque no escala con el sistema, es una opción buena solo para sistemas pequeños con bajo uso de recursos o entornos de prueba.

## Enlaces

* [Relacionado con] [ADR-02: Selección de PostgreSQL como Motor de Bases de Datos](./adr-02-almacenar-datos.md)
* [Relacionado con] [ADR-03: Selección de Go como Lenguaje y Runtime del Backend](./adr-03-definir-runtime.md)
* [Relacionado con] [ADR-08: Encolar Tareas Pesadas](./adr-08-encolar-tareas.md)
* [Relacionado con] [ADR-15: Orquestar Contenedores para Escalado Horizontal](./adr-15-orquestar-contenedores.md)
* [Requiere] Definición de la estrategia de invalidación de caché ante edición de lotes publicados
* [Requiere] Definición del comportamiento de degradación controlada ante indisponibilidad de Redis