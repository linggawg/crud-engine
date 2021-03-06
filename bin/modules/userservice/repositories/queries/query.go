package queries

import (
	"context"
	models "engine/bin/modules/userservice/models/domain"
)

type UserServicePostgre interface {
	GetByServiceIDAndUserId(ctx context.Context, serviceId, userId string) (userService *models.UserService, err error)
}
