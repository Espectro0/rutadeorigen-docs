# Autorizar Acciones

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-05

**Historia técnica:** Definir el modelo de autorización que decide, una vez identificado el usuario, qué acciones puede ejecutar y sobre qué recursos.

## Contexto y planteamiento del problema

Autenticar a un usuario solo responde "quién es"; el sistema además necesita decidir qué puede hacer. Un usuario autenticado no debe poder consultar información perteneciente a otra cuenta, es decir, la autorización depende no solo del rol del usuario sino de si el recurso le pertenece. 

¿Qué modelo de autorización decide, en cada operación, si un usuario puede ejecutarla sobre un recurso concreto, cubriendo tanto el control por rol como el control por propiedad del recurso?

## Impulsores de decisión

* QS-03 Intento de modificar información pública sin autorización
* QS-16 Solicitud de información privada ajena
* QS-17 Ejecución de funciones fuera del rol asignado

## Opciones consideradas

* RBAC (roles predefinidos)
* ABAC (autorización basada en atributos: rol + propiedad del recurso + contexto)
* Autorización delegada al proveedor de identidad externo (scopes/permissions de OAuth2)

## Resultado de la decisión

Opción elegida: **ABAC (autorización basada en atributos)**, porque es la única opción que cubre directamente un bloqueo del acceso a información que, aunque el usuario esté autenticado y tenga el rol adecuado, no le pertenece. El rol del usuario se conserva como uno de los atributos evaluados, pero la decisión de autorización se calcula combinando ese rol con la propiedad del recurso solicitado y el contexto de la operación, en lugar de basarse únicamente en el rol.

### Consecuencias positivas

* Cubre explícitamente el control de propiedad de recursos exigido, no solo el control por rol.
* Un mismo mecanismo de autorización resuelve ambos casos, sin necesitar dos sistemas distintos.
* Queda disponible para incorporar más atributos de contexto a futuro sin rediseñar el modelo.

### Consecuencias negativas

* Es más compleja de diseñar, implementar y auditar que un RBAC puro, al depender de la evaluación de múltiples atributos por solicitud.
* Requiere que cada verificación de autorización tenga acceso a los datos de propiedad del recurso, no solo al rol del usuario autenticado.
* Puede representar un esfuerzo de implementación mayor al estrictamente necesario si, en la práctica, la mayoría de las reglas terminan dependiendo solo del rol.

## Pros y contras de las opciones

### RBAC (roles predefinidos)

Cada usuario tiene un rol fijo, y un middleware verifica en cada endpoint si el rol tiene permiso para la operación solicitada.

* Bien, porque es un modelo simple y bien entendido, fácil de auditar de forma centralizada.
* Bien, porque encaja directamente con los roles ya identificados en el proyecto.
* Malo, porque no distingue si un recurso pertenece al usuario autenticado.
* Malo, porque necesitaría un mecanismo adicional para resolver la verificación de propiedad, duplicando lógica de autorización.

### ABAC (autorización basada en atributos: rol + propiedad del recurso + contexto)

La decisión de autorización se calcula combinando el rol del usuario, la propiedad del recurso solicitado y el contexto de la operación.

* Bien, porque cubre en un mismo mecanismo tanto el control por rol, como el control por propiedad del recurso.
* Bien, porque queda disponible para incorporar más atributos de contexto a futuro sin rediseñar el modelo.
* Malo, porque es más compleja de diseñar, implementar y auditar que RBAC puro.
* Malo, porque cada verificación necesita acceso a los datos de propiedad del recurso, no solo al rol del usuario.

### Autorización delegada al proveedor de identidad externo (scopes/permissions de OAuth2)

El mismo proveedor de identidad gestiona los permisos del usuario mediante scopes, y el backend los interpreta.

* Bien, porque reutiliza el proveedor ya elegido, sin un componente adicional que operar.
* Malo, porque los proveedores de identidad externos modelan bien la autenticación, pero no reglas de negocio específicas como la propiedad de un recurso, que dependen de datos que solo tiene el backend.
* Malo, porque no resuelve por sí sola, el caso que exige conocer la relación entre el usuario y el recurso concreto almacenado en la base de datos propia.

## Enlaces

* [Relacionado con] [ADR-09: Autenticar Usuarios](./adr-09-autenticar-usuarios.md)
