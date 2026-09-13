# PoC 06: Procesamiento de Archivos antes de Almacenarlos

Prueba de concepto (PoC) que implementa el procesamiento básico de archivos (validación y optimización de imágenes) en **Go**, como primer paso hacia lo definido en el ADR-19. Para esta etapa, en vez de integrar la cola de RabbitMQ y MinIO, se usan carpetas locales `input` y `output`, de forma que la lógica de procesamiento en sí pueda validarse de forma aislada antes de conectarla al resto del pipeline.

## Arquitectura y Justificación

El ADR-19 decidió orquestar herramientas de código abierto especializadas por tipo de operación, en vez de reimplementarlas o depender de un servicio administrado. Esta PoC valida esa idea en su forma más simple: la lógica se separa en tres responsabilidades — el orquestador (`main.go`) que recorre los archivos, el validador que detecta el tipo real del archivo y su tamaño, y el procesador de imágenes que las optimiza. Así, el resto de la aplicación no depende de los detalles internos de cómo se procesa cada tipo de archivo, y se pueden ir sumando más procesadores (video, escaneo de seguridad) sin tocar el resto del flujo.

Por ahora, el procesamiento completo (optimizar, limpiar metadatos) solo está implementado para **imágenes**; los documentos y los videos se validan y se copian tal cual, dejando explícito en el código que la transcodificación con ffmpeg y el escaneo con ClamAV quedan para una siguiente iteración, ya que requieren binarios externos instalados.

### Estructura del proyecto

```
poc-06-procesamiento-archivos/
├── cmd/
│   └── main.go                    # Recorre input/, procesa cada archivo según su tipo y escribe el resultado en output/
├── internal/
│   └── procesamiento/
│       ├── validar.go              # Detecta el tipo real del archivo y valida su tamaño
│       └── imagen.go               # Redimensiona y comprime imágenes; al recodificarlas elimina sus metadatos EXIF
├── input/                           # Carpeta de entrada: aquí se colocan los archivos de prueba
├── output/                          # Carpeta de salida: aquí quedan los archivos ya procesados
├── go.mod
└── go.sum
```

## Guía de Levantamiento

### 1. Instalar dependencias

```bash
go get github.com/disintegration/imaging
go mod tidy
```

### 2. Ejecución

Coloca uno o más archivos de prueba en la carpeta `input/` (fotos, documentos, etc.) y luego:

```bash
go run ./cmd
```

No se necesita ningún archivo `.env`: todo el procesamiento ocurre sobre las carpetas locales.

## Comportamiento Observado

- **Rechazo por tamaño:** un archivo que supera el límite configurado (20 MB) se rechaza de inmediato, sin procesarlo, con el motivo impreso en consola.
- **Optimización de imágenes:** una foto de prueba de 3000px de ancho y 92 KB se redimensionó a 1600px y quedó en 26 KB tras la compresión — una reducción de más del 70% en el tamaño del archivo.
- **Validación por contenido:** el tipo de archivo se detecta a partir de sus primeros bytes (no de su extensión), por lo que un archivo renombrado con una extensión falsa no engaña al validador.
- **Documentos sin procesar todavía:** un documento de texto se valida correctamente y se copia sin cambios a `output/`, marcado en consola como pendiente de transcodificación/escaneo.

## Stack Tecnológico

- **Lenguaje:** Go (Golang)
- **Procesamiento de imágenes:** `disintegration/imaging`
- **Detección de tipo de archivo:** `net/http` (`http.DetectContentType`)
