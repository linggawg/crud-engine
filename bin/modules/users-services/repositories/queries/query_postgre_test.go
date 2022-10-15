package queries_test

import (
	"context"
	models "engine/bin/modules/users-services/models/domain"
	"engine/bin/modules/users-services/repositories/queries"
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
var usersServices = &models.UsersServices{
	ID:         random,
	UserID:     random,
	ServiceID:  random,
	CreatedAt:  null.TimeFrom(time.Now()),
	CreatedBy:  &random,
	ModifiedAt: null.TimeFrom(time.Now()),
	ModifiedBy: &random,
}

func TestFindOneByServiceIDAndUserId(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock := NewMock()

		query := "SELECT id, user_id, service_id, created_at, created_by, modified_at, modified_by FROM users_services WHERE service_id = $1 AND user_id = $2"

		rows := sqlmock.NewRows([]string{"id", "user_id", "service_id", "created_at", "created_by", "modified_at", "modified_by"}).
			AddRow(usersServices.ID, usersServices.UserID, usersServices.ServiceID, usersServices.CreatedAt, usersServices.CreatedBy, usersServices.ModifiedAt, usersServices.ModifiedBy)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(usersServices.ID, usersServices.UserID).WillReturnRows(rows)

		entity, err := queries.NewUsersServicesQuery(db).FindOneByServiceIDAndUserId(context.TODO(), usersServices.ID, usersServices.UserID)
		assert.NotNil(t, entity)
		assert.NoError(t, err)
	})
	t.Run("Error", func(t *testing.T) {
		db, mock := NewMock()

		query := "SELECT id, user_id, service_id, created_at, created_by, modified_at, modified_by FROM users_services WHERE service_id = $1 AND user_id = $2"

		rows := sqlmock.NewRows([]string{"id", "user_id", "service_id", "created_at", "created_by", "modified_at", "modified_by"})
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(usersServices.ID, usersServices.UserID).WillReturnRows(rows)

		entity, err := queries.NewUsersServicesQuery(db).FindOneByServiceIDAndUserId(context.TODO(), usersServices.ID, usersServices.UserID)
		assert.Empty(t, entity)
		assert.Error(t, err)
	})
	t.Run("DB no connection", func(t *testing.T) {
		db, _ := NewMock()

		db.Close()
		entity, err := queries.NewUsersServicesQuery(db).FindOneByServiceIDAndUserId(context.TODO(), usersServices.ID, usersServices.UserID)
		assert.Empty(t, entity)
		assert.Error(t, err)
		assert.EqualError(t, err, "error establishing a database connection")
	})
}

func TestFindOneByServiceUrlAndUserIdAndMethod(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock := NewMock()

		query := "SELECT us.id, us.user_id, us.service_id, us.created_at, us.created_by, us.modified_at, us.modified_by FROM users_services us JOIN services s ON s.id = us.service_id WHERE s.service_url = $1 AND us.user_id = $2 AND s.method = $3"

		rows := sqlmock.NewRows([]string{"id", "user_id", "service_id", "created_at", "created_by", "modified_at", "modified_by"}).
			AddRow(usersServices.ID, usersServices.UserID, usersServices.ServiceID, usersServices.CreatedAt, usersServices.CreatedBy, usersServices.ModifiedAt, usersServices.ModifiedBy)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("service_url", usersServices.UserID, "method").WillReturnRows(rows)

		entity, err := queries.NewUsersServicesQuery(db).FindOneByServiceUrlAndUserIdAndMethodAndQueryIsNull(context.TODO(), "service_url", usersServices.UserID, "method")
		assert.NotNil(t, entity)
		assert.NoError(t, err)
	})
	t.Run("Error", func(t *testing.T) {
		db, mock := NewMock()

		query := "SELECT us.id, us.user_id, us.service_id, us.created_at, us.created_by, us.modified_at, us.modified_by FROM users_services us JOIN services s ON s.id = us.service_id WHERE s.service_url = $1 AND us.user_id = $2 AND s.method = $3"

		rows := sqlmock.NewRows([]string{"id", "user_id", "service_id", "created_at", "created_by", "modified_at", "modified_by"})

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("service_url", usersServices.UserID, "method").WillReturnRows(rows)

		entity, err := queries.NewUsersServicesQuery(db).FindOneByServiceUrlAndUserIdAndMethodAndQueryIsNull(context.TODO(), "service_url", usersServices.UserID, "method")
		assert.Empty(t, entity)
		assert.Error(t, err)
	})
	t.Run("DB no connection", func(t *testing.T) {
		db, _ := NewMock()

		db.Close()
		entity, err := queries.NewUsersServicesQuery(db).FindOneByServiceUrlAndUserIdAndMethodAndQueryIsNull(context.TODO(), "service_url", usersServices.UserID, "method")
		assert.Empty(t, entity)
		assert.Error(t, err)
		assert.EqualError(t, err, "error establishing a database connection")
	})
}
