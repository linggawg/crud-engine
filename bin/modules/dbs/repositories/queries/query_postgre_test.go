package queries_test

import (
	"context"
	models "engine/bin/modules/dbs/models/domain"
	"engine/bin/modules/dbs/repositories/queries"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func NewMock() (*sqlx.DB, sqlmock.Sqlmock) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	db := sqlx.NewDb(mockDb, "sqlmock")
	return db, mock
}

var random = uuid.New().String()
var dbs = &models.Dbs{
	ID:         random,
	AppID:      random,
	Name:       random,
	Host:       random,
	Port:       9999,
	Username:   random,
	Password:   &random,
	Dialect:    random,
	CreatedAt:  null.TimeFrom(time.Now()),
	CreatedBy:  &random,
	ModifiedAt: null.TimeFrom(time.Now()),
	ModifiedBy: &random,
}

func TestFindOneByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock := NewMock()

		query := "SELECT id, app_id, name, host, port, username, password, dialect FROM dbs WHERE id = $1"

		rows := sqlmock.NewRows([]string{"id", "app_id", "name", "host", "port", "username", "password", "dialect", "created_at", "created_by", "modified_at", "modified_by"}).
			AddRow(dbs.ID, dbs.AppID, dbs.Name, dbs.Host, dbs.Port, dbs.Username, dbs.Password, dbs.Dialect, dbs.CreatedAt, dbs.CreatedBy, dbs.ModifiedAt, dbs.ModifiedBy)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(dbs.ID).WillReturnRows(rows)

		entity, err := queries.NewDbsQuery(db).FindOneByID(context.TODO(), dbs.ID)
		assert.NotNil(t, entity)
		assert.NoError(t, err)
	})
	t.Run("Error", func(t *testing.T) {
		db, mock := NewMock()

		query := "SELECT id, app_id, name, host, port, username, password, dialect FROM dbs WHERE id = $1"

		rows := sqlmock.NewRows([]string{"id", "app_id", "name", "host", "port", "username", "password", "dialect", "created_at", "created_by", "modified_at", "modified_by"})
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(dbs.ID).WillReturnRows(rows)

		entity, err := queries.NewDbsQuery(db).FindOneByID(context.TODO(), dbs.ID)
		assert.Empty(t, entity)
		assert.Error(t, err)
	})
	t.Run("DB no connection", func(t *testing.T) {
		db, _ := NewMock()

		db.Close()
		entity, err := queries.NewDbsQuery(db).FindOneByID(context.TODO(), dbs.ID)
		assert.Empty(t, entity)
		assert.Error(t, err)
		assert.EqualError(t, err, "error establishing a database connection")
	})
}
