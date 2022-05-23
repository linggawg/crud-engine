package handler

import (
	"crud-engine/pkg/utils"
	"fmt"
	"log"
	"net/http"

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
	table := c.Param("table")
	db := h.db

	value := c.Param("value")
	field := c.QueryParam("field_id")
	sqlStatement := "DELETE FROM " + table + " WHERE " + field + " ='" + value + "'"
	fmt.Println(sqlStatement)

	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}

	result, err := stmt.Exec()
	if err != nil {
		log.Println(err)
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}

	resultId, _ := result.LastInsertId()
	message := "successfully delete " + table + " with " + field + " " + value
	return utils.Response(resultId, message, http.StatusOK, c)
}
