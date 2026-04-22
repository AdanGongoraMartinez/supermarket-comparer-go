# Supermarket Comparer API

REST API para comparar precios de productos de supermercado.

## Tecnologias

- **Go 1.21**
- **GORM** (ORM)
- **PostgreSQL** (Base de datos)
- **joho/godotenv** (Variables de entorno)

## Instalacion

```bash
# 1. Copiar archivo de configuracion
cp .env.example .env

# 2. Iniciar PostgreSQL con Docker
docker-compose up -d

# 3. Ejecutar el servidor
go run cmd/server/main.go
```

El servidor iniciara en `http://localhost:3000`.

## Comandos

```bash
# Ejecutar servidor
go run cmd/server/main.go

# Ejecutar tests
go test ./...

# Compilar binario
go build -o server ./cmd/server

# Verificar con go vet
go vet ./...
```

## Endpoints

| Metodo | Ruta               | Descripcion                                                  |
| ------ | ------------------ | ------------------------------------------------------------ |
| GET    | `/health`          | Health check                                                 |
| POST   | `/products`        | Crear producto                                               |
| GET    | `/products`        | Buscar productos (query: `name`, `categoryId`, `activeOnly`) |
| GET    | `/products/{id}`   | Obtener producto por ID                                      |
| DELETE | `/products/{id}`   | Desactivar producto                                          |
| POST   | `/categories`      | Crear categoria                                              |
| GET    | `/categories`      | Buscar categorias (query: `name`)                            |
| GET    | `/categories/{id}` | Obtener categoria por ID                                     |
| DELETE | `/categories/{id}` | Eliminar categoria                                           |

## Implementado

- Gestion de productos (CRUD parcial)
- Gestion de categorias (CRUD completo)
- Conexion a PostgreSQL con auto-migracion
- Manejo de errores estructurado
- Tests unitarios con repositories falsos (no funcionales)

## No implementado

- Supermercados/Tiendas
- Precios de productos
- Comparacion de precios entre supermercados
- Autenticacion y autorizacion
- Paginacion
- Endpoints de actualizacion (PUT/PATCH)
- Documentacion de la API (Swagger/OpenAPI)
- Tests de integracion
- Despliegue (Dockerfile, CI/CD)
