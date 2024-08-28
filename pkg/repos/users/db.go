// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.3.0-wicked-fork

package users

import (
	"github.com/stumble/dcache"
	"github.com/stumble/wpgx"
)

// BeforeDump allows you to edit result before dump.
type BeforeDump func(m *User)

type CacheQuerierConn interface {
	GetCache() *dcache.DCache
	GetConn() wpgx.WQuerier
}

type CacheWGConn interface {
	GetCache() *dcache.DCache
	GetConn() wpgx.WGConn
}

func New(db wpgx.WGConn, cache *dcache.DCache) *Queries {
	return &Queries{db: db, cache: cache}
}

type Queries struct {
	db    wpgx.WGConn
	cache *dcache.DCache
}

var _ CacheWGConn = (*Queries)(nil)

func (q *Queries) GetCache() *dcache.DCache {
	return q.cache
}

func (q *Queries) GetConn() wpgx.WGConn {
	return q.db
}

func (q *Queries) AsReadOnly() *ReadOnlyQueries {
	return &ReadOnlyQueries{
		db:    q.db,
		cache: q.cache,
	}
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

func (q *Queries) UseReplica(replicaQuerier wpgx.WQuerier) *ReadOnlyQueries {
	return &ReadOnlyQueries{
		db:    replicaQuerier,
		cache: q.cache,
	}
}

type ReadOnlyQueries struct {
	db    wpgx.WQuerier
	cache *dcache.DCache
}

var _ CacheQuerierConn = (*ReadOnlyQueries)(nil)

func (q *ReadOnlyQueries) WithCache(cache *dcache.DCache) *ReadOnlyQueries {
	return &ReadOnlyQueries{
		db:    q.db,
		cache: cache,
	}
}

func (q *ReadOnlyQueries) GetCache() *dcache.DCache {
	return q.cache
}

func (q *ReadOnlyQueries) GetConn() wpgx.WQuerier {
	return q.db
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
