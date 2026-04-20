package products

import (
	"supermarket-comparer-go/internal/core"
	"supermarket-comparer-go/internal/entities"
)

type ProductRepository interface {
	Create(input CreateProductInput) *core.Result[entities.Product]
	FindByID(id string) *core.Result[entities.Product]
	FindByName(name string) *core.Result[[]entities.Product]
	Search(filters ProductSearchFilters) *core.Result[[]entities.Product]
	Deactivate(id string) *core.Result[any]
}