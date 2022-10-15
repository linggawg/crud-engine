package usecases

import (
	"context"
	models "engine/bin/modules/services/models/domain"
	"engine/bin/pkg/utils"
)

type CommandUsecase interface {
	DeleteByServiceUrl(ctx context.Context, payload models.ServicesRequest) utils.Result
}
