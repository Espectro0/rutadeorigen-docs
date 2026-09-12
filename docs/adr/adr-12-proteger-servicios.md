# Proteger Servicios Externos

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-05

**Historia técnica:** Definir el mecanismo que evita que la falla de un servicio externo afecte las funcionalidades propias del sistema que no dependen de él.

## Contexto y planteamiento del problema

Ante la indisponibilidad de un servicio externo, al menos el 95% de las funcionalidades propias del sistema permanezcan operativas. La plataforma ya depende de varios servicios externos por decisiones previas, y cada llamada a uno de ellos que no maneje su indisponibilidad de forma controlada puede degradar o bloquear operaciones que, en realidad, no dependen de ese servicio.

¿Qué mecanismo evita que la falla o lentitud de un servicio externo se propague hacia las funcionalidades propias del sistema que no dependen de él?

## Impulsores de decisión

* QS-10 Fallo de un servicio externo

## Opciones consideradas

* Circuit Breaker + timeouts
* Proxy/service Mesh dedicado
* Reintentos simples con backoff, sin Circuit Breaker

## Resultado de la decisión

Opción elegida: **Circuit Breaker + timeouts**, porque protege cada llamada a un servicio externo directamente en el código, sin añadir un componente de infraestructura adicional como un proxy o service mesh dedicado, cuya complejidad operativa no se justifica para el alcance del proyecto.

### Consecuencias positivas

* Corta las llamadas hacia un servicio externo caído una vez se supera el umbral de fallos configurado, evitando saturar las solicitudes que probablemente fallarán o tardarán en agotar su timeout.
* No introduce infraestructura adicional que operar.
* Permite que las funcionalidades propias que no dependen del servicio externo caído sigan operando.

### Consecuencias negativas

* Cada servicio externo requiere su propia configuración de umbrales de fallos y timeouts dentro del código, en lugar de una política centralizada.
* La lógica de protección queda dispersa en cada punto del código que llama a un servicio externo, si no se encapsula de forma consistente.
* No desacopla completamente la política de resiliencia del código de negocio, como sí lo haría un proxy dedicado.

## Pros y contras de las opciones

### Circuit Breaker + timeouts

Una librería propia del backend envuelve cada llamada a un servicio externo, cortando las llamadas tras superar un umbral de fallos.

* Bien, porque no requiere infraestructura adicional, integrándose directamente en el código donde se llama al servicio externo.
* Bien, porque da control total sobre los umbrales y timeouts de cada servicio externo.
* Malo, porque cada servicio externo requiere su propia configuración dentro del código del backend.
* Malo, porque la lógica de protección puede quedar dispersa si no se encapsula de forma consistente.

### Proxy/service mesh dedicado que aplique circuit breaking de forma centralizada

Un componente de infraestructura intercepta todas las llamadas salientes y aplica circuit breaking de forma centralizada, sin tocar el código del backend.

* Bien, porque centraliza la política de resiliencia, desacoplada del código de negocio.
* Bien, porque es reutilizable para cualquier llamada saliente nueva sin modificar el backend.
* Malo, porque introduce un componente de infraestructura adicional que operar y mantener.
* Malo, porque su complejidad es considerable frente al alcance del proyecto y su equipo de una sola persona.

### Reintentos simples con backoff, sin circuit breaker dedicado

Cada llamada a un servicio externo se reintenta con backoff exponencial ante un fallo, sin un mecanismo que corte las llamadas de forma anticipada.

* Bien, porque su implementación es mínima y reutiliza patrones ya usados para la pérdida de conexión.
* Malo, porque no corta las llamadas hacia un servicio externo caído: sigue intentando y esperando el timeout de cada solicitud.
* Malo, porque no evita que una caída prolongada del servicio externo termine saturando al backend con solicitudes en espera.

## Enlaces

* [Relacionado con] [ADR-05: Reintentar y Persistir Temporalmente ante Pérdida de Conexión](./adr-05-reintentar-conexion.md)
* [Relacionado con] [ADR-09: Autenticar Usuarios](./adr-09-autenticar-usuarios.md)
* [Relacionado con] [ADR-11: Auditar Cambios](./adr-11-auditar-cambios.md)
* [Relacionado con] [ADR-14: Integrar Asistente de IA de Trazabilidad](./adr-14-integrar-asistente-ia.md)
* [Relacionado con] [ADR-18: Notificar Invitaciones a una Organización](./adr-18-notificar-invitaciones.md)
