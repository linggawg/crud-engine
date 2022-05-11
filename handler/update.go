package handler

import (
	"crud-engine/config"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Update(c echo.Context) error {
	//table := strings.ReplaceAll(c.Path(), "/", "")
	db := config.CreateCon()

	//id := c.Param("id")
	sqlStatement := "UPDATE table_name SET column1=value, column2=value2, WHERE some_column=some_value"

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
