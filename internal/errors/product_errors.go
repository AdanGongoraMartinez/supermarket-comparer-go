package errors

type InvalidProductNameError struct {
	name       string
	statusCode int
}

func (e *InvalidProductNameError) Error() string {
	return "invalid product name: " + e.name
}

func (e *InvalidProductNameError) GetStatusCode() int { return e.statusCode }

func NewInvalidProductNameError(name string) *InvalidProductNameError {
	return &InvalidProductNameError{name: name, statusCode: 400}
}

type InvalidProductIDError struct {
	ID         string
	statusCode int
}

func (e *InvalidProductIDError) Error() string {
	return "invalid product id: " + e.ID
}

func (e *InvalidProductIDError) GetStatusCode() int { return e.statusCode }

func NewInvalidProductIDError(id string) *InvalidProductIDError {
	return &InvalidProductIDError{ID: id, statusCode: 400}
}

type ProductNotFoundError struct {
	ID         string
	statusCode int
}

func (e *ProductNotFoundError) Error() string {
	return "product not found: " + e.ID
}

func (e *ProductNotFoundError) GetStatusCode() int { return e.statusCode }

func NewProductNotFoundError(id string) *ProductNotFoundError {
	return &ProductNotFoundError{ID: id, statusCode: 404}
}

type ProductAlreadyExistsError struct {
	Name         string
	Presentation string
	statusCode   int
}

func (e *ProductAlreadyExistsError) Error() string {
	return "product already exists: " + e.Name + " (" + e.Presentation + ")"
}

func (e *ProductAlreadyExistsError) GetStatusCode() int { return e.statusCode }

func NewProductAlreadyExistsError(name, presentation string) *ProductAlreadyExistsError {
	return &ProductAlreadyExistsError{Name: name, Presentation: presentation, statusCode: 409}
}
