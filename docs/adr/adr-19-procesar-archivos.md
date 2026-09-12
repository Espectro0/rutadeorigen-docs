# Procesar Archivos antes de Almacenarlos en MinIO

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-12

**Historia técnica:** Definir cómo se implementa el procesamiento que un worker ejecuta sobre las fotos, documentos y videos de evidencia antes de que su versión final quede en MinIO. La ejecución asíncrona de ese worker mediante una cola ya quedó resuelta en el ADR-08; este ADR se enfoca solo en qué hace el worker con el archivo, no en cómo se dispara.

## Contexto y planteamiento del problema

Los archivos que se cargan como evidencia (fotos, documentos y videos) llegan en su formato y tamaño originales. Antes de reemplazar el objeto en MinIO por su versión final, el worker necesita: reducir su tamaño, validar que el tipo y tamaño de archivo sean los esperados, escanear que no contengan software malicioso, y eliminar metadatos que no deberían exponerse.

## Impulsores de decisión

* QS-11 Almacenamiento de archivos de un proceso de trazabilidad

## Opciones consideradas

* Orquestar herramientas de código abierto especializadas por tipo de operación (ej. ffmpeg para video, una librería de imágenes para fotos, ClamAV para escaneo de seguridad)
* Delegar el procesamiento a un servicio administrado de medios en la nube (ej. Cloudinary, imgix, AWS Elemental MediaConvert)
* Implementar cada operación de procesamiento desde cero, sin librerías ni herramientas externas

## Resultado de la decisión

Opción elegida: **Orquestar herramientas de código abierto especializadas por tipo de operación**, porque reutiliza soluciones maduras y gratuitas para cada tarea, evitando tanto el costo recurrente de un servicio administrado como el esfuerzo inviable de reimplementar esas operaciones desde cero, consistente con el criterio de presupuesto limitado ya usado en el ADR-07.

### Consecuencias positivas

* Cada operación se apoya en una herramienta madura y especializada en esa tarea concreta, en vez de código propio para algo tan complejo como un codec de video.
* Sin costo recurrente por licencia ni por volumen procesado, coherente con el resto de decisiones de infraestructura del proyecto.
* Las herramientas se actualizan de forma independiente entre sí y de la lógica propia del worker.

### Consecuencias negativas

* El worker debe integrar y mantener varias herramientas distintas en vez de una sola solución integrada.
* El rendimiento y los recursos que consume el procesamiento quedan bajo responsabilidad de dimensionar el propio servidor del worker, sin la elasticidad de un servicio administrado.
* Cada herramienta puede introducir sus propias vulnerabilidades o requerir actualizaciones de seguridad que hay que rastrear por separado.

## Pros y contras de las opciones

### Orquestar herramientas de código abierto especializadas por tipo de operación

El worker no reimplementa compresión, transcodificación ni escaneo de malware; para cada operación invoca una herramienta o librería madura y especializada en esa tarea concreta, y coordina el resultado.

* Bien, porque reutiliza herramientas maduras, probadas por años en producción, en vez de reimplementar algoritmos de compresión o codecs de video.
* Bien, porque son de código abierto y autoalojadas, sin costo recurrente por licencia ni por volumen procesado.
* Bien, porque cada herramienta se actualiza de forma independiente sin tener que tocar la lógica propia del worker.
* Malo, porque cada herramienta tiene su propia forma de configurarse e instalarse, sumando piezas distintas que mantener en vez de una sola solución integrada.
* Malo, porque el rendimiento y los recursos (CPU/memoria) que consume cada herramienta quedan bajo responsabilidad del propio servidor donde corre el worker, sin la elasticidad de un servicio administrado.

### Delegar el procesamiento a un servicio administrado de medios en la nube

El worker sube el archivo original a un servicio externo especializado en procesamiento de imágenes/video, que devuelve la versión optimizada, validada o escaneada.

* Bien, porque ofrece procesamiento de nivel productivo (transformaciones, transcodificación) sin mantener infraestructura ni herramientas propias.
* Bien, porque escala automáticamente con el volumen de archivos, sin dimensionar servidores propios.
* Malo, porque genera un costo recurrente por volumen procesado y/o almacenado, contrario al presupuesto limitado del proyecto (mismo criterio que descartó S3/GCS en el ADR-07).
* Malo, porque no todos estos servicios cubren a la vez las cuatro necesidades (optimización, validación, escaneo de seguridad y limpieza de metadatos); normalmente hay que combinar varios proveedores.
* Malo, porque el archivo original viaja fuera de la infraestructura propia para ser procesado, algo a evaluar si contiene datos sensibles antes de la limpieza de metadatos.

### Implementar cada operación de procesamiento desde cero, sin librerías ni herramientas externas

El equipo escribe su propia lógica de compresión de imágenes, transcodificación de video, detección de malware y limpieza de metadatos.

* Bien, porque no depende de herramientas ni servicios externos, con control total sobre el comportamiento exacto de cada paso.
* Malo, porque reimplementar compresión de imágenes, codecs de video o detección de malware es un esfuerzo enorme e inviable para un equipo de una sola persona, con un resultado casi siempre peor que el de herramientas maduras con años de desarrollo.
* Malo, porque el mantenimiento y la corrección de errores de seguridad (ej. en el escaneo de malware) quedarían enteramente bajo responsabilidad propia, sin el respaldo de una comunidad o proveedor.

## Enlaces

* [Relacionado con] [ADR-07: Almacenar Evidencias como Objetos](./adr-07-almacenar-evidencias.md)
* [Relacionado con] [ADR-08: Encolar Tareas Pesadas](./adr-08-encolar-tareas.md)
