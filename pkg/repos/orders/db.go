// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.2.1-wicked-fork

package orders

import (
	"github.com/stumble/dcache"
	"github.com/stumble/wpgx"
)

// BeforeDump allows you to edit result before dump.
type BeforeDump func(m *Order)

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
CREATE TABLE IF NOT EXISTS orders (
   id         INT         GENERATED ALWAYS AS IDENTITY,
   user_id    INT         references Users(ID) ON DELETE SET NULL,
   book_id    INT         references Items(ID) ON DELETE SET NULL,
   price      BIGINT      NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   is_deleted BOOLEAN     NOT NULL,
   CONSTRAINT orders_id_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS orders_item_id_idx ON orders (ItemID);
`
