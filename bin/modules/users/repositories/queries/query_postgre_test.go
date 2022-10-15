package queries_test

import (
	"context"
	models "engine/bin/modules/users/models/domain"
	"engine/bin/modules/users/repositories/queries"
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
var users = &models.Users{
	ID:         random,
	RoleID:     random,
	Username:   random,
	Password:   random,
	CreatedAt:  null.TimeFrom(time.Now()),
	CreatedBy:  &random,
	ModifiedAt: null.TimeFrom(time.Now()),
	ModifiedBy: &random,
}

func TestFindOneByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock := NewMock()

		query := `
		SELECT
			id,
			role_id,
			username,
			password,
			created_at,
			created_by,
			modified_at,
			modified_by
		FROM
			users
		WHERE id = $1
			`

		rows := sqlmock.NewRows([]string{"id", "role_id", "username", "password", "created_at", "created_by", "modified_at", "modified_by"}).
			AddRow(users.ID, users.RoleID, users.Username, users.Password, users.CreatedAt, users.CreatedBy, users.ModifiedAt, users.ModifiedBy)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(users.ID).WillReturnRows(rows)

		entity, err := queries.NewUsersQuery(db).FindOneByID(context.TODO(), users.ID)
		assert.NotNil(t, entity)
		assert.NoError(t, err)
	})
	t.Run("Error", func(t *testing.T) {
		db, mock := NewMock()

		query := `
		SELECT
			id,
			role_id,
			username,
			password,
			created_at,
			created_by,
			modified_at,
			modified_by
		FROM
			users
		WHERE id = $1
			`

		rows := sqlmock.NewRows([]string{"id", "role_id", "username", "password", "created_at", "created_by", "modified_at", "modified_by"})
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(users.ID).WillReturnRows(rows)

		entity, err := queries.NewUsersQuery(db).FindOneByID(context.TODO(), users.ID)
		assert.Empty(t, entity)
		assert.Error(t, err)
	})
	t.Run("DB no connection", func(t *testing.T) {
		db, _ := NewMock()

		db.Close()
		entity, err := queries.NewUsersQuery(db).FindOneByID(context.TODO(), users.Username)
		assert.Empty(t, entity)
		assert.Error(t, err)
		assert.EqualError(t, err, "error establishing a database connection")
	})
}

func TestFindOneByUsername(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock := NewMock()

		query := `
		SELECT
			id,
			role_id,
			username,
			password,
			created_at,
			created_by,
			modified_at,
			modified_by
		FROM
			users
		WHERE username = $1
			`

		rows := sqlmock.NewRows([]string{"id", "role_id", "username", "password", "created_at", "created_by", "modified_at", "modified_by"}).
			AddRow(users.ID, users.RoleID, users.Username, users.Password, users.CreatedAt, users.CreatedBy, users.ModifiedAt, users.ModifiedBy)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(users.Username).WillReturnRows(rows)

		entity, err := queries.NewUsersQuery(db).FindOneByUsername(context.TODO(), users.Username)
		assert.NotNil(t, entity)
		assert.NoError(t, err)
	})
	t.Run("Error", func(t *testing.T) {
		db, mock := NewMock()

		query := `
		SELECT
			id,
			role_id,
			username,
			password,
			created_at,
			created_by,
			modified_at,
			modified_by
		FROM
			users
		WHERE username = $1
			`

		rows := sqlmock.NewRows([]string{"id", "role_id", "username", "password", "created_at", "created_by", "modified_at", "modified_by"})
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(users.Username).WillReturnRows(rows)

		entity, err := queries.NewUsersQuery(db).FindOneByUsername(context.TODO(), users.Username)
		assert.Empty(t, entity)
		assert.Error(t, err)
	})
	t.Run("DB no connection", func(t *testing.T) {
		db, _ := NewMock()

		db.Close()
		entity, err := queries.NewUsersQuery(db).FindOneByUsername(context.TODO(), users.Username)
		assert.Empty(t, entity)
		assert.Error(t, err)
		assert.EqualError(t, err, "error establishing a database connection")
	})
}
