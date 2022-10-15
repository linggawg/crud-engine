package commands

import (
	"context"
	models "engine/bin/modules/users-services/models/domain"
)

type UsersServicesPostgre interface {
	InsertOne(ctx context.Context, usersServices *models.UsersServices) error
	DeleteByUsersIdAndServiceUrl(ctx context.Context, usersId, serviceUrl string) error
}
