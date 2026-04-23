package errors

type InvalidCategoryNameError struct {
	statusCode int
}

func (e *InvalidCategoryNameError) Error() string {
	return "invalid category name: name is required"
}

func (e *InvalidCategoryNameError) GetStatusCode() int { return e.statusCode }

func NewInvalidCategoryNameError(name string) *InvalidCategoryNameError {
	return &InvalidCategoryNameError{statusCode: 400}
}

type CategoryNotFoundError struct {
	ID         string
	statusCode int
}

func (e *CategoryNotFoundError) Error() string {
	return "category not found: " + e.ID
}

func (e *CategoryNotFoundError) GetStatusCode() int { return e.statusCode }

func NewCategoryNotFoundError(name string) *CategoryNotFoundError {
	return &CategoryNotFoundError{statusCode: 404}
}

type CategoryAlreadyExistsError struct {
	Name       string
	statusCode int
}

func (e *CategoryAlreadyExistsError) Error() string {
	return "category already exists: " + e.Name
}

func (e *CategoryAlreadyExistsError) GetStatusCode() int { return e.statusCode }

func NewCategoryAlreadyExistsError(name string) *CategoryAlreadyExistsError {
	return &CategoryAlreadyExistsError{statusCode: 409}
}

type InvalidCategoryIDError struct {
	ID         string
	statusCode int
}

func (e *InvalidCategoryIDError) Error() string {
	return "invalid category id: " + e.ID
}

func (e *InvalidCategoryIDError) GetStatusCode() int { return e.statusCode }

func NewInvalidCategoryIDError(name string) *InvalidCategoryIDError {
	return &InvalidCategoryIDError{statusCode: 409}
}

