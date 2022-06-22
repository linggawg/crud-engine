package usecases

import (
	"context"
	"crud-engine/modules/userservice/repositories/queries"
	httpError "crud-engine/pkg/http-error"
	"crud-engine/pkg/utils"
)

type UserServiceQueryUsecase struct {
	UserServicePostgreQuery queries.UserServicePostgre
}

func NewQueryUsecase(UserServicePostgreQuery queries.UserServicePostgre) *UserServiceQueryUsecase {
	return &UserServiceQueryUsecase{
		UserServicePostgreQuery: UserServicePostgreQuery,
	}
}

func (u UserServiceQueryUsecase) GetByServiceIDAndUserId(ctx context.Context, serviceId, userId string) utils.Result {
	var result utils.Result

	dbs, err := u.UserServicePostgreQuery.GetByServiceIDAndUserId(ctx, serviceId, userId)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Data user service tidak ditemukan"
		result.Error = errObj
		return result
	}

	result.Data = dbs
	return result
}
