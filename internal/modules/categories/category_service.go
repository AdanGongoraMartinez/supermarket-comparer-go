package categories

import (
	"supermarket-comparer-go/internal/core"
	"supermarket-comparer-go/internal/entities"
	"supermarket-comparer-go/internal/errors"
)

type CategoryService struct {
	repo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(input CreateCategoryInput) *core.Result[entities.Category] {
	if err := s.validateCreateInput(input); err != nil {
		return core.Fail[entities.Category](err)
	}

	result := s.repo.FindByName(input.Name)
	if !result.IsSuccess() {
		return core.Fail[entities.Category](result.GetError())
	}

	existingCategories := result.GetValue()
	for _, c := range existingCategories {
		if c.Name == input.Name {
			return core.Fail[entities.Category](&errors.CategoryAlreadyExistsError{Name: input.Name})
		}
	}

	createResult := s.repo.Create(input)
	if !createResult.IsSuccess() {
		return core.Fail[entities.Category](createResult.GetError())
	}

	return core.Ok(createResult.GetValue())
}

func (s *CategoryService) GetCategoryByID(id string) *core.Result[entities.Category] {
	if !core.IsValidUUIDString(id) {
		return core.Fail[entities.Category](&errors.CategoryNotFoundError{ID: id})
	}

	result := s.repo.FindByID(id)
	if !result.IsSuccess() {
		return core.Fail[entities.Category](result.GetError())
	}

	return core.Ok(result.GetValue())
}

func (s *CategoryService) SearchCategories(filters CategorySearchFilters) *core.Result[[]entities.Category] {
	result := s.repo.Search(filters)
	if !result.IsSuccess() {
		return core.Fail[[]entities.Category](result.GetError())
	}

	return core.Ok(result.GetValue())
}

func (s *CategoryService) DeleteCategory(id string) *core.Result[any] {
	if !core.IsValidUUIDString(id) {
		return core.Fail[any](&errors.CategoryNotFoundError{ID: id})
	}

	existingResult := s.repo.FindByID(id)
	if !existingResult.IsSuccess() {
		return core.Fail[any](existingResult.GetError())
	}

	result := s.repo.Delete(id)
	if !result.IsSuccess() {
		return core.Fail[any](result.GetError())
	}

	return core.Ok[any](nil)
}

func (s *CategoryService) validateCreateInput(input CreateCategoryInput) error {
	if input.Name == "" {
		return &errors.InvalidCategoryNameError{}
	}
	return nil
}