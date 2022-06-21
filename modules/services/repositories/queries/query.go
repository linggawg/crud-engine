package queries

import (
	"context"
	models "crud-engine/modules/services/models/domain"
)

type ServicesPostgre interface {
	GetByServiceUrlAndMethod(ctx context.Context, serviceUrl, method string) (services *models.Services, err error)
	GetByServiceDefinitionAndMethod(ctx context.Context, serviceDefinition, method string) (services *models.Services, err error)
}
