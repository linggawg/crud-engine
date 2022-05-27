package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"os"
	"strings"
)

func getPrimaryKey(db *sqlx.DB, table string, c echo.Context) (p *PrimaryKey, err error) {
	var primarykey PrimaryKey
	if os.Getenv("DB_DIALECT") == "mysql" {
		var keyMysql KeyMySQL
		sql := "show columns from " + table + " where `Key` = 'PRI'"
		err = db.GetContext(c.Request().Context(), &keyMysql, sql)
		if err != nil {
			return nil, err
		}
		primarykey.column = keyMysql.Field
		if strings.Contains(keyMysql.TypeData, "int") {
			primarykey.format = "int"
		} else {
			primarykey.format = "varchar"
		}
	} else if os.Getenv("DB_DIALECT") == "postgres" {
		var keyPostgres KeyPostgres
		sql := "SELECT pg_attribute.attname, format_type(pg_attribute.atttypid, pg_attribute.atttypmod) " +
			"FROM pg_index, pg_class, pg_attribute, pg_namespace " +
			"WHERE pg_class.oid = '" + table + "'::regclass " +
			"AND indrelid = pg_class.oid " +
			"AND pg_class.relnamespace = pg_namespace.oid " +
			"AND pg_attribute.attrelid = pg_class.oid " +
			"AND pg_attribute.attnum = any(pg_index.indkey) AND indisprimary"
		err = db.GetContext(c.Request().Context(), &keyPostgres, sql)
		if err != nil {
			return nil, err
		}
		primarykey.column = keyPostgres.Attname
		if strings.Contains(keyPostgres.FormatType, "int") {
			primarykey.format = "int"
		} else {
			primarykey.format = "varchar"
		}
	}
	return &primarykey, err
}

type PrimaryKey struct {
	column string
	format string
}

type KeyPostgres struct {
	Attname    string `db:"attname"`
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
