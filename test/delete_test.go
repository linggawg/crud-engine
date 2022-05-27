package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"crud-engine/handler"
	conn "crud-engine/pkg/database"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	responseDeleteJSON := `{"success":true,"data":0,"message":"successfully delete province with id 99","code":200}
`
	sqlx, err := conn.InitSqlx(GlobalEnv.SQLXDatabase)
	if err != nil {
		panic(err)
	}
	log.Println("Database successfully initialized")
	
	q := make(url.Values)
	q.Set("field_id", "id")

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("table", "value")
	c.SetParamValues("province", "99")

	if assert.NoError(t, handler.New(sqlx).Delete(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, responseDeleteJSON, rec.Body.String())
		log.Println("Delete handler test success")
	}
}