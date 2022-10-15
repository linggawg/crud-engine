package queries

import (
	"context"
	models "engine/bin/modules/users-services/models/domain"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type UsersServicesPostgreQuery struct {
	db *sqlx.DB
}

func NewUsersServicesQuery(db *sqlx.DB) *UsersServicesPostgreQuery {
	return &UsersServicesPostgreQuery{db}
}

func (s *UsersServicesPostgreQuery) FindOneByServiceIDAndUserId(ctx context.Context, serviceId, userId string) (usersServices *models.UsersServices, err error) {
	var userservice models.UsersServices
	query := `
	SELECT
		id,
		user_id,
		service_id,
		created_at,
		created_by,
		modified_at,
		modified_by
	FROM
		users_services
	WHERE service_id = $1 AND user_id = $2
	`

	err = s.db.PingContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error establishing a database connection")
	}

	err = s.db.GetContext(ctx, &userservice, query, serviceId, userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &userservice, nil
}

func (s *UsersServicesPostgreQuery) FindOneByServiceUrlAndUserIdAndMethodAndQueryIsNull(ctx context.Context, serviceUrl, userId, method string) (usersServices *models.UsersServices, err error) {
	var userservice models.UsersServices
	query := `
	SELECT
		us.id,
		us.user_id,
		us.service_id,
		us.created_at,
		us.created_by,
		us.modified_at,
		us.modified_by
	FROM
		users_services us
	JOIN services s ON s.id = us.service_id
	WHERE s.service_url = $1 AND us.user_id = $2 AND s.method = $3 AND s.query_id IS NULL
	`

	err = s.db.PingContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error establishing a database connection")
	}

	err = s.db.GetContext(ctx, &userservice, query, serviceUrl, userId, method)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &userservice, nil
}
