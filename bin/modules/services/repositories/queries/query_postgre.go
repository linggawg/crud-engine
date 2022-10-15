package queries

import (
	"context"
	models "engine/bin/modules/services/models/domain"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type ServicesPostgreQuery struct {
	db *sqlx.DB
}

func NewServicesQuery(db *sqlx.DB) *ServicesPostgreQuery {
	return &ServicesPostgreQuery{db}
}

func (s *ServicesPostgreQuery) FindOneByServiceUrlAndMethodAndQueryIsNull(ctx context.Context, serviceUrl, method string) (services *models.Services, err error) {
	var service models.Services
	query := `
	SELECT
		id,
		db_id,
		query_id,
		service_url,
		method,
		created_at,
		created_by,
		modified_at,
		modified_by
	FROM
		services
	WHERE service_url = $1 AND method = $2 AND query_id IS NULL
		`

	err = s.db.PingContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error establishing a database connection")
	}

	err = s.db.GetContext(ctx, &service, query, serviceUrl, method)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &service, nil
}

func (s *ServicesPostgreQuery) FindOneByQueryID(ctx context.Context, queryID string) (services *models.Services, err error) {
	var service models.Services
	query := `
	SELECT
		id,
		db_id,
		query_id,
		service_url,
		method,
		created_at,
		created_by,
		modified_at,
		modified_by
	FROM
		services
	WHERE query_id = $1
	`

	err = s.db.PingContext(ctx)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error establishing a database connection")
	}

	err = s.db.GetContext(ctx, &service, query, queryID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &service, nil
}
