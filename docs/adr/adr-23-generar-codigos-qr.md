# Generar el Código QR de cada Lote

* **Estado:** Reemplazado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-20

**Historia técnica:** Definir dónde y con qué herramienta se genera la imagen del código QR de cada lote, distinto del mecanismo que garantiza que la URL codificada en él siga resolviendo indefinidamente, ya resuelto en el ADR-16.

## Contexto y planteamiento del problema

Cada lote de café necesita un código QR físico, imprimible en el empaque, que codifique la URL pública de su página de trazabilidad. Ese QR se genera una sola vez, en algún momento del ciclo de vida del lote, y una vez impreso no puede corregirse: necesita suficiente corrección de errores para tolerar el desgaste y la suciedad propios de un empaque físico en circulación. El ADR-08 ya definió que la generación del sello digital de verificación del lote ocurre de forma asíncrona sobre RabbitMQ; la generación del QR es una operación relacionada pero distinta, que este ADR resuelve por separado.

¿Dónde y con qué herramienta se genera la imagen del código QR de cada lote, garantizando suficiente corrección de errores para su impresión física?

## Impulsores de decisión

* QS-01 Consumidor consulta la información de trazabilidad de un café
* QS-12 Generación de sello digital
* TC-02 Permanencia de QR

## Opciones consideradas

* Backend en Go, en el mismo worker que genera el sello digital
* Cliente (frontend) bajo demanda
* Servicio externo (API de terceros)

## Resultado de la decisión

Opción elegida: **Servicio externo (API de terceros)**, concretamente **QuickChart** (`api.quickchart.io/qr`), porque es la opción más simple de integrar (una sola solicitud HTTP con la URL del lote como parámetro, sin sumar una librería ni lógica de generación de imágenes al backend), y porque, al ser un proyecto de código abierto, deja abierta la posibilidad de auto-alojarlo más adelante sin cambiar la integración ya construida, de forma similar a como el ADR-07 dejó abierta la migración de MinIO hacia otro proveedor compatible con S3.

### Consecuencias positivas

* Es la opción de menor esfuerzo de implementación: una sola solicitud HTTP, sin sumar una librería de generación de QR ni lógica de rasterizado al backend.
* Al ser de código abierto, existe una ruta clara de auto-alojamiento a futuro si el volumen de lotes o la dependencia de un tercero externo lo justifican, sin rediseñar la integración.
* No consume recursos de cómputo propios (CPU/memoria) del worker del backend ni del navegador del cliente para generar la imagen.

### Consecuencias negativas

* Introduce una dependencia externa para un componente crítico del sistema, en contra del patrón que el resto del proyecto ha seguido consistentemente de preferir herramientas propias o autoalojadas cuando existe una alternativa madura, como MinIO sobre S3 o Retraced sobre WorkOS.
* La disponibilidad de QuickChart queda fuera del control del proyecto: una caída del servicio externo justo en el momento de sellar un lote bloquearía la generación de su QR, salvo que se proteja la llamada con el Circuit Breaker ya definido en el ADR-12.
* No deja, por sí sola, un artefacto propio y versionado del QR exacto generado en cada momento; el resultado de la API debe descargarse y guardarse explícitamente en MinIO (ADR-07) si se quiere conservar un registro auditable.

## Pros y contras de las opciones

### Backend en Go (worker del sello digital)

El propio worker que calcula el sello digital de verificación, ya ejecutándose de forma asíncrona sobre RabbitMQ (ADR-08), genera también la imagen del QR con una librería como `go-qrcode`.

* Bien, porque queda como parte del mismo punto de generación que el sello digital, sin coordinar dos procesos distintos.
* Bien, porque no depende de la disponibilidad de un tercero externo para completar el sellado de un lote.
* Bien, porque el artefacto generado puede guardarse directamente en MinIO junto con el resto de evidencias del lote, con control total sobre su formato y nivel de corrección de errores.
* Malo, porque suma una responsabilidad más al worker del sello digital, aunque sea una operación liviana.

### Cliente (frontend) bajo demanda

El navegador genera la imagen del QR en el momento en que alguien lo visualiza o descarga desde el panel, con una librería JavaScript como `qrcode.react`.

* Bien, porque no consume ningún recurso del backend para generarlo.
* Bien, porque es la opción más simple de implementar dentro del frontend ya elegido (React) en el ADR-17.
* Malo, porque depende de que el frontend arme correctamente, en cada renderizado, la URL versionada ya definida en el ADR-16.
* Malo, porque no deja un artefacto fijo y auditable del QR exacto que terminó impreso; si un consumidor reporta un QR roto, no hay un registro histórico de qué imagen se generó en ese momento.

### Servicio externo (API de terceros)

Un proveedor externo recibe la URL del lote y devuelve la imagen del QR ya generada, sin lógica propia de generación en el backend ni en el cliente.

* Bien, porque es la opción de menor esfuerzo de implementación, con una sola solicitud HTTP.
* Bien, porque no consume recursos de cómputo propios para generar la imagen.
* Bien, porque, en el caso de QuickChart, al ser de código abierto deja abierta la ruta de auto-alojarlo después sin rediseñar la integración.
* Malo, porque introduce una dependencia externa para un componente crítico, en contra del patrón de infraestructura propia seguido en el resto del proyecto.
* Malo, porque la disponibilidad del QR de un lote queda ligada a la disponibilidad de un tercero, salvo que se proteja con el Circuit Breaker del ADR-12.

## Enlaces

* [Reemplazado por] [ADR-25: Migrar la Generación de Códigos QR al Servicio de QR](./adr-25-migrar-generacion-qr.md)
* [Relacionado con] [ADR-16: Versionar las URLs de los Códigos QR](./adr-16-versionar-urls.md)
* [Relacionado con] [ADR-08: Encolar Tareas Pesadas](./adr-08-encolar-tareas.md)
* [Relacionado con] [ADR-07: Almacenar Evidencias como Objetos](./adr-07-almacenar-evidencias.md)
* [Relacionado con] [ADR-12: Proteger Servicios Externos](./adr-12-proteger-servicios.md)
