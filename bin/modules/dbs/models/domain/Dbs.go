package models

import (
	"gopkg.in/guregu/null.v4"
)

type Dbs struct {
	ID         string    `db:"id" json:"id"`
	AppID      string    `db:"app_id" json:"app_id"`
	Name       string    `db:"name" json:"name"`
	Host       string    `db:"host" json:"host"`
	Port       int       `db:"port" json:"port"`
	Username   string    `db:"username" json:"username"`
	Password   *string   `db:"password" json:"password"`
	Dialect    string    `db:"dialect" json:"dialect"`
	CreatedAt  null.Time `db:"created_at" json:"created_at"`
	CreatedBy  *string   `db:"created_by" json:"created_by"`
	ModifiedAt null.Time `db:"modified_at" json:"modified_at"`
	ModifiedBy *string   `db:"modified_by" json:"modified_by"`
}

type DbsList []Dbs
