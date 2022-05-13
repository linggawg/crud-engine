package handler

import (
	"crud-engine/config"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Get(c echo.Context) error {
	db := config.CreateCon()
	var sqlStatement string

	isQuery := c.QueryParam("isQuery")
	if isQuery == "true" {
		sqlStatement = c.Param("table")
	}else{
		table := c.Param("table")
		primaryKey, err := getPrimaryKey(db, table, c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	
		// pageNo := c.QueryParam("pageNo")
		pageSize := ""
		if c.QueryParam("pageSize") != "" {
			pageSize = " LIMIT " + c.QueryParam("pageSize")
		}

		isDistinct := ""
		if c.QueryParam("isDistinct") == "true" {
			isDistinct = "DISTINCT "
		}
	
		query := ""
		if c.QueryParam("query") != "" {
			query = " WHERE " + c.QueryParam("query")
		}
		
		colls := c.QueryParam("colls")
		if colls == "" {
			colls = "*"
		}
	
		sortBy := " ORDER BY " 
		if c.QueryParam("sortBy") == "" {
			sortBy += primaryKey+", asc" 
		} else {
			sortBy += c.QueryParam("sortBy")
		}
		sqlStatement = "SELECT " + isDistinct + colls +" FROM " + table + query + sortBy + pageSize
	}

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
