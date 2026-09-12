# Notificar Invitaciones a una Organización

* **Estado:** Propuesto
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-08

**Historia técnica:** Definir el mecanismo que notifica a una persona que fue invitada a unirse a una organización (finca, tostador o marca) dentro de la plataforma, incluyendo el caso en que esa persona todavía no tiene una cuenta.

## Contexto y planteamiento del problema

La plataforma permite el ingreso a varias organizaciones, lo cual asume que existe una forma de invitar a una persona a unirse a una de ellas. E una funcionalidad nueva que surge directamente de administrar quién pertenece a cada organización. La mayoría de las personas invitadas todavía no tiene cuenta, así que la notificación no puede depender de que la persona ya esté usando el sistema.

¿Cómo se notifica a una persona que fue invitada a unirse a una organización, incluyendo el caso en que todavía no tiene cuenta en la plataforma?

## Impulsores de decisión

* TC-07 Aislamiento entre organizaciones

## Opciones consideradas

* Email transaccional mediante un proveedor especializado (Resend)
* Notificaciones in-app + email combinados
* Servicio externo multicanal (Novu) que orquesta email, push e in-app

## Resultado de la decisión

Opción elegida: **Email transaccional mediante Resend**, porque es el único canal que funciona sin que el invitado tenga cuenta previa ni la plataforma instalada: el correo llega con un enlace para aceptar la invitación y crear su cuenta en el mismo flujo. El envío se encola de forma asíncrona sobre RabbitMQ (ya adoptado), sin bloquear la solicitud de quien invita, y Resend queda protegido como cualquier otro servicio externo mediante el Circuit Breaker ya definido.

### Consecuencias positivas

* Funciona incluso si la persona invitada nunca ha usado la plataforma, sin depender de que abra la aplicación para enterarse.
* No requiere conexión persistente ni una app instalada, a diferencia de una notificación push.
* Se integra directamente con la infraestructura ya decidida: se encola sobre RabbitMQ y se protege con el Circuit Breaker.
* Es la opción de menor costo y complejidad operativa entre las consideradas.

### Consecuencias negativas

* Es un único canal: si el correo cae en spam o el invitado no lo revisa, no hay un mecanismo alternativo (como un centro de notificaciones in-app) que lo respalde.
* Introduce una dependencia externa más y un costo recurrente por envío, sumado a los servicios externos que ya tiene la plataforma.
* No deja una notificación persistente dentro de la plataforma para cuando la persona sí tenga cuenta y quiera revisar invitaciones pasadas o pendientes.

## Pros y contras de las opciones

### Email transaccional mediante un proveedor especializado (Resend)

Un proveedor externo envía el correo de invitación con un enlace para aceptar y, si aplica, crear la cuenta.

* Bien, porque es el único canal que llega a alguien sin cuenta previa en la plataforma.
* Bien, porque no depende de conexión persistente ni de una app instalada.
* Bien, porque es la opción más simple y económica de implementar y operar.
* Malo, porque es un único canal, sin respaldo si el correo no se revisa o cae en spam.
* Malo, porque no deja un historial de invitaciones visible dentro de la plataforma.

### Notificaciones in-app + email combinados

Una tabla propia en PostgreSQL guarda la notificación, visible en un centro de notificaciones al entrar a la plataforma, y en paralelo se envía un correo de refuerzo.

* Bien, porque cubre tanto a quien no tiene cuenta (por email) como a un usuario ya activo (por el centro de notificaciones).
* Bien, porque deja un historial de invitaciones consultable dentro de la plataforma.
* Malo, porque son dos mecanismos que construir y mantener en lugar de uno.
* Malo, porque el canal in-app no aporta nada al caso más común de esta ADR: alguien que todavía no tiene cuenta.

### Servicio externo multicanal (Novu) que orquesta email, push e in-app

Un proveedor externo centraliza las plantillas y decide por qué canal (email, push, in-app) notificar, dejando esa lógica fuera del backend propio.

* Bien, porque facilita agregar más canales (ej. push móvil) a futuro sin rediseñar el mecanismo de notificaciones.
* Bien, porque centraliza las plantillas de todos los canales en un solo lugar.
* Malo, porque suma una dependencia externa más a una plataforma que ya depende de varios servicios (RabbitMQ, WorkOS, el proveedor OIDC, el LLM externo).
* Malo, porque resuelve un problema (orquestar múltiples canales) que todavía no existe: hoy solo se necesita un canal, email.

## Enlaces

* [Relacionado con] [ADR-08: Encolar Tareas Pesadas](./adr-08-encolar-tareas.md)
* [Relacionado con] [ADR-09: Autenticar Usuarios](./adr-09-autenticar-usuarios.md)
* [Relacionado con] [ADR-12: Proteger Servicios Externos](./adr-12-proteger-servicios.md)
