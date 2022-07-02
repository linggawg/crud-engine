package commands

import (
	"context"
	models "engine/bin/modules/users/models/domain"
)

type UsersPostgre interface {
	InsertOne(ctx context.Context, users *models.Users) error
}
