# Exponer los Servicios mediante un Reverse Proxy

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-20

**Historia técnica:** Definir el mecanismo que recibe el tráfico externo del clúster de k3s (ADR-15), termina TLS, y lo enruta hacia el backend y los paneles administrativos de la infraestructura propia (MinIO, Retraced, Grafana), sin interrumpir las páginas públicas de trazabilidad durante una actualización.

## Contexto y planteamiento del problema

El ADR-15 decidió desplegar el backend en k3s con autoescalado horizontal, dejando explícitamente abierto un ecosistema de ingress controllers para resolver cómo entra el tráfico externo, sin definir cuál usar. Además de las réplicas del backend, el clúster también aloja varios paneles administrativos propios — la consola de MinIO (ADR-07), el panel de Retraced (ADR-20) y Grafana (ADR-13) — que deben quedar expuestos de forma controlada y con TLS, sin que cada uno gestione su propio certificado por separado.

¿Qué mecanismo recibe el tráfico externo del clúster, termina TLS de forma centralizada, y lo enruta hacia el servicio interno correcto, sin interrumpir el servicio durante actualizaciones?

## Impulsores de decisión

* QS-01 Consumidor consulta la información de trazabilidad de un café
* QS-13 Solicitudes concurrentes masivas
* TC-10 Prácticas DevOps (despliegue sin interrumpir páginas públicas ni QR ya impresos)

## Opciones consideradas

* Traefik (incluido por defecto en k3s)
* Nginx Ingress Controller
* Caddy (vía caddy-ingress-controller)

## Resultado de la decisión

Opción elegida: **Traefik**, porque viene incluido por defecto en la distribución k3s ya elegida en el ADR-15, sin exigir un paso adicional de instalación ni la desactivación de un componente ya presente, y porque gestiona certificados TLS automáticos contra Let's Encrypt de forma nativa, sin depender de instalar y mantener `cert-manager` como componente aparte. Además, sus middlewares propios permiten aplicar headers de seguridad y límites de tasa de solicitudes directamente sobre las rutas expuestas, sin sumar un componente de infraestructura adicional solo para eso.

### Consecuencias positivas

* No agrega ningún paso de instalación: ya está presente desde que se desplegó k3s en el ADR-15.
* TLS automático nativo contra Let's Encrypt, sin depender de `cert-manager` como pieza separada que mantener.
* Sus middlewares (rate limiting, headers de seguridad, compresión, auth básica) cubren buena parte de la protección perimetral básica sin sumar un componente dedicado.
* Detecta automáticamente los servicios de Kubernetes mediante los recursos IngressRoute/Ingress, sin configuración manual de balanceo.

### Consecuencias negativas

* Al venir "gratis" con k3s, es fácil dejarlo sin configurar con criterio (TLS, middlewares) en vez de tratarlo como una decisión explícita, precisamente lo que este ADR busca evitar.
* Su ecosistema de plugins y su comunidad, aunque activos, son más pequeños que los de Nginx Ingress Controller, el estándar de facto fuera de k3s.
* Migrar a otro ingress controller más adelante implica reescribir las reglas de enrutamiento en un formato distinto, ya que Traefik usa su propio CRD (IngressRoute) además del Ingress estándar.

## Pros y contras de las opciones

### Traefik

Ingress controller incluido por defecto en k3s, con gestión nativa de TLS automático y middlewares propios para rate limiting y headers de seguridad.

* Bien, porque no exige ningún paso de instalación adicional sobre lo ya desplegado en el ADR-15.
* Bien, porque resuelve TLS automático sin depender de `cert-manager` como componente separado.
* Bien, porque sus middlewares cubren protección perimetral básica sin sumar infraestructura nueva.
* Malo, porque su ecosistema y comunidad son más pequeños que los de Nginx Ingress Controller.
* Malo, porque su CRD propio (IngressRoute) no es completamente portable si se migra a otro ingress controller a futuro.

### Nginx Ingress Controller

El ingress controller más adoptado del ecosistema Kubernetes en general, con la documentación y comunidad más extensas de las tres opciones.

* Bien, porque es el estándar de facto, con la mayor cantidad de documentación y soluciones a problemas específicos disponibles.
* Bien, porque usa el recurso Ingress estándar de Kubernetes, más portable entre distintos proveedores/clústeres.
* Malo, porque k3s no lo trae instalado: exige desactivar Traefik explícitamente e instalarlo aparte, sumando un paso operativo.
* Malo, porque TLS automático requiere instalar y mantener `cert-manager` como componente adicional, a diferencia de Traefik.

### Caddy

Servidor web e ingress controller (vía `caddy-ingress-controller`) con la configuración más simple y legible de las tres, y HTTPS automático como su característica de diseño insignia.

* Bien, porque su configuración declarativa (Caddyfile) es la más simple y legible de las tres opciones.
* Bien, porque el HTTPS automático es su característica central desde el diseño, no un añadido.
* Malo, porque su integración con Kubernetes es la menos madura de las tres, con menor adopción y menos ejemplos disponibles.
* Malo, porque, igual que Nginx Ingress Controller, exige desactivar Traefik e instalar este componente aparte.

## Enlaces

* [Relacionado con] [ADR-15: Orquestar Contenedores para Escalado Horizontal](./adr-15-orquestar-contenedores.md)
* [Relacionado con] [ADR-07: Almacenar Evidencias como Objetos](./adr-07-almacenar-evidencias.md)
* [Relacionado con] [ADR-13: Monitorear el Sistema](./adr-13-monitorear-sistema.md)
* [Relacionado con] [ADR-20: Migrar la Auditoría de Cambios a Retraced Autoalojado](./adr-20-migrar-auditoria-retraced.md)
