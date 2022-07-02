package queries

import (
	"context"
	models "engine/bin/modules/userservice/models/domain"

	"github.com/jmoiron/sqlx"
)

type UserServicePostgreQuery struct {
	db *sqlx.DB
}

func NewUserServiceQuery(db *sqlx.DB) *UserServicePostgreQuery {
	return &UserServicePostgreQuery{db}
}

func (s *UserServicePostgreQuery) GetByServiceIDAndUserId(ctx context.Context, serviceId, userId string) (userServices *models.UserService, err error) {
	var userService models.UserService
	query := `
	SELECT
		id,
		user_id,
		service_id
	FROM
		user_service
	WHERE service_id = $1 AND user_id = $2
		`
	err = s.db.GetContext(ctx, &userService, query, serviceId, userId)
	if err != nil {
		return nil, err
	}
	return &userService, nil
}
