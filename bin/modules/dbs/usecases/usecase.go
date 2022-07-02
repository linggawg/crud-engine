package usecases

import (
	"context"
	"engine/bin/pkg/utils"
)

type QueryUsecase interface {
	GetDbsConnection(ctx context.Context, userId, method, serviceUrl string, isQuery bool) utils.Result
}
