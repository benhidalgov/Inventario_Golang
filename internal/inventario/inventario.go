package inventario

import (
	"errors"
	"sync"
)

// Producto representa un artículo del inventario.
// Los tags `json` controlan cómo se serializa/deserializa en JSON.
type Producto struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Stock  int    `json:"stock"`
}

// Inventario almacena los productos de forma segura para uso concurrente.
// El campo `productos` es privado (minúscula) para forzar que todo acceso
// pase por los métodos del struct, que usan el mutex internamente.
// Si fuera público, cualquier parte del código podría leer/escribir el mapa
// sin bloquear el mutex, causando race conditions.
type Inventario struct {
	mu   sync.Mutex
	repo *SQLiteRepository // Ahora el inventario conoce la base de datos
}

// CompraRequest define la estructura del JSON que recibimos en /comprar.
type CompraRequest struct {
	ProductoID int `json:"producto_id"`
	Cantidad   int `json:"cantidad"`
}

// Errores específicos del dominio.
// Usar errores tipados (con errors.New o variables centinela) permite
// que el caller los identifique con errors.Is() y tome decisiones
// (por ejemplo, devolver un código HTTP distinto según el tipo de error).
var (
	ErrProductoNoEncontrado = errors.New("producto no encontrado")
	ErrStockInsuficiente    = errors.New("stock insuficiente")
)

// NewInventario crea un Inventario conectado a la base de datos.
// Recibe el repositorio SQLite ya inicializado.
func NewInventario(r *SQLiteRepository) *Inventario {
	return &Inventario{
		repo: r,
	}
}

// AgregarProducto inserta o reemplaza un producto en el inventario.
// Usa el mutex para que solo una goroutine escriba a la vez.
// Ahora delega la operación a la base de datos SQLite.
func (inv *Inventario) AgregarProducto(p Producto) error {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	// Delegamos al repositorio que hace el INSERT en la DB
	return inv.repo.AgregarProductoDB(p)
}

// ObtenerProductos devuelve todos los productos de la base de datos.
// El mutex protege la lectura para evitar conflictos con escrituras
// concurrentes desde otras goroutines.
func (inv *Inventario) ObtenerProductos() (map[int]Producto, error) {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	// Delegamos al repositorio que hace el SELECT en la DB
	return inv.repo.ObtenerProductosDB()
}

// ActualizarStock modifica el stock de un producto por la cantidad indicada.
// Valores negativos restan stock (compra), positivos lo suman (reposición).
// El mutex protege la operación para evitar race conditions.
func (inv *Inventario) ActualizarStock(id int, cantidad int) error {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	// Llamamos a la base de datos de forma segura
	return inv.repo.ActualizarStockDB(id, cantidad)
}
