package categories

import (
	"supermarket-comparer-go/internal/core"
	"supermarket-comparer-go/internal/entities"
)

type CategoryRepository interface {
	Create(input CreateCategoryInput) *core.Result[entities.Category]
	FindByID(id string) *core.Result[entities.Category]
	FindByName(name string) *core.Result[[]entities.Category]
	Search(filters CategorySearchFilters) *core.Result[[]entities.Category]
	Delete(id string) *core.Result[any]
}