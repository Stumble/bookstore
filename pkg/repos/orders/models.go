// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.3.0-wicked-fork

package orders

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type BookCategory string

const (
	BookCategoryComputerScience BookCategory = "computer_science"
	BookCategoryPhilosophy      BookCategory = "philosophy"
	BookCategoryComic           BookCategory = "comic"
)

func (e *BookCategory) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = BookCategory(s)
	case string:
		*e = BookCategory(s)
	default:
		return fmt.Errorf("unsupported scan type for BookCategory: %T", src)
	}
	return nil
}

type NullBookCategory struct {
	BookCategory BookCategory
	Valid        bool // Valid is true if BookCategory is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullBookCategory) Scan(value interface{}) error {
	if value == nil {
		ns.BookCategory, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.BookCategory.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullBookCategory) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.BookCategory), nil
}

func (e BookCategory) Valid() bool {
	switch e {
	case BookCategoryComputerScience,
		BookCategoryPhilosophy,
		BookCategoryComic:
		return true
	}
	return false
}

func AllBookCategoryValues() []BookCategory {
	return []BookCategory{
		BookCategoryComputerScience,
		BookCategoryPhilosophy,
		BookCategoryComic,
	}
}

type Order struct {
	ID        int32     `json:"id"`
	UserID    *int32    `json:"user_id"`
	BookID    *int32    `json:"book_id"`
	Price     int64     `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	IsDeleted bool      `json:"is_deleted"`
}
