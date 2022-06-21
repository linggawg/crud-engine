package usecases

import (
	"context"
	"crud-engine/pkg/utils"
)

type QueryUsecase interface {
	GetByServiceUrlAndMethod(ctx context.Context, serviceUrl, method string) utils.Result
	GetByServiceDefinitionAndMethod(ctx context.Context, serviceDefinition, method string) utils.Result
}
