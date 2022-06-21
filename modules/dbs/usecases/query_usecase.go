package usecases

import (
	"context"
	"crud-engine/modules/dbs/repositories/queries"
	httpError "crud-engine/pkg/http-error"
	"crud-engine/pkg/utils"
)

type DbsQueryUsecase struct {
	DbsPostgreQuery queries.DbsPostgre
}

func NewQueryUsecase(DbsPostgreQuery queries.DbsPostgre) *DbsQueryUsecase {
	return &DbsQueryUsecase{
		DbsPostgreQuery: DbsPostgreQuery,
	}
}

func (u DbsQueryUsecase) GetByID(ctx context.Context, id string) utils.Result {
	var result utils.Result

	dbs, err := u.DbsPostgreQuery.GetByID(ctx, id)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = "Data dbs tidak ditemukan"
		result.Error = errObj
		return result
	}

	result.Data = dbs
	return result
}
