// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.2.1-wicked-fork

package users

import (
	"github.com/stumble/dcache"
	"github.com/stumble/wpgx"
)

// BeforeDump allows you to edit result before dump.
type BeforeDump func(m *User)

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
CREATE TABLE IF NOT EXISTS users (
   id          INT          GENERATED ALWAYS AS IDENTITY,
   name        VARCHAR(255) NOT NULL,
   metadata    JSON,
   image       TEXT         NOT NULL,
   created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
   CONSTRAINT users_id_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS users_created_at_idx
    ON Users (CreatedAt);

CREATE UNIQUE INDEX IF NOT EXISTS users_lower_name_key
    ON Users ((lower(Name))) INCLUDE (ID);
`
