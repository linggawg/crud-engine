package queries

import (
	"context"
	models "engine/bin/modules/roles/models/domain"
)

type RolesPostgre interface {
	FindOneByID(ctx context.Context, id string) (roles *models.Roles, err error)
}
