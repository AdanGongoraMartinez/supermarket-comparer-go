package categories

import (
	"gorm.io/gorm"
	"supermarket-comparer-go/internal/core"
	"supermarket-comparer-go/internal/database"
	"supermarket-comparer-go/internal/entities"
	"supermarket-comparer-go/internal/errors"
)

type CategoryRepositoryImpl struct{}

func NewCategoryRepository() *CategoryRepositoryImpl {
	return &CategoryRepositoryImpl{}
}

func (r *CategoryRepositoryImpl) Create(input CreateCategoryInput) *core.Result[entities.Category] {
	category := entities.CategoryModel{
		Name: input.Name,
	}

	result := database.DB.Create(&category)
	if result.Error != nil {
		return core.Fail[entities.Category](errors.NewDatabaseError("failed to create category", result.Error))
	}

	return core.Ok(r.mapModelToEntity(category))
}

func (r *CategoryRepositoryImpl) FindByID(id string) *core.Result[entities.Category] {
	var category entities.CategoryModel
	result := database.DB.First(&category, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return core.Fail[entities.Category](&errors.CategoryNotFoundError{ID: id})
		}
		return core.Fail[entities.Category](errors.NewDatabaseError("failed to find category", result.Error))
	}

	return core.Ok(r.mapModelToEntity(category))
}

func (r *CategoryRepositoryImpl) FindByName(name string) *core.Result[[]entities.Category] {
	var categories []entities.CategoryModel
	result := database.DB.Where("name = ?", name).Order("name").Find(&categories)
	if result.Error != nil {
		return core.Fail[[]entities.Category](errors.NewDatabaseError("failed to find categories", result.Error))
	}

	entitiesList := make([]entities.Category, len(categories))
	for i, c := range categories {
		entitiesList[i] = r.mapModelToEntity(c)
	}

	return core.Ok(entitiesList)
}

func (r *CategoryRepositoryImpl) Search(filters CategorySearchFilters) *core.Result[[]entities.Category] {
	var categories []entities.CategoryModel
	query := database.DB.Model(&entities.CategoryModel{})

	if filters.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filters.Name+"%")
	}

	result := query.Find(&categories)
	if result.Error != nil {
		return core.Fail[[]entities.Category](errors.NewDatabaseError("failed to search categories", result.Error))
	}

	entitiesList := make([]entities.Category, len(categories))
	for i, c := range categories {
		entitiesList[i] = r.mapModelToEntity(c)
	}

	return core.Ok(entitiesList)
}

func (r *CategoryRepositoryImpl) Delete(id string) *core.Result[any] {
	result := database.DB.Delete(&entities.CategoryModel{}, "id = ?", id)
	if result.Error != nil {
		return core.Fail[any](errors.NewDatabaseError("failed to delete category", result.Error))
	}

	if result.RowsAffected == 0 {
		return core.Fail[any](&errors.CategoryNotFoundError{ID: id})
	}

	return core.Ok[any](nil)
}

func (r *CategoryRepositoryImpl) mapModelToEntity(model entities.CategoryModel) entities.Category {
	return entities.Category{
		BaseEntity: entities.BaseEntity{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
		},
		Name: model.Name,
	}
}