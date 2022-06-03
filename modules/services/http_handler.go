package services

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

func (s *HttpSqlx) GetByServiceUrlAndMethod(ctx context.Context, serviceUrl, method string) (services *models.Services, err error) {
	var service models.Services
	query := `
	SELECT
		id,
		db_id,
		service_url,
		method,
		service_definition,
		is_query
	FROM
		services
	WHERE service_url = $1 AND method = $2
		`
	err = s.db.GetContext(ctx, &service, query, serviceUrl, method)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (s *HttpSqlx) GetByServiceDefinitionAndMethod(ctx context.Context, serviceDefinition, method string) (services *models.Services, err error) {
	var service models.Services
	query := `
	SELECT
		id,
		db_id,
		service_url,
		method,
		service_definition,
		is_query
	FROM
		services
	WHERE is_query = true AND service_definition = $1 AND method = $2
		`
	err = s.db.GetContext(ctx, &service, query, serviceDefinition, method)
	if err != nil {
		return nil, err
	}
	return &service, nil
}
