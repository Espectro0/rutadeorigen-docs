# Selección de Go como Lenguaje y Runtime del Backend

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-04

**Historia técnica:** Selección del lenguaje de programación y runtime sobre el cual se construye el backend del sistema, base de la cual dependen las demás decisiones de arquitectura.

## Contexto y planteamiento del problema

*Ruta de Origen* necesita atender consultas públicas de trazabilidad de forma masiva y concurrente, registrar procesos de trazabilidad en campo casi de inmediato, generar sellos digitales de verificación de forma asíncrona, e integrar en el futuro sensores IoT que envían lecturas de forma concurrente.

¿Qué lenguaje/runtime sostiene mejor cientos o miles de operaciones concurrentes con baja latencia, se despliega de forma simple y permite fácilmente sacar avances por módulos, y cuenta con un ecosistema maduro para las decisiones que dependen de él?

## Impulsores de decisión

* QS-01 Consumidor consulta la información de trazabilidad de un café
* QS-09 Registro inmediato de procesos
* QS-12 Generación de sello digital
* QS-13 Solicitudes concurrentes masivas
* RF-13 Integración con Sensores IoT
* RF-14 Registro de Datos IoT

## Opciones consideradas

* Go
* Node.js / TypeScript
* Python (FastAPI / Django)
* Java (Spring Boot)

## Resultado de la decisión

Opción elegida: **Go**, porque ofrece concurrencia nativa mediante goroutines lo que permite atender los picos y las conexiones simultáneas sin la sobrecarga de un modelo de hilos pesados ni el cuello de botella de un solo hilo de ejecución. Al ser un lenguaje compilado a un binario estático, se despliega de forma simple y económica, y cuenta con un ecosistema maduro que sostiene directamente las decisiones tomadas.

### Consecuencias positivas

* Concurrencia nativa sin necesidad de librerías externas de manejo asíncrono.
* Bajo consumo de memoria y CPU por instancia, reduciendo el costo de infraestructura bajo picos de tráfico.
* Compilación a un único binario estático, sin runtime externo, simplificando el despliegue y el pipeline de CI/CD.
* Tipado estático que detecta errores en tiempo de compilación, reduciendo el riesgo de fallas en producción.
* Ecosistema maduro de librerías ya adoptado en decisiones previas.

### Consecuencias negativas

* Menor velocidad de desarrollo inicial comparado con lenguajes dinámicos por el tipado estático y la verbosidad del manejo de errores.
* Ecosistema de librerías para integraciones de IA menos maduro que en Python.
* Curva de aprendizaje adicional.

## Pros y contras de las opciones

### Go

Lenguaje compilado, de tipado estático, diseñado por Google con soporte nativo de concurrencia mediante goroutines y canales.

* Bien, porque las goroutines son livianas frente a un hilo de sistema operativo, permitiendo miles de operaciones concurrentes sin degradar el rendimiento.
* Bien, porque compila a un único binario estático, sin dependencias de runtime externo, simplificando el despliegue en un equipo de una sola persona.
* Bien, porque su bajo consumo de memoria/CPU reduce el costo de infraestructura, alineado con el presupuesto limitado del proyecto.
* Bien, porque cuenta con un ecosistema maduro de librerías.
* Bien, porque su tipado estático detecta errores en tiempo de compilación, reduciendo el riesgo de fallas en producción.
* Malo, porque su desarrollo inicial es más lento que en lenguajes dinámicos, por la verbosidad de su manejo explícito de errores.
* Malo, porque su ecosistema de librerías para integraciones de IA es menos maduro que el de Python.

### Node.js / TypeScript

Runtime de JavaScript basado en un bucle de eventos de un solo hilo, con TypeScript como capa de tipado estático opcional.

* Bien, porque su modelo asíncrono no bloqueante también sostiene muchas conexiones simultáneas de I/O.
* Bien, porque comparte el mismo lenguaje (JavaScript/TypeScript) con un eventual frontend, reduciendo el cambio de contexto para un equipo pequeño.
* Bien, porque cuenta con un ecosistema amplio de librerías para integraciones rápidas, incluida IA.
* Malo, porque su modelo de un solo hilo obliga a delegar el trabajo verdaderamente pesado a procesos adicionales, añadiendo complejidad que en Go se resuelve con una goroutine simple.
* Malo, porque el tipado de TypeScript se pierde en tiempo de ejecución, permitiendo errores de tipo que Go evita en producción.
* Malo, porque su consumo de memoria por instancia es mayor que el de un binario compilado de Go, elevando el costo de infraestructura bajo picos de tráfico.

### Python (FastAPI / Django)

Lenguaje interpretado, de tipado dinámico, con un ecosistema muy maduro para integraciones de datos e inteligencia artificial.

* Bien, porque su ecosistema de IA/ML es el más maduro de las cuatro opciones.
* Bien, porque su curva de aprendizaje y velocidad de desarrollo inicial son altas.
* Malo, porque su modelo de concurrencia real está limitado por el GIL (Global Interpreter Lock), obligando a usar múltiples procesos para aprovechar varios núcleos, lo que consume más memoria por unidad de concurrencia que las goroutines.
* Malo, porque al ser interpretado, su rendimiento base por petición es menor que el de un binario compilado.
* Malo, porque su tipado dinámico traslada más errores de tipos a tiempo de ejecución, un riesgo mayor para un equipo de una sola persona sin red de otros revisores.

### Java (Spring Boot)

Lenguaje compilado a bytecode, ejecutado sobre la máquina virtual de Java (JVM), con tipado estático y un ecosistema empresarial extenso.

* Bien, porque ofrece tipado estático y un ecosistema maduro de librerías empresariales, similar en robustez al de Go.
* Bien, porque su modelo de hilos, aunque más pesado que las goroutines, es suficientemente maduro para sostener concurrencia moderada.
* Malo, porque la JVM tiene un consumo de memoria base considerablemente mayor que un binario de Go, elevando el costo de infraestructura del proyecto.
* Malo, porque un hilo de Java es órdenes de magnitud más pesado que una goroutine, por lo que sostener miles de conexiones concurrentes requiere más memoria y ajuste fino del pool de hilos.
* Malo, porque su tiempo de arranque y su curva de configuración son mayores, lo que no favorece la agilidad de despliegue de un equipo de una sola persona.

## Enlaces

* [Relacionado con] [ADR-01: Uso de Redis como memoria caché](./adr-01-gestionar-cache.md)
* [Relacionado con] [ADR-02: Selección de PostgreSQL como Motor de Bases de Datos](./adr-02-almacenar-datos.md)
