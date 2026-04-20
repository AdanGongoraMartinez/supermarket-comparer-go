package categories

type CreateCategoryInput struct {
	Name string `json:"name"`
}

type CategorySearchFilters struct {
	Name string `json:"name,omitempty"`
}