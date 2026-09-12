# Seleccionar el Framework de Frontend

* **Estado:** Propuesto
* **Decidentes:** Juan Esteban Jaramillo Ramírez
* **Fecha:** 2026-09-05

**Historia técnica:** Selección del framework que construye la interfaz consumida por productores, tostadores, marcas y consumidores finales, incluyendo las páginas públicas de trazabilidad accedidas mediante código QR.

## Contexto y planteamiento del problema

La interfaz debe funcionar en dispositivos de gama baja, debe estar preparada desde el inicio para múltiples idiomas pensando en mercados de exportación, y debe cumplir criterios de accesibilidad para personas con distintos niveles de habilidad digital o diversidad funcional. Además, las páginas públicas de trazabilidad consultadas por consumidores mediante QR deben responder con baja latencia y buen posicionamiento en buscadores.

¿Qué framework de frontend sostiene una interfaz accesible e internacionalizada desde el inicio, con un rendimiento aceptable en dispositivos de gama baja tanto para el panel de los actores de la cadena como para las páginas públicas de trazabilidad?

## Impulsores de decisión

* QS-05 Usuario nuevo utiliza las funciones principales
* QS-07 Acceso desde diferentes dispositivos
* QS-21 Reproducción de contenido multimedia
* TC-03 Dispositivos de gama baja
* TC-05 Internacionalización
* TC-06 Accesibilidad

## Opciones consideradas

* React + Next.js
* Vue + Nuxt
* Svelte + SvelteKit
* Astro

## Resultado de la decisión

Opción elegida: **React + Next.js**, porque cuenta con el ecosistema más grande y maduro de librerías de internacionalización (`next-intl`, `react-i18next`) y de componentes accesibles, y porque el renderizado híbrido de Next.js (SSR/SSG) permite que las páginas públicas de trazabilidad respondan rápido en dispositivos de gama baja y posicionen bien en buscadores, sin renunciar a una experiencia totalmente interactiva en el panel de productores, tostadores y marcas.

### Consecuencias positivas

* Cuenta con el ecosistema de librerías de internacionalización y accesibilidad más grande y maduro entre las opciones consideradas.
* El renderizado híbrido de Next.js permite servir las páginas públicas de QR pre-renderizadas, mejorando el tiempo de carga en dispositivos de gama baja y el posicionamiento en buscadores.
* Su comunidad amplia facilita encontrar soluciones a problemas específicos y, a futuro, incorporar más personas al equipo.

### Consecuencias negativas

* El bundle y el tiempo de ejecución de React son, por defecto, más pesados que los de Svelte, exigiendo disciplina propia (code-splitting, evitar librerías de componentes pesadas) para realmente servir bien a los dispositivos de gama baja que exige TC-03.
* Requiere más configuración y boilerplate que Vue para un equipo pequeño.
* Para las páginas públicas de solo lectura, envía más JavaScript al cliente del que enviaría por defecto una arquitectura basada en islands como Astro.

## Pros y contras de las opciones

### React + Next.js

Biblioteca de interfaz de usuario de amplia adopción, con Next.js como framework que añade enrutamiento, renderizado híbrido (SSR/SSG) y una convención de proyecto completa.

* Bien, porque tiene el ecosistema de librerías de i18n y accesibilidad más grande y maduro de las cuatro opciones.
* Bien, porque su renderizado híbrido favorece tanto el rendimiento en dispositivos de gama baja como el posicionamiento en buscadores de las páginas públicas.
* Bien, porque su comunidad amplia facilita la resolución de problemas y el crecimiento futuro del equipo.
* Malo, porque su bundle y tiempo de ejecución son, por defecto, más pesados que los de Svelte.
* Malo, porque exige más configuración y boilerplate que Vue.

### Vue + Nuxt

Framework progresivo con una curva de aprendizaje más suave, y Nuxt como su equivalente a Next.js para enrutamiento y renderizado híbrido.

* Bien, porque su curva de aprendizaje es más suave que la de React para un equipo pequeño.
* Bien, porque su bundle base es, en general, algo más liviano que el de React.
* Bien, porque también ofrece renderizado híbrido (SSR/SSG) mediante Nuxt.
* Malo, porque su ecosistema de librerías de i18n y accesibilidad, aunque maduro, es más pequeño que el de React.
* Malo, porque su comunidad y oferta de talento son menores que las de React.

### Svelte + SvelteKit

Framework que compila los componentes a JavaScript mínimo en tiempo de compilación, sin virtual DOM, con SvelteKit como su capa de enrutamiento y renderizado.

* Bien, porque es la opción más liviana en tiempo de ejecución de las cuatro, favoreciendo directamente a los dispositivos de gama baja.
* Bien, porque su sintaxis es simple y reduce el código repetitivo frente a React.
* Malo, porque su ecosistema de librerías de internacionalización y accesibilidad es considerablemente más reducido que el de React o Vue.
* Malo, porque su comunidad más pequeña dificulta encontrar soluciones a problemas poco comunes.

### Astro

Framework orientado a contenido mayormente estático, que envía cero JavaScript por defecto y solo hidrata componentes interactivos puntuales ("islands").

* Bien, porque es la opción con menor JavaScript enviado al cliente para páginas de solo lectura, como la página pública de trazabilidad consultada por QR.
* Bien, porque favorece directamente el rendimiento en dispositivos de gama baja para ese tipo de páginas.
* Malo, porque su modelo de islands es menos adecuado para un panel interno con mucha interactividad, como el de productores, tostadores y marcas.
* Malo, porque su ecosistema y comunidad son más pequeños que los de React.

## Enlaces

* Ninguno — esta es la primera decisión específica de frontend dentro del conjunto de ADR.
