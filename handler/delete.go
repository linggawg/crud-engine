package handler

import (
	"crud-engine/config"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Delete(c echo.Context) error {
	table := c.Param("table")
	db := config.CreateCon()

	id := c.Param("id")
	key, err := getPrimaryKey(db, table, c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sqlStatement := "DELETE FROM " + table + " WHERE " + key + " ='" + id + "'"

	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	result, err := stmt.Exec()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	resultId, _ := result.LastInsertId()
	return c.JSON(http.StatusOK, resultId)
}
