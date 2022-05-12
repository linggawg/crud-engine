package handler

import (
	"crud-engine/config"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Delete(c echo.Context) error {
	table := c.Param("table")
	db := config.CreateCon()

	id := c.Param("id")
	key, err := getPrimaryKey(db, table, c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sqlStatement := "DELETE FROM " + table + " WHERE " + key + " = " + id

	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	result, err := stmt.Exec()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

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
