package services

import (
	"context"
	"crud-engine/modules/models"

	"github.com/jmoiron/sqlx"
)

type HttpSqlx struct {
	db *sqlx.DB
}

func (s *HttpSqlx) GetByServiceUrlAndMethod(ctx context.Context, serviceUrl string, method string) (services *models.Services, err error) {
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