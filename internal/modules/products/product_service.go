package products

import (
	"supermarket-comparer-go/internal/core"
	"supermarket-comparer-go/internal/entities"
	"supermarket-comparer-go/internal/errors"
)

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(input CreateProductInput) (entities.Product, error) {
	if err := s.validateCreateInput(input); err != nil {
		return entities.Product{}, err
	}

	existingProducts, err := s.repo.FindByName(input.Name)
	if err != nil {
		return entities.Product{}, err
	}

	brand := ""
	if input.Brand != nil {
		brand = *input.Brand
	}
	presentation := ""
	if input.Presentation != nil {
		presentation = *input.Presentation
	}

	for _, p := range existingProducts {
		if p.Name == input.Name && p.Brand == brand && p.Presentation == presentation {
			return entities.Product{}, &errors.ProductAlreadyExistsError{
				Name:         input.Name,
				Presentation: presentation,
			}
		}
	}

	return s.repo.Create(input)
}

func (s *ProductService) GetProductByID(id string) (entities.Product, error) {
	if !core.IsValidUUIDString(id) {
		return entities.Product{}, errors.NewInvalidProductIDError(id)
	}

	return s.repo.FindByID(id)
}

func (s *ProductService) SearchProducts(filters ProductSearchFilters) ([]entities.Product, error) {
	return s.repo.Search(filters)
}

func (s *ProductService) DeactivateProduct(id string) error {
	if !core.IsValidUUIDString(id) {
		return errors.NewInvalidProductIDError(id)
	}

	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Deactivate(id)
}

func (s *ProductService) validateCreateInput(input CreateProductInput) error {
	if input.Name == "" {
		return errors.NewInvalidProductNameError(input.Name)
	}
	return nil
}
