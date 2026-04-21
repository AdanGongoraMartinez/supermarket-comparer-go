package categories

import (
	"testing"

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

func (r *FakeCategoryRepository) Create(input CreateCategoryInput) (entities.Category, error) {
	if r.findErr != nil {
		return entities.Category{}, r.findErr
	}
	category := entities.Category{
		BaseEntity: entities.BaseEntity{
			ID: "test-id",
		},
		Name: input.Name,
	}
	r.categories = append(r.categories, category)
	return category, nil
}

func (r *FakeCategoryRepository) FindByID(id string) (entities.Category, error) {
	if r.findErr != nil {
		return entities.Category{}, r.findErr
	}
	for _, c := range r.categories {
		if c.ID == id {
			return c, nil
		}
	}
	return entities.Category{}, &errors.CategoryNotFoundError{ID: id}
}

func (r *FakeCategoryRepository) FindByName(name string) ([]entities.Category, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	var filtered []entities.Category
	for _, c := range r.categories {
		if c.Name == name {
			filtered = append(filtered, c)
		}
	}
	return filtered, nil
}

func (r *FakeCategoryRepository) Search(filters CategorySearchFilters) ([]entities.Category, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	return r.categories, nil
}

func (r *FakeCategoryRepository) Delete(id string) error {
	if r.findErr != nil {
		return r.findErr
	}
	for i, c := range r.categories {
		if c.ID == id {
			r.categories = append(r.categories[:i], r.categories[i+1:]...)
			return nil
		}
	}
	return &errors.CategoryNotFoundError{ID: id}
}

func TestCreateCategory_Success(t *testing.T) {
	repo := NewFakeCategoryRepository()
	service := NewCategoryService(repo)

	input := CreateCategoryInput{
		Name: "Dairy",
	}

	_, err := service.CreateCategory(input)

	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestCreateCategory_InvalidName(t *testing.T) {
	repo := NewFakeCategoryRepository()
	service := NewCategoryService(repo)

	input := CreateCategoryInput{
		Name: "",
	}

	_, err := service.CreateCategory(input)

	if err == nil {
		t.Error("expected failure for empty name")
	}
}

func TestGetCategoryByID_NotFound(t *testing.T) {
	repo := NewFakeCategoryRepository()
	service := NewCategoryService(repo)

	_, err := service.GetCategoryByID("non-existent-id")

	if err == nil {
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

	categories, err := service.SearchCategories(filters)

	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}

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

	err := service.DeleteCategory("550e8400-e29b-41d4-a716-446655440000")

	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestDeleteCategory_InvalidID(t *testing.T) {
	repo := NewFakeCategoryRepository()
	service := NewCategoryService(repo)

	err := service.DeleteCategory("invalid-id")

	if err == nil {
		t.Error("expected failure for invalid UUID")
	}
}