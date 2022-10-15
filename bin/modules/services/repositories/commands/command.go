package commands

import (
	"context"
	models "engine/bin/modules/services/models/domain"
)

type ServicesPostgre interface {
	InsertOne(ctx context.Context, services *models.Services) error
	DeleteByServiceUrl(ctx context.Context, serviceUrl string) error
}
