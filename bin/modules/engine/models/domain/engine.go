package models

import (
	dbsModels "engine/bin/modules/dbs/models/domain"
	resourcesMappingModels "engine/bin/modules/resources-mapping/models/domain"
	"engine/bin/pkg/token"
)

type PrimaryKey struct {
	Column string
	Format string
}

type KeyPostgres struct {
	AttName    string `db:"attname"`
	FormatType string `db:"format_type"`
}

type KeyMySQL struct {
	Field       string  `db:"Field"`
	TypeData    string  `db:"Type"`
	NullData    string  `db:"Null"`
	Key         string  `db:"Key"`
	DefaultData *string `db:"Default"`
	Extra       *string `db:"Extra"`
}

type InformationSchema struct {
	ColumName  string `db:"column_name"`
	IsNullable string `db:"is_nullable"`
}

type GetList struct {
	Page       *int        `json:"pageNo"`
	Size       *int        `json:"pageSize"`
	Sort       string      `json:"sortBy"`
	IsDistinct bool        `json:"isDistinct"`
	Filter     string      `json:"filter"`
	Columns    string      `json:"columns"`
	Key        string      `json:"key"`
	Opts       token.Claim `json:"opts"`
}

type EngineRequest struct {
	Table   string                 `json:"table"`
	FieldId string                 `json:"field_id"`
	Value   string                 `json:"value"`
	Data    map[string]interface{} `json:"data"`
	Opts    token.Claim            `json:"opts"`
}

type EngineConfig struct {
	Dbs                  dbsModels.Dbs
	ResourcesMappingList resourcesMappingModels.ResourcesMappingList
}
