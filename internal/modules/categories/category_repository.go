package categories

import (
	"supermarket-comparer-go/internal/entities"
)

type CategoryRepository interface {
	Create(input CreateCategoryInput) (entities.Category, error)
	FindByID(id string) (entities.Category, error)
	FindByName(name string) ([]entities.Category, error)
	Search(filters CategorySearchFilters) ([]entities.Category, error)
	Delete(id string) error
}