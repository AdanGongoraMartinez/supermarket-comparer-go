package entities

type Product struct {
	BaseEntity
	Name        string  `json:"name"`
	Brand       string  `json:"brand,omitempty"`
	Presentation string `json:"presentation,omitempty"`
	Barcode    string  `json:"barcode,omitempty"`
	CategoryID string  `json:"categoryId,omitempty"`
	Active     bool    `json:"active"`
}