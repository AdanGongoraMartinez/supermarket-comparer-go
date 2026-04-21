package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

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

	// if err := database.CloseDB(); err != nil {
	// 	log.Fatalf("Failed to close database: %v", err)
	// }

	productRepo := products.NewProductRepository()
	productService := products.NewProductService(productRepo)

	categoryRepo := categories.NewCategoryRepository()
	categoryService := categories.NewCategoryService(categoryRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.Handle("/products", products.ProductHandler(productService))
	mux.Handle("/categories", categories.CategoryHandler(categoryService))

	port := getAvailablePort(os.Getenv("PORT"))

	log.Printf("Server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getAvailablePort(preferred string) string {
	ports := []string{"3000", "3001", "3002", "3003", "3004", "8080", "8081"}

	if preferred != "" {
		if !isPortInUse(preferred) {
			return preferred
		}
		ports = append([]string{preferred}, ports...)
	}

	for _, p := range ports {
		if !isPortInUse(p) {
			return p
		}
	}

	log.Fatalf("No available port found")
	return ""
}

func isPortInUse(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return true
	}
	ln.Close()
	time.Sleep(10 * time.Millisecond)
	return false
}
