package handler

import (
	"crud-engine/pkg/utils"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// Delete DeleteDate godoc
// @Summary      Delete Data
// @Description  Delete data by ID (primary key)
// @Tags         CrudEngine
// @Accept       json
// @Produce      json
// @Param        table   path    string  true  "Table Name"
// @Param        id   path    string  true  "Primary Key"
// @Security Authorization
// @Success      200  {object} utils.BaseWrapperModel
// @Router       /sql/{table}/{id} [delete]
func (h *HttpSqlx) Delete(c echo.Context) error {
	errorMessage := os.Getenv("DELETE_ERROR_MESSAGE")
	table := c.Param("table")
	db := h.db

	value := c.Param("value")
	field := c.QueryParam("field_id")
	informationSchemas, err := sqlIsNullable(db, table, os.Getenv("DB_DIALECT"), c)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}
	isFoundField := false
	for _, i := range informationSchemas {
		if i.ColumName == field {
			isFoundField = true
			break
		}
	}
	if !isFoundField {
		errorMessage += ", field_id '" + field + "' is not found"
		return utils.Response(nil, errorMessage, http.StatusBadRequest, c)
	}
	sqlStatement := "DELETE FROM " + table + " WHERE " + field + " ='" + value + "'"

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
	message := "successfully delete " + table + " with " + field + " " + value
	return utils.Response(resultId, message, http.StatusOK, c)
}
