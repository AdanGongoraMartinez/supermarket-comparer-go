package errors

type InvalidCategoryNameError struct{}

func (e *InvalidCategoryNameError) Error() string {
	return "invalid category name: name is required"
}

type CategoryNotFoundError struct {
	ID string
}

func (e *CategoryNotFoundError) Error() string {
	return "category not found: " + e.ID
}

type CategoryAlreadyExistsError struct {
	Name string
}

func (e *CategoryAlreadyExistsError) Error() string {
	return "category already exists: " + e.Name
}