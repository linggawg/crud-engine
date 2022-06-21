package usecases

import (
	"context"
	"crud-engine/pkg/utils"
)

type QueryUsecase interface {
	GetByID(ctx context.Context, id string) utils.Result
}
