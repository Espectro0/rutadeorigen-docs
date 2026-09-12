# Versionar las URLs de los Códigos QR

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-05

**Historia técnica:** Definir el mecanismo que garantiza que un código QR impreso en una bolsa de café siga funcionando indefinidamente, incluso si el dominio o la estructura de URLs del sistema cambian en el futuro.

## Contexto y planteamiento del problema

Un código QR impreso en un empaque físico no puede modificarse ni reemplazarse una vez que la bolsa está en manos de un consumidor o una tienda. Cualquier cambio futuro de dominio o de la estructura de URLs de las páginas públicas de trazabilidad no puede invalidar un QR ya impreso: una redirección rota rompe directamente la promesa de valor del producto ante el consumidor final.

¿Qué mecanismo permite que la URL codificada en un QR ya impreso siga resolviendo al recurso correcto, incluso después de futuros cambios de dominio o de estructura de URLs?

## Impulsores de decisión

* TC-02 Permanencia de QR
* QS-01 Consumidor consulta la información de trazabilidad de un café

## Opciones consideradas

* Identificador opaco (UUID) más una capa de resolución interna
* Subdominio dedicado y estable exclusivamente para QR
* Tabla de redirecciones versionadas

## Resultado de la decisión

Opción elegida: **Tabla de redirecciones versionadas**, porque no exige introducir un identificador opaco distinto del que ya identifica un lote o producto, ni un subdominio separado del resto de la plataforma: cada vez que cambia la estructura de una URL o el dominio, se guarda un registro que mapea la URL anterior hacia la ubicación vigente, y el sistema resuelve automáticamente cualquier URL histórica válida mediante una redirección permanente (301) hacia su destino actual.

### Consecuencias positivas

* No exige cambiar el esquema de identificadores de lote/producto ya definido, ni introducir una capa adicional de indirección para todo QR nuevo.
* Un cambio de dominio o de estructura de URLs deja de romper los QR ya impresos: basta con registrar el mapeo correspondiente.
* La redirección permanente (301) conserva, de paso, el posicionamiento en buscadores de las páginas públicas ya indexadas.

### Consecuencias negativas

* La tabla de redirecciones crece indefinidamente y nunca puede purgarse por completo, porque un QR impreso hace años puede seguir escaneándose en cualquier momento.
* Depende de un proceso humano disciplinado: si al cambiar una estructura de URLs alguien olvida registrar la redirección correspondiente, ese conjunto de QR queda roto sin que el sistema lo detecte por sí solo.
* No protege por sí sola contra cambios que alteren el formato o la longitud del identificador subyacente del lote/producto; ese es un riesgo que debe gestionarse aparte, con disciplina en el diseño del identificador.

## Pros y contras de las opciones

### Identificador opaco (UUID) más una capa de resolución interna

El QR codifica únicamente un identificador estable y sin significado (UUID); una capa de resolución interna decide, en el momento de la consulta, a qué recurso apunta ese identificador.

* Bien, porque el identificador nunca necesita cambiar, sin importar cuántas veces cambie la estructura interna del sistema.
* Bien, porque desacopla por completo la URL pública de cualquier detalle interno de implementación.
* Malo, porque exige introducir y mantener un esquema de identificadores opacos independiente del que ya usa el modelo de datos para lotes/productos.
* Malo, porque cualquier QR ya impreso con la estructura de URL anterior a esta decisión igual necesitaría un mecanismo de redirección para no quedar roto.

### Subdominio dedicado y estable exclusivamente para QR

Un subdominio fijo (por ejemplo, `qr.rutadeorigen.co`) se reserva exclusivamente para resolver códigos QR, separado del dominio principal usado para el resto de la plataforma.

* Bien, porque aísla los QR de cualquier cambio futuro de dominio o de marca en el resto del sitio.
* Bien, porque es simple de comunicar y de razonar: un solo dominio dedicado, nunca reutilizado para otro fin.
* Malo, porque no resuelve por sí solo un cambio futuro de estructura de URLs dentro de ese mismo subdominio.
* Malo, porque introduce una dependencia adicional de la gestión de DNS de ese subdominio específico, que igual debe mantenerse indefinidamente.

### Tabla de redirecciones versionadas

Cada cambio de estructura de URLs o de dominio se registra como un mapeo desde la URL anterior hacia la vigente, resuelto mediante una redirección permanente.

* Bien, porque no exige un esquema de identificadores nuevo ni un subdominio separado del resto de la plataforma.
* Bien, porque conserva el posicionamiento en buscadores de las páginas públicas ya indexadas, gracias a la redirección 301.
* Malo, porque la tabla de redirecciones crece indefinidamente y depende de que el equipo recuerde registrar cada cambio.
* Malo, porque no cubre, por sí sola, un cambio en el formato del identificador subyacente del lote/producto.

## Enlaces

* [Relacionado con] [ADR-01: Uso de Redis como memoria caché](./adr-01-gestionar-cache.md)
* [Relacionado con] [ADR-02: Selección de PostgreSQL como Motor de Bases de Datos](./adr-02-almacenar-datos.md)
* [Requiere] Definición de un proceso operativo que garantice registrar cada redirección al cambiar la estructura de URLs
