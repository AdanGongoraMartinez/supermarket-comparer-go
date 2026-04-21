package products

import (
	"supermarket-comparer-go/internal/entities"
)

type ProductRepository interface {
	Create(input CreateProductInput) (entities.Product, error)
	FindByID(id string) (entities.Product, error)
	FindByName(name string) ([]entities.Product, error)
	Search(filters ProductSearchFilters) ([]entities.Product, error)
	Deactivate(id string) error
}