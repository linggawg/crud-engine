package usecases

import (
	"context"
	models "engine/bin/modules/users-services/models/domain"
	"engine/bin/pkg/utils"
)

type CommandUsecase interface {
	InsertAllByServices(ctx context.Context, payload models.UsersServicesRequest) utils.Result
	DeleteByUsersIdAndServiceUrl(ctx context.Context, payload models.UsersServicesRequest) utils.Result
}
