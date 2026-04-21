package products

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"supermarket-comparer-go/internal/core"
)

func ProductHandler(service *ProductService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/products":
			createProduct(w, r, service)
		case r.Method == "GET" && r.URL.Path == "/products":
			searchProducts(w, r, service)
		case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/products/"):
			id := strings.TrimPrefix(r.URL.Path, "/products/")
			getProductByID(w, r, service, id)
		case r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, "/products/"):
			id := strings.TrimPrefix(r.URL.Path, "/products/")
			deleteProduct(w, r, service, id)
		default:
			http.NotFound(w, r)
		}
	})
}

func createProduct(w http.ResponseWriter, r *http.Request, service *ProductService) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	var input CreateProductInput
	if err := json.Unmarshal(body, &input); err != nil {
		sendError(w, http.StatusBadRequest, "invalid json")
		return
	}

	value, err := service.CreateProduct(input)
	response := core.HandleResult(value, err, 201)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func searchProducts(w http.ResponseWriter, r *http.Request, service *ProductService) {
	query := r.URL.Query()

	filters := ProductSearchFilters{
		Name:       query.Get("name"),
		CategoryID: query.Get("categoryId"),
		ActiveOnly: query.Get("activeOnly") != "false",
	}

	value, err := service.SearchProducts(filters)
	response := core.HandleResult(value, err, 200)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func getProductByID(w http.ResponseWriter, r *http.Request, service *ProductService, id string) {
	value, err := service.GetProductByID(id)
	response := core.HandleResult(value, err, 200)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func deleteProduct(w http.ResponseWriter, r *http.Request, service *ProductService, id string) {
	err := service.DeactivateProduct(id)
	response := core.HandleEmptyResult(err, 204)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	if response.StatusCode != 204 {
		json.NewEncoder(w).Encode(response)
	}
}

func sendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(core.APIResponse{
		Success: false,
		Error:   message,
	})
}