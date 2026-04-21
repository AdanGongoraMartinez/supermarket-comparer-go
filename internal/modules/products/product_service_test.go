package products

import (
	"testing"

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

func (r *FakeProductRepository) Create(input CreateProductInput) (entities.Product, error) {
	if r.findErr != nil {
		return entities.Product{}, r.findErr
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
	return product, nil
}

func (r *FakeProductRepository) FindByID(id string) (entities.Product, error) {
	if r.findErr != nil {
		return entities.Product{}, r.findErr
	}
	for _, p := range r.products {
		if p.ID == id {
			return p, nil
		}
	}
	return entities.Product{}, &errors.ProductNotFoundError{ID: id}
}

func (r *FakeProductRepository) FindByName(name string) ([]entities.Product, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	var filtered []entities.Product
	for _, p := range r.products {
		if p.Name == name {
			filtered = append(filtered, p)
		}
	}
	return filtered, nil
}

func (r *FakeProductRepository) Search(filters ProductSearchFilters) ([]entities.Product, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	if len(r.products) == 0 {
		return nil, &errors.ProductNotFoundError{ID: filters.Name}
	}
	return r.products, nil
}

func (r *FakeProductRepository) Deactivate(id string) error {
	if r.findErr != nil {
		return r.findErr
	}
	for i, p := range r.products {
		if p.ID == id {
			r.products[i].Active = false
			return nil
		}
	}
	return &errors.ProductNotFoundError{ID: id}
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

	_, err := service.CreateProduct(input)

	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestCreateProduct_InvalidName(t *testing.T) {
	repo := NewFakeProductRepository()
	service := NewProductService(repo)

	input := CreateProductInput{
		Name: "",
	}

	_, err := service.CreateProduct(input)

	if err == nil {
		t.Error("expected failure for empty name")
	}
}

func TestGetProductByID_NotFound(t *testing.T) {
	repo := NewFakeProductRepository()
	service := NewProductService(repo)

	_, err := service.GetProductByID("non-existent-id")

	if err == nil {
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

	products, err := service.SearchProducts(filters)

	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}

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

	err := service.DeactivateProduct("550e8400-e29b-41d4-a716-446655440000")

	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestDeactivateProduct_InvalidID(t *testing.T) {
	repo := NewFakeProductRepository()
	service := NewProductService(repo)

	err := service.DeactivateProduct("invalid-id")

	if err == nil {
		t.Error("expected failure for invalid UUID")
	}
}