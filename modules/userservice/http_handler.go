package userservice

import (
	"context"
	"crud-engine/modules/models"

	"github.com/jmoiron/sqlx"
)

type HttpSqlx struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *HttpSqlx {
	return &HttpSqlx{db}
}

func (s *HttpSqlx) GetByServiceIDAndUserId(ctx context.Context, serviceId string, userId string) (userservices *models.UserService, err error) {
	var userservice models.UserService
	query := `
	SELECT
		id,
		user_id,
		service_id
	FROM
		user_service
	WHERE service_id = $1 AND user_id = $2
		`
	err = s.db.GetContext(ctx, &userservice, query, serviceId, userId)
	if err != nil {
		return nil, err
	}
	return &userservice, nil
}
