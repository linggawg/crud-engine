package models

import (
	"engine/bin/pkg/token"

	"gopkg.in/guregu/null.v4"
)

type Users struct {
	ID         string    `db:"id" json:"id"`
	RoleID     string    `db:"role_id" json:"role_id"`
	Username   string    `db:"username" json:"username"`
	Password   string    `db:"password" json:"password,omitempty"`
	CreatedAt  null.Time `db:"created_at" json:"created_at"`
	CreatedBy  *string   `db:"created_by" json:"created_by"`
	ModifiedAt null.Time `db:"modified_at" json:"modified_at"`
	ModifiedBy *string   `db:"modified_by" json:"modified_by"`
}

type UsersList []Users

type ReqLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Duration int    `json:"duration" validate:"required"`
}

type ReqUser struct {
	Username string      `json:"username" validate:"required"`
	Password string      `json:"password" validate:"required"`
	Opts     token.Claim `json:"opts"`
	RoleID   string      `json:"role_id" validate:"required"`
}
