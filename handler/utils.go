package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func getPrimaryKey(db *sqlx.DB, table string, c echo.Context) (key string, err error) {
	var keyMysql KeyMySQL
	sql := "show columns from " + table + " where `Key` = 'PRI'"
	err = db.GetContext(c.Request().Context(), &keyMysql, sql)
	if err != nil {
		return "", err
	}
	return keyMysql.Field, nil
}

type KeyMySQL struct {
	Field       string  `db:"Field"`
	TypeData    string  `db:"Type"`
	NullData    string  `db:"Null"`
	Key         string  `db:"Key"`
	DefaultData *string `db:"Default"`
	Extra       *string `db:"Extra"`
}
