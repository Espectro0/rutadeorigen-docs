# Reintentar y Persistir Temporalmente ante Pérdida de Conexión

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-05

**Historia técnica:** Definir cómo el sistema conserva la información diligenciada y reintenta el envío cuando un usuario pierde la conexión durante un registro, y cómo se define la política general de reintentos ante fallos transitorios.

## Contexto y planteamiento del problema

Los productores registran procesos de trazabilidad frecuentemente desde el campo, con conectividad limitada o intermitente. Si la conexión se pierde a mitad de un registro, el sistema debe conservar lo diligenciado y permitir completar el envío al recuperar la señal, sin generar registros corruptos ni obligar al usuario a repetir todo el formulario desde cero. 

## Impulsores de decisión

* QS-02 Pérdida de conexión durante un registro
* QS-10 Fallo de un servicio externo

## Opciones consideradas

* Persistencia local
* Background Sync API (Service Worker) con sincronización automática en segundo plano
* Reintento inmediato sin persistencia local

## Resultado de la decisión

Opción elegida: **Persistencia local**, porque conserva la información diligenciada y permite completar el envío al recuperar conexión sin depender de una API de navegador con soporte limitado, dado que la plataforma debe funcionar en cualquier dispositivo. El reintento se dispara cuando el usuario permanece en la aplicación o la reabre, lo cual es consistente con el escenario real de un productor que reintenta al recuperar señal en su finca. La misma política de espera entre reintentos se reutiliza para las llamadas del backend a dependencias externas antes de que el circuit breaker decida abrir el circuito.

### Consecuencias positivas

* La información diligenciada no se pierde ante un corte de conexión.
* No depende de una API de navegador con soporte limitado, funcionando en cualquier dispositivo.
* La misma política de reintentos se reutiliza tanto en el cliente como en las llamadas del backend a servicios externos, evitando definir el criterio dos veces.

### Consecuencias negativas

* Si el usuario cierra la aplicación o el navegador y no vuelve a abrirla, el reintento no se dispara automáticamente en segundo plano.
* Requiere código adicional en el cliente para detectar la reconexión, gestionar el backoff y evitar reenvíos duplicados.

## Pros y contras de las opciones

### Persistencia local

El formulario se guarda localmente en el dispositivo del usuario a medida que se diligencia; al detectar que la conexión se recuperó, el cliente reintenta el envío con una espera creciente entre intentos.

* Bien, porque es relativamente simple de implementar y no requiere infraestructura adicional.
* Bien, porque funciona en cualquier navegador moderno, sin depender de APIs con soporte limitado.
* Bien, porque conserva los datos diligenciados y permite completar el registro tras recuperar conexión.
* Malo, porque si el usuario cierra la pestaña o la aplicación, el reintento no se dispara hasta que la vuelva a abrir.
* Malo, porque gestionar manualmente la detección de conexión y el backoff añade código adicional en el cliente.

### Background Sync API (Service Worker)

Un Service Worker registra una tarea de sincronización que el navegador ejecuta automáticamente en segundo plano en cuanto detecta conexión disponible, incluso si la pestaña ya se cerró.

* Bien, porque sincroniza en segundo plano sin necesidad de que el usuario reabra la aplicación.
* Bien, porque delega en el navegador la detección de reconexión, en vez de implementarla manualmente.
* Malo, porque su soporte de navegadores es limitado, lo que choca con el requisito de funcionar en cualquier dispositivo.
* Malo, porque añade la complejidad de registrar y mantener un Service Worker, mayor carga para un equipo de una sola persona.

### Reintento inmediato sin persistencia local

El cliente simplemente reintenta enviar la solicitud inmediatamente tras un fallo, sin guardar el formulario en ningún almacenamiento local.

* Bien, porque es trivial de implementar, sin código adicional de persistencia.
* Malo, porque si la conexión se pierde y el usuario cierra o recarga la aplicación, toda la información diligenciada se pierde.
* Malo, porque no controla el estado del envío, arriesgando registros duplicados si el reintento se dispara más de una vez.

## Enlaces

* [Relacionado con] [ADR-12: Proteger Servicios Externos](./adr-12-proteger-servicios.md)
