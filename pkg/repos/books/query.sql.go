// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.1.11-dirty-wicked-fork
// source: query.sql

package books

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type BulkInsertByCopyfromParams struct {
	Name        string
	Description string
	Metadata    []byte
	Category    BookCategory
	Price       float32
}

const getAllBooks = `-- name: GetAllBooks :many
SELECT id, name, description, metadata, category, price, dummy_field, created_at, updated_at FROM books
`

// -- cache : 10m
func (q *Queries) GetAllBooks(ctx context.Context) ([]Book, error) {
	q.db.CountIntent("books.GetAllBooks")
	dbRead := func() (any, time.Duration, error) {
		cacheDuration := time.Duration(time.Millisecond * 600000)
		rows, err := q.db.WQuery(ctx, "books.GetAllBooks", getAllBooks)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()
		var items []Book
		for rows.Next() {
			var i *Book = new(Book)
			if err := rows.Scan(
				&i.ID,
				&i.Name,
				&i.Description,
				&i.Metadata,
				&i.Category,
				&i.Price,
				&i.DummyField,
				&i.CreatedAt,
				&i.UpdatedAt,
			); err != nil {
				return nil, 0, err
			}
			items = append(items, *i)
		}
		if err := rows.Err(); err != nil {
			return nil, 0, err
		}
		return items, cacheDuration, nil
	}
	if q.cache == nil {
		items, _, err := dbRead()
		return items.([]Book), err
	}
	var items []Book
	err := q.cache.GetWithTtl(ctx, "books:GetAllBooks:", &items, dbRead, false, false)
	if err != nil {
		return nil, err
	}

	return items, err
}

const getAllBooks2 = `-- name: GetAllBooks2 :many
SELECT id, name, description, metadata, category, price, dummy_field, created_at, updated_at FROM books
`

// -- cache : 10m
func (q *Queries) GetAllBooks2(ctx context.Context) ([]Book, error) {
	q.db.CountIntent("books.GetAllBooks2")
	dbRead := func() (any, time.Duration, error) {
		cacheDuration := time.Duration(time.Millisecond * 600000)
		rows, err := q.db.WQuery(ctx, "books.GetAllBooks2", getAllBooks2)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()
		var items []Book
		for rows.Next() {
			var i *Book = new(Book)
			if err := rows.Scan(
				&i.ID,
				&i.Name,
				&i.Description,
				&i.Metadata,
				&i.Category,
				&i.Price,
				&i.DummyField,
				&i.CreatedAt,
				&i.UpdatedAt,
			); err != nil {
				return nil, 0, err
			}
			items = append(items, *i)
		}
		if err := rows.Err(); err != nil {
			return nil, 0, err
		}
		return items, cacheDuration, nil
	}
	if q.cache == nil {
		items, _, err := dbRead()
		return items.([]Book), err
	}
	var items []Book
	err := q.cache.GetWithTtl(ctx, "books:GetAllBooks2:", &items, dbRead, false, false)
	if err != nil {
		return nil, err
	}

	return items, err
}

const getBookByID = `-- name: GetBookByID :one
SELECT id, name, description, metadata, category, price, dummy_field, created_at, updated_at FROM books WHERE id = $1
`

// -- cache : 10m
func (q *Queries) GetBookByID(ctx context.Context, id int32) (*Book, error) {
	q.db.CountIntent("books.GetBookByID")
	dbRead := func() (any, time.Duration, error) {
		cacheDuration := time.Duration(time.Millisecond * 600000)
		row := q.db.WQueryRow(ctx, "books.GetBookByID", getBookByID, id)
		var i *Book = new(Book)
		err := row.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Metadata,
			&i.Category,
			&i.Price,
			&i.DummyField,
			&i.CreatedAt,
			&i.UpdatedAt,
		)
		if err == pgx.ErrNoRows {
			return (*Book)(nil), cacheDuration, nil
		}
		return i, cacheDuration, err
	}
	if q.cache == nil {
		i, _, err := dbRead()
		return i.(*Book), err
	}

	var i *Book
	err := q.cache.GetWithTtl(ctx, "books:GetBookByID:"+hashIfLong(fmt.Sprintf("%+v", id)), &i, dbRead, false, false)
	if err != nil {
		return nil, err
	}

	return i, err
}

const getBookByNameMaybe = `-- name: GetBookByNameMaybe :many
SELECT id, name, description, metadata, category, price, dummy_field, created_at, updated_at FROM books WHERE
  name LIKE coalesce($1, name)
`

// -- cache : 10m
func (q *Queries) GetBookByNameMaybe(ctx context.Context, name *string) ([]Book, error) {
	q.db.CountIntent("books.GetBookByNameMaybe")
	dbRead := func() (any, time.Duration, error) {
		cacheDuration := time.Duration(time.Millisecond * 600000)
		rows, err := q.db.WQuery(ctx, "books.GetBookByNameMaybe", getBookByNameMaybe, name)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()
		var items []Book
		for rows.Next() {
			var i *Book = new(Book)
			if err := rows.Scan(
				&i.ID,
				&i.Name,
				&i.Description,
				&i.Metadata,
				&i.Category,
				&i.Price,
				&i.DummyField,
				&i.CreatedAt,
				&i.UpdatedAt,
			); err != nil {
				return nil, 0, err
			}
			items = append(items, *i)
		}
		if err := rows.Err(); err != nil {
			return nil, 0, err
		}
		return items, cacheDuration, nil
	}
	if q.cache == nil {
		items, _, err := dbRead()
		return items.([]Book), err
	}
	var items []Book
	err := q.cache.GetWithTtl(ctx, "books:GetBookByNameMaybe:"+hashIfLong(fmt.Sprintf("%+v", ptrStr(name))), &items, dbRead, false, false)
	if err != nil {
		return nil, err
	}

	return items, err
}

const getBookBySpec = `-- name: GetBookBySpec :many
SELECT id, name, description, metadata, category, price, dummy_field, created_at, updated_at FROM books WHERE
  name LIKE coalesce($1, name) AND
  price = coalesce($2, price) AND
  ($3::int is NULL or dummy_field = $3)
`

type GetBookBySpecParams struct {
	Name  *string
	Price *float32
	Dummy *int32
}

// CacheKey - cache key
func (arg GetBookBySpecParams) CacheKey() string {
	prefix := "books:GetBookBySpec:"
	return prefix + hashIfLong(fmt.Sprintf("%+v,%+v,%+v", ptrStr(arg.Name), ptrStr(arg.Price), ptrStr(arg.Dummy)))
}

// -- cache : 10m
func (q *Queries) GetBookBySpec(ctx context.Context, arg GetBookBySpecParams) ([]Book, error) {
	q.db.CountIntent("books.GetBookBySpec")
	dbRead := func() (any, time.Duration, error) {
		cacheDuration := time.Duration(time.Millisecond * 600000)
		rows, err := q.db.WQuery(ctx, "books.GetBookBySpec", getBookBySpec, arg.Name, arg.Price, arg.Dummy)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()
		var items []Book
		for rows.Next() {
			var i *Book = new(Book)
			if err := rows.Scan(
				&i.ID,
				&i.Name,
				&i.Description,
				&i.Metadata,
				&i.Category,
				&i.Price,
				&i.DummyField,
				&i.CreatedAt,
				&i.UpdatedAt,
			); err != nil {
				return nil, 0, err
			}
			items = append(items, *i)
		}
		if err := rows.Err(); err != nil {
			return nil, 0, err
		}
		return items, cacheDuration, nil
	}
	if q.cache == nil {
		items, _, err := dbRead()
		return items.([]Book), err
	}
	var items []Book
	err := q.cache.GetWithTtl(ctx, arg.CacheKey(), &items, dbRead, false, false)
	if err != nil {
		return nil, err
	}

	return items, err
}

const insert = `-- name: Insert :exec
INSERT INTO books (
   name, description, metadata, category, price
) VALUES (
  $1, $2, $3, $4, $5
)
`

type InsertParams struct {
	Name        string
	Description string
	Metadata    []byte
	Category    BookCategory
	Price       float32
}

// -- invalidate : [GetAllBooks, GetAllBooks2]
func (q *Queries) Insert(ctx context.Context, arg InsertParams) error {
	_, err := q.db.WExec(ctx, "books.Insert", insert,
		arg.Name,
		arg.Description,
		arg.Metadata,
		arg.Category,
		arg.Price,
	)
	if err != nil {
		return err
	}
	// invalidate
	_ = q.db.PostExec(func() error {
		var anyErr error
		{
			key := "books:GetAllBooks:"
			err = q.cache.Invalidate(ctx, key)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msgf(
					"Failed to invalidate: %s", key)
				anyErr = err
			}
		}
		{
			key := "books:GetAllBooks2:"
			err = q.cache.Invalidate(ctx, key)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msgf(
					"Failed to invalidate: %s", key)
				anyErr = err
			}
		}
		return anyErr
	})
	return nil
}

const insertAndReturnID = `-- name: InsertAndReturnID :one
INSERT INTO books (
   name, description, metadata, category, price
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING id
`

type InsertAndReturnIDParams struct {
	Name        string
	Description string
	Metadata    []byte
	Category    BookCategory
	Price       float32
}

func (q *Queries) InsertAndReturnID(ctx context.Context, arg InsertAndReturnIDParams) (*int32, error) {
	row := q.db.WQueryRow(ctx, "books.InsertAndReturnID", insertAndReturnID,
		arg.Name,
		arg.Description,
		arg.Metadata,
		arg.Category,
		arg.Price,
	)
	var id *int32 = new(int32)
	err := row.Scan(id)
	if err == pgx.ErrNoRows {
		return (*int32)(nil), nil
	} else if err != nil {
		return nil, err
	}

	return id, err
}

const insertWithInvalidate = `-- name: InsertWithInvalidate :exec
INSERT INTO books (
   id, name, description, metadata, category, dummy_field, price
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
`

type InsertWithInvalidateParams struct {
	ID          int32
	Name        string
	Description string
	Metadata    []byte
	Category    BookCategory
	DummyField  *int32
	Price       float32
}

// -- invalidate : [GetBookByNameMaybe, GetBookBySpec]
func (q *Queries) InsertWithInvalidate(ctx context.Context, arg InsertWithInvalidateParams, getBookByNameMaybe **string, getBookBySpec *GetBookBySpecParams) error {
	_, err := q.db.WExec(ctx, "books.InsertWithInvalidate", insertWithInvalidate,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.Metadata,
		arg.Category,
		arg.DummyField,
		arg.Price,
	)
	if err != nil {
		return err
	}
	// invalidate
	_ = q.db.PostExec(func() error {
		var anyErr error
		{
			if getBookByNameMaybe != nil {
				key := "books:GetBookByNameMaybe:" + hashIfLong(fmt.Sprintf("%+v", ptrStr((*getBookByNameMaybe))))
				err = q.cache.Invalidate(ctx, key)
				if err != nil {
					log.Ctx(ctx).Error().Err(err).Msgf(
						"Failed to invalidate: %s", key)
					anyErr = err
				}
			}
		}
		{
			if getBookBySpec != nil {
				key := (*getBookBySpec).CacheKey()
				err = q.cache.Invalidate(ctx, key)
				if err != nil {
					log.Ctx(ctx).Error().Err(err).Msgf(
						"Failed to invalidate: %s", key)
					anyErr = err
				}
			}
		}
		return anyErr
	})
	return nil
}

const listByCategory = `-- name: ListByCategory :many
SELECT id, name, description, metadata, category, price, dummy_field, created_at, updated_at
FROM
  books
WHERE
  category = $1 AND id > $2
ORDER BY
  id
LIMIT $3
`

type ListByCategoryParams struct {
	Category BookCategory
	After    int32
	First    int32
}

func (q *Queries) ListByCategory(ctx context.Context, arg ListByCategoryParams) ([]Book, error) {
	q.db.CountIntent("books.ListByCategory")
	rows, err := q.db.WQuery(ctx, "books.ListByCategory", listByCategory, arg.Category, arg.After, arg.First)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Book
	for rows.Next() {
		var i *Book = new(Book)
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Metadata,
			&i.Category,
			&i.Price,
			&i.DummyField,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, *i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, err
}

const partialUpdateByID = `-- name: PartialUpdateByID :exec
UPDATE books
SET
  description = coalesce($1, description),
  metadata = coalesce($2, metadata),
  price = coalesce($3, price),
  dummy_field = coalesce($4, dummy_field),
  updated_at = NOW()
WHERE
  id = $5
`

type PartialUpdateByIDParams struct {
	Description *string
	Meta        []byte
	Price       *float32
	DummyField  *int32
	ID          int32
}

func (q *Queries) PartialUpdateByID(ctx context.Context, arg PartialUpdateByIDParams) error {
	_, err := q.db.WExec(ctx, "books.PartialUpdateByID", partialUpdateByID,
		arg.Description,
		arg.Meta,
		arg.Price,
		arg.DummyField,
		arg.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

const refreshIDSerial = `-- name: RefreshIDSerial :exec
SELECT setval(seq_name, (SELECT MAX(id) FROM books)+1, false)
FROM PG_GET_SERIAL_SEQUENCE('books', 'id') as seq_name
`

func (q *Queries) RefreshIDSerial(ctx context.Context) error {
	_, err := q.db.WExec(ctx, "books.RefreshIDSerial", refreshIDSerial)
	if err != nil {
		return err
	}

	return nil
}

const searchBooks = `-- name: SearchBooks :many
SELECT id, name, description, metadata, category, price, dummy_field, created_at, updated_at FROM books WHERE name like $1
`

// -- cache : 10m
func (q *Queries) SearchBooks(ctx context.Context, name string) ([]Book, error) {
	q.db.CountIntent("books.SearchBooks")
	dbRead := func() (any, time.Duration, error) {
		cacheDuration := time.Duration(time.Millisecond * 600000)
		rows, err := q.db.WQuery(ctx, "books.SearchBooks", searchBooks, name)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()
		var items []Book
		for rows.Next() {
			var i *Book = new(Book)
			if err := rows.Scan(
				&i.ID,
				&i.Name,
				&i.Description,
				&i.Metadata,
				&i.Category,
				&i.Price,
				&i.DummyField,
				&i.CreatedAt,
				&i.UpdatedAt,
			); err != nil {
				return nil, 0, err
			}
			items = append(items, *i)
		}
		if err := rows.Err(); err != nil {
			return nil, 0, err
		}
		return items, cacheDuration, nil
	}
	if q.cache == nil {
		items, _, err := dbRead()
		return items.([]Book), err
	}
	var items []Book
	err := q.cache.GetWithTtl(ctx, "books:SearchBooks:"+hashIfLong(fmt.Sprintf("%+v", name)), &items, dbRead, false, false)
	if err != nil {
		return nil, err
	}

	return items, err
}

const updateBookByID = `-- name: UpdateBookByID :exec
UPDATE books
SET
  description = $1, metadata = $2, price = $3, updated_at = NOW()
WHERE
  id = $4
`

type UpdateBookByIDParams struct {
	Description string
	Meta        []byte
	Price       float32
	ID          int32
}

// -- invalidate : [GetBookByID]
func (q *Queries) UpdateBookByID(ctx context.Context, arg UpdateBookByIDParams, getBookByID *int32) error {
	_, err := q.db.WExec(ctx, "books.UpdateBookByID", updateBookByID,
		arg.Description,
		arg.Meta,
		arg.Price,
		arg.ID,
	)
	if err != nil {
		return err
	}
	// invalidate
	_ = q.db.PostExec(func() error {
		var anyErr error
		{
			if getBookByID != nil {
				key := "books:GetBookByID:" + hashIfLong(fmt.Sprintf("%+v", (*getBookByID)))
				err = q.cache.Invalidate(ctx, key)
				if err != nil {
					log.Ctx(ctx).Error().Err(err).Msgf(
						"Failed to invalidate: %s", key)
					anyErr = err
				}
			}
		}
		return anyErr
	})
	return nil
}

//// auto generated functions

func (q *Queries) Dump(ctx context.Context, beforeDump ...BeforeDump) ([]byte, error) {
	sql := "SELECT id,name,description,metadata,category,price,dummy_field,created_at,updated_at FROM \"books\" ORDER BY id,name,description,price,dummy_field,created_at,updated_at ASC;"
	rows, err := q.db.WQuery(ctx, "books.Dump", sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Book
	for rows.Next() {
		var v Book
		if err := rows.Scan(&v.ID, &v.Name, &v.Description, &v.Metadata, &v.Category, &v.Price, &v.DummyField, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		for _, applyBeforeDump := range beforeDump {
			applyBeforeDump(&v)
		}
		items = append(items, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	bytes, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (q *Queries) Load(ctx context.Context, data []byte) error {
	sql := "INSERT INTO \"books\" (id,name,description,metadata,category,price,dummy_field,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);"
	rows := make([]Book, 0)
	err := json.Unmarshal(data, &rows)
	if err != nil {
		return err
	}
	for _, row := range rows {
		_, err := q.db.WExec(ctx, "books.Load", sql, row.ID, row.Name, row.Description, row.Metadata, row.Category, row.Price, row.DummyField, row.CreatedAt, row.UpdatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func hashIfLong(v string) string {
	if len(v) > 64 {
		hash := sha256.Sum256([]byte(v))
		return "h(" + hex.EncodeToString(hash[:]) + ")"
	}
	return v
}

func ptrStr[T any](v *T) string {
	if v == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%+v", *v)
}

// eliminate unused error
var _ = log.Logger
var _ = fmt.Sprintf("")
var _ = time.Now()
var _ = json.RawMessage{}
var _ = sha256.Sum256(nil)
var _ = hex.EncodeToString(nil)
