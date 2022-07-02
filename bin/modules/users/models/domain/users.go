package models

import (
	"gopkg.in/guregu/null.v4"
)

type Users struct {
	ID         string    `db:"id" json:"id"`
	Username   string    `db:"username" json:"username"`
	Email      string    `db:"email" json:"email"`
	Password   string    `db:"password" json:"password,omitempty"`
	CreatedAt  null.Time `db:"created_at" json:"created_at"`
	CreatedBy  string    `db:"created_by" json:"created_by"`
	ModifiedAt null.Time `db:"modified_at" json:"modified_at"`
	ModifiedBy *string   `db:"modified_by" json:"modified_by"`
}

type UsersList []Users

type ReqLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ReqUser struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	UserId   string `json:"userid"`
}
