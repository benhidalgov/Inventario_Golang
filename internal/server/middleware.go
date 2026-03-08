// Paquete server: contiene la lógica del servidor HTTP,
// incluyendo middlewares y configuración de rutas.
package server

import (
	"log"
	"net/http"
	"time"
)

// ========================================================================
// ¿Qué es un Middleware?
// ========================================================================
// Un middleware es una función que "envuelve" a un handler HTTP.
// Se ejecuta ANTES y/o DESPUÉS del handler real, permitiendo agregar
// lógica transversal (logging, seguridad, etc.) sin modificar los handlers.
//
// El patrón es siempre el mismo:
//   func MiMiddleware(next http.HandlerFunc) http.HandlerFunc {
//       return func(w http.ResponseWriter, r *http.Request) {
//           // Lógica ANTES del handler
//           next(w, r)  // Llamar al handler real
//           // Lógica DESPUÉS del handler
//       }
//   }
//
// Se pueden encadenar: Logger(Recovery(miHandler))
// El request pasa por Logger → Recovery → miHandler
// ========================================================================

// Logger es un middleware que registra información de cada request HTTP.
// Recibe el handler "next" (el siguiente en la cadena) y devuelve un
// nuevo handler que lo envuelve con lógica de logging.
//
// Firma: func(http.HandlerFunc) http.HandlerFunc
//   - Entrada: el handler que viene después en la cadena
//   - Salida:  un nuevo handler que agrega logging
func Logger(next http.HandlerFunc) http.HandlerFunc {
	// Retornamos un closure (función anónima) que captura "next"
	// del scope exterior. Este closure ES el nuevo handler.
	return func(w http.ResponseWriter, r *http.Request) {
		// time.Now() captura el instante antes de procesar el request.
		start := time.Now()

		//  IMPORTANTE: Llamamos a next() para que el request llegue
		// al handler real. Sin esta línea, el request se "traga" aquí
		// y el cliente nunca recibe respuesta del handler.
		next(w, r)

		// time.Since(start) calcula cuánto tardó el handler en responder.
		// Se ejecuta DESPUÉS de que next() termina.
		// Formato del log: | MÉTODO | RUTA | IP_CLIENTE | DURACIÓN |
		log.Printf("| %s | %s | %s | %v |",
			r.Method,          // GET, POST, PUT, DELETE, etc.
			r.URL.Path,        // La ruta solicitada, ej: "/inventario"
			r.RemoteAddr,      // IP y puerto del cliente, ej: "127.0.0.1:54321"
			time.Since(start)) // Duración, ej: "1.234ms"
	}
}

// Recovery es un middleware que protege contra panics.
// Si un handler tiene un panic (error fatal), Recovery lo "atrapa"
// y devuelve un error 500 en vez de hacer crash todo el servidor.
//
// Sin este middleware, un panic en cualquier handler tumbaría
// el servidor entero, afectando a TODOS los usuarios.
func Recovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// defer: esta función anónima se ejecutará AL FINAL de este handler,
		// incluso si hay un panic. Es como un "finally" en otros lenguajes.
		defer func() {
			// recover() captura el valor del panic (si hubo uno).
			// Si no hubo panic, retorna nil y no entra al if.
			if err := recover(); err != nil {
				// Logueamos el panic para debugging.
				// %v imprime el valor del panic (puede ser string, error, etc.)
				log.Printf("[PANIC] %v", err)

				// Respondemos al cliente con un error 500 genérico.
				// No revelamos detalles del panic al cliente por seguridad.
				http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
			}
		}() // Los () al final ejecutan la función defer inmediatamente cuando toca

		// Ejecutamos el handler real. Si este hace panic,
		// el defer de arriba lo atrapa automáticamente.
		next(w, r)
	}
}
