package handler

import (
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
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sqlStatement := "DELETE FROM " + table + " WHERE " + key + " ='" + id + "'"

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
