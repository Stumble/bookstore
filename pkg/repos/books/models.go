// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0-65-ge5c2ed73-dirty-wicked-fork

package books

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Category string

const (
	CategoryComputerScience Category = "computer_science"
	CategoryPhilosophy      Category = "philosophy"
	CategoryComic           Category = "comic"
)

func (e *Category) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Category(s)
	case string:
		*e = Category(s)
	default:
		return fmt.Errorf("unsupported scan type for Category: %T", src)
	}
	return nil
}

type NullCategory struct {
	Category Category
	Valid    bool // Valid is true if Category is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCategory) Scan(value interface{}) error {
	if value == nil {
		ns.Category, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Category.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCategory) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Category), nil
}

func (e Category) Valid() bool {
	switch e {
	case CategoryComputerScience,
		CategoryPhilosophy,
		CategoryComic:
		return true
	}
	return false
}

func AllCategoryValues() []Category {
	return []Category{
		CategoryComputerScience,
		CategoryPhilosophy,
		CategoryComic,
	}
}

type Book struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Metadata    []byte         `json:"metadata"`
	Category    interface{}    `json:"category"`
	Price       pgtype.Numeric `json:"price"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
