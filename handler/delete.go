package handler

import (
	"crud-engine/config"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Delete(c echo.Context) error {
	table := c.Param("table")
	db := config.CreateCon()

	id := c.Param("id")
	sqlStatement := "DELETE FROM " + table + " WHERE id = " + id

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
