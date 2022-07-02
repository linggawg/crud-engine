package models

import (
	"gopkg.in/guregu/null.v4"
)

type Services struct {
	ID                string    `db:"id" json:"id"`
	DbID              string    `db:"db_id" json:"db_id"`
	Method            string    `db:"method" json:"method"`
	ServiceUrl        *string   `db:"service_url" json:"service_url"`
	ServiceDefinition *string   `db:"service_definition" json:"service_definition"`
	IsQuery           bool      `db:"is_query" json:"is_query"`
	CreatedAt         null.Time `db:"created_at" json:"created_at"`
	CreatedBy         string    `db:"created_by" json:"created_by"`
	ModifiedAt        null.Time `db:"modified_at" json:"modified_at"`
	ModifiedBy        string    `db:"modified_by" json:"modified_by"`
}

// Services list
type ServicesList []Services
