# Code Review - Supermarket Comparer API

## Resumen

Se analizaron todos los archivos del proyecto identificando problemas de código sin usar, malas prácticas e inconsistencias.

---

## 1. Código Sin Usar (Unused Code)

### 1.1 Funciones no utilizadas en repository

**Archivo**: `internal/modules/products/product_repository_impl.go`

| Función                        | Líneas  | Descripción                                  |
| ------------------------------ | ------- | -------------------------------------------- |
| `FindByNameAndBrand`           | 129-154 | Definida pero nunca llamada                  |
| `containsBrandAndPresentation` | 156-176 | Helper no usado, definida pero nunca llamada |

### 1.3 Función no llamada

## 2. Malas Prácticas

### 2.1 Nivel de logging demasiado verboso

**Ubicación**: `internal/database/database.go:34`

```go
Logger: logger.Default.LogMode(logger.Info),
```

**Problema**: En producción, esto genera太多的 logs. Debe usarse `logger.Silent` o `logger.Warn`.

### 2.2 Sin graceful shutdown

**Ubicación**: `cmd/server/main.go:49`

```go
if err := http.ListenAndServe(":"+port, mux); err != nil {
```

**Problema**: El servidor no maneja señales SIGTERM ni cierre gracefully de conexiones activas.

### 2.3 Sleep innecesario

**Ubicación**: `cmd/server/main.go:80`

```go
time.Sleep(10 * time.Millisecond)
```

**Problema**: Sleep de 10ms después de cerrar el listener es innecesario.

### 2.4 Manejo de errores genérico

**Ubicación**: `internal/core/api_response.go:46-51`

```go
func getErrorStatus(err error) int {
    if err == nil {
        return 500
    }
    return 400  // Siempre 400 para cualquier error
}
```

**Problema**: Todos los errores retornan 400, incluyendo 404 (not found), 500 (db errors), etc.

### 2.5 Falta validación de Foreign Key

**Ubicación**: `internal/modules/products/product_service.go`
**Problema**: Al crear un producto con `CategoryID`, no se valida que la categoría existe realmente.

### 2.6 Validación inconsistente de ID

**Ubicación**: `internal/modules/categories/category_service.go:49`

```go
if !core.IsValidUUIDString(id) {
    return &errors.CategoryNotFoundError{ID: id}  // Retorna "not found" para ID inválido
}
```

**Problema**:返错误 "not found" cuando el ID es inválido - debería ser "invalid ID error".

---

## 3. Inconsistencias

### 3.1 DELETE inconsistente entre módulos

| Módulo     | Operación | Comportamiento                |
| ---------- | --------- | ----------------------------- |
| Products   | DELETE    | Soft delete (desactiva)       |
| Categories | DELETE    | Hard delete (elimina físicos) |

**Impacto**: Inconsistencia para el usuario - DELETE behaves differently según el recurso.

### 3.2 Error types inconsistentes

| Servicio        | ID Inválido             | Error Retornado                                      |
| --------------- | ----------------------- | ---------------------------------------------------- |
| ProductService  | `InvalidProductIDError` | ✅ Correcto                                          |
| CategoryService | `CategoryNotFoundError` | ❌ Incorrecto (debería ser `InvalidCategoryIDError`) |

### 3.3 Ausencia de Update (PUT/PATCH)

**Archivo**: README.md:72

- No hay endpoints de actualización implementados para ningún recurso

### 3.4 No hay validación de contenido

**Archivos Handlers**: `product_handler.go`, `category_handler.go`

- Usan `io.ReadAll` + `json.Unmarshal` sin validación de schema
- No hay max request body size limit

### 3.5 Nil pointer potencial

**Ubicación**: `internal/modules/products/product_repository_impl.go:39`

```go
result := database.DB.First(&product, "id = ?", id)
```

**Problema**: Si `id` es nil o vacío, la query puede fallar silenciosamente.

---

## 5. Recomendaciones de Prioridad

| Prioridad | Issue              | Reparación                                         |
| --------- | ------------------ | -------------------------------------------------- |
| Alta      | Código sin usar    | Eliminar `FindByNameAndBrand`                      |
| Alta      | Error handling 400 | Diferenciar códigos: 400 vs 404 vs 500             |
| Media     | Logging verboso    | Cambiar a `logger.Warn` o `logger.Silent`          |
| Media     | Soft/Hard delete   | Estandarizar comportamiento de DELETE              |
| Media     | Invalid ID error   | CategoryService debe retornar error de ID inválido |
| Baja      | Graceful shutdown  | Agregar manejo de señales                          |
| Baja      | Hardcoded ports    | Mover a constants o config                         |
| Baja      | Validación FK      | Agregar validación de CategoryID existe            |

---

## 6. Métricas del Proyecto

| Métrica                      | Valor |
| ---------------------------- | ----- |
| Total archivos .go           | 22    |
| Líneas de código (aprox)     | ~900  |
| Archivos con código sin usar | 3     |
| Funciones duplicadas         | 1     |
| Inconsistencias encontradas  | 5     |
| Malas prácticas              | 6     |

---

_Review generado el 22/04/2026_

