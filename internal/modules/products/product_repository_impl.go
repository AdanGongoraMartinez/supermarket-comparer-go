package products

import (
	"strings"

	"gorm.io/gorm"
	"supermarket-comparer-go/internal/core"
	"supermarket-comparer-go/internal/database"
	"supermarket-comparer-go/internal/entities"
	"supermarket-comparer-go/internal/errors"
)

type ProductRepositoryImpl struct{}

func NewProductRepository() *ProductRepositoryImpl {
	return &ProductRepositoryImpl{}
}

func (r *ProductRepositoryImpl) Create(input CreateProductInput) *core.Result[entities.Product] {
	product := entities.ProductModel{
		Name:         input.Name,
		Brand:        input.Brand,
		Presentation: input.Presentation,
		Barcode:      input.Barcode,
		CategoryID:  input.CategoryID,
		Active:      true,
	}

	result := database.DB.Create(&product)
	if result.Error != nil {
		return core.Fail[entities.Product](errors.NewDatabaseError("failed to create product", result.Error))
	}

	return core.Ok(r.mapModelToEntity(product))
}

func (r *ProductRepositoryImpl) FindByID(id string) *core.Result[entities.Product] {
	var product entities.ProductModel
	result := database.DB.First(&product, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return core.Fail[entities.Product](&errors.ProductNotFoundError{ID: id})
		}
		return core.Fail[entities.Product](errors.NewDatabaseError("failed to find product", result.Error))
	}

	return core.Ok(r.mapModelToEntity(product))
}

func (r *ProductRepositoryImpl) FindByName(name string) *core.Result[[]entities.Product] {
	var products []entities.ProductModel
	result := database.DB.Where("name = ?", name).Order("name").Find(&products)
	if result.Error != nil {
		return core.Fail[[]entities.Product](errors.NewDatabaseError("failed to find products", result.Error))
	}

	entitiesList := make([]entities.Product, len(products))
	for i, p := range products {
		entitiesList[i] = r.mapModelToEntity(p)
	}

	return core.Ok(entitiesList)
}

func (r *ProductRepositoryImpl) Search(filters ProductSearchFilters) *core.Result[[]entities.Product] {
	var products []entities.ProductModel
	query := database.DB.Model(&entities.ProductModel{})

	if filters.ActiveOnly {
		query = query.Where("active = ?", true)
	}

	if filters.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filters.Name+"%")
	}

	if filters.CategoryID != "" {
		query = query.Where("category_id = ?", filters.CategoryID)
	}

	result := query.Find(&products)
	if result.Error != nil {
		return core.Fail[[]entities.Product](errors.NewDatabaseError("failed to search products", result.Error))
	}

	if len(products) == 0 {
		return core.Fail[[]entities.Product](&errors.ProductNotFoundError{ID: filters.Name})
	}

	entitiesList := make([]entities.Product, len(products))
	for i, p := range products {
		entitiesList[i] = r.mapModelToEntity(p)
	}

	return core.Ok(entitiesList)
}

func (r *ProductRepositoryImpl) Deactivate(id string) *core.Result[any] {
	result := database.DB.Model(&entities.ProductModel{}).Where("id = ?", id).Update("active", false)
	if result.Error != nil {
		return core.Fail[any](errors.NewDatabaseError("failed to deactivate product", result.Error))
	}

	if result.RowsAffected == 0 {
		return core.Fail[any](&errors.ProductNotFoundError{ID: id})
	}

	return core.Ok[any](nil)
}

func (r *ProductRepositoryImpl) mapModelToEntity(model entities.ProductModel) entities.Product {
	return entities.Product{
		BaseEntity: entities.BaseEntity{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
		},
		Name:         model.Name,
		Brand:        stringOrNil(model.Brand),
		Presentation: stringOrNil(model.Presentation),
		Barcode:      stringOrNil(model.Barcode),
		CategoryID:  stringOrNil(model.CategoryID),
		Active:      model.Active,
	}
}

func stringOrNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (r *ProductRepositoryImpl) FindByNameAndBrand(name, brand, presentation string) *core.Result[entities.Product] {
	var product entities.ProductModel
	query := database.DB.Where("name = ?", name)

	if brand != "" {
		query = query.Where("brand = ?", brand)
	} else {
		query = query.Where("brand IS NULL")
	}

	if presentation != "" {
		query = query.Where("presentation = ?", presentation)
	} else {
		query = query.Where("presentation IS NULL")
	}

	result := query.First(&product)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return core.Fail[entities.Product](nil)
		}
		return core.Fail[entities.Product](errors.NewDatabaseError("failed to find product", result.Error))
	}

	return core.Ok(r.mapModelToEntity(product))
}

func containsBrandAndPresentation(products []entities.ProductModel, name, brand, presentation string) bool {
	for _, p := range products {
		matchName := strings.EqualFold(p.Name, name)
		matchBrand := false
		if brand == "" {
			matchBrand = p.Brand == nil
		} else {
			matchBrand = p.Brand != nil && *p.Brand == brand
		}
		matchPresentation := false
		if presentation == "" {
			matchPresentation = p.Presentation == nil
		} else {
			matchPresentation = p.Presentation != nil && *p.Presentation == presentation
		}
		if matchName && matchBrand && matchPresentation {
			return true
		}
	}
	return false
}