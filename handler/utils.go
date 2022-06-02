package handler

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"os"
	"strings"
)

type InformationSchema struct {
	ColumName  string `db:"column_name"`
	IsNullable string `db:"is_nullable"`
}

func SetQuery(query string) (newQuery string) {
	switch os.Getenv("DB_DIALECT") {
	case "postgres":
		count := strings.Count(query, "?")
		for i := 1; i <= count; i++ {
			query = strings.Replace(query, "?", fmt.Sprintf("$%d", i), 1)
		}
		return query
	default:
		return query
	}
}

func sqlIsNullable(db *sqlx.DB, table string, c echo.Context) (informationSchema []InformationSchema, err error) {
	var (
		is      []InformationSchema
		sql     string
		args    []interface{}
		dialect = os.Getenv("DB_DIALECT")
	)

	if dialect == "mysql" {
		sql = "SELECT column_name, is_nullable " +
			"FROM information_schema.COLUMNS " +
			"WHERE table_schema = DATABASE () AND table_name = ? " +
			"ORDER BY table_name, ordinal_position;"
	} else if dialect == "postgres" {
		sql = "SELECT column_name, is_nullable " +
			"FROM information_schema.columns " +
			"WHERE table_schema = current_schema() " +
			"AND table_name = $1 " +
			"order by table_name,ordinal_position; "
	}
	args = append(args, table)
	err = db.SelectContext(c.Request().Context(), &is, sql, args...)
	if err != nil {
		return nil, err
	}
	return is, nil
}
func getPrimaryKey(db *sqlx.DB, table string, c echo.Context) (p *PrimaryKey, err error) {
	var (
		primarykey PrimaryKey
		dialect    = os.Getenv("DB_DIALECT")
	)

	if dialect == "mysql" {
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
	} else if dialect == "postgres" {
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
