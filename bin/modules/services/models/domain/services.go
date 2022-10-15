package models

import (
	"engine/bin/pkg/token"

	"gopkg.in/guregu/null.v4"
)

type Services struct {
	ID                string    `db:"id" json:"id"`
	DbID              string    `db:"db_id" json:"db_id"`
	QueryID           *string   `db:"query_id" json:"query_id"`
	Method            string    `db:"method" json:"method"`
	ServiceUrl        *string   `db:"service_url" json:"service_url"`
	CreatedAt         null.Time `db:"created_at" json:"created_at"`
	CreatedBy         *string    `db:"created_by" json:"created_by"`
	ModifiedAt        null.Time `db:"modified_at" json:"modified_at"`
	ModifiedBy        *string    `db:"modified_by" json:"modified_by"`
}

type ServicesList []Services

type ServicesRequest struct {
	ServiceUrl 	string		`json:"service_url" validate:"required"`
	Opts 	 	token.Claim `json:"opts"`
}