package categories

import (
	"testing"

	"supermarket-comparer-go/internal/core"
	"supermarket-comparer-go/internal/entities"
	"supermarket-comparer-go/internal/errors"
)

type FakeCategoryRepository struct {
	categories []entities.Category
	findErr   error
}

func NewFakeCategoryRepository() *FakeCategoryRepository {
	return &FakeCategoryRepository{categories: []entities.Category{}}
}

func (r *FakeCategoryRepository) Create(input CreateCategoryInput) *core.Result[entities.Category] {
	if r.findErr != nil {
		return core.Fail[entities.Category](r.findErr)
	}
	category := entities.Category{
		BaseEntity: entities.BaseEntity{
			ID: "test-id",
		},
		Name: input.Name,
	}
	r.categories = append(r.categories, category)
	return core.Ok(category)
}

func (r *FakeCategoryRepository) FindByID(id string) *core.Result[entities.Category] {
	if r.findErr != nil {
		return core.Fail[entities.Category](r.findErr)
	}
	for _, c := range r.categories {
		if c.ID == id {
			return core.Ok(c)
		}
	}
	return core.Fail[entities.Category](&errors.CategoryNotFoundError{ID: id})
}

func (r *FakeCategoryRepository) FindByName(name string) *core.Result[[]entities.Category] {
	if r.findErr != nil {
		return core.Fail[[]entities.Category](r.findErr)
	}
	var filtered []entities.Category
	for _, c := range r.categories {
		if c.Name == name {
			filtered = append(filtered, c)
		}
	}
	return core.Ok(filtered)
}

func (r *FakeCategoryRepository) Search(filters CategorySearchFilters) *core.Result[[]entities.Category] {
	if r.findErr != nil {
		return core.Fail[[]entities.Category](r.findErr)
	}
	return core.Ok(r.categories)
}

func (r *FakeCategoryRepository) Delete(id string) *core.Result[any] {
	if r.findErr != nil {
		return core.Fail[any](r.findErr)
	}
	for i, c := range r.categories {
		if c.ID == id {
			r.categories = append(r.categories[:i], r.categories[i+1:]...)
			return core.Ok[any](nil)
		}
	}
	return core.Fail[any](&errors.CategoryNotFoundError{ID: id})
}

func TestCreateCategory_Success(t *testing.T) {
	repo := NewFakeCategoryRepository()
	service := NewCategoryService(repo)

	input := CreateCategoryInput{
		Name: "Dairy",
	}

	result := service.CreateCategory(input)

	if !result.IsSuccess() {
		t.Errorf("expected success, got error: %v", result.GetError())
	}

	value := result.GetValue()
	if value.Name != "Dairy" {
		t.Errorf("expected name Dairy, got %s", value.Name)
	}
}

func TestCreateCategory_InvalidName(t *testing.T) {
	repo := NewFakeCategoryRepository()
	service := NewCategoryService(repo)

	input := CreateCategoryInput{
		Name: "",
	}

	result := service.CreateCategory(input)

	if result.IsSuccess() {
		t.Error("expected failure for empty name")
	}

	err := result.GetError()
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestGetCategoryByID_NotFound(t *testing.T) {
	repo := NewFakeCategoryRepository()
	service := NewCategoryService(repo)

	result := service.GetCategoryByID("non-existent-id")

	if result.IsSuccess() {
		t.Error("expected failure for non-existent category")
	}
}

func TestSearchCategories_Success(t *testing.T) {
	repo := NewFakeCategoryRepository()
	repo.categories = []entities.Category{
		{BaseEntity: entities.BaseEntity{ID: "1"}, Name: "Dairy"},
		{BaseEntity: entities.BaseEntity{ID: "2"}, Name: "Bakery"},
	}
	service := NewCategoryService(repo)

	filters := CategorySearchFilters{
		Name: "",
	}

	result := service.SearchCategories(filters)

	if !result.IsSuccess() {
		t.Errorf("expected success, got error: %v", result.GetError())
	}

	categories := result.GetValue()
	if len(categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(categories))
	}
}

func TestDeleteCategory_Success(t *testing.T) {
	repo := NewFakeCategoryRepository()
	repo.categories = []entities.Category{
		{BaseEntity: entities.BaseEntity{ID: "550e8400-e29b-41d4-a716-446655440000"}, Name: "Dairy"},
	}
	service := NewCategoryService(repo)

	result := service.DeleteCategory("550e8400-e29b-41d4-a716-446655440000")

	if !result.IsSuccess() {
		t.Errorf("expected success, got error: %v", result.GetError())
	}
}

func TestDeleteCategory_InvalidID(t *testing.T) {
	repo := NewFakeCategoryRepository()
	service := NewCategoryService(repo)

	result := service.DeleteCategory("invalid-id")

	if result.IsSuccess() {
		t.Error("expected failure for invalid UUID")
	}
}