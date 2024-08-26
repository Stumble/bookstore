// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.2.2-wicked-fork

package revenues

import (
	"github.com/stumble/dcache"
	"github.com/stumble/wpgx"
)

// BeforeDump allows you to edit result before dump.
type BeforeDump func(m *ByBookRevenue)

func New(db wpgx.WGConn, cache *dcache.DCache) *Queries {
	return &Queries{db: db, cache: cache}
}

type Queries struct {
	db    wpgx.WGConn
	cache *dcache.DCache
}

func (q *Queries) WithTx(tx *wpgx.WTx) *Queries {
	return &Queries{
		db:    tx,
		cache: q.cache,
	}
}

func (q *Queries) WithCache(cache *dcache.DCache) *Queries {
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
