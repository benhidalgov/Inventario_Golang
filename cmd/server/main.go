package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"tienda/internal/inventario"
)

func main() {
	// 1. Creamos el repositorio SQLite (se conecta a la DB y crea la tabla)
	repo, err := inventario.NewSQLiteRepository("inventario.db")
	if err != nil {
		log.Fatal("Error al crear el repositorio:", err)
	}

	// 2. Creamos el inventario pasándole el repositorio
	inv := inventario.NewInventario(repo)

	// 3. Agregamos un producto de ejemplo (ahora puede fallar, así que chequeamos err)
	err = inv.AgregarProducto(inventario.Producto{ID: 1, Nombre: "Escritorio", Stock: 10})
	if err != nil {
		log.Fatal("Error al agregar producto:", err)
	}

	// GET /inventario - Devuelve todos los productos en formato JSON.
	http.HandleFunc("/inventario", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// ObtenerProductos ahora devuelve (map, error)
		productos, err := inv.ObtenerProductos()
		if err != nil {
			http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(productos)
	})

	// POST /comprar - Procesa una compra restando stock.
	http.HandleFunc("/comprar", procesarCompra(inv))

	log.Println("El servidor corre en el puerto 8080")

	// log.Fatal() hace dos cosas:
	// 1. Imprime el error con timestamp (como log.Println)
	// 2. Llama a os.Exit(1) para terminar el programa
	// Sin esto, si el puerto está ocupado, ListenAndServe retorna un error
	// que se descartaría silenciosamente y el programa terminaría sin aviso.
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// procesarCompra retorna un http.HandlerFunc (closure) que captura `inv`.
// Este patrón permite inyectar dependencias en handlers HTTP sin usar
// variables globales. El handler resultante tiene acceso a `inv` porque
// Go captura la variable del scope exterior en el closure.
func procesarCompra(inv *inventario.Inventario) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validamos el método HTTP. Solo aceptamos POST porque una compra
		// modifica estado (no es idempotente), y GET es para consultas.
		if r.Method != http.MethodPost {
			http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
			return
		}

		var req inventario.CompraRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Json invalido", http.StatusBadRequest)
			return
		}

		// Intentamos actualizar el stock (restamos con valor negativo).
		err = inv.ActualizarStock(req.ProductoID, -req.Cantidad)
		if err != nil {
			// Usamos errors.Is() para distinguir el tipo de error y devolver
			// el código HTTP apropiado:
			// - 404 Not Found: el producto no existe
			// - 409 Conflict: el producto existe pero no hay stock suficiente
			// - 500 Internal Server Error: error inesperado
			// Esto es mejor que devolver siempre 404, porque le da al cliente
			// información útil para saber qué salió mal.
			switch {
			case errors.Is(err, inventario.ErrProductoNoEncontrado):
				http.Error(w, err.Error(), http.StatusNotFound)
			case errors.Is(err, inventario.ErrStockInsuficiente):
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Stock actualizado correctamente para el producto %d", req.ProductoID)
	}
}
