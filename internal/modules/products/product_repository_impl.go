package products

import (
	"strings"

	"supermarket-comparer-go/internal/database"
	"supermarket-comparer-go/internal/entities"
	"supermarket-comparer-go/internal/errors"

	"gorm.io/gorm"
)

type ProductRepositoryImpl struct{}

func NewProductRepository() *ProductRepositoryImpl {
	return &ProductRepositoryImpl{}
}

func (r *ProductRepositoryImpl) Create(input CreateProductInput) (entities.Product, error) {
	product := entities.ProductModel{
		Name:         input.Name,
		Brand:        input.Brand,
		Presentation: input.Presentation,
		Barcode:      input.Barcode,
		CategoryID:   input.CategoryID,
		Active:       true,
	}

	result := database.DB.Create(&product)
	if result.Error != nil {
		return entities.Product{}, errors.NewDatabaseError("failed to create product", result.Error)
	}

	return r.mapModelToEntity(product), nil
}

func (r *ProductRepositoryImpl) FindByID(id string) (entities.Product, error) {
	var product entities.ProductModel
	result := database.DB.First(&product, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return entities.Product{}, &errors.ProductNotFoundError{ID: id}
		}
		return entities.Product{}, errors.NewDatabaseError("failed to find product", result.Error)
	}

	return r.mapModelToEntity(product), nil
}

func (r *ProductRepositoryImpl) FindByName(name string) ([]entities.Product, error) {
	var products []entities.ProductModel
	result := database.DB.Where("name = ?", name).Order("name").Find(&products)
	if result.Error != nil {
		return nil, errors.NewDatabaseError("failed to find products", result.Error)
	}

	entitiesList := make([]entities.Product, len(products))
	for i, p := range products {
		entitiesList[i] = r.mapModelToEntity(p)
	}

	return entitiesList, nil
}

func (r *ProductRepositoryImpl) Search(filters ProductSearchFilters) ([]entities.Product, error) {
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
		return nil, errors.NewDatabaseError("failed to search products", result.Error)
	}

	entitiesList := make([]entities.Product, len(products))
	for i, p := range products {
		entitiesList[i] = r.mapModelToEntity(p)
	}

	return entitiesList, nil
}

func (r *ProductRepositoryImpl) Deactivate(id string) error {
	result := database.DB.Model(&entities.ProductModel{}).Where("id = ?", id).Update("active", false)
	if result.Error != nil {
		return errors.NewDatabaseError("failed to deactivate product", result.Error)
	}

	if result.RowsAffected == 0 {
		return &errors.ProductNotFoundError{ID: id}
	}

	return nil
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
		CategoryID:   stringOrNil(model.CategoryID),
		Active:       model.Active,
	}
}

func stringOrNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (r *ProductRepositoryImpl) FindByNameAndBrand(name, brand, presentation string) (entities.Product, error) {
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
			return entities.Product{}, nil
		}
		return entities.Product{}, errors.NewDatabaseError("failed to find product", result.Error)
	}

	return r.mapModelToEntity(product), nil
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
