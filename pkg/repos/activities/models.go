// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.2.2-wicked-fork

package activities

import (
	"time"
)

type Activity struct {
	ID        int32     `json:"id"`
	Action    string    `json:"action"`
	Parameter *string   `json:"parameter"`
	CreatedAt time.Time `json:"created_at"`
}
