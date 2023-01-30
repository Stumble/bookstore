// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.0.1-wicked-fork

package books

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stumble/wpgx"
)

type WGConn interface {
	WQuery(
		ctx context.Context, name string, unprepared string, args ...interface{}) (pgx.Rows, error)
	WQueryRow(
		ctx context.Context, name string, unprepared string, args ...interface{}) pgx.Row
	WExec(
		ctx context.Context, name string, unprepared string, args ...interface{}) (pgconn.CommandTag, error)

	PostExec(f wpgx.PostExecFunc) error
}

type ReadWithTtlFunc = func() (any, time.Duration, error)

// BeforeDump allows you to edit result before dump.
type BeforeDump func(m *Book)

type Cache interface {
	GetWithTtl(
		ctx context.Context, key string, target any,
		readWithTtl ReadWithTtlFunc, noCache bool, noStore bool) error
	Set(ctx context.Context, key string, val any, ttl time.Duration) error
	Invalidate(ctx context.Context, key string) error
}

func New(db WGConn, cache Cache) *Queries {
	return &Queries{db: db, cache: cache}
}

type Queries struct {
	db    WGConn
	cache Cache
}

func (q *Queries) WithTx(tx *wpgx.WTx) *Queries {
	return &Queries{
		db:    tx,
		cache: q.cache,
	}
}

func (q *Queries) WithCache(cache Cache) *Queries {
	return &Queries{
		db:    q.db,
		cache: cache,
	}
}

var Schema = `
CREATE TYPE book_category AS ENUM (
    'computer_science',
    'philosophy',
    'comic'
);

CREATE TABLE IF NOT EXISTS books (
   id            SERIAL              NOT NULL,
   name          VARCHAR(255)        NOT NULL,
   description   VARCHAR(255)        NOT NULL,
   metadata      JSON,
   category      book_category       NOT NULL,
   price         REAL                NOT NULL,
   dummy_field   INT,
   created_at    TIMESTAMP           NOT NULL DEFAULT NOW(),
   updated_at    TIMESTAMP           NOT NULL DEFAULT NOW(),
   CONSTRAINT books_id_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS books_name_idx ON books (name);
CREATE INDEX IF NOT EXISTS books_category_id_idx ON books (category, id);
`