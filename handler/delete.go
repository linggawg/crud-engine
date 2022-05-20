package handler

import (
	"crud-engine/pkg/utils"
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
// @Success      200  {object} map[string]interface{}
// @Router       /{table}/{id} [delete]
func (h *HttpSqlx) Delete(c echo.Context) error {
	table := c.Param("table")
	db := h.db

	id := c.Param("id")
	key, err := getPrimaryKey(db, table, c)
	if err != nil {
		log.Println(err)
		return utils.Response(nil, err.Error(), http.StatusBadRequest, c)
	}

	sqlStatement := "DELETE FROM " + table + " WHERE " + key + " ='" + id + "'"

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
	message := "successfully update " + table + " with Id " + id
	return utils.Response(resultId, message, http.StatusOK, c)
}
