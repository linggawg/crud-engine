package queries

import (
	"context"
	models "engine/bin/modules/users/models/domain"
)

type UsersPostgre interface {
	FindOneByID(ctx context.Context, id string) (users *models.Users, err error)
	FindOneByUsername(ctx context.Context, username string) (users *models.Users, err error)
}
