package models

import (
	"gopkg.in/guregu/null.v4"
)

// User struct
type User struct {
	ID         string    `db:"id" json:"id"`
	Username   string    `db:"username" json:"username"`
	Email      string    `db:"email" json:"email"`
	Password   string    `db:"password" json:"password,omitempty"`
	CreatedAt  null.Time `db:"created_at" json:"created_at"`
	CreatedBy  string    `db:"created_by" json:"created_by"`
	ModifiedAt null.Time `db:"modified_at" json:"modified_at"`
	ModifiedBy *string   `db:"modified_by" json:"modified_by"`
}

type Apps struct {
	ID         string    `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	CreatedAt  null.Time `db:"created_at" json:"created_at"`
	CreatedBy  string    `db:"created_by" json:"created_by"`
	ModifiedAt null.Time `db:"modified_at" json:"modified_at"`
	ModifiedBy string    `db:"modified_by" json:"modified_by"`
}
