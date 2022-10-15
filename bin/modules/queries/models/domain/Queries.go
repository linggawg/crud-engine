package models

import (
	"gopkg.in/guregu/null.v4"
)

type Queries struct {
	ID         		string    `db:"id" json:"id"`
	Key         	string    `db:"key" json:"key"`
	QueryDefinition	string    `db:"query_definition" json:"query_definition"`
	CreatedAt  		null.Time `db:"created_at" json:"created_at"`
	CreatedBy  		*string   `db:"created_by" json:"created_by"`
	ModifiedAt	 	null.Time `db:"modified_at" json:"modified_at"`
	ModifiedBy 		*string   `db:"modified_by" json:"modified_by"`
}

type QueriesList []Queries
