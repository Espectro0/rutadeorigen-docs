# Integrar Asistente de IA de Trazabilidad

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-05

**Historia técnica:** Definir cómo se implementa el asistente de inteligencia artificial que da recomendaciones a productores, tostadores y marcas a partir de la información de trazabilidad ya registrada en la plataforma.

## Contexto y planteamiento del problema

Un asistente de IA que recomiende a partir de datos de trazabilidad ya subidos a la plataforma, donde lo difícil no es utilizarla, es que el asistente genere información sin respaldo en los datos reales del productor, tostador o marca. Construir o entrenar propiamente un modelo es demasiado costoso, tiene un gran tiempo de desarrollo y una curva de aprendizaje muy grande.

¿Cómo se implementa el asistente de IA de forma que sus recomendaciones estén respaldadas por los datos de trazabilidad reales, sin construir ni entrenar un modelo propio?

## Impulsores de decisión

* RF-21 Asistente IA de Trazabilidad

## Opciones consideradas

* Integración con LLM externo vía, utilizando RAG
* Modelo de lenguaje autoalojado corriendo en la propia infraestructura
* Asistente basado en reglas/respuestas predefinidas

## Resultado de la decisión

Opción elegida: **Integración con LLM externo vía, utilizando RAG**, porque permite ofrecer una experiencia conversacional real sin entrenar ni operar un modelo propio, y el RAG es la forma para que el asistente no genere información sin respaldo en los datos ya registrados.

### Consecuencias positivas

* Ofrece una experiencia conversacional de alta calidad sin entrenar ni operar un modelo propio.
* El RAG ancla las respuestas del asistente en los datos de trazabilidad reales.
* Se integra como cualquier otro servicio externo.

### Consecuencias negativas

* Genera un costo recurrente por uso.
* Requiere enviar parte de la información de trazabilidad al proveedor externo del modelo para darle contexto, lo cual debe manejarse con cuidado frente a la información privada de productores, tostadores y marcas.
* Depende de la disponibilidad del proveedor externo del LLM.

## Pros y contras de las opciones

### Integración con LLM externo vía, utilizando RAG

El backend consulta un LLM externo mediante API, recuperando primero los datos de trazabilidad relevantes (RAG) para fundamentar la respuesta del asistente.

* Bien, porque no requiere entrenar ni operar un modelo propio.
* Bien, porque el RAG ancla las respuestas en los datos reales.
* Bien, porque se integra igual que otros servicios externos.
* Malo, porque genera un costo recurrente por uso.
* Malo, porque implica enviar parte de la información de trazabilidad a un tercero.

### Modelo de lenguaje autoalojado corriendo en la propia infraestructura

Un modelo de lenguaje de código abierto se despliega y ejecuta dentro de la infraestructura propia del proyecto.

* Bien, porque no genera costo por token ni envía datos a terceros.
* Bien, porque da control total sobre el modelo y sus datos.
* Malo, porque requiere cómputo (GPU) significativamente mayor al resto de la infraestructura del proyecto.
* Malo, porque implica una carga operativa alta (ajustar y mantener el modelo).

### Asistente basado en reglas/respuestas predefinidas

El asistente responde mediante un conjunto de reglas o respuestas predefinidas, tipo FAQ o árbol de decisión, sin un modelo generativo.

* Bien, porque no genera costo recurrente ni depende de un proveedor externo.
* Bien, porque sus respuestas son totalmente predecibles y controladas.
* Malo, porque no ofrece una experiencia conversacional real ni se adapta a preguntas no anticipadas.
* Malo, porque se aleja del propósito de un "asistente de IA".

## Enlaces

* [Relacionado con] [ADR-12: Proteger Servicios Externos](./adr-12-proteger-servicios.md)
