package queries_test

import (
	"context"
	models "engine/bin/modules/queries/models/domain"
	"engine/bin/modules/queries/repositories/queries"
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
var queriis = &models.Queries{
	ID:              random,
	Key:             random,
	QueryDefinition: random,
	CreatedAt:       null.TimeFrom(time.Now()),
	CreatedBy:       &random,
	ModifiedAt:      null.TimeFrom(time.Now()),
	ModifiedBy:      &random,
}

func TestFindOneByKey(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock := NewMock()

		query := "SELECT id, key, query_definition, created_at, created_by, modified_at, modified_by FROM queries WHERE key = $1"

		rows := sqlmock.NewRows([]string{"id", "key", "query_definition", "created_at", "created_by", "modified_at", "modified_by"}).
			AddRow(queriis.ID, queriis.Key, queriis.QueryDefinition, queriis.CreatedAt, queriis.CreatedBy, queriis.ModifiedAt, queriis.ModifiedBy)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(queriis.Key).WillReturnRows(rows)

		entity, err := queries.NewQueriesQuery(db).FindOneByKey(context.TODO(), queriis.Key)
		assert.NotNil(t, entity)
		assert.NoError(t, err)
	})
	t.Run("Error", func(t *testing.T) {
		db, mock := NewMock()

		query := "SELECT id, key, query_definition, created_at, created_by, modified_at, modified_by FROM queries WHERE key = $1"

		rows := sqlmock.NewRows([]string{"id", "key", "query_definition", "created_at", "created_by", "modified_at", "modified_by"})
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(queriis.Key).WillReturnRows(rows)

		entity, err := queries.NewQueriesQuery(db).FindOneByKey(context.TODO(), queriis.Key)
		assert.Empty(t, entity)
		assert.Error(t, err)
	})
	t.Run("DB no connection", func(t *testing.T) {
		db, _ := NewMock()

		db.Close()
		entity, err := queries.NewQueriesQuery(db).FindOneByKey(context.TODO(), queriis.Key)
		assert.Empty(t, entity)
		assert.Error(t, err)
		assert.EqualError(t, err, "error establishing a database connection")
	})
}
