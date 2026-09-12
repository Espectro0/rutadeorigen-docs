# Autenticar Usuarios

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-05

**Historia técnica:** Definir el mecanismo de autenticación de los usuarios registrados que operan sobre la información de trazabilidad.

## Contexto y planteamiento del problema

Construir y mantener un sistema de autenticación propio (hashing de contraseñas, recuperación de cuenta, factores adicionales, detección de contexto sospechoso) implica una carga de seguridad considerable para un equipo de una sola persona, donde cualquier error de implementación compromete directamente la confianza de productores, marcas y consumidores en la plataforma.

¿Quién y cómo verifica la identidad de un usuario antes de permitirle operar sobre la información de trazabilidad?

## Impulsores de decisión

* QS-03 Intento de modificar información pública sin autorización
* QS-20 Inicio de sesión desde contexto sospechoso

## Opciones consideradas

* JWT (tokens autocontenidos, firmados por el propio backend)
* Sesiones de servidor almacenadas en Redis
* Autenticación delegada a un proveedor externo (OAuth2/OIDC con un identity provider)

## Resultado de la decisión

Opción elegida: **Autenticación delegada a un proveedor externo (OAuth2/OIDC)**, porque traslada la complejidad y el riesgo de seguridad de construir un sistema de autenticación propio a un proveedor especializado, reduciendo la superficie de error de un equipo de una sola persona. El backend consume la identidad verificada por el proveedor mediante el protocolo estándar OIDC, sin almacenar credenciales sensibles propias.

### Consecuencias positivas

* Delega a un proveedor especializado la complejidad de seguridad, reduciendo el código propio a mantener.
* El backend no almacena ni gestiona contraseñas, reduciendo el impacto de una eventual brecha de seguridad propia.
* Se apoya en un protocolo estándar (OIDC), con bibliotecas y documentación ampliamente disponibles para su integración en Go.

### Consecuencias negativas

* Introduce un costo recurrente y una dependencia de un tercero para una funcionalidad crítica del sistema.
* La disponibilidad de la autenticación queda ligada a la disponibilidad del proveedor externo, un riesgo que debe gestionarse junto con la estrategia de resiliencia ante fallos de servicios externos.
* Menor control sobre el detalle exacto de las políticas de seguridad, que dependen de lo que el proveedor exponga.

## Pros y contras de las opciones

### JWT (tokens autocontenidos, firmados por el propio backend)

El backend emite tokens firmados que contienen la identidad y los permisos del usuario, verificables sin consultar un almacén central.

* Bien, porque no requiere estado en el servidor, encajando con las múltiples instancias del escalado horizontal.
* Bien, porque es un estándar ampliamente soportado en clientes web y móviles.
* Malo, porque revocar un token antes de su expiración requiere un mecanismo adicional.
* Malo, porque implica construir y mantener toda la lógica de autenticación de forma propia.

### Sesiones de servidor almacenadas en Redis

Cada sesión autenticada se guarda en Redis, y el servidor la consulta para validar cada solicitud.

* Bien, porque permite revocación inmediata: basta con eliminar la sesión almacenada.
* Bien, porque reutiliza el Redis ya adoptado, sin un nuevo componente de infraestructura.
* Malo, porque cada solicitud autenticada depende de una consulta adicional a Redis, que ya cumple múltiples roles.
* Malo, porque igual que la opción de JWT, implica construir y mantener toda la lógica de autenticación propia.

### Autenticación delegada a un proveedor externo (OAuth2/OIDC con un identity provider)

Un proveedor de identidad externo gestiona el registro, la autenticación y la detección de accesos sospechosos, y el backend consume la identidad verificada mediante OIDC.

* Bien, porque delega la complejidad y el riesgo de seguridad a un proveedor especializado.
* Bien, porque reduce el código propio de autenticación que un equipo de una sola persona debe mantener y auditar.
* Malo, porque introduce un costo recurrente y una dependencia crítica de un tercero.
* Malo, porque reduce el control sobre el detalle exacto de las políticas de seguridad aplicadas.

## Enlaces

* [Relacionado con] [ADR-01: Uso de Redis como memoria caché](./adr-01-gestionar-cache.md)
* [Relacionado con] [ADR-10: Autorizar Acciones](./adr-10-autorizar-acciones.md)
* [Relacionado con] [ADR-12: Proteger Servicios Externos](./adr-12-proteger-servicios.md)
* [Relacionado con] [ADR-18: Notificar Invitaciones a una Organización](./adr-18-notificar-invitaciones.md)
