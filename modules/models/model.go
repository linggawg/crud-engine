package models

import "gopkg.in/guregu/null.v4"

// User struct
type User struct {
	ID           int64     `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	Password     string    `db:"password" json:"password,omitempty"`
	CreatedAt    null.Time `db:"created_at" json:"created_at"`
	CreatedBy    string    `db:"created_by" json:"created_by"`
	UpdatedAt    null.Time `db:"updated_at" json:"updated_at"`
	LastUpdateBy *string   `db:"last_update_by" json:"last_update_by"`
	IsDeleted    bool      `db:"is_deleted" json:"is_deleted"`
}
