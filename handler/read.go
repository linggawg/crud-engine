package handler

import (
	"crud-engine/config"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Read(c echo.Context) error {
	table := c.Param("table")
	db := config.CreateCon()

	sqlStatement := "SELECT * FROM " + table

	rows, err := db.QueryContext(c.Request().Context(), sqlStatement)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}

	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, string(jsonData))
}
