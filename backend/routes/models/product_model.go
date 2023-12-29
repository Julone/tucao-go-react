package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
)

// Product struct to describe product object.
type Product struct {
	ID            uuid.UUID    `gorm:"column:id" json:"id" `
	UserID        uuid.UUID    `gorm:"column:user_id" json:"user_id" `
	UserInfo      User         `gorm:"foreignKey:ID" json:"user_info" validate:"-"`
	Title         string       `gorm:"column:title" json:"title" validate:"required,lte=255"`
	Author        string       `gorm:"column:author" json:"author" validate:"required,lte=255"`
	ProductStatus int          `gorm:"column:product_status" json:"product_status" validate:"required,len=1"`
	ProductAttrs  ProductAttrs `gorm:"column:product_attrs;type:json" json:"product_attrs"`
	BaseDbTime
}

type ProductRepo struct {
	Curd[Product]
}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{}
}

// ProductAttrs struct to describe product attributes.
type ProductAttrs struct {
	Picture     string `json:"picture"`
	Description string `json:"description"`
	Rating      int    `json:"rating" validate:"-"`
}

// Value make the ProductAttrs struct implement the driver.Valuer interface.
// This method simply returns the JSON-encoded representation of the struct.
func (b ProductAttrs) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Scan make the ProductAttrs struct implement the sql.Scanner interface.
// This method simply decodes a JSON-encoded value into the struct fields.
func (b *ProductAttrs) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(j, &b)
}
