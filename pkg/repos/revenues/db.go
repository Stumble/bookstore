// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.1.0-wicked-fork

package revenues

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
type BeforeDump func(m *ByBookRevenue)

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
CREATE MATERIALIZED VIEW IF NOT EXISTS by_book_revenues AS
  SELECT
    books.id,
    books.name,
    books.category,
    books.price,
    books.created_at,
    sum(orders.price) AS total,
    sum(
      CASE WHEN
        (orders.created_at > now() - interval '30 day')
      THEN orders.price ELSE 0 END
    ) AS last30d
  FROM
    books
    LEFT JOIN orders ON books.id = orders.book_id
  GROUP BY
      books.id;

CREATE UNIQUE INDEX IF NOT EXISTS v_books_id_unique_idx
  ON v_books (id);

CREATE UNIQUE INDEX IF NOT EXISTS v_books_total_volume_idx
  ON v_books (total);

CREATE UNIQUE INDEX IF NOT EXISTS v_books_total_last_30d_volume_idx
  ON v_books (last30d);
`
