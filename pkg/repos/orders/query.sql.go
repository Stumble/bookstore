// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.0.0-wicked-fork
// source: query.sql

package orders

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

const bulkUpdate = `-- name: BulkUpdate :exec
UPDATE orders
SET
  price=temp.price,
  book_id=temp.book_id
FROM
  (
    SELECT
      UNNEST($1::int[]) as id,
      UNNEST($2::bigint[]) as price,
      UNNEST($3::int[]) as book_id
  ) AS temp
WHERE
  orders.id=temp.id
`

type BulkUpdateParams struct {
	ID     []int32
	Price  []int64
	BookID []int32
}

func (q *Queries) BulkUpdate(ctx context.Context, arg BulkUpdateParams) error {
	_, err := q.db.WExec(ctx, "BulkUpdate", bulkUpdate, arg.ID, arg.Price, arg.BookID)
	if err != nil {
		return err
	}

	return nil
}

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO orders (
  user_id, book_id, price, is_deleted
) VALUES (
  $1, $2, $3, FALSE
)
RETURNING id, user_id, book_id, price, created_at, is_deleted
`

type CreateAuthorParams struct {
	UserID *int32
	BookID *int32
	Price  int64
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (*Order, error) {
	row := q.db.WQueryRow(ctx, "CreateAuthor", createAuthor, arg.UserID, arg.BookID, arg.Price)
	var i *Order = new(Order)
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.BookID,
		&i.Price,
		&i.CreatedAt,
		&i.IsDeleted,
	)
	if err == pgx.ErrNoRows {
		return (*Order)(nil), nil
	} else if err != nil {
		return nil, err
	}

	return i, err
}

const deleteOrder = `-- name: DeleteOrder :exec
UPDATE orders
SET
  is_deleted = TRUE
WHERE
  id = $1
`

func (q *Queries) DeleteOrder(ctx context.Context, id int32) error {
	_, err := q.db.WExec(ctx, "DeleteOrder", deleteOrder, id)
	if err != nil {
		return err
	}

	return nil
}

const getOrderByID = `-- name: GetOrderByID :one
SELECT
  orders.ID,
  orders.user_id,
  orders.book_id,
  orders.created_at,
  users.name AS user_name,
  users.image AS user_thumbnail,
  books.name AS book_name,
  books.price As book_price,
  books.metadata As book_metadata
FROM
  orders
  INNER JOIN books ON orders.book_id = books.id
  INNER JOIN users ON orders.user_id = users.id
WHERE
  orders.is_deleted = FALSE
`

type GetOrderByIDRow struct {
	ID            int32
	UserID        *int32
	BookID        *int32
	CreatedAt     time.Time
	UserName      string
	UserThumbnail string
	BookName      string
	BookPrice     pgtype.Numeric
	BookMetadata  []byte
}

// -- cache : 10m
func (q *Queries) GetOrderByID(ctx context.Context) (*GetOrderByIDRow, error) {
	dbRead := func() (any, time.Duration, error) {
		cacheDuration := time.Duration(time.Millisecond * 600000)
		row := q.db.WQueryRow(ctx, "GetOrderByID", getOrderByID)
		var i *GetOrderByIDRow = new(GetOrderByIDRow)
		err := row.Scan(
			&i.ID,
			&i.UserID,
			&i.BookID,
			&i.CreatedAt,
			&i.UserName,
			&i.UserThumbnail,
			&i.BookName,
			&i.BookPrice,
			&i.BookMetadata,
		)
		if err == pgx.ErrNoRows {
			return (*GetOrderByIDRow)(nil), cacheDuration, nil
		}
		return i, cacheDuration, err
	}
	if q.cache == nil {
		i, _, err := dbRead()
		return i.(*GetOrderByIDRow), err
	}

	var i *GetOrderByIDRow
	err := q.cache.GetWithTtl(ctx, "orders:GetOrderByID", &i, dbRead, false, false)
	if err != nil {
		return nil, err
	}

	return i, err
}

const listOrdersByGender = `-- name: ListOrdersByGender :many
WITH users_by_gender AS (
  SELECT id, name, metadata, image, created_at FROM users WHERE users.metadata->>'gender' = $3::text
)
SELECT id, user_id, book_id, price, created_at, is_deleted FROM orders
WHERE
  user_id IN (SELECT id FROM users_by_gender) AND orders.id > $1
LIMIT $2
`

type ListOrdersByGenderParams struct {
	After  int32
	First  int32
	Gender string
}

// CacheKey - cache key
func (arg ListOrdersByGenderParams) CacheKey() string {
	prefix := "orders:ListOrdersByGender:"
	return prefix + fmt.Sprintf("%+v,%+v,%+v", arg.After, arg.First, arg.Gender)
}

// -- cache : 1m
// This is just an example for using type annotation for JSON field and 'with clause'.
func (q *Queries) ListOrdersByGender(ctx context.Context, arg ListOrdersByGenderParams) ([]Order, error) {
	dbRead := func() (any, time.Duration, error) {
		cacheDuration := time.Duration(time.Millisecond * 60000)
		rows, err := q.db.WQuery(ctx, "ListOrdersByGender", listOrdersByGender, arg.After, arg.First, arg.Gender)
		if err != nil {
			return nil, 0, err
		}
		defer rows.Close()
		var items []Order
		for rows.Next() {
			var i *Order = new(Order)
			if err := rows.Scan(
				&i.ID,
				&i.UserID,
				&i.BookID,
				&i.Price,
				&i.CreatedAt,
				&i.IsDeleted,
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
		return items.([]Order), err
	}
	var items []Order
	err := q.cache.GetWithTtl(ctx, arg.CacheKey(), &items, dbRead, false, false)
	if err != nil {
		return nil, err
	}

	return items, err
}

const listOrdersByUser = `-- name: ListOrdersByUser :many
SELECT id, user_id, book_id, price, created_at, is_deleted FROM orders
WHERE
  user_id = $1 AND created_at < $2
ORDER BY created_at DESC
LIMIT $3
`

type ListOrdersByUserParams struct {
	UserID *int32
	After  time.Time
	First  int32
}

func (q *Queries) ListOrdersByUser(ctx context.Context, arg ListOrdersByUserParams) ([]Order, error) {
	rows, err := q.db.WQuery(ctx, "ListOrdersByUser", listOrdersByUser, arg.UserID, arg.After, arg.First)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Order
	for rows.Next() {
		var i *Order = new(Order)
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.BookID,
			&i.Price,
			&i.CreatedAt,
			&i.IsDeleted,
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

const listOrdersByUserAndBook = `-- name: ListOrdersByUserAndBook :many
SELECT id, user_id, book_id, price, created_at, is_deleted FROM orders
WHERE
  (user_id, book_id) IN (
  SELECT
    UNNEST($1::int[]),
    UNNEST($2::int[])
)
`

type ListOrdersByUserAndBookParams struct {
	UserID []int32
	BookID []int32
}

func (q *Queries) ListOrdersByUserAndBook(ctx context.Context, arg ListOrdersByUserAndBookParams) ([]Order, error) {
	rows, err := q.db.WQuery(ctx, "ListOrdersByUserAndBook", listOrdersByUserAndBook, arg.UserID, arg.BookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Order
	for rows.Next() {
		var i *Order = new(Order)
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.BookID,
			&i.Price,
			&i.CreatedAt,
			&i.IsDeleted,
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

//// auto generated functions

func (q *Queries) Dump(ctx context.Context, beforeDump ...BeforeDump) ([]byte, error) {
	sql := "SELECT id,user_id,book_id,price,created_at,is_deleted FROM orders ORDER BY id,user_id,book_id,price,created_at,is_deleted ASC;"
	rows, err := q.db.WQuery(ctx, "Dump", sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Order
	for rows.Next() {
		var v Order
		if err := rows.Scan(&v.ID, &v.UserID, &v.BookID, &v.Price, &v.CreatedAt, &v.IsDeleted); err != nil {
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
	sql := "INSERT INTO orders (id,user_id,book_id,price,created_at,is_deleted) VALUES ($1,$2,$3,$4,$5,$6);"
	rows := make([]Order, 0)
	err := json.Unmarshal(data, &rows)
	if err != nil {
		return err
	}
	for _, row := range rows {
		_, err := q.db.WExec(ctx, "Load", sql, row.ID, row.UserID, row.BookID, row.Price, row.CreatedAt, row.IsDeleted)
		if err != nil {
			return err
		}
	}
	return nil
}

// eliminate unused error
var _ = log.Logger
var _ = fmt.Sprintf("")
var _ = time.Now()
var _ = json.RawMessage{}
