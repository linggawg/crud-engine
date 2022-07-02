package models

import "engine/bin/pkg/token"

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
	IsQuery    bool        `json:"isQuery"`
	IsDistinct bool        `json:"isDistinct"`
	Query      string      `json:"query"`
	Colls      string      `json:"colls"`
	Opts       token.Claim `json:"opts"`
}

type EngineRequest struct {
	Table   string                 `json:"table"`
	FieldId string                 `json:"field_id"`
	Value   string                 `json:"value"`
	Data    map[string]interface{} `json:"data"`
	Opts    token.Claim            `json:"opts"`
}
