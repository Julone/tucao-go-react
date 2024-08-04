package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
)

// LogRecord struct to describe product object.
type LogRecord struct {
	ID              uuid.UUID      `gorm:"column:id;type:bigint;not null;primaryKey;auto_increment" json:"id" `
	UserID          string         `gorm:"column:user_id" json:"user_id" `
	Title           string         `gorm:"column:title" json:"title" validate:"required,lte=255"`
	Author          string         `gorm:"column:author" json:"author" validae:"required,lte=255"`
	LogRecordStatus int            `gorm:"column:log_status" json:"log_status" validate:"required,len=1"`
	LogRecordAttrs  LogRecordAttrs `gorm:"column:log_attrs;type:json" json:"log_attrs"`
	BaseDbTime
}

type LogRecordRepo struct {
	Curd[LogRecord]
}

func NewLogRecordRepo() *LogRecordRepo {
	return &LogRecordRepo{}
}

// LogRecordAttrs struct to describe product attributes.
type LogRecordAttrs struct {
	Picture     string `json:"picture"`
	Description string `json:"description"`
	Rating      int    `json:"rating" validate:"-"`
}

// Value make the LogRecordAttrs struct implement the driver.Valuer interface.
// This method simply returns the JSON-encoded representation of the struct.
func (b LogRecordAttrs) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Scan make the LogRecordAttrs struct implement the sql.Scanner interface.
// This method simply decodes a JSON-encoded value into the struct fields.
func (b *LogRecordAttrs) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(j, &b)
}
