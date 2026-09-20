# Internacionalizar la Interfaz con un Catálogo de Mensajes

* **Estado:** Aprobado
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-20

**Historia técnica:** Definir la librería y el mecanismo de catálogo de mensajes que traduce la interfaz de Ruta de Origen a distintos idiomas, cumpliendo el requisito de internacionalización ya fijado como criterio de selección de frontend en el ADR-17, sin haber definido todavía cuál usar en la práctica.

## Contexto y planteamiento del problema

El ADR-17 exige que la interfaz esté "preparada desde el inicio para múltiples idiomas pensando en mercados de exportación" (TC-05), y usó justamente la disponibilidad de librerías de internacionalización maduras (`next-intl`, `react-i18next`) como uno de los criterios para elegir React + Next.js sobre las demás opciones. Falta decidir en cuál de esas librerías se apoya el catálogo de mensajes: dónde viven las traducciones, cómo se cargan por idioma, y cómo se manejan casos como plurales, género gramatical o formatos de fecha/número que varían entre mercados de exportación.

¿Qué librería y mecanismo de catálogo de mensajes traduce la interfaz a distintos idiomas, cubriendo tanto el panel interno como las páginas públicas de trazabilidad, sin renunciar a las ventajas de rendimiento (SSR/SSG) ya decididas en el ADR-17?

## Impulsores de decisión

* TC-05 Internacionalización
* TC-06 Accesibilidad
* QS-05 Usuario nuevo utiliza las funciones principales

## Opciones consideradas

* next-intl
* react-i18next + next-i18next
* FormatJS / react-intl

## Resultado de la decisión

Opción elegida: **react-i18next + next-i18next**, porque reutiliza el ecosistema i18next, el más maduro y con más plugins del mercado JavaScript (detección automática de idioma, carga diferida de namespaces, herramientas de extracción de claves), coherente con el criterio ya usado en el ADR-17 de priorizar el ecosistema más grande disponible.

Riesgo conocido: a la fecha de este ADR el proyecto todavía no fijó si el frontend se construirá sobre el App Router o el Pages Router de Next.js. `next-i18next` está documentado oficialmente como una solución pensada para el Pages Router, con soporte limitado y no oficial del App Router — si el proyecto adopta el App Router, esta decisión debe revisarse antes de implementarse.

### Consecuencias positivas

* Reutiliza el ecosistema i18next, el más maduro y con más plugins de detección de idioma, carga de namespaces y herramientas de extracción de claves.
* Su comunidad amplia facilita resolver problemas específicos de internacionalización poco comunes (pluralización compleja, formatos regionales).
* No depende de convenciones nuevas para el equipo: es la opción más documentada en tutoriales y ejemplos de Next.js con i18n.

### Consecuencias negativas

* `next-i18next` está documentado oficialmente como pensado para el Pages Router; su soporte del App Router es limitado y no oficial, lo que representa un riesgo si el proyecto adopta el App Router más adelante sin revisar esta decisión.
* Requiere más configuración inicial (archivos de configuración de i18next, estructura de namespaces) que next-intl para el caso simple.
* No es type-safe por defecto: una clave de traducción faltante o mal escrita solo se detecta en tiempo de ejecución, a diferencia de next-intl con TypeScript.

## Pros y contras de las opciones

### next-intl

Librería diseñada específicamente para el App Router moderno de Next.js, con soporte de ICU MessageFormat y tipado estático de las claves de traducción.

* Bien, porque está diseñada específicamente para el App Router, sin las limitaciones de soporte que sí tiene next-i18next.
* Bien, porque es type-safe con TypeScript: una traducción faltante se detecta en tiempo de compilación.
* Bien, porque su configuración inicial es más simple que la de i18next para el caso estándar.
* Malo, porque es la más nueva de las tres, con menos años de adopción en producción y una comunidad más pequeña.
* Malo, porque su ecosistema de plugins es más reducido que el de i18next.

### react-i18next + next-i18next

El ecosistema i18next, el más maduro y extendido del mercado JavaScript, con next-i18next como su integración para Next.js.

* Bien, porque su ecosistema de plugins (detección de idioma, carga diferida, herramientas de extracción) es el más grande y maduro de las tres opciones.
* Bien, porque su comunidad amplia facilita resolver casos poco comunes de pluralización o formatos regionales.
* Malo, porque next-i18next está pensado oficialmente para el Pages Router, con soporte limitado y no oficial del App Router.
* Malo, porque no es type-safe por defecto: una clave faltante solo se detecta en tiempo de ejecución.

### FormatJS / react-intl

Librería de internacionalización con soporte de ICU MessageFormat, usada en el ejemplo oficial de internacionalización de Next.js.

* Bien, porque también soporta ICU MessageFormat, cubriendo plurales, género gramatical y formatos regionales de forma estándar.
* Bien, porque es la librería usada en el ejemplo oficial de Next.js para internacionalización.
* Malo, porque su integración con el App Router es manual, siguiendo ejemplos de la comunidad, sin un paquete dedicado como next-intl.
* Malo, porque su API es más verbosa que la de next-intl para el caso de uso simple.

## Enlaces

* [Relacionado con] [ADR-17: Seleccionar el Framework de Frontend](./adr-17-seleccionar-frontend.md)
