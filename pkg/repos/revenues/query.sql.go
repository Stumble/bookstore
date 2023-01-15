// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0-67-g8e483d59-wicked-fork
// source: query.sql

package revenues

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

const getTopItems = `-- name: GetTopItems :many
select id, name, category, price, created_at, total, last30d from by_book_revenues
order by
  total
limit 3
`

func (q *Queries) GetTopItems(ctx context.Context) ([]ByBookRevenue, error) {
	rows, err := q.db.WQuery(ctx, "GetTopItems", getTopItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ByBookRevenue
	for rows.Next() {
		var i *ByBookRevenue = new(ByBookRevenue)
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Category,
			&i.Price,
			&i.CreatedAt,
			&i.Total,
			&i.Last30d,
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

const refresh = `-- name: Refresh :exec
REFRESH MATERIALIZED VIEW CONCURRENTLY by_book_revenues
`

func (q *Queries) Refresh(ctx context.Context) error {
	_, err := q.db.WExec(ctx, "Refresh", refresh)
	if err != nil {
		return err
	}

	return nil
}

//// auto generated functions

func (q *Queries) Dump(ctx context.Context, beforeDump ...BeforeDump) ([]byte, error) {
	sql := "SELECT id,name,category,price,created_at,total,last30d FROM by_book_revenues ORDER BY id,name,category,price,created_at,total,last30d ASC;"
	rows, err := q.db.WQuery(ctx, "Dump", sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ByBookRevenue
	for rows.Next() {
		var v ByBookRevenue
		if err := rows.Scan(&v.ID, &v.Name, &v.Category, &v.Price, &v.CreatedAt, &v.Total, &v.Last30d); err != nil {
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
	sql := "INSERT INTO by_book_revenues (id,name,category,price,created_at,total,last30d) VALUES ($1,$2,$3,$4,$5,$6,$7);"
	rows := make([]ByBookRevenue, 0)
	err := json.Unmarshal(data, &rows)
	if err != nil {
		return err
	}
	for _, row := range rows {
		_, err := q.db.WExec(ctx, "Load", sql, row.ID, row.Name, row.Category, row.Price, row.CreatedAt, row.Total, row.Last30d)
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
