# Orquestar Contenedores para Escalado Horizontal

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-05

**Historia técnica:** Definir el mecanismo de orquestación que despliega el backend en múltiples instancias, permitiendo escalarlo horizontalmente sin depender de una sola instancia y sin interrumpir las páginas públicas de trazabilidad ni los códigos QR ya impresos durante una actualización.

## Contexto y planteamiento del problema

Varias decisiones ya tomadas asumen que el backend corre en múltiples instancias coordinadas entre sí, sin garantía de que dos solicitudes consecutivas sean atendidas por la misma. Además, el sistema debe poder actualizarse con frecuencia sin dejar de atender las páginas públicas de los lotes ni los QR ya impresos en bolsas físicas.

¿Qué mecanismo despliega y coordina múltiples instancias del backend, permitiendo escalarlas horizontalmente y actualizarlas sin tiempo de inactividad?

## Impulsores de decisión

* QS-01 Consumidor consulta la información de trazabilidad de un café
* QS-13 Solicitudes concurrentes masivas
* TC-10 Prácticas DevOps (despliegue sin interrumpir páginas públicas ni QR ya impresos)

## Opciones consideradas

* Kubernetes autoalojado (k3s)
* Docker Swarm
* Servicio de contenedores administrado (Cloud Run / AWS Fargate)
* HashiCorp Nomad

## Resultado de la decisión

Opción elegida: **Kubernetes (k3s)**, porque ofrece autoescalado horizontal nativo (HPA) que reacciona directamente a métricas, sin necesitar un mecanismo propio adicional para decidir cuándo agregar o quitar instancias, y porque al ser el estándar de facto de la industria cuenta con la documentación, las herramientas y la comunidad más extensas de las opciones consideradas. La distribución ligera (k3s) reduce parte de la carga operativa que implicaría una instalación completa de Kubernetes, sin renunciar a su ecosistema.

### Consecuencias positivas

* Autoescalado horizontal nativo (HPA) que reacciona a métricas de CPU/memoria sin configuración adicional, y que puede extenderse a las métricas ya definidas en Prometheus/Grafana mediante un adaptador estándar.
* Estándar de facto de la industria, con la documentación, las herramientas y la comunidad más extensas de las opciones consideradas.
* Soporta actualizaciones progresivas (rolling updates) sin detener el servicio, cumpliendo la exigencia de no interrumpir las páginas públicas ni los QR ya impresos.
* Deja abierto un ecosistema amplio (ingress controllers, autoescalado avanzado, service mesh) para crecer sin migrar de herramienta de orquestación.

### Consecuencias negativas

* Implica una curva de aprendizaje y una carga operativa considerable para un equipo de una sola persona, incluso en su distribución ligera.
* Su superficie de configuración (manifiestos, controladores, red interna) es mucho mayor que la que exige el resto de la infraestructura ya adoptada.
* Aunque k3s reduce el peso de administrar el plano de control frente a una instalación completa de Kubernetes, ese plano de control sigue siendo responsabilidad propia, a diferencia de un servicio administrado.
* Extender el autoescalado más allá de CPU/memoria hacia las métricas de negocio ya definidas en Prometheus/Grafana requiere instalar y mantener un adaptador adicional (`prometheus-adapter`).

## Pros y contras de las opciones

### Kubernetes (k3s)

Distribución ligera de Kubernetes, autoalojada, con autoescalado horizontal nativo (HPA) y el ecosistema de orquestación más extenso del mercado.

* Bien, porque su autoescalado horizontal es nativo, reaccionando directamente a métricas sin un mecanismo propio adicional.
* Bien, porque es el estándar de facto de la industria, con la documentación y las herramientas más extensas disponibles.
* Malo, porque implica una curva de aprendizaje y una carga operativa considerable.
* Malo, porque su superficie de configuración es mucho mayor que la que exige el resto de la infraestructura ya adoptada.

### Docker Swarm

Herramienta de orquestación nativa de Docker, que coordina múltiples instancias de los mismos contenedores ya definidos en Docker Compose.

* Bien, porque reutiliza directamente la sintaxis de Docker ya conocida, sin una herramienta de orquestación distinta que aprender.
* Bien, porque no tiene costo de licencia ni ata el despliegue a un proveedor cloud específico.
* Bien, porque soporta actualizaciones progresivas sin detener el servicio.
* Malo, porque no ofrece autoescalado automático nativo, a diferencia de Kubernetes o un servicio administrado.
* Malo, porque su comunidad y ecosistema de herramientas son considerablemente más pequeños que los de Kubernetes.

### Servicio de contenedores administrado (Cloud Run / AWS Fargate)

Un proveedor cloud administra el escalado, la infraestructura subyacente y el enrutamiento de los contenedores, cobrando por uso.

* Bien, porque el autoescalado es completamente automático y administrado por el proveedor, sin infraestructura propia que operar.
* Bien, porque reduce a cero la carga operativa de mantener un clúster o nodos propios.
* Malo, porque genera un costo recurrente por uso, en lugar de una alternativa de código abierto sin costo de licencia.
* Malo, porque introduce una dependencia directa de un proveedor cloud específico, similar a la que ya se evitó deliberadamente al elegir MinIO en lugar de un almacenamiento de objetos administrado.

### HashiCorp Nomad

Orquestador de cargas de trabajo de propósito más general que Kubernetes, capaz de coordinar contenedores y otros tipos de procesos.

* Bien, porque su modelo es más simple de operar que Kubernetes.
* Bien, porque no está limitado únicamente a contenedores, dejando abierta la posibilidad de orquestar otro tipo de cargas a futuro.
* Malo, porque su adopción y comunidad son menores que las de Kubernetes o Docker Swarm, con menos documentación disponible para resolver problemas específicos.
* Malo, porque introduce una herramienta adicional distinta de Docker, sin la ventaja de continuidad que sí ofrece Docker Swarm.

## Enlaces

* [Relacionado con] [ADR-01: Uso de Redis como memoria caché](./adr-01-gestionar-cache.md)
* [Relacionado con] [ADR-13: Monitorear el Sistema](./adr-13-monitorear-sistema.md)
* [Requiere] Instalación y configuración de un adaptador de métricas (`prometheus-adapter`) para que el autoescalado horizontal reaccione a las métricas ya definidas en Prometheus/Grafana, más allá de CPU/memoria
