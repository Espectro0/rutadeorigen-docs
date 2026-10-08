# Automatizar la Integración y el Despliegue Continuos

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-10-08

**Historia técnica:** Definir cómo se organiza el código de los servicios definidos en el ADR-24, cómo se integran los cambios y cómo llegan a producción en el clúster de k3s (ADR-15), de forma automática y sin interrumpir las páginas públicas de trazabilidad ni los códigos QR ya impresos.

## Contexto y planteamiento del problema

El ADR-24 dividió el backend en seis microservicios que deben poder desplegarse de forma independiente, y el ADR-15 definió k3s como plataforma de orquestación, con actualizaciones progresivas sin tiempo de inactividad. El TC-10 exige además que la plataforma pueda actualizarse con frecuencia, con pruebas, despliegues y monitoreo automatizados. Sin una estrategia definida, cada servicio se construiría y desplegaría de forma manual, con el riesgo de publicar en producción un cambio sin validar, de desplegar versiones incompatibles entre servicios o de no poder volver con rapidez a una versión anterior.

¿Cómo se organiza el código, se valida cada cambio y se despliega automáticamente cada servicio en k3s, permitiendo volver a una versión anterior ante una falla?

## Impulsores de decisión

* TC-10 Prácticas DevOps
* TC-02 Permanencia de QR
* TC-09 Prácticas de código limpio
* QS-01 Consumidor consulta la información de trazabilidad de un café
* BC-06 Costo de operación sostenible

## Opciones consideradas

* GitHub Actions para CI y despliegue GitOps con Argo CD
* GitHub Actions para CI y despliegue directo desde un runner con Helm
* Forgejo autoalojado con Forgejo Actions y despliegue GitOps con Flux
* GitLab CI con sus propios runners

## Resultado de la decisión

Opción elegida: **GitHub Actions para CI y despliegue GitOps con Argo CD**, porque separa la integración, que corre en la infraestructura de GitHub sin consumir recursos propios, del despliegue, que el propio clúster de k3s jala desde Git. Así el clúster nunca recibe conexiones entrantes desde el CI, el estado de producción queda descrito en el repositorio en todo momento, y volver a una versión anterior se reduce a revertir un commit o usar el rollback de Argo CD.

La estrategia completa queda así:

* **Repositorio:** un monorepo con los servicios (`services/`), las librerías compartidas (`libs/`), el frontend (`web/`) y los manifiestos de despliegue (`deploy/`).
* **Ramificación:** *trunk-based development*. `main` es la única rama permanente y siempre debe estar desplegable; los cambios entran por ramas cortas mediante pull request.
* **Pipelines por servicio:** cada servicio tiene su propio workflow de GitHub Actions, filtrado por la ruta de su código, de modo que un cambio en un servicio solo construye y despliega ese servicio. Un cambio en `libs/` reconstruye todos los servicios que dependen de esa librería.
* **Validaciones en cada pull request:** lint (`golangci-lint` y ESLint), pruebas unitarias y de integración, verificación de las reglas de Clean Architecture del ADR-24 (`go-arch-lint`) y escaneo de vulnerabilidades de dependencias e imágenes (Trivy). Un pull request que no pasa todas las validaciones no puede mezclarse.
* **Imágenes:** cada servicio se construye como imagen de contenedor y se publica en GitHub Container Registry (GHCR), etiquetada con el SHA del commit. Nunca se despliega una etiqueta mutable como `latest`.
* **Manifiestos:** un chart de Helm común, con un archivo de valores por servicio dentro de `deploy/`.
* **Entorno:** únicamente producción.
* **Disparador:** cada merge a `main` despliega en producción. Tras publicar la imagen, el workflow actualiza su etiqueta en los valores de Helm del servicio mediante un commit automático, y Argo CD sincroniza el cambio en el clúster.
* **Migraciones de base de datos:** se ejecutan como un Job de Kubernetes en la fase *PreSync* de Argo CD, antes de actualizar el servicio; si la migración falla, el despliegue no continúa. Las migraciones deben ser compatibles hacia atrás (*expand/contract*), porque durante la actualización progresiva conviven pods de la versión anterior y de la nueva.

### Consecuencias positivas

* Cada servicio se valida y despliega de forma independiente y automática, cumpliendo el TC-10.
* El clúster jala los cambios desde Git: no necesita exponer su API ni recibir conexiones entrantes desde el CI, algo especialmente relevante en infraestructura propia.
* El estado de producción queda descrito por completo en el repositorio: cualquier cambio queda versionado y auditable, y volver a una versión anterior es revertir un commit o usar el rollback de Argo CD.
* Las imágenes etiquetadas con el SHA del commit permiten saber con exactitud qué código corre en producción y redesplegar cualquier versión anterior.
* La integración corre en máquinas de GitHub, sin consumir los recursos del servidor propio.
* Las migraciones en *PreSync* evitan que una versión nueva arranque contra un esquema de base de datos que todavía no tiene.

### Consecuencias negativas

* Sin un entorno previo, cada merge a `main` llega directamente a producción: un error que no detecten las pruebas afecta a los usuarios reales, y la única protección es la calidad de las validaciones, las sondas de salud de k3s y el rollback.
* Las validaciones no incluyen pruebas de compatibilidad de los contratos de eventos entre servicios: un cambio incompatible en un evento solo se detectaría en producción, cuando un consumidor falle al procesarlo.
* Argo CD es un componente más que instalar, actualizar y monitorear dentro del clúster.
* El repositorio, el CI y el registro de imágenes quedan concentrados en GitHub: una caída o un cambio en sus condiciones afecta todo el flujo de entrega, aunque no a los servicios que ya están en producción.
* Las migraciones compatibles hacia atrás obligan a dividir ciertos cambios de esquema en varios despliegues sucesivos.
* Los commits automáticos del CI para actualizar las etiquetas de las imágenes agregan ruido al historial de `main`.

## Pros y contras de las opciones

### GitHub Actions para CI y despliegue GitOps con Argo CD

GitHub Actions valida y construye las imágenes; Argo CD, instalado en k3s, sincroniza el clúster con los manifiestos descritos en Git.

* Bien, porque el clúster jala los cambios y no necesita recibir conexiones entrantes desde el CI.
* Bien, porque el estado de producción queda versionado en Git, con rollback inmediato y un panel que muestra qué está desplegado.
* Bien, porque GitHub Actions y GHCR se integran de forma nativa, sin credenciales adicionales.
* Malo, porque suma Argo CD como componente que operar dentro del clúster.
* Malo, porque concentra la entrega en GitHub como proveedor.

### GitHub Actions para CI y despliegue directo desde un runner con Helm

Un runner autoalojado con acceso al clúster ejecuta `helm upgrade` al final de cada workflow.

* Bien, porque no suma ningún componente adicional al clúster.
* Bien, porque el flujo completo queda en un solo lugar, dentro del workflow.
* Malo, porque el runner necesita credenciales con permisos de escritura sobre el clúster.
* Malo, porque el estado real del clúster puede desviarse de lo descrito en Git sin que nada lo detecte ni lo corrija.

### Forgejo autoalojado con Forgejo Actions y despliegue GitOps con Flux

Repositorio, CI y registro autoalojados, con Flux sincronizando el clúster desde Git.

* Bien, porque elimina la dependencia de GitHub y sigue el patrón de infraestructura propia del proyecto.
* Bien, porque Flux es más liviano que Argo CD.
* Malo, porque suma Forgejo, sus runners y su registro como servicios propios que operar, respaldar y asegurar.
* Malo, porque los runners consumirían los recursos del mismo servidor que atiende producción.

### GitLab CI con sus propios runners

GitLab como plataforma integrada de repositorio, CI y registro, en su versión gestionada o autoalojada.

* Bien, porque integra repositorio, CI, registro y entornos en una sola plataforma.
* Malo, porque la versión autoalojada es pesada en recursos para un servidor propio.
* Malo, porque la versión gestionada tiene límites de minutos de CI más restrictivos en su plan gratuito.

## Enlaces

* [Relacionado con] [ADR-24: Definir arquitectura concreta para el Backend](./adr-24-arquitectura-backend.md)
* [Relacionado con] [ADR-15: Orquestar Contenedores para Escalado Horizontal](./adr-15-orquestar-contenedores.md)
* [Relacionado con] [ADR-21: Exponer los Servicios mediante un Reverse Proxy](./adr-21-exponer-servicios.md)
* [Relacionado con] [ADR-02: Selección de PostgreSQL como Motor de Bases de Datos](./adr-02-almacenar-datos.md)
* [Relacionado con] [ADR-13: Monitorear el Sistema](./adr-13-monitorear-sistema.md)
