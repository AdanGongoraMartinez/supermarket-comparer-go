package categories

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"supermarket-comparer-go/internal/core"
)

func CategoryHandler(service *CategoryService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/categories":
			createCategory(w, r, service)
		case r.Method == "GET" && r.URL.Path == "/categories":
			searchCategories(w, r, service)
		case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/categories/"):
			id := strings.TrimPrefix(r.URL.Path, "/categories/")
			getCategoryByID(w, r, service, id)
		case r.Method == "DELETE" && strings.HasPrefix(r.URL.Path, "/categories/"):
			id := strings.TrimPrefix(r.URL.Path, "/categories/")
			deleteCategory(w, r, service, id)
		default:
			http.NotFound(w, r)
		}
	})
}

func createCategory(w http.ResponseWriter, r *http.Request, service *CategoryService) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	var input CreateCategoryInput
	if err := json.Unmarshal(body, &input); err != nil {
		sendError(w, http.StatusBadRequest, "invalid json")
		return
	}

	result := service.CreateCategory(input)
	response := core.HandleResult(result, 201)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func searchCategories(w http.ResponseWriter, r *http.Request, service *CategoryService) {
	query := r.URL.Query()

	filters := CategorySearchFilters{
		Name: query.Get("name"),
	}

	result := service.SearchCategories(filters)
	response := core.HandleResult(result, 200)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func getCategoryByID(w http.ResponseWriter, r *http.Request, service *CategoryService, id string) {
	result := service.GetCategoryByID(id)
	response := core.HandleResult(result, 200)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func deleteCategory(w http.ResponseWriter, r *http.Request, service *CategoryService, id string) {
	result := service.DeleteCategory(id)
	response := core.HandleEmptyResult(result, 204)
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
		Error:  message,
	})
}