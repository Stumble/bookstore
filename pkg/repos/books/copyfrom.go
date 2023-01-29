// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.0.0-wicked-fork
// source: copyfrom.go

package books

import (
	"context"
)

// iteratorForBulkInsert implements pgx.CopyFromSource.
type iteratorForBulkInsert struct {
	rows                 []BulkInsertParams
	skippedFirstNextCall bool
}

func (r *iteratorForBulkInsert) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForBulkInsert) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].Name,
		r.rows[0].Description,
		r.rows[0].Metadata,
		r.rows[0].Category,
		r.rows[0].Price,
	}, nil
}

func (r iteratorForBulkInsert) Err() error {
	return nil
}

func (q *Queries) BulkInsert(ctx context.Context, arg []BulkInsertParams) (int64, error) {
	return q.db.WCopyFrom(ctx, "BulkInsert", []string{"books"}, []string{"name", "description", "metadata", "category", "price"}, &iteratorForBulkInsert{rows: arg})
}
