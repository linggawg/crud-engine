package queries

import (
	"context"
	models "crud-engine/modules/services/models/domain"
	"github.com/jmoiron/sqlx"
)

type ServicesPostgreQuery struct {
	db *sqlx.DB
}

func NewServicesQuery(db *sqlx.DB) *ServicesPostgreQuery {
	return &ServicesPostgreQuery{db}
}

func (s *ServicesPostgreQuery) GetByServiceUrlAndMethod(ctx context.Context, serviceUrl, method string) (services *models.Services, err error) {
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

func (s *ServicesPostgreQuery) GetByServiceDefinitionAndMethod(ctx context.Context, serviceDefinition, method string) (services *models.Services, err error) {
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
