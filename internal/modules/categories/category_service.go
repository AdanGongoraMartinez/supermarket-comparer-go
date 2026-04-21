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

func (s *CategoryService) CreateCategory(input CreateCategoryInput) (entities.Category, error) {
	if err := s.validateCreateInput(input); err != nil {
		return entities.Category{}, err
	}

	existingCategories, err := s.repo.FindByName(input.Name)
	if err != nil {
		return entities.Category{}, err
	}

	for _, c := range existingCategories {
		if c.Name == input.Name {
			return entities.Category{}, &errors.CategoryAlreadyExistsError{Name: input.Name}
		}
	}

	return s.repo.Create(input)
}

func (s *CategoryService) GetCategoryByID(id string) (entities.Category, error) {
	if !core.IsValidUUIDString(id) {
		return entities.Category{}, &errors.CategoryNotFoundError{ID: id}
	}

	return s.repo.FindByID(id)
}

func (s *CategoryService) SearchCategories(filters CategorySearchFilters) ([]entities.Category, error) {
	return s.repo.Search(filters)
}

func (s *CategoryService) DeleteCategory(id string) error {
	if !core.IsValidUUIDString(id) {
		return &errors.CategoryNotFoundError{ID: id}
	}

	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *CategoryService) validateCreateInput(input CreateCategoryInput) error {
	if input.Name == "" {
		return &errors.InvalidCategoryNameError{}
	}
	return nil
}