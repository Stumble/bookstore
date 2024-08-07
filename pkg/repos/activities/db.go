// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.2.1-wicked-fork

package activities

import (
	"github.com/stumble/dcache"
	"github.com/stumble/wpgx"
)

// BeforeDump allows you to edit result before dump.
type BeforeDump func(m *Activity)

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
CREATE TABLE IF NOT EXISTS activities (
   id            INT                 GENERATED ALWAYS AS IDENTITY,
   action        VARCHAR(255)        NOT NULL,
   parameter     TEXT,
   created_at    TIMESTAMPTZ           NOT NULL DEFAULT NOW(),
   CONSTRAINT activities_id_pkey PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS activities_action_idx ON activities (action);
`
