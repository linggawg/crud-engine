package handler

import (
	"crud-engine/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func Create(c echo.Context) error {
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

	fmt.Println(columns)
	fmt.Println(values)
	sqlStatement := getColumn(table, columns) + getValue(values)
	fmt.Println(sqlStatement)

	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	result, err := stmt.Exec()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	log.Println(result.LastInsertId())
	return c.JSON(http.StatusOK, values)
}

func getColumn(table string, columns string) string {
	sqlStatement := "INSERT " + table + " (" + columns + ") "
	return sqlStatement
}

func getValue(values string) string {
	sqlStatement := "VALUES (" + values + ")"
	return sqlStatement
}
