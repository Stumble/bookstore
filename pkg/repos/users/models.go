// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v2.1.7-wicked-fork

package users

import (
	"time"
)

type User struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	Metadata  []byte    `json:"metadata"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}
