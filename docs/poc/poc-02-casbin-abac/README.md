# PoC 02: Sistema de Autorización Basada en Atributos (ABAC)

Prueba de concepto (PoC) que implementa un control de acceso mediante **ABAC (Attribute-Based Access Control)** utilizando **Casbin** y **Go**. Evalúa permisos comparando atributos del usuario (rol, negocio) contra atributos del recurso (dueño, estado), en vez de depender de una tabla fija de roles y permisos.

> A diferencia de RBAC, que nos obliga a tener permisos fijos por rol, ABAC evalúa atrivutos en tiempo de evaluación, evitando el colapso de los roles cuando el acceso depende del contexto y no solamente de "quién es" el usuario.

## Arquitectura y Justificación

Un modelo de permisos codificado a mano esparcido por el código de negocio se vuelve difícil de mantener y de auditar a medida que crecen las reglas. **Casbin** desacopla esa lógica: la política de acceso vive en un archivo de configuración (`model.conf`) separado del código, evaluado por un motor de reglas independiente.

Se usa el modo de Casbin sin políticas persistidas (CSV, Bases de Datos) se compara directamente los atributos de los modelos que ya existen en el dominio.

### Estructura del proyecto

```
poc-02-casbin-abac/
├── cmd/
│   └── main.go            # Casos de prueba
├── internal/
│   ├── auth/
│   │   ├── enforcer.go    # Carga el modelo y construye el enforcer de Casbin
│   │   └── model.conf     # Modelo ABAC
│   └── domain/
│       └── domain.go      # Estructuras de usuario y lote
├── go.mod
└── go.sum
```
### Reglas de Autorización (`internal/auth/model.conf`)

| Rol | Acción | Condición |
| :--- | :---: | :--- |
| `admin` | `edit`, `read` | Acceso total. |
| `producer` | `edit`, `read` | Solo si tiene pertenece al mismo negocio; no puede editar si el lote ya está publicado. |
| `consumer` | `read` | Solo si el lote está publicado. |

## Guía de Levantamiento

### 1. Instalar dependencias

```bash
go get github.com/casbin/casbin/v2
go mod tidy
```

### 2. Ejecución

```bash
go run ./cmd
```

> No requiere levantar ningún servicio externo. El motor de autorización corre dentro del binario.

## Comportamiento Esperado

- **Acceso total de un administrador:** cualquier acción sobre cualquier lote es permitida sin condiciones.
- **Aislamiento por negocio:** un productor no puede editar ni leer lotes que pertenezcan a otro negocio.
- **Bloqueo por estado:** un lote en estado publicado deja de ser editable, incluso para el dueño del negocio.
- **Lectura pública controlada:** un consumidor solo puede leer lotes ya publicados.
- **Trazabilidad:** cada caso imprime el resultado obtenido junto al esperado, permitiendo detectar de inmediato cualquier regla que no se comporte como se diseñó.

> Los resultados de los test realizados para el PoC fueron los siguientes: 

```shell
[OK] Admin edits any batch                                   -> allowed=true  [expected=true ]
[OK] Producer edits their own draft batch                    -> allowed=true  [expected=true ]
[OK] Producer edits their own published batch                -> allowed=false [expected=false]
[OK] Producer edits another producer's draft batch           -> allowed=false [expected=false]
[OK] Consumer edits any batch                                -> allowed=false [expected=false]

All tests passed.
```

## Stack Tecnológico

- **Lenguaje:** Go (Golang)
- **Motor de Autorización:** Casbin (`casbin/v2`) - modelo ABAC