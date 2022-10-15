package usecases

import (
	"context"
	"engine/bin/pkg/utils"
)

type QueryUsecase interface {
	GetDbsConnection(ctx context.Context, userId, method, serviceUrl, key string) utils.Result
}
