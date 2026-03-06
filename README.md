# 🏪 Tienda - API de Inventario

> Microservicio de inventario desarrollado en Go con enfoque en alto rendimiento, concurrencia segura y persistencia con SQLite.

## 🛠️ Tech Stack

| Tecnología | Uso |
|---|---|
| **Go 1.21+** | Lenguaje principal |
| **SQLite** | Base de datos embebida con `database/sql` |
| **sync.Mutex** | Concurrencia segura para manejo de stock |
| **Makefile** | Automatización de tareas |

## 🛡️ Key Features

- **Prevención de SQL Injection** — Consultas parametrizadas con placeholders (`?`)
- **Thread-Safety** — Protección contra race conditions con `sync.Mutex`
- **Arquitectura Limpia** — Separación entre `cmd/` (punto de entrada) e `internal/` (lógica y datos)
- **Patrón Repository** — Capa de acceso a datos desacoplada de la lógica de negocio

## 📁 Estructura del Proyecto

```
├── cmd/
│   └── server/
│       └── main.go          # Punto de entrada, rutas HTTP
├── internal/
│   └── inventario/
│       ├── inventario.go    # Lógica de negocio (mutex, validaciones)
│       └── repository.go    # Acceso a datos (SQLite)
├── makefile                 # Comandos de automatización
├── go.mod
└── go.sum
```

## 🚥 Quick Start

```bash
# Clonar el repositorio
git clone https://github.com/tu-usuario/tienda-inventario.git
cd tienda-inventario

# Ejecutar el servidor (compila y corre)
make run

# O compilar el binario para producción
make build
./tienda-server.exe
```

El servidor arranca en `http://localhost:8080`

## 📡 API Endpoints

### `GET /inventario`
Devuelve todos los productos.

```bash
curl http://localhost:8080/inventario
```

**Respuesta:**
```json
{
  "1": {"id": 1, "nombre": "Escritorio", "stock": 10}
}
```

### `POST /comprar`
Realiza una compra restando stock.

```bash
curl -X POST http://localhost:8080/comprar \
  -H "Content-Type: application/json" \
  -d '{"producto_id": 1, "cantidad": 3}'
```

**Respuestas:**
| Código | Significado |
|---|---|
| `200 OK` | Stock actualizado correctamente |
| `404 Not Found` | Producto no encontrado |
| `409 Conflict` | Stock insuficiente |

## ⚙️ Comandos del Makefile

| Comando | Descripción |
|---|---|
| `make run` | Compila y ejecuta el servidor |
| `make build` | Genera el binario de producción |
| `make clean` | Elimina archivos generados |
| `make test` | Ejecuta los tests |
