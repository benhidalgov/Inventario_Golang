package inventario

import (
	"sync"
	"testing"
)

func TestActualizarStock(t *testing.T) {
	//Preparacion
	//se usa repositorio en meomoria o un mock o un test por asi decirlo para no ensuciar la veradera DB
	repo, _ := NewSQLiteRepository(":memory:")
	inv := NewInventario(repo)

	//Se inserta producto de prueba
	query := "INSERT INTO productos (id, nombre, stock) VALUES (1, 'Silla test', 10)"
	repo.db.Exec(query)

	err := inv.ActualizarStock(1, -3)

	if err != nil {
		t.Errorf("No se esperaba erorr, pero se obtuvo: %v", err)
	}

}

func TestConcurrenciaStock(t *testing.T) {
	repo, _ := NewSQLiteRepository(":memory:")
	inv := NewInventario(repo)

	stockInicial := 100
	repo.db.Exec("INSERT INTO productos (id, nombre, stock) VALUES (1, 'Producto concurrente', ?)", stockInicial)

	var wg sync.WaitGroup
	numCompras := 50

	for i := 0; i < numCompras; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			inv.ActualizarStock(1, -1)
		}()
	}
	wg.Wait()
}
