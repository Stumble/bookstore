// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.1.4-wicked-fork
// source: query.sql

package users

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
)

const complicated = `-- name: Complicated :one
WITH RECURSIVE fibonacci(n,x,y) AS (
	SELECT
    	1 AS n ,
  		0 :: int AS x,
    	1 :: int AS y
  	UNION ALL
  	SELECT
    	n + 1 AS n,
  		y AS x,
    	x + y AS y
  	FROM fibonacci
  	WHERE n < $1::int
	)
SELECT
	x
FROM fibonacci
`

// -- cache : 1m
// example of sqlc cannot handle recursive query.
func (q *Queries) Complicated(ctx context.Context, n int32) (*int32, error) {
	dbRead := func() (any, time.Duration, error) {
		cacheDuration := time.Duration(time.Millisecond * 60000)
		row := q.db.WQueryRow(ctx, "users.Complicated", complicated, n)
		var x *int32 = new(int32)
		err := row.Scan(x)
		if err == pgx.ErrNoRows {
			return (*int32)(nil), cacheDuration, nil
		}
		return x, cacheDuration, err
	}
	if q.cache == nil {
		x, _, err := dbRead()
		return x.(*int32), err
	}

	var x *int32
	err := q.cache.GetWithTtl(ctx, "users:Complicated:"+hashIfLong(fmt.Sprintf("%+v", n)), &x, dbRead, false, false)
	if err != nil {
		return nil, err
	}

	return x, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO Users (
  name, metadata, image
) VALUES (
  $1, $2, $3
)
RETURNING id, name, metadata, image, created_at
`

type CreateUserParams struct {
	Name     string
	Metadata []byte
	Image    string
}

// -- invalidate : [GetUserByID, GetUserByName]
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams, getUserByID *int32, getUserByName *string) (*User, error) {
	row := q.db.WQueryRow(ctx, "users.CreateUser", createUser, arg.Name, arg.Metadata, arg.Image)
	var i *User = new(User)
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Metadata,
		&i.Image,
		&i.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return (*User)(nil), nil
	} else if err != nil {
		return nil, err
	}

	// invalidate
	_ = q.db.PostExec(func() error {
		var anyErr error
		if getUserByID != nil {
			key := "users:GetUserByID:" + hashIfLong(fmt.Sprintf("%+v", (*getUserByID)))
			err = q.cache.Invalidate(ctx, key)
			if err != nil {
				log.Error().Err(err).Msgf(
					"Failed to invalidate: %s", key)
				anyErr = err
			}
		}
		if getUserByName != nil {
			key := "users:GetUserByName:" + hashIfLong(fmt.Sprintf("%+v", (*getUserByName)))
			err = q.cache.Invalidate(ctx, key)
			if err != nil {
				log.Error().Err(err).Msgf(
					"Failed to invalidate: %s", key)
				anyErr = err
			}
		}
		return anyErr
	})
	return i, err
}

const deleteBadUsers = `-- name: DeleteBadUsers :execresult
DELETE FROM Users
WHERE NAME LIKE $1
`

// -- invalidate : [GetUserByID]
func (q *Queries) DeleteBadUsers(ctx context.Context, name string, getUserByID *int32) (pgconn.CommandTag, error) {
	rv, err := q.db.WExec(ctx, "users.DeleteBadUsers", deleteBadUsers, name)
	if err != nil {
		return rv, err
	}
	// invalidate
	_ = q.db.PostExec(func() error {
		var anyErr error
		if getUserByID != nil {
			key := "users:GetUserByID:" + hashIfLong(fmt.Sprintf("%+v", (*getUserByID)))
			err = q.cache.Invalidate(ctx, key)
			if err != nil {
				log.Error().Err(err).Msgf(
					"Failed to invalidate: %s", key)
				anyErr = err
			}
		}
		return anyErr
	})
	return rv, nil
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM Users
WHERE id = $1
`

// -- invalidate : [GetUserByID, GetUserByName]
func (q *Queries) DeleteUser(ctx context.Context, id int32, getUserByID *int32, getUserByName *string) error {
	_, err := q.db.WExec(ctx, "users.DeleteUser", deleteUser, id)
	if err != nil {
		return err
	}
	// invalidate
	_ = q.db.PostExec(func() error {
		var anyErr error
		if getUserByID != nil {
			key := "users:GetUserByID:" + hashIfLong(fmt.Sprintf("%+v", (*getUserByID)))
			err = q.cache.Invalidate(ctx, key)
			if err != nil {
				log.Error().Err(err).Msgf(
					"Failed to invalidate: %s", key)
				anyErr = err
			}
		}
		if getUserByName != nil {
			key := "users:GetUserByName:" + hashIfLong(fmt.Sprintf("%+v", (*getUserByName)))
			err = q.cache.Invalidate(ctx, key)
			if err != nil {
				log.Error().Err(err).Msgf(
					"Failed to invalidate: %s", key)
				anyErr = err
			}
		}
		return anyErr
	})
	return nil
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, metadata, image, created_at FROM users
WHERE id = $1 LIMIT 1
`

// -- cache : 30s
func (q *Queries) GetUserByID(ctx context.Context, id int32) (*User, error) {
	dbRead := func() (any, time.Duration, error) {
		cacheDuration := time.Duration(time.Millisecond * 30000)
		row := q.db.WQueryRow(ctx, "users.GetUserByID", getUserByID, id)
		var i *User = new(User)
		err := row.Scan(
			&i.ID,
			&i.Name,
			&i.Metadata,
			&i.Image,
			&i.CreatedAt,
		)
		if err == pgx.ErrNoRows {
			return (*User)(nil), cacheDuration, nil
		}
		return i, cacheDuration, err
	}
	if q.cache == nil {
		i, _, err := dbRead()
		return i.(*User), err
	}

	var i *User
	err := q.cache.GetWithTtl(ctx, "users:GetUserByID:"+hashIfLong(fmt.Sprintf("%+v", id)), &i, dbRead, false, false)
	if err != nil {
		return nil, err
	}

	return i, err
}

const getUserByName = `-- name: GetUserByName :one
SELECT id, name, metadata, image, created_at FROM users
WHERE Name = $1 LIMIT 1
`

// -- cache : 5m
func (q *Queries) GetUserByName(ctx context.Context, name string) (*User, error) {
	dbRead := func() (any, time.Duration, error) {
		cacheDuration := time.Duration(time.Millisecond * 300000)
		row := q.db.WQueryRow(ctx, "users.GetUserByName", getUserByName, name)
		var i *User = new(User)
		err := row.Scan(
			&i.ID,
			&i.Name,
			&i.Metadata,
			&i.Image,
			&i.CreatedAt,
		)
		if err == pgx.ErrNoRows {
			return (*User)(nil), cacheDuration, nil
		}
		return i, cacheDuration, err
	}
	if q.cache == nil {
		i, _, err := dbRead()
		return i.(*User), err
	}

	var i *User
	err := q.cache.GetWithTtl(ctx, "users:GetUserByName:"+hashIfLong(fmt.Sprintf("%+v", name)), &i, dbRead, false, false)
	if err != nil {
		return nil, err
	}

	return i, err
}

const incorrectQuery = `-- name: IncorrectQuery :one
SELECT id, name, metadata, image, created_at FROM users
WHERE Name = $1 LIMIT 1
`

// -- cache : 5m
func (q *Queries) IncorrectQuery(ctx context.Context, namePointer *string) (*User, error) {
	dbRead := func() (any, time.Duration, error) {
		cacheDuration := time.Duration(time.Millisecond * 300000)
		row := q.db.WQueryRow(ctx, "users.IncorrectQuery", incorrectQuery, namePointer)
		var i *User = new(User)
		err := row.Scan(
			&i.ID,
			&i.Name,
			&i.Metadata,
			&i.Image,
			&i.CreatedAt,
		)
		if err == pgx.ErrNoRows {
			return (*User)(nil), cacheDuration, nil
		}
		return i, cacheDuration, err
	}
	if q.cache == nil {
		i, _, err := dbRead()
		return i.(*User), err
	}

	var i *User
	err := q.cache.GetWithTtl(ctx, "users:IncorrectQuery:"+hashIfLong(fmt.Sprintf("%+v", ptrStr(namePointer))), &i, dbRead, false, false)
	if err != nil {
		return nil, err
	}

	return i, err
}

const listUserNames = `-- name: ListUserNames :many
SELECT id, name FROM users
WHERE id > $1
ORDER BY id
LIMIT $2
`

type ListUserNamesParams struct {
	After int32
	First int32
}

type ListUserNamesRow struct {
	ID   int32
	Name string
}

func (q *Queries) ListUserNames(ctx context.Context, arg ListUserNamesParams) ([]ListUserNamesRow, error) {
	rows, err := q.db.WQuery(ctx, "users.ListUserNames", listUserNames, arg.After, arg.First)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUserNamesRow
	for rows.Next() {
		var i *ListUserNamesRow = new(ListUserNamesRow)
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, *i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, name, metadata, image, created_at FROM users
WHERE id > $1
ORDER BY id
LIMIT $2
`

type ListUsersParams struct {
	After int32
	First int32
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.WQuery(ctx, "users.ListUsers", listUsers, arg.After, arg.First)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i *User = new(User)
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Metadata,
			&i.Image,
			&i.CreatedAt,
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

const updateMetaByID = `-- name: UpdateMetaByID :execrows
UPDATE users
SET
  metadata = $1
WHERE
  id = $2
`

type UpdateMetaByIDParams struct {
	Metadata []byte
	ID       int32
}

func (q *Queries) UpdateMetaByID(ctx context.Context, arg UpdateMetaByIDParams) (int64, error) {
	result, err := q.db.WExec(ctx, "users.UpdateMetaByID", updateMetaByID, arg.Metadata, arg.ID)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

const updateNameByID = `-- name: UpdateNameByID :one
UPDATE users
SET
  name = $1
WHERE
  id = $2
RETURNING ID
`

type UpdateNameByIDParams struct {
	Name string
	ID   int32
}

func (q *Queries) UpdateNameByID(ctx context.Context, arg UpdateNameByIDParams) (*int32, error) {
	row := q.db.WQueryRow(ctx, "users.UpdateNameByID", updateNameByID, arg.Name, arg.ID)
	var id *int32 = new(int32)
	err := row.Scan(id)
	if err == pgx.ErrNoRows {
		return (*int32)(nil), nil
	} else if err != nil {
		return nil, err
	}

	return id, err
}

const updateUserGrade = `-- name: UpdateUserGrade :execrows
UPDATE users
  SET metadata = jsonb_set(metadata, '{grade}', $1::text, true)
WHERE
  Name LIKE $2
`

type UpdateUserGradeParams struct {
	Grade string
	Name  string
}

// -- invalidate : [GetUserByID]
func (q *Queries) UpdateUserGrade(ctx context.Context, arg UpdateUserGradeParams, getUserByID *int32) (int64, error) {
	result, err := q.db.WExec(ctx, "users.UpdateUserGrade", updateUserGrade, arg.Grade, arg.Name)
	if err != nil {
		return 0, err
	}
	// invalidate
	_ = q.db.PostExec(func() error {
		var anyErr error
		if getUserByID != nil {
			key := "users:GetUserByID:" + hashIfLong(fmt.Sprintf("%+v", (*getUserByID)))
			err = q.cache.Invalidate(ctx, key)
			if err != nil {
				log.Error().Err(err).Msgf(
					"Failed to invalidate: %s", key)
				anyErr = err
			}
		}
		return anyErr
	})
	return result.RowsAffected(), nil
}

const upsertUsers = `-- name: UpsertUsers :exec
insert into users
  (name, metadata, image)
select
        unnest($1::VARCHAR(255)[]),
        unnest($2::JSON[]),
        unnest($3::TEXT[])
on conflict ON CONSTRAINT users_lower_name_key do
update set
    metadata = excluded.metadata,
    image = excluded.image
`

type UpsertUsersParams struct {
	Name     []string
	Metadata [][]byte
	Image    []string
}

func (q *Queries) UpsertUsers(ctx context.Context, arg UpsertUsersParams) error {
	_, err := q.db.WExec(ctx, "users.UpsertUsers", upsertUsers, arg.Name, arg.Metadata, arg.Image)
	if err != nil {
		return err
	}

	return nil
}

//// auto generated functions

func (q *Queries) Dump(ctx context.Context, beforeDump ...BeforeDump) ([]byte, error) {
	sql := "SELECT id,name,metadata,image,created_at FROM users ORDER BY id,name,image,created_at ASC;"
	rows, err := q.db.WQuery(ctx, "users.Dump", sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var v User
		if err := rows.Scan(&v.ID, &v.Name, &v.Metadata, &v.Image, &v.CreatedAt); err != nil {
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
	sql := "INSERT INTO users (id,name,metadata,image,created_at) VALUES ($1,$2,$3,$4,$5);"
	rows := make([]User, 0)
	err := json.Unmarshal(data, &rows)
	if err != nil {
		return err
	}
	for _, row := range rows {
		_, err := q.db.WExec(ctx, "users.Load", sql, row.ID, row.Name, row.Metadata, row.Image, row.CreatedAt)
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
