package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"crud-engine/handler"
	conn "crud-engine/pkg/database"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPut(t *testing.T) {
	responsePutJSON := `{"success":true,"data":0,"message":"successfully update province with id 99","code":200}
`
	putJSON := `{"created_by": "system","created_date": "2022-05-27 09:45:06","id": "99","is_deleted": "0","name": "TEST PUT PROVINCE"}`

	sqlx, err := conn.InitSqlx(GlobalEnv.SQLXDatabase)
	if err != nil {
		panic(err)
	}
	log.Println("Database successfully initialized")
	
	q := make(url.Values)
	q.Set("field_id", "id")

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/?"+q.Encode(), strings.NewReader(putJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("table", "value")
	c.SetParamValues("province", "99")

	if assert.NoError(t, handler.New(sqlx).Put(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, responsePutJSON, rec.Body.String())
		log.Println("Put handler test success")
	}
}