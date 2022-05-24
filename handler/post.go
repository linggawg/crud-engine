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

// Post UpdateData godoc
// @Summary      Insert Data
// @Description  Insert data by column name in format JSON
// @Tags         CrudEngine
// @Accept       json
// @Produce      json
// @Param        table   path    string  true  "Table Name"
// @Param		 insertRequest body map[string]interface{} true "JSON request body based on column name"
// @Security Authorization
// @Success      200  {object} utils.BaseWrapperModel
// @Router       /sql/{table} [post]
func (h *HttpSqlx) Post(c echo.Context) error {
	errorMessage := os.Getenv("POST_ERROR_MESSAGE")
	table := c.Param("table")
	db := h.db

	var jsonBody map[string]interface{}
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
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
	sqlStatement := "INSERT INTO " + table + " (" + columns + ") VALUES (" + values + ")"

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
	return utils.Response(resultId, "successfully insert "+table, http.StatusOK, c)
}
