package usecases

import (
	"context"
	"crud-engine/pkg/utils"
)

type QueryUsecase interface {
	GetByServiceIDAndUserId(ctx context.Context, serviceId, userId string) utils.Result
}
