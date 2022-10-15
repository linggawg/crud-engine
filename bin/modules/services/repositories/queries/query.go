package queries

import (
	"context"
	models "engine/bin/modules/services/models/domain"
)

type ServicesPostgre interface {
	FindOneByQueryID(ctx context.Context, queryID string) (services *models.Services, err error)
	FindOneByServiceUrlAndMethodAndQueryIsNull(ctx context.Context, serviceUrl, method string) (services *models.Services, err error)
}
