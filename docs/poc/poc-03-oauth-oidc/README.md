# PoC 03: Autenticación OAuth2/OIDC con WorkOS AuthKit

Prueba de concepto (PoC) que implementa autenticación de usuarios mediante **OAuth2/OIDC** usando **WorkOS AuthKit**. Simula el login de un usuario a Ruta de Origen a través de la pantalla hospedada de AuthKit, con un backend en Go que gestiona la sesión y un frontend básico que refleja el estado de autenticación.

## Arquitectura y Justificación

Delegar la autenticación a **WorkOS AuthKit** evita construir y mantener un sistema propio de login y ofrece un flujo OAuth2/OIDC hospedado listo para producción, con hasta 1M de usuarios activos mensuales gratis. El backend solo se encarga de iniciar el flujo, intercambiar el código de autorización por el usuario autenticado y mantener una sesión propia, sin necesidad de una base de datos para el alcance de este PoC. La lógica se separa en tres responsabilidades: construcción del cliente, manejo de sesión y los handlers HTTP. El frontend es vanilla HTML/CSS/JS, sin frameworks, ya que solo necesita reflejar dos estados: sesión iniciada o no.

### Estructura del proyecto

```
poc-03-oauth-oidc/
├── backend/
│   ├── cmd/
│   │   └── main.go            # Carga variables de entorno, arranca el servidor y sirve el frontend
│   ├── internal/
│   │   └── auth/
│   │       ├── client.go       # Construye el cliente de WorkOS a partir de las variables de entorno
│   │       ├── handlers.go     # Handlers de /login, /callback, /api/me y /api/logout
│   │       └── session.go      # Sesión en memoria + cookie HttpOnly
│   ├── go.mod
│   └── go.sum
└── frontend/
    ├── index.html               # Estructura base de la página
    ├── styles.css               # Estilos de la tarjeta de login/perfil
    └── app.js                   # Consume /api/me y pinta login o perfil
```

### Parámetros de Configuración (`.env.example`)

Todos los parámetros de configuración estan estructurados dentro del `.env`, por lo que puedes ver una plantilla en la carpeta de este PoC. (`./backend/.env.example`)

## Guía de Levantamiento

### 1. Configurar variables de entorno

Crear un archivo `.env` dentro de `backend/` con `WORKOS_API_KEY`, `WORKOS_CLIENT_ID` y `WORKOS_REDIRECT_URI` (por defecto `http://localhost:8080/callback`). Se cargan automáticamente al inicio.

### 2. Instalar dependencias

```bash
go get github.com/workos/workos-go/v10
go get github.com/joho/godotenv
go mod tidy
```

### 3. Ejecución

```bash
cd backend
go run ./cmd
```

El servidor sirve tanto la API (`/login`, `/callback`, `/api/me`, `/api/logout`) como el frontend estático desde `../frontend`. Abrir `http://localhost:8080`.

## Comportamiento Observado

- **Login vía AuthKit:** al entrar sin sesión, la página muestra un botón que redirige a la pantalla hospedada de AuthKit (email/contraseña) usando `GetAuthKitAuthorizationURL` con `provider=authkit`.
- **Intercambio de código:** tras autenticarse, WorkOS redirige a `/callback` con un `code`, que el backend intercambia por el usuario autenticado (`AuthenticateWithCode`).
- **Sesión propia:** el usuario queda guardado en un mapa en memoria, identificado por una cookie `HttpOnly`, sin depender de las cookies de sesión de WorkOS.
- **Frontend dinámico:** `app.js` consulta `/api/me` al cargar; si hay sesión pinta el perfil (avatar, nombre, correo) y un botón de cerrar sesión, si no, pinta el botón de login.
- **Logout:** `/api/logout` destruye la sesión local; no cierra sesión en WorkOS (no era necesario para el alcance del PoC).

## Stack Tecnológico

- **Lenguaje:** Go (Golang)
- **Identity Provider:** WorkOS AuthKit (`workos-go/v10`)
- **Configuración:** `godotenv`
- **Frontend:** HTML, CSS y JavaScript vanilla
