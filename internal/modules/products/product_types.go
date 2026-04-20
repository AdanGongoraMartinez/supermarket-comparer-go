package products

type CreateProductInput struct {
	Name         string  `json:"name"`
	Brand        *string `json:"brand,omitempty"`
	Presentation *string `json:"presentation,omitempty"`
	Barcode      *string `json:"barcode,omitempty"`
	CategoryID  *string `json:"categoryId,omitempty"`
}

type ProductSearchFilters struct {
	Name       string `json:"name,omitempty"`
	CategoryID string `json:"categoryId,omitempty"`
	ActiveOnly bool   `json:"activeOnly"`
}