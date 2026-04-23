package main

import (
	"log"
	"net/http"
	"os"

	"supermarket-comparer-go/internal/database"
	"supermarket-comparer-go/internal/modules/categories"
	"supermarket-comparer-go/internal/modules/products"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	productRepo := products.NewProductRepository()
	productService := products.NewProductService(productRepo)

	categoryRepo := categories.NewCategoryRepository()
	categoryService := categories.NewCategoryService(categoryRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.Handle("/products/", products.ProductHandler(productService))
	mux.Handle("/categories/", categories.CategoryHandler(categoryService))

	port := os.Getenv("PORT")

	log.Printf("Server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
