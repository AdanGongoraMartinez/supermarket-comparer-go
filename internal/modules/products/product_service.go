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

func (s *ProductService) CreateProduct(input CreateProductInput) *core.Result[entities.Product] {
	if err := s.validateCreateInput(input); err != nil {
		return core.Fail[entities.Product](err)
	}

	existingResult := s.repo.FindByName(input.Name)
	if !existingResult.IsSuccess() {
		return core.Fail[entities.Product](existingResult.GetError())
	}

	existingProducts := existingResult.GetValue()
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
			return core.Fail[entities.Product](&errors.ProductAlreadyExistsError{
				Name:        input.Name,
				Presentation: presentation,
			})
		}
	}

	createResult := s.repo.Create(input)
	if !createResult.IsSuccess() {
		return core.Fail[entities.Product](createResult.GetError())
	}

	return core.Ok(createResult.GetValue())
}

func (s *ProductService) GetProductByID(id string) *core.Result[entities.Product] {
	if !core.IsValidUUIDString(id) {
		return core.Fail[entities.Product](&errors.InvalidProductIDError{ID: id})
	}

	result := s.repo.FindByID(id)
	if !result.IsSuccess() {
		return core.Fail[entities.Product](result.GetError())
	}

	return core.Ok(result.GetValue())
}

func (s *ProductService) SearchProducts(filters ProductSearchFilters) *core.Result[[]entities.Product] {
	result := s.repo.Search(filters)
	if !result.IsSuccess() {
		return core.Fail[[]entities.Product](result.GetError())
	}

	return core.Ok(result.GetValue())
}

func (s *ProductService) DeactivateProduct(id string) *core.Result[any] {
	if !core.IsValidUUIDString(id) {
		return core.Fail[any](&errors.InvalidProductIDError{ID: id})
	}

	existingResult := s.repo.FindByID(id)
	if !existingResult.IsSuccess() {
		return core.Fail[any](existingResult.GetError())
	}

	result := s.repo.Deactivate(id)
	if !result.IsSuccess() {
		return core.Fail[any](result.GetError())
	}

	return core.Ok[any](nil)
}

func (s *ProductService) validateCreateInput(input CreateProductInput) error {
	if input.Name == "" {
		return &errors.InvalidProductNameError{Name: input.Name}
	}
	return nil
}