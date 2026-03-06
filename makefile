# Variables
BINARY_NAME=tienda-server.exe
MAIN_PATH=./cmd/server/main.go

# Comando por defecto: compila y corre el servidor
run:
	go run $(MAIN_PATH)

# Compila el binario para producción
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Limpia los archivos generados
clean:
	go clean
	del /f $(BINARY_NAME) 2>nul || exit /b 0

# Ejecuta los tests (cuando los tengamos)
test:
	go test ./... -v