# Almacenar Evidencias como Objetos

* **Estado:** Propuesto
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-05

**Historia técnica:** Definir dónde y cómo se almacenan los archivos de evidencia que respaldan cada proceso de trazabilidad, separando ese almacenamiento de la base de datos relacional.

## Contexto y planteamiento del problema

Cada proceso de trazabilidad se respalda con evidencias multimedia cargadas principalmente desde el campo, con conectividad intermitente. Guardar esos archivos directamente en PostgreSQL no es viable. Se necesita un mecanismo de almacenamiento de objetos, propio y separado de la base de datos relacional, que identifique cada archivo de forma única, conserve sus metadatos y garantice que el 98% de las cargas se completen correctamente sin pérdida de archivos ya confirmados como almacenados.

## Impulsores de decisión

* QS-11 Almacenamiento de archivos de un proceso de trazabilidad

## Opciones consideradas

* MinIO (almacenamiento de objetos autoalojado, compatible con S3)
* Amazon S3
* Google Cloud Storage
* Cloudflare R2

## Resultado de la decisión

Opción elegida: **MinIO**, porque al ser autoalojado no genera un costo recurrente en un proyecto de presupuesto limitado y equipo de una sola persona, y al implementar la misma API que S3 permite usar desde el inicio el mismo SDK y el mismo modelo de buckets/objetos que un proveedor cloud, dejando abierta la posibilidad de migrar la carga hacia un servicio compatible con S3 más adelante si el proyecto crece, sin tener que rediseñar la capa de almacenamiento.

### Consecuencias positivas

* Sin costo recurrente de almacenamiento en la nube durante el desarrollo y la operación inicial del proyecto.
* Compatibilidad con la API de S3 desde el día uno: el código de la aplicación queda desacoplado del proveedor concreto.
* Se despliega junto al resto de la infraestructura propia, sin depender de una cuenta ni credenciales de un tercero.
* Ruta de migración clara: mover los objetos a un proveedor compatible con S3 no exige cambiar el SDK ni el modelo de datos, solo la configuración de conexión.

### Consecuencias negativas

* La disponibilidad, el respaldo y el escalado del propio almacenamiento quedan bajo responsabilidad del proyecto, sin el SLA de un proveedor administrado.
* Requiere configuración y mantenimiento operativo adicional que un servicio administrado resolvería de forma nativa.
* El servidor donde corre MinIO debe dimensionarse para el crecimiento del volumen de evidencias, a diferencia de un servicio cloud con capacidad prácticamente ilimitada.

## Pros y contras de las opciones

### MinIO

Servicio de almacenamiento de objetos de código abierto, autoalojado, que implementa la misma API que Amazon S3.

* Bien, porque no tiene costo de licencia ni de uso, ajustándose al presupuesto limitado del proyecto.
* Bien, porque al ser compatible con S3 no genera dependencia de código hacia un proveedor cloud específico.
* Bien, porque se despliega junto al resto de la infraestructura propia, sin depender de una cuenta externa.
* Malo, porque la disponibilidad y los respaldos del almacenamiento son responsabilidad del proyecto, sin SLA de un tercero.
* Malo, porque el escalado depende de dimensionar manualmente el servidor donde corre, a diferencia de un servicio administrado.

### Amazon S3

Servicio de almacenamiento de objetos administrado por Amazon Web Services, estándar de facto del mercado.

* Bien, porque ofrece una durabilidad y disponibilidad muy altas sin que el proyecto administre infraestructura.
* Bien, porque cuenta con SDKs maduros para Go y una amplia documentación e integraciones.
* Malo, porque genera un costo recurrente desde el primer archivo almacenado, algo difícil de sostener en un proyecto de presupuesto limitado.
* Malo, porque introduce una dependencia directa de un proveedor cloud específico.

### Google Cloud Storage

Servicio de almacenamiento de objetos administrado por Google Cloud.

* Bien, porque ofrece garantías de disponibilidad y durabilidad similares a S3, sin administración de infraestructura propia.
* Malo, porque genera un costo recurrente igual que S3, sin una ventaja técnica clara sobre las demás opciones para este proyecto.
* Malo, porque introduce un segundo proveedor cloud distinto al resto de la infraestructura, sin justificación técnica suficiente.

### Cloudflare R2

Servicio de almacenamiento de objetos administrado por Cloudflare, compatible con la API de S3 y sin costo de egreso de datos.

* Bien, porque no cobra por la salida de datos, lo que favorece servir evidencias multimedia a consumidores públicos vía QR.
* Bien, porque al ser compatible con S3 permitiría una migración similar a la que ya habilita MinIO.
* Malo, porque sigue siendo un servicio administrado por un tercero, con menor recorrido y adopción que Amazon S3.
* Malo, porque tampoco elimina el costo por almacenamiento y por operaciones, un factor sensible dado el presupuesto del proyecto.

## Enlaces

* [Relacionado con] [ADR-02: Selección de PostgreSQL como Motor de Bases de Datos](./adr-02-almacenar-datos.md)
* [Relacionado con] [ADR-06: Respaldar y Recuperar Datos ante Desastres](./adr-06-respaldar-datos.md)
