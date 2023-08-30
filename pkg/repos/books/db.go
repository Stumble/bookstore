// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.1.7-wicked-fork

package books

import (
	"github.com/stumble/dcache"
	"github.com/stumble/wpgx"
)

// BeforeDump allows you to edit result before dump.
type BeforeDump func(m *Book)

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
   created_at    TIMESTAMPTZ           NOT NULL DEFAULT NOW(),
   updated_at    TIMESTAMPTZ           NOT NULL DEFAULT NOW(),
   CONSTRAINT books_id_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS books_name_idx ON books (name);
CREATE INDEX IF NOT EXISTS books_category_id_idx ON books (category, id);
`
