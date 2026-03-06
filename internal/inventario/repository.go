package inventario

import (
	"database/sql"

	// El guion bajo (_) significa "importar solo por sus efectos secundarios".
	// Este paquete registra el driver "sqlite3" en database/sql mediante su
	// función init(). Sin esto, sql.Open("sqlite3", ...) fallaría.
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteRepository encapsula la conexión a la base de datos SQLite.
// Es nuestra capa de acceso a datos (patrón Repository).
type SQLiteRepository struct {
	db *sql.DB // conexión a la base de datos
}

// NewSQLiteRepository crea una nueva instancia del repositorio.
// Recibe la ruta al archivo .db de SQLite y devuelve el repositorio
// o un error si algo falla.
func NewSQLiteRepository(path string) (*SQLiteRepository, error) {
	// sql.Open no abre la conexión de inmediato, solo la configura.
	// "sqlite3" es el nombre del driver registrado por go-sqlite3.
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	// Usamos backticks (`) para strings multi-línea en Go (raw strings).
	// CREATE TABLE IF NOT EXISTS crea la tabla solo si no existe aún,
	// evitando errores al ejecutar el programa más de una vez.
	query := `CREATE TABLE IF NOT EXISTS productos (
		id INTEGER PRIMARY KEY,
		nombre TEXT NOT NULL,
		stock INTEGER NOT NULL
	)`

	// db.Exec ejecuta una consulta que no devuelve filas (INSERT, UPDATE, CREATE...).
	// Ignoramos el primer valor de retorno (_) porque no necesitamos el resultado.
	_, err = db.Exec(query)
	return &SQLiteRepository{db: db}, err
}

// ActualizarStockDB modifica el stock de un producto en la base de datos.
// - id: el ID del producto a actualizar.
// - cambio: la cantidad a sumar (positivo) o restar (negativo) al stock.
func (r *SQLiteRepository) ActualizarStockDB(id, cambio int) error {
	// Los signos ? son placeholders para prevenir SQL injection.
	// SQLite los reemplaza de forma segura con los valores que pasamos después.
	query := `UPDATE productos SET stock = stock + ? WHERE id = ?`

	// db.Exec ejecuta la query y devuelve un Result con info de las filas afectadas.
	// Pasamos "cambio" y "id" en el orden en que aparecen los ? en la query.
	res, err := r.db.Exec(query, cambio, id)
	if err != nil {
		return err
	}

	// RowsAffected() nos dice cuántas filas fueron modificadas.
	// Si es 0, significa que no existe ningún producto con ese ID.
	filas, _ := res.RowsAffected()
	if filas == 0 {
		return ErrProductoNoEncontrado
	}
	return nil
}

// AgregarProductoDB inserta un producto en la base de datos.
// Si ya existe un producto con el mismo ID, lo reemplaza (REPLACE).
func (r *SQLiteRepository) AgregarProductoDB(p Producto) error {
	// INSERT OR REPLACE: si el ID ya existe, borra la fila vieja e inserta la nueva.
	// Los ? son placeholders que previenen SQL injection.
	query := `INSERT OR REPLACE INTO productos (id, nombre, stock) VALUES (?, ?, ?)`

	_, err := r.db.Exec(query, p.ID, p.Nombre, p.Stock)
	return err
}

// ObtenerProductosDB lee todos los productos de la base de datos
// y los devuelve como un mapa [id]Producto.
func (r *SQLiteRepository) ObtenerProductosDB() (map[int]Producto, error) {
	// SELECT * trae todas las filas de la tabla.
	query := `SELECT id, nombre, stock FROM productos`

	// db.Query (a diferencia de db.Exec) se usa para consultas que SÍ devuelven filas.
	// Devuelve un *sql.Rows que debemos recorrer e ir cerrando al finalizar.
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	// defer rows.Close() garantiza que liberamos los recursos de la consulta
	// cuando la función termine, incluso si hay un error en el medio.
	defer rows.Close()

	// Creamos el mapa donde guardaremos los productos.
	productos := make(map[int]Producto)

	// rows.Next() avanza al siguiente resultado. Devuelve false cuando
	// no quedan más filas o si ocurre un error.
	for rows.Next() {
		var p Producto
		// rows.Scan() lee los valores de la fila actual y los guarda
		// en las variables que le pasamos (por puntero con &).
		// El orden debe coincidir con el SELECT: id, nombre, stock.
		err := rows.Scan(&p.ID, &p.Nombre, &p.Stock)
		if err != nil {
			return nil, err
		}
		productos[p.ID] = p
	}

	return productos, nil
}
