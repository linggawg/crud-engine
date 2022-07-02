package queries

import (
	"context"
	models "engine/bin/modules/users/models/domain"
)

type UsersPostgre interface {
	GetByID(ctx context.Context, id string) (users *models.Users, err error)
	GetByEmail(ctx context.Context, email string) (user *models.Users, err error)
}
