package products

import (
	"testing"

	"supermarket-comparer-go/internal/core"
	"supermarket-comparer-go/internal/entities"
	"supermarket-comparer-go/internal/errors"
)

type FakeProductRepository struct {
	products []entities.Product
	findErr  error
}

func NewFakeProductRepository() *FakeProductRepository {
	return &FakeProductRepository{products: []entities.Product{}}
}

func (r *FakeProductRepository) Create(input CreateProductInput) *core.Result[entities.Product] {
	if r.findErr != nil {
		return core.Fail[entities.Product](r.findErr)
	}
	product := entities.Product{
		BaseEntity: entities.BaseEntity{
			ID: "test-id",
		},
		Name:         input.Name,
		Brand:       strPtrToStr(input.Brand),
		Presentation: strPtrToStr(input.Presentation),
		Active:      true,
	}
	r.products = append(r.products, product)
	return core.Ok(product)
}

func (r *FakeProductRepository) FindByID(id string) *core.Result[entities.Product] {
	if r.findErr != nil {
		return core.Fail[entities.Product](r.findErr)
	}
	for _, p := range r.products {
		if p.ID == id {
			return core.Ok(p)
		}
	}
	return core.Fail[entities.Product](&errors.ProductNotFoundError{ID: id})
}

func (r *FakeProductRepository) FindByName(name string) *core.Result[[]entities.Product] {
	if r.findErr != nil {
		return core.Fail[[]entities.Product](r.findErr)
	}
	var filtered []entities.Product
	for _, p := range r.products {
		if p.Name == name {
			filtered = append(filtered, p)
		}
	}
	return core.Ok(filtered)
}

func (r *FakeProductRepository) Search(filters ProductSearchFilters) *core.Result[[]entities.Product] {
	if r.findErr != nil {
		return core.Fail[[]entities.Product](r.findErr)
	}
	if len(r.products) == 0 {
		return core.Fail[[]entities.Product](&errors.ProductNotFoundError{ID: filters.Name})
	}
	return core.Ok(r.products)
}

func (r *FakeProductRepository) Deactivate(id string) *core.Result[any] {
	if r.findErr != nil {
		return core.Fail[any](r.findErr)
	}
	for i, p := range r.products {
		if p.ID == id {
			r.products[i].Active = false
			return core.Ok[any](nil)
		}
	}
	return core.Fail[any](&errors.ProductNotFoundError{ID: id})
}

func strPtrToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func strPtr(s string) *string {
	return &s
}

func TestCreateProduct_Success(t *testing.T) {
	repo := NewFakeProductRepository()
	service := NewProductService(repo)

	input := CreateProductInput{
		Name:         "Milk",
		Brand:        strPtr("Parmalat"),
		Presentation: strPtr("1L"),
	}

	result := service.CreateProduct(input)

	if !result.IsSuccess() {
		t.Errorf("expected success, got error: %v", result.GetError())
	}

	value := result.GetValue()
	if value.Name != "Milk" {
		t.Errorf("expected name Milk, got %s", value.Name)
	}
}

func TestCreateProduct_InvalidName(t *testing.T) {
	repo := NewFakeProductRepository()
	service := NewProductService(repo)

	input := CreateProductInput{
		Name: "",
	}

	result := service.CreateProduct(input)

	if result.IsSuccess() {
		t.Error("expected failure for empty name")
	}

	err := result.GetError()
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestGetProductByID_NotFound(t *testing.T) {
	repo := NewFakeProductRepository()
	service := NewProductService(repo)

	result := service.GetProductByID("non-existent-id")

	if result.IsSuccess() {
		t.Error("expected failure for non-existent product")
	}
}

func TestSearchProducts_Success(t *testing.T) {
	repo := NewFakeProductRepository()
	repo.products = []entities.Product{
		{BaseEntity: entities.BaseEntity{ID: "1"}, Name: "Milk"},
		{BaseEntity: entities.BaseEntity{ID: "2"}, Name: "Bread"},
	}
	service := NewProductService(repo)

	filters := ProductSearchFilters{
		Name:       "",
		ActiveOnly: true,
	}

	result := service.SearchProducts(filters)

	if !result.IsSuccess() {
		t.Errorf("expected success, got error: %v", result.GetError())
	}

	products := result.GetValue()
	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
}

func TestDeactivateProduct_Success(t *testing.T) {
	repo := NewFakeProductRepository()
	repo.products = []entities.Product{
		{BaseEntity: entities.BaseEntity{ID: "550e8400-e29b-41d4-a716-446655440000"}, Name: "Milk", Active: true},
	}
	service := NewProductService(repo)

	result := service.DeactivateProduct("550e8400-e29b-41d4-a716-446655440000")

	if !result.IsSuccess() {
		t.Errorf("expected success, got error: %v", result.GetError())
	}
}

func TestDeactivateProduct_InvalidID(t *testing.T) {
	repo := NewFakeProductRepository()
	service := NewProductService(repo)

	result := service.DeactivateProduct("invalid-id")

	if result.IsSuccess() {
		t.Error("expected failure for invalid UUID")
	}
}