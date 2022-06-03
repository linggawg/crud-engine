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

type Services struct {
	ID                string    `db:"id" json:"id"`
	DbID              string    `db:"db_id" json:"db_id"`
	Method            string    `db:"method" json:"method"`
	ServiceUrl        string    `db:"service_url" json:"service_url"`
	ServiceDefinition string    `db:"service_definition" json:"service_definition"`
	IsQuery           bool      `db:"is_query" json:"is_query"`
	CreatedAt         null.Time `db:"created_at" json:"created_at"`
	CreatedBy         string    `db:"created_by" json:"created_by"`
	ModifiedAt        null.Time `db:"modified_at" json:"modified_at"`
	ModifiedBy        string    `db:"modified_by" json:"modified_by"`
}

type UserService struct {
	ID         string    `db:"id" json:"id"`
	UserID     string    `db:"user_id" json:"user_id"`
	ServiceID  string    `db:"service_id" json:"service_id"`
	CreatedAt  null.Time `db:"created_at" json:"created_at"`
	CreatedBy  string    `db:"created_by" json:"created_by"`
	ModifiedAt null.Time `db:"modified_at" json:"modified_at"`
	ModifiedBy string    `db:"modified_by" json:"modified_by"`
}
