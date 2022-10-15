package models

import (
	"gopkg.in/guregu/null.v4"
)

type ResourcesMapping struct {
	ID           string    `db:"id" json:"id"`
	ServiceId    string    `db:"service_id" json:"service_id"`
	SourceOrigin string    `db:"source_origin" json:"source_origin"`
	SourceAlias  string    `db:"source_alias" json:"source_alias"`
	CreatedAt    null.Time `db:"created_at" json:"created_at"`
	CreatedBy    *string   `db:"created_by" json:"created_by"`
	ModifiedAt   null.Time `db:"modified_at" json:"modified_at"`
	ModifiedBy   *string   `db:"modified_by" json:"modified_by"`
}

type ResourcesMappingList []ResourcesMapping
