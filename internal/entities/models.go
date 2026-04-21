package entities

import "time"

type CategoryModel struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"type:text;not null" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

type ProductModel struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name         string    `gorm:"type:text;not null" json:"name"`
	Brand        *string   `gorm:"type:text" json:"brand,omitempty"`
	Presentation *string   `gorm:"type:text" json:"presentation,omitempty"`
	Barcode      *string   `gorm:"type:text" json:"barcode,omitempty"`
	CategoryID   *string   `gorm:"type:uuid" json:"categoryId,omitempty"`
	Active       bool      `gorm:"type:boolean;default:true" json:"active"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (CategoryModel) TableName() string {
	return "categories"
}

func (ProductModel) TableName() string {
	return "products"
}

