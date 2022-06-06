package userservice_test

import (
	"context"
	"crud-engine/modules/models"
	_userservice "crud-engine/modules/userservice"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"testing"
	"time"
)

var mockUserService = models.UserServices{
	models.UserService{
		ID:         "53b4aac9-fae2-4f54-b7c1-f88d8caef32f",
		UserID:     "424cf355-51cb-4325-a085-71b0b8ac2a38",
		ServiceID:  "ddfc24fe-7eb5-4c4d-88a5-2001bc16e76e",
		CreatedAt:  null.TimeFrom(time.Now()),
		CreatedBy:  "424cf355-51cb-4325-a085-71b0b8ac2a38",
		ModifiedAt: null.TimeFrom(time.Now()),
		ModifiedBy: "424cf355-51cb-4325-a085-71b0b8ac2a38",
	},
	models.UserService{
		ID:         "1f70c357-cfe0-45fa-9fec-281b37c1600c",
		UserID:     "93147c1d-197e-4556-9d24-1992196aaa03",
		ServiceID:  "6c34520b-d554-47b8-9fd2-f84bb6a042e6",
		CreatedAt:  null.TimeFrom(time.Now()),
		CreatedBy:  "93147c1d-197e-4556-9d24-1992196aaa03",
		ModifiedAt: null.TimeFrom(time.Now()),
		ModifiedBy: "93147c1d-197e-4556-9d24-1992196aaa03",
	},
}

func TestGetByServiceIDAndUserId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "service_id", "created_at", "created_by", "modified_at", "modified_by"})
	for _, user := range mockUserService {
		rows.AddRow(user.ID, user.UserID, user.ServiceID, user.CreatedAt, user.CreatedBy, user.ModifiedAt, user.ModifiedBy)
	}

	query := "SELECT id, user_id, service_id FROM user_service WHERE service_id = \\$1 AND user_id = \\$2"
	mock.ExpectQuery(query).WithArgs(mockUserService[0].ServiceID, mockUserService[0].UserID).WillReturnRows(rows)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	a := _userservice.New(sqlxDB)

	mUserServices, err := a.GetByServiceIDAndUserId(context.TODO(), mockUserService[0].ServiceID, mockUserService[0].UserID)
	assert.NoError(t, err)
	assert.NotNil(t, mUserServices)
}

func TestGetByServiceIDAndUserIdNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	rows := sqlmock.NewRows([]string{"id", "user_id", "service_id", "created_at", "created_by", "modified_at", "modified_by"})
	query := "SELECT id, user_id, service_id FROM user_service WHERE service_id = \\$1 AND user_id = \\$2"
	mock.ExpectQuery(query).WithArgs(mockUserService[0].ServiceID, mockUserService[0].UserID).WillReturnRows(rows)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	a := _userservice.New(sqlxDB)

	mUserServices, err := a.GetByServiceIDAndUserId(context.TODO(), mockUserService[0].ServiceID, mockUserService[0].UserID)
	assert.Empty(t, mUserServices)
	assert.Error(t, err)
	assert.Equal(t, err, sql.ErrNoRows)
}
