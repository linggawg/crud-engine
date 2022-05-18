package handler

import (
	"crud-engine/config"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// Put UpdateData godoc
// @Summary      Update Data
// @Description  Update data by ID (primary key) and data by column name in format JSON
// @Tags         CrudEngine
// @Accept       json
// @Produce      json
// @Param        table   path    string  true  "Table Name"
// @Param        id   path    string  true  "Primary Key"
// @Param		 updateRequest body map[string]interface{} true "JSON request body based on column name"
// @Success      200  {object} map[string]interface{}
// @Router       /{table}/{id} [put]
func Put(c echo.Context) error {
	table := c.Param("table")
	id := c.Param("id")
	db := config.CreateCon()

	var jsonBody map[string]interface{}
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	primaryKey, err := getPrimaryKey(db, table, c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var setData string
	for key := range jsonBody {
		setData += key + "='" + jsonBody[key].(string) + "', "
	}
	setData = strings.TrimRight(setData, ", ")

	sqlStatement := "UPDATE " + table + " SET " + setData + " WHERE " + primaryKey + "='" + id + "'"

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
