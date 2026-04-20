package errors

type InvalidProductNameError struct {
	Name string
}

func (e *InvalidProductNameError) Error() string {
	return "invalid product name: " + e.Name
}

type InvalidProductIDError struct {
	ID string
}

func (e *InvalidProductIDError) Error() string {
	return "invalid product id: " + e.ID
}

type ProductNotFoundError struct {
	ID string
}

func (e *ProductNotFoundError) Error() string {
	return "product not found: " + e.ID
}

type ProductAlreadyExistsError struct {
	Name        string
	Presentation string
}

func (e *ProductAlreadyExistsError) Error() string {
	return "product already exists: " + e.Name + " (" + e.Presentation + ")"
}