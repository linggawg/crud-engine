package handler

import (
	"crud-engine/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"
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
// @Security Authorization
// @Success      200  {object} utils.BaseWrapperModel
// @Router       /sql/{table}/{id} [put]
func (h *HttpSqlx) Put(c echo.Context) error {
	errorMessage := os.Getenv("PUT_ERROR_MESSAGE")
	table := c.Param("table")
	db := h.db

	var jsonBody map[string]interface{}
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	var setData string
	for key := range jsonBody {
		setData += key + "='" + jsonBody[key].(string) + "', "
	}
	setData = strings.TrimRight(setData, ", ")

	value := c.Param("value")
	field := c.QueryParam("field_id")
	sqlStatement := "UPDATE " + table + " SET " + setData + " WHERE " + field + " ='" + value + "'"

	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	result, err := stmt.Exec()
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}

	resultId, _ := result.LastInsertId()
	message := "successfully update "  + table + " with " + field + " " + value
	return utils.Response(resultId, message, http.StatusOK, c)
}
