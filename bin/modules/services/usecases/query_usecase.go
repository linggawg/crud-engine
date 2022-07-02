package usecases

import (
	"context"
	"engine/bin/modules/services/repositories/queries"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/utils"
)

type ServicesQueryUsecase struct {
	ServicesPostgreQuery queries.ServicesPostgre
}

func NewQueryUsecase(ServicesPostgreQuery queries.ServicesPostgre) *ServicesQueryUsecase {
	return &ServicesQueryUsecase{
		ServicesPostgreQuery: ServicesPostgreQuery,
	}
}

func (u ServicesQueryUsecase) GetByServiceUrlAndMethod(ctx context.Context, serviceUrl, method string) utils.Result {
	var result utils.Result

	service, err := u.ServicesPostgreQuery.GetByServiceUrlAndMethod(ctx, serviceUrl, method)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Data service url tidak ditemukan"
		result.Error = errObj
		return result
	}

	result.Data = service
	return result
}

func (u ServicesQueryUsecase) GetByServiceDefinitionAndMethod(ctx context.Context, serviceDefinition, method string) utils.Result {
	var result utils.Result

	service, err := u.ServicesPostgreQuery.GetByServiceDefinitionAndMethod(ctx, serviceDefinition, method)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Data service definition tidak ditemukan"
		result.Error = errObj
		return result
	}

	result.Data = service
	return result
}
