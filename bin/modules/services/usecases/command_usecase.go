package usecases

import (
	"context"
	models "engine/bin/modules/services/models/domain"
	"engine/bin/modules/services/repositories/commands"
	"engine/bin/modules/services/repositories/queries"
	httpError "engine/bin/pkg/http-error"
	"engine/bin/pkg/utils"
	"strings"
)

type ServicesCommandUsecase struct {
	ServicesPostgreCommand commands.ServicesPostgre
	ServicesPostgreQuery   queries.ServicesPostgre
}

func NewServicesCommandUsecase(ServicesPostgreCommand commands.ServicesPostgre, ServicesPostgreQuery queries.ServicesPostgre) *ServicesCommandUsecase {
	return &ServicesCommandUsecase{
		ServicesPostgreCommand: ServicesPostgreCommand,
		ServicesPostgreQuery:   ServicesPostgreQuery,
	}
}

func (s ServicesCommandUsecase) DeleteByServiceUrl(ctx context.Context, payload models.ServicesRequest) utils.Result {
	var result utils.Result
	if !strings.EqualFold(payload.Opts.RoleName, "admin") {
		errObj := httpError.NewUnauthorized()
		errObj.Message = "unauthorized access"
		result.Error = errObj
		return result
	}

	err := s.ServicesPostgreCommand.DeleteByServiceUrl(ctx, payload.ServiceUrl)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = err.Error()
		result.Error = errObj
		return result
	}
	return result
}
