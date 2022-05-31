package userservice

import (
	"context"
	"crud-engine/modules/models"

	"github.com/jmoiron/sqlx"
)

type HttpSqlx struct {
	db *sqlx.DB
}

func (s *HttpSqlx) GetByServiceID(ctx context.Context, id string) (userservices *models.UserService, err error) {
	var userservice models.UserService
	query := `
	SELECT
		id,
		user_id,
		service_id
	FROM
		user_service
	WHERE service_id = $1
		`
	err = s.db.GetContext(ctx, &userservice, query, id)
	if err != nil {
		return nil, err
	}
	return &userservice, nil
}

func (s *HttpSqlx) GetByUserID(ctx context.Context, id string) (userservices *models.UserService, err error) {
	var userservice models.UserService
	query := `
	SELECT
		id,
		user_id,
		service_id
	FROM
		user_service
	WHERE user_id = $1
		`
	err = s.db.GetContext(ctx, &userservice, query, id)
	if err != nil {
		return nil, err
	}
	return &userservice, nil
}