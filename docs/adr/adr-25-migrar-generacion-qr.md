# Migrar la Generación de Códigos QR al Servicio de QR

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-10-08

**Historia técnica:** Reemplazar la decisión del ADR-23 de generar la imagen del código QR mediante la API externa de QuickChart por una generación propia dentro del servicio de QR definido en el ADR-24, eliminando la dependencia de un tercero en un componente crítico del sistema.

## Contexto y planteamiento del problema

El ADR-23 eligió QuickChart por ser la opción de menor esfuerzo de integración, pero reconoció tres consecuencias negativas: introduce una dependencia externa para un componente crítico, en contra del patrón de infraestructura propia seguido en el resto del proyecto (MinIO sobre S3, Retraced sobre WorkOS); la disponibilidad del QR de un lote queda ligada a la de un tercero; y no deja por sí sola un registro reproducible de la imagen generada. Posteriormente, el ADR-24 definió un servicio de QR independiente, responsable de los códigos, la tabla de redirecciones y la resolución de escaneos. Con un servicio propio dedicado a los QR, generar la imagen dentro de él deja de sumar responsabilidades a otro componente y elimina la necesidad de llamar a un tercero.

¿Dónde y con qué herramienta se genera ahora la imagen del código QR, sin depender de un servicio externo y garantizando suficiente corrección de errores para su impresión física?

## Impulsores de decisión

* QS-01 Consumidor consulta la información de trazabilidad de un café
* QS-10 Fallo de un servicio externo
* TC-02 Permanencia de QR
* BC-06 Costo de operación sostenible

## Opciones consideradas

* Mantener QuickChart como API externa (decisión original del ADR-23)
* Generar la imagen dentro del servicio de QR con una librería de Go
* Autoalojar QuickChart dentro del clúster de k3s
* Generar la imagen en el cliente (frontend) bajo demanda — ya descartada en el ADR-23

## Resultado de la decisión

Opción elegida: **Generar la imagen dentro del servicio de QR con una librería de Go**, porque elimina la dependencia de un tercero en un componente crítico sin sumar un servicio adicional que operar, y porque la generación de un QR es una operación liviana que el propio servicio responsable de los códigos puede resolver directamente, consistente con el patrón de infraestructura propia seguido en el resto del proyecto.

La generación sigue estas reglas:

* **Librería:** una librería de Go que exponga la matriz de módulos del código (por ejemplo, `skip2/go-qrcode`), a partir de la cual el servicio renderiza los formatos de salida.
* **Formatos:** SVG, vectorial y apto para imprenta en cualquier tamaño, y PNG, para la vista previa en el panel y el uso en web.
* **Corrección de errores:** nivel **Q (25 %)**, como equilibrio entre la tolerancia al desgaste y la suciedad propios de un empaque físico y la densidad del código, que permite imprimirlo en tamaños pequeños.
* **Momento de generación:** bajo demanda, cuando se solicita la imagen desde el panel o la API, sin almacenarla de forma permanente.
* **Caché:** la imagen generada se guarda en Redis (ADR-01), con una llave compuesta por el contenido codificado, el formato, el tamaño y el nivel de corrección de errores, de modo que solicitudes repetidas no la regeneren.
* **Reproducibilidad:** como la generación es determinista para un mismo contenido y unos mismos parámetros, el servicio registra junto a cada código la versión de la librería y los parámetros de generación, lo que permite reconstruir exactamente la imagen entregada en cualquier momento sin conservar el archivo.

### Consecuencias positivas

* Elimina la dependencia de un tercero externo en un componente crítico: la generación del QR ya no puede fallar por una caída de QuickChart, ni necesita protegerse con el Circuit Breaker del ADR-12.
* No suma ningún servicio de infraestructura adicional que desplegar ni mantener: la generación vive dentro del servicio de QR ya definido en el ADR-24.
* Da control total sobre el formato, el tamaño y el nivel de corrección de errores de la imagen.
* Generar bajo demanda con caché en Redis evita ocupar almacenamiento permanente en MinIO por cada código, alineado con el control de costos del BC-06 y el BC-07.
* Registrar la versión de la librería y los parámetros resuelve la falta de un registro auditable señalada en el ADR-23, sin guardar cada imagen.

### Consecuencias negativas

* Suma a mantener una librería de generación y el código que renderiza SVG y PNG a partir de la matriz del QR, frente a una sola solicitud HTTP con QuickChart.
* La reproducibilidad depende de conservar la versión exacta de la librería y los parámetros: si se actualiza la librería y cambia la forma de construir la matriz, una regeneración podría producir una imagen distinta a la entregada originalmente, aunque codifique el mismo contenido.
* Ante una pérdida de la caché en Redis, las imágenes se regeneran en las solicitudes siguientes, consumiendo CPU del servicio de QR de forma puntual.
* La generación consume recursos de cómputo propios del servicio de QR, aunque sea una operación liviana.

## Pros y contras de las opciones

### Mantener QuickChart como API externa

Conservar la decisión original del ADR-23: el servicio de QR solicita la imagen a `api.quickchart.io/qr` con la URL del código como parámetro.

* Bien, porque no exige ningún cambio sobre lo ya decidido.
* Bien, porque es la opción con menor código propio: una sola solicitud HTTP.
* Malo, porque mantiene una dependencia externa en un componente crítico, en contra del patrón de infraestructura propia del proyecto.
* Malo, porque la disponibilidad de la generación del QR sigue ligada a la de un tercero.

### Generar la imagen dentro del servicio de QR con una librería de Go

El servicio de QR genera directamente la imagen a partir del contenido del código, usando una librería de Go.

* Bien, porque elimina por completo la dependencia de un tercero para generar el QR.
* Bien, porque no suma ningún servicio adicional que operar: la generación queda dentro del servicio responsable de los códigos.
* Bien, porque da control total sobre el formato y el nivel de corrección de errores.
* Malo, porque suma una librería y el código de renderizado que mantener dentro del servicio.

### Autoalojar QuickChart dentro del clúster de k3s

Desplegar QuickChart, que es de código abierto, como un servicio propio dentro del clúster, conservando la misma integración HTTP del ADR-23.

* Bien, porque elimina la dependencia de un tercero sin cambiar la forma de integración ya construida.
* Bien, porque era la ruta de evolución que el propio ADR-23 dejó abierta.
* Malo, porque suma un servicio más que desplegar, actualizar y monitorear solo para una operación liviana.
* Malo, porque QuickChart es una herramienta general de gráficos, mucho más grande de lo que se necesita para generar códigos QR.

### Generar la imagen en el cliente (frontend) bajo demanda

El navegador genera la imagen en el momento en que alguien la visualiza o descarga desde el panel, con una librería de JavaScript.

* Bien, porque no consume recursos del backend.
* Malo, porque deja la generación fuera del servicio responsable de los códigos, duplicando la lógica de qué URL codifica cada QR.
* Malo, porque no permite entregar la imagen a través de la API a integraciones distintas del panel, como una imprenta o un tostador que genere sus etiquetas automáticamente.

## Enlaces

* [Reemplaza a] [ADR-23: Generar el Código QR de cada Lote](./adr-23-generar-codigos-qr.md)
* [Relacionado con] [ADR-24: Definir arquitectura concreta para el Backend](./adr-24-arquitectura-backend.md)
* [Relacionado con] [ADR-01: Uso de Redis como memoria caché](./adr-01-gestionar-cache.md)
* [Relacionado con] [ADR-16: Versionar las URLs de los Códigos QR](./adr-16-versionar-urls.md)
* [Relacionado con] [ADR-12: Proteger Servicios Externos](./adr-12-proteger-servicios.md)
