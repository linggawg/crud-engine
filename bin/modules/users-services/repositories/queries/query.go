package queries

import (
	"context"
	models "engine/bin/modules/users-services/models/domain"
)

type UsersServicesPostgre interface {
	FindOneByServiceIDAndUserId(ctx context.Context, serviceId, userId string) (usersServices *models.UsersServices, err error)
	FindOneByServiceUrlAndUserIdAndMethodAndQueryIsNull(ctx context.Context, serviceUrl, userId, method string) (usersServices *models.UsersServices, err error)
}
