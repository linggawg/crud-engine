package handler

import (
	"crud-engine/config"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func Post(c echo.Context) error {
	table := c.Param("table")
	db := config.CreateCon()

	var jsonBody map[string]interface{}
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var columns string
	var values string
	for key := range jsonBody {
		columns += key + ", "
		values += "'"
		values += jsonBody[key].(string)
		values += "', "
	}
	columns = strings.TrimRight(columns, ", ")
	values = strings.TrimRight(values, ", ")
	sqlStatement := "INSERT " + table + " (" + columns + ") VALUES (" + values + ")"

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