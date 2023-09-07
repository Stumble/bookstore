// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.1.8-wicked-fork
// source: query.sql

package activities

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

const insert = `-- name: Insert :exec
INSERT INTO activities (
   action, parameter, created_at
) VALUES (
  $1, $2, NOW()
)
`

type InsertParams struct {
	Action    string
	Parameter *string
}

func (q *Queries) Insert(ctx context.Context, arg InsertParams) error {
	_, err := q.db.WExec(ctx, "activities.Insert", insert, arg.Action, arg.Parameter)
	if err != nil {
		return err
	}

	return nil
}

//// auto generated functions

func (q *Queries) Dump(ctx context.Context, beforeDump ...BeforeDump) ([]byte, error) {
	sql := "SELECT id,action,parameter,created_at FROM activities ORDER BY id,action,parameter,created_at ASC;"
	rows, err := q.db.WQuery(ctx, "activities.Dump", sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Activity
	for rows.Next() {
		var v Activity
		if err := rows.Scan(&v.ID, &v.Action, &v.Parameter, &v.CreatedAt); err != nil {
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
	sql := "INSERT INTO activities (id,action,parameter,created_at) VALUES ($1,$2,$3,$4);"
	rows := make([]Activity, 0)
	err := json.Unmarshal(data, &rows)
	if err != nil {
		return err
	}
	for _, row := range rows {
		_, err := q.db.WExec(ctx, "activities.Load", sql, row.ID, row.Action, row.Parameter, row.CreatedAt)
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
